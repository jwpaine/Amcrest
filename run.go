package main

import (
	amcrest "amcrest/Amcrest"
	"fmt"
	_ "fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func saveFrame(camera *amcrest.Camera, filename string) {
	frame, err := camera.GetSnapshot()

	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	// Write the []byte directly to the file
	_, err = file.Write(frame)
	if err != nil {
		log.Fatalf("Error saving image to file: %v", err)
	}
}

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	// // Use the strings directly as raw keys
	uri := os.Getenv("URI")
	loginname := os.Getenv("LOGINNAME")
	password := os.Getenv("PASSWORD")

	camera := amcrest.Init(uri, loginname, password)

	for i := 0; i < 1; i++ {
		filename := fmt.Sprintf("frame_%d.jpg", i)
		saveFrame(camera, filename)
	}

}
