package main

import (
	"fmt"
	"log"

	"tuber/pkg/apply"
	"tuber/pkg/yamldownloader"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println("no .env file found")
	}

	yamls, err := yamldownloader.FindLayer()

	if err != nil {
		log.Fatal(err)
	}

	out, err := apply.Apply(yamls)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", out)
}
