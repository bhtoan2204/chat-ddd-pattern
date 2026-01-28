package main

import (
	"fmt"
	"go-socket/scaffold/generator"
	"go-socket/scaffold/models"
	"log"
)

const API_SPEC_DIR = "scaffold/api"

func main() {
	apiSpec, err := models.LoadAPISpecDir(API_SPEC_DIR)
	if err != nil {
		log.Fatalf("Failed to load API spec: %v", err)
	}
	msg, err := generator.GenerateRequest(apiSpec.Endpoints)
	if err != nil {
		log.Fatalf("Failed to generate requests: %v", err)
	}
	fmt.Println(msg)
	msg, err = generator.GenerateResponse(apiSpec.Endpoints)
	if err != nil {
		log.Fatalf("Failed to generate responses: %v", err)
	}
	fmt.Println(msg)
	msg, err = generator.GenerateHandler(apiSpec.Endpoints)
	if err != nil {
		log.Fatalf("Failed to generate handlers: %v", err)
	}
	fmt.Println(msg)
	msg, err = generator.GenerateRouting(apiSpec)
	if err != nil {
		log.Fatalf("Failed to generate routing: %v", err)
	}
	fmt.Println(msg)
	msg, err = generator.GenerateRegistry(apiSpec)
	if err != nil {
		log.Fatalf("Failed to generate registry: %v", err)
	}
	fmt.Println(msg)
}
