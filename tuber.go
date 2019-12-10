package main

import (
	"fmt"
	"log"

	"tuber/pkg/apply"
	"tuber/pkg/listen"
	"tuber/pkg/yamldownloader"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	listen.Listen(func(event *listen.RegistryEvent, err error) {
		spew.Dump(event)
	})

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
