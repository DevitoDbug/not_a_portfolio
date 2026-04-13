---
title: From TCP Streams to Redis Commands (Understanding RESP)
date: 8 April, 2026
---

Computers don’t understand “words”, only bytes
It is easy to forget this when writing code, but computers don't strings like: 
```text
SET name david
```

To us, this is a command. 
To a computer, its just a sequence of bytes.
The big question is, how does a "command" travel from the Redis CLI to the Redis server 
and end up stored in memory?

Just like many other technologies, Redis uses a client server architecture. In our case, the 
**Redis CLI** is the client and the **Redis Server** is the server. That means commands 
must travel over a network. But networks don't send words, they send bytes. 

This creates a huge problem:
How does the server know where a command starts, where it ends and what each part means?

Before Redis can actually store anything, it must first **parse** the incoming data into a structure
(usually a data type) it understands. In my case, building the redis server in **GO** meant converting 
the incoming byte stream into a Go struct.

Ladies and gentlemen, this is exactly where **RESP (Redis Serialization Protocol)** comes in.

# The core problem
TCP sends a continuous stream of bytes — there are no message boundaries. When data is sent over TCP, 
you have no guarantee that what you receive corresponds to what was sent in a single write. 
The data may arrive in chunks, be split arbitrarily, or even be combined with other messages.

For example, if you send:

```text
    Hello world
```

The receiver might get:
- "Hello world" (best case)
- "Hel", "lo ", "world"
- or even partial fragments at unpredictable intervals
This behavior is entirely controlled by the operating system and the network stack, completely unpredictable from our servers view.

This leads to a few fundamental problems:
- Where does one command end?
- How many arguments are there?
- What type is each value?


