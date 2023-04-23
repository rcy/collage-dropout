package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
)

const (
	canvasWidth  = 800
	canvasHeight = 800
	maxImageSize = 200
	maxRotate    = 20
)

func main() {
	rand.Seed(time.Now().UnixNano())

	imagesDir := "./images"
	images := []string{}

	files, err := ioutil.ReadDir(imagesDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".jpg" ||
			filepath.Ext(file.Name()) == ".jpeg" ||
			filepath.Ext(file.Name()) == ".png" ||
			filepath.Ext(file.Name()) == ".gif" {
			images = append(images, filepath.Join(imagesDir, file.Name()))
		}
	}

	createCollage(images)
}

func createCollage(images []string) {
	canvas := imaging.New(canvasWidth, canvasHeight, color.White)

	for _, imagePath := range images {
		img, err := imaging.Open(imagePath)
		if err != nil {
			continue
		}

		aspectRatio := float64(img.Bounds().Size().X) / float64(img.Bounds().Size().Y)
		imageWidth := math.Min(float64(maxImageSize), float64(img.Bounds().Size().X))
		imageHeight := imageWidth / aspectRatio

		x := rand.Intn(canvasWidth - int(imageWidth))
		y := rand.Intn(canvasHeight - int(imageHeight))
		rotate := rand.Float64()*float64(maxRotate*2) - float64(maxRotate)

		img = imaging.Resize(img, int(imageWidth), int(imageHeight), imaging.Lanczos)
		img = imaging.Rotate(img, rotate, color.Transparent)

		canvas = imaging.Paste(canvas, img, image.Pt(x, y))
	}

	file, err := os.Create("collage.png")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	err = png.Encode(file, canvas)
	if err != nil {
		panic(err)
	}

	fmt.Println("Collage saved as collage.png")
}
