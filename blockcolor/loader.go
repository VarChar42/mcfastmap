package blockcolor

import (
	"bufio"
	"image/color"
	"os"
	"strconv"
	"strings"
)

type ColorSchema map[string]color.RGBA

func LoadSchema(path string) *ColorSchema {
	file, err := os.Open(path)

	if err != nil {
		return nil
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	schema := make(ColorSchema)

	for scanner.Scan() {
		data := strings.Split(scanner.Text(), " ")

		r, _ := strconv.Atoi(data[1])
		g, _ := strconv.Atoi(data[2])
		b, _ := strconv.Atoi(data[3])

		rgb := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}

		schema[data[0]] = rgb
	}

	return &schema
}
