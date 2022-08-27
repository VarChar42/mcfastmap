package stitch

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strconv"
	"strings"
)

func Stitch(path string) { // TODO: needs cleanup

	padding := 5
	files, _ := os.ReadDir(path)

	minX := ^int(0)
	maxX := 0

	minZ := ^int(0)
	maxZ := 0

	for _, file := range files {
		data := strings.Split(file.Name(), ".")
		if len(data) != 4 || data[3] != "png" {
			continue
		}

		x, _ := strconv.Atoi(data[1])
		z, _ := strconv.Atoi(data[2])

		minX = min(x, minX)
		maxX = max(x, maxX)

		minZ = min(z, minZ)
		maxZ = max(z, maxZ)

	}

	x := (maxX-minX)*(512+padding) + 512
	y := (maxZ-minZ)*(512+padding) + 512

	total := image.Point{x, y}

	fmt.Printf("total: %v\n", total)

	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, total})

	for _, file := range files {
		data := strings.Split(file.Name(), ".")
		if len(data) != 4 || data[3] != "png" {
			continue
		}

		x, _ := strconv.Atoi(data[1])
		z, _ := strconv.Atoi(data[2])

		sourceFile, _ := os.Open(path + file.Name())
		source, _ := png.Decode(sourceFile)

		pX := (x - minX) * (512 + padding)
		pY := (z - minZ) * (512 + padding)

		rec := image.Rectangle{image.Point{pX, pY}, image.Point{pX + 512, pY + 512}}

		draw.Draw(img, rec, source, image.Point{0, 0}, draw.Over)

	}

	finalFile, _ := os.Create("final.png")

	png.Encode(finalFile, img)
	fmt.Printf("%v %v\n", x, y)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
