package main

import (
	"fmt"
	"log"

	"{{module_path}}/internal/{{project_name}}"
)

func main() {
	fmt.Println("Welcome to {{project_name}}!")
	fmt.Println("Version:", {{project_name}}.Version)

	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	// Your code here
	fmt.Println("{{project_name}} is running successfully!")
	return nil
}
