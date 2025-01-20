package main

import (
	amcrest "amcrest/Amcrest"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
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

	fmt.Printf("%s %s %s\n", uri, loginname, password)

	camera := amcrest.Init(uri, loginname, password)

	frame_changed := []byte{}

	for i := 0; i < 10; i++ {
		// filename := fmt.Sprintf("frame_%d.jpg", i)
		// saveFrame(camera, filename)

		frame, err := camera.GetSnapshot()

		if err != nil {
			log.Fatal(err)
		}

		if len(frame_changed) == 0 {
			fmt.Println("Saving first frame")
			frame_changed = frame
			continue
		}

		// Decode as JPEG
		saved, err := jpeg.Decode(bytes.NewReader(frame_changed))
		if err != nil {
			log.Fatalf("Error decoding saved JPEG image: %v", err)
		}
		current, err := jpeg.Decode(bytes.NewReader(frame))
		if err != nil {
			log.Fatalf("Error decoding current JPEG image: %v", err)
		}

		fmt.Println("Images successfully decoded as JPEG")

		// bounds := img.Bounds()
		// fmt.Printf("Image dimensions: %dx%d\n", bounds.Dx(), bounds.Dy())
		compareImages(saved, current)

	}

}

func compareImages(img1, img2 image.Image) {

	if img1.Bounds() != img2.Bounds() {
		fmt.Println("not same bounds")
	} else {
		fmt.Println("same bounds")
	}

	// Compare pictures by calculating histogram for each color value?
	bounds := img1.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			colors := img1.At(x, y)
			r1, g1, b1, a1 := colors.RGBA()

		}
	}

	// // Images are identical
	// return true
}
