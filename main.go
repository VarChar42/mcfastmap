package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/VarChar42/mcfastmap/blockcolor"
	"github.com/VarChar42/mcfastmap/render"
)

func main() {

	colorSchema := *blockcolor.LoadSchema("filtered.txt")

	world := "world"

	files, _ := os.ReadDir(world + "/region/")

	waitChan := make(chan struct{}, 20)

	var wg sync.WaitGroup

	for _, file := range files {
		data := strings.Split(file.Name(), ".")
		if len(data) != 4 || data[3] != "mca" {
			continue
		}

		x, _ := strconv.Atoi(data[1])
		z, _ := strconv.Atoi(data[2])

		waitChan <- struct{}{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			render.Render(world, x, z, colorSchema)
			<-waitChan
		}()

	}

	fmt.Printf("Waiting\n")
	wg.Wait()
}
