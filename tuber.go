package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/joho/godotenv"
	"tuber/yamldownloader"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	yamls, err := yamldownloader.FindLayer()

	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("kubectl", "apply", "-f", "-")
	// cmd := exec.Command("cat")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		lastIndex := len(yamls) - 1
		for i, yaml := range yamls {
			io.WriteString(stdin, yaml.Content)
			if i < lastIndex {
				io.WriteString(stdin, "---\n")
			}
		}
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", out)
}
