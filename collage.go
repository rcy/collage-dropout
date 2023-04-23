package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

func main() {
	dir := "images"
	output := "collage.png"
	width, height := 1600, 1200

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	rand.Seed(time.Now().UnixNano())
	dc := gg.NewContext(width, height)

	for _, file := range files {
		if !file.IsDir() {
			imgPath := filepath.Join(dir, file.Name())
			img, err := gg.LoadImage(imgPath)
			if err != nil {
				fmt.Println("Error loading image:", err)
				continue
			}

			rotateAngle := rand.Float64()*20 - 10
			rotatedImg := imaging.Rotate(img, rotateAngle, image.Transparent)

			x := rand.Float64() * float64(width-rotatedImg.Bounds().Dx())
			y := rand.Float64() * float64(height-rotatedImg.Bounds().Dy())

			dc.DrawImage(rotatedImg, int(x), int(y))
		}
	}

	err = dc.SavePNG(output)
	if err != nil {
		fmt.Println("Error saving collage image:", err)
		return
	}
	fmt.Println("Collage image saved as", output)
}
