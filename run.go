package main

import (
	amcrest "amcrest/Amcrest"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func writeFrame(image []byte) {

	file, err := os.Create(time.Now().String() + ".jpg")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	// Write the []byte directly to the file
	_, err = file.Write(image)
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
	frame_saved := []byte{}

	threshold := 2000
	maxPercent := 5.5

	for i := 0; i < 100; i++ {

		frame_current, err := camera.GetSnapshot()

		if err != nil {
			log.Fatal(err)
		}

		if len(frame_saved) == 0 {
			fmt.Println("Saving first frame")
			frame_saved = frame_current
			continue
		}

		// Decode as JPEG
		saved, err := jpeg.Decode(bytes.NewReader(frame_saved))
		if err != nil {
			log.Fatalf("Error decoding saved JPEG image: %v", err)
		}
		current, err := jpeg.Decode(bytes.NewReader(frame_current))
		if err != nil {
			log.Fatalf("Error decoding current JPEG image: %v", err)
		}

		if Change(saved, current, uint32(threshold), float64(maxPercent)) {
			fmt.Printf("Images differ > %f percent\n", maxPercent)
			writeFrame(frame_current)
			frame_saved = frame_current

		}

	}

}

func Change(img1, img2 image.Image, threshold uint32, maxPercent float64) bool {

	if img1.Bounds() != img2.Bounds() {
		fmt.Println("not same bounds")
		return false
	}

	// Compare pictures by calculating the difference in channel averages between images
	bounds := img1.Bounds()
	changedPixels := 0
	totalPixels := (bounds.Max.X - bounds.Min.X) * (bounds.Max.Y - bounds.Min.Y)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {

			r1, g1, b1, a1 := img1.At(x, y).RGBA()
			r2, g2, b2, a2 := img2.At(x, y).RGBA()

			// Check if any channel difference exceeds the threshold
			if absDiff(r1, r2) > threshold || absDiff(g1, g2) > threshold || absDiff(b1, b2) > threshold || absDiff(a1, a2) > threshold {
				changedPixels++
			}
		}
	}
	percentageChanged := float64(changedPixels) / float64(totalPixels) * 100
	fmt.Printf("Percentage changed: %f\n", percentageChanged)
	return percentageChanged > float64(maxPercent)

}

func absDiff(a, b uint32) uint32 {
	if a > b {
		return a - b
	}
	return b - a
}
