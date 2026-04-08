package main

import (
	"fmt"

	"github.com/DevitoDbug/portfolio/internals/server"
)

func main() {
	port := ":8081"
	server := server.NewServer(port)

	err := server.StartServer()
	if err != nil {
		fmt.Printf("failed to start server. Error: %v\n", err)
	}

	// blogPath := "./internals/blogs/"
	// files, err := os.ReadDir(blogPath)
	// if err != nil {
	// 	fmt.Printf("Error reading file is: %v\n", err)
	// 	os.Exit(1)
	// }
	//
	// for _, file := range files {
	// 	fmt.Println(file.Name())
	//
	// 	filePath := blogPath + file.Name()
	// 	fileContent, err := os.ReadFile(filePath)
	// 	if err != nil {
	// 		fmt.Printf("error reading file %v. Error: %v\n", filePath, err)
	// 		os.Exit(1)
	// 	}
	//
	// 	fmt.Println("file content is: ")
	// 	fmt.Printf("%s\n", fileContent)
	//
	// 	var buf bytes.Buffer
	// 	thing := goldmark.New()
	// 	err = thing.Convert(fileContent, &buf)
	// 	if err != nil {
	// 		fmt.Printf("error converting md to html %v.\nError: %v\n", filePath, err)
	// 		os.Exit(1)
	// 	}
	// 	fmt.Println(buf.String())
	// }
}
