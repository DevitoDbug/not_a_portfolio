---
title: From TCP Streams to Redis Commands (Understanding RESP)
date: 8 April, 2026
---

Computers don't understand "words", only bytes.

It is easy to forget this when writing code, but to a computer, a command like:

```text
SET name david
```

is just a sequence of bytes. The big question is: how does a command travel from the Redis CLI to the Redis server and end up stored in memory?

Redis uses a client-server architecture. The **Redis CLI** is the client, the **Redis Server** is the server, and commands must travel between them over a network. But networks don't send words — they send bytes. This creates a fundamental problem: how does the server know where a command starts, where it ends, and what each part means?

Before Redis can store anything, it must first **parse** the incoming bytes into a structure it understands. When I built my own Redis server in Go, that meant converting a raw byte stream into a Go struct.

This is exactly where **RESP (Redis Serialization Protocol)** comes in.

# The Core Problem

TCP sends a continuous stream of bytes with no message boundaries. When data travels over TCP, you have no guarantee that what you receive corresponds to what was sent in a single write. Data may arrive in chunks, be split arbitrarily, or even be combined with other messages.

Send this:

```text
Hello world
```

The receiver might get:
- `"Hello world"` — all at once (best case)
- `"Hel"`, `"lo "`, `"world"` — in three separate reads
- or arbitrary fragments at unpredictable intervals

The operating system and network stack control this entirely. From your server's perspective, it is completely unpredictable.

This creates three fundamental questions:
1. Where does one command end and the next begin?
2. How many arguments does a command have?
3. What type is each value?

# Why Serialization

Serialization is the process of turning structured data into bytes so it can travel over a network. The important part that's easy to miss: **you must also be able to reverse it**. Whatever you send must be reconstructable on the other side with the exact same meaning — this is called deserialization.

Consider sending two commands back-to-back:

```text
SET name David
SET age 20
```

Both get converted to bytes and sent together. The server receives a blob of bytes and must figure out: where does the first command end? Where does the second begin? Without a clear answer, the server might try to interpret `David SET age` as a single command.

We need a contract. A set of rules that both the client and server agree on, so the sender can encode structure into the bytes, and the receiver can reliably decode it.

That contract is RESP.

# Building Intuition Step by Step

Let's start with a simple command:

```text
SET name david
```

**Step 1 — Raw text (ambiguous)**

The server receives bytes. It has no idea where `SET` ends and `name` begins, or whether `david` is one argument or two. Spaces and newlines don't mean anything at the byte level.

**Step 2 — Represent it as an array**

We know this is a list of three strings: `["SET", "name", "david"]`. If only we could communicate that structure explicitly.

**Step 3 — Add lengths**

What if we told the server the length of each part before sending it? `SET` is 3 bytes, `name` is 4 bytes, `david` is 5 bytes. Now the server knows exactly how many bytes to read for each piece.

**Step 4 — RESP format**

RESP formalizes this idea. Here's what our command looks like in RESP:

```text
*3
$3 SET
$4 name
$5 david
```

- `*3` means "this is an array of 3 elements"
- `$3` means "the next element is a string of 3 bytes"

**Step 5 — Wire format**

RESP uses `\r\n` (a carriage return followed by a newline) as a delimiter between parts. The actual bytes sent over the wire look like this:

```text
*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$5\r\ndavid\r\n
```

The server reads `*3\r\n`, knows to expect 3 elements, then reads each `$<length>\r\n<data>\r\n` block in sequence. No ambiguity.

# RESP Data Types

RESP uses the **first byte** of each value to communicate its type. Five types cover everything Redis needs:

| Prefix | Type          | Example                  |
|--------|---------------|--------------------------|
| `*`    | Array         | `*3\r\n...`              |
| `$`    | Bulk String   | `$5\r\nhello\r\n`        |
| `+`    | Simple String | `+OK\r\n`                |
| `:`    | Integer       | `:42\r\n`                |
| `-`    | Error         | `-ERR unknown command\r\n` |

**Arrays (`*`)** are how commands are sent. Every Redis command arrives as an array of bulk strings.

**Bulk Strings (`$`)** carry actual data. They always include the byte length upfront, so the parser knows exactly how many bytes to read — no guessing.

**Simple Strings (`+`)** are short, length-free strings. They can't contain `\r\n` inside them, which is fine because they're only used for simple fixed responses.

**Integers (`:`)** send a number directly in the text, terminated by `\r\n`.

**Errors (`-`)** look like simple strings but signal that something went wrong.

# Why Two String Types?

It might seem redundant to have both Simple Strings and Bulk Strings. They serve different roles:

**Simple Strings** are for control messages — fast, minimal, and predictable. When you run `SET name david`, the server responds with `+OK\r\n`. It's always `OK`, always the same two bytes, so there's no need for a length prefix. The constraint is that they can't contain `\r\n` inside them (because `\r\n` is how the parser knows the string is done).

**Bulk Strings** are for actual data. When you run `GET name`, the server returns `$5\r\ndavid\r\n`. The length prefix `$5` lets the parser read exactly 5 bytes, meaning the value could be anything — spaces, newlines, binary data — and it would still parse correctly.

**Simple = control messages. Bulk = data.**

# How Parsing Works

The parser reads one byte, uses it to decide the type, then delegates to the appropriate handler. Here's the mental model:

```text
Read first byte:
  '*' → read count, parse that many elements recursively
  '$' → read length, read that many bytes
  '+' → read until \r\n
  ':' → read until \r\n, convert to integer
  '-' → read until \r\n, treat as error
```

Walking through our example:

1. Read `*` → this is an array
2. Read `3\r\n` → expect 3 elements
3. Read `$` → next element is a bulk string
4. Read `3\r\n` → it's 3 bytes long
5. Read `SET\r\n` → first element is `"SET"`
6. Repeat for `name` and `david`

At the end, you have a structured command: `["SET", "name", "david"]`.

# Common Pitfalls

**Assuming TCP is message-based.** It isn't. A single `read()` call might return half a command, or two commands merged together. Your parser must handle this. Never assume one network read equals one message.

**Partial reads.** If you expect 5 bytes and only 3 arrive, you must buffer the data and wait for the rest. Ignoring this causes corrupted parses.

**CRLF mistakes.** RESP uses `\r\n`, not just `\n`. Searching for only `\n` will include the `\r` in your data, which will silently corrupt everything.

**Off-by-one errors.** If `$5` says 5 bytes, read exactly 5 bytes — then consume the trailing `\r\n` separately. Reading 7 bytes to "include the delimiter" shifts every subsequent read and breaks the parser.

# The Mental Model

- **TCP** is a raw byte stream with no structure
- **RESP** is a set of rules that impose structure on that stream
- **The parser** reads those rules and reconstructs the original command

RESP solves exactly the problems TCP creates: it defines boundaries, communicates types, and makes the format self-describing. The server never has to guess — every byte has a defined meaning based on what came before it.

# Conclusion

What starts as a simple `SET name david` in your terminal goes through a precise journey:

1. Split into an array of strings
2. Each string gets prefixed with its type and length
3. Parts are joined with `\r\n` delimiters
4. Bytes travel over TCP
5. The server's parser reads the first byte, determines the type, and reconstructs the command

The design is intentional. RESP is simple to implement, easy to debug (it's human-readable), and fast to parse. No external dependencies, no complex encoding — just a few prefixes and delimiters.

If you're building any kind of networked system that needs a custom protocol, this is the kind of thinking you'll need: not "how do I send data", but "how does the receiver know what they're looking at?"

Try building a RESP parser yourself. It's one of those exercises that makes TCP — and protocols in general — click in a way that reading about them never quite does.
