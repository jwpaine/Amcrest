package main

import (
	amcrest "amcrest/Amcrest"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	// Use the strings directly as raw keys
	uri := os.Getenv("URI")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	camera := amcrest.Init(uri)

	err = camera.LoadAuth(username, password)

	if err != nil {
		log.Fatal(err)
	}

	frame, err := camera.GetSnapshot()

	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("./test.jpg")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	// Wrap the []byte in an io.Reader using bytes.NewReader
	reader := bytes.NewReader(frame)

	// Write response body to file
	_, err = io.Copy(file, reader)
	if err != nil {
		log.Fatalf("Error saving response to file: %v", err)
	}

	fmt.Println("Successfully saved file!")

}
