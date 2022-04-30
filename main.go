package main

import (
	"github.com/3n0ugh/gen-populus/generator"
	"log"
)

func main() {
	err := generator.Generate(100)
	if err != nil {
		log.Println(err)
	}
}
