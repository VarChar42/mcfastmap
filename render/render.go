package render

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"sort"

	"github.com/Tnze/go-mc/level"
	"github.com/Tnze/go-mc/level/block"
	"github.com/Tnze/go-mc/save"
	"github.com/Tnze/go-mc/save/region"
	"github.com/VarChar42/mcfastmap/blockcolor"
)

func Render(world string, xMca int, zMca int, schema blockcolor.ColorSchema) {

	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{512, 512}})

	regionFile := fmt.Sprintf("%s/region/r.%d.%d.mca", world, xMca, zMca)

	fmt.Printf("regionFile: %v\n", regionFile)

	r, err := region.Open(regionFile)
	if err != nil {
		//panic(err)
		return
	}
	defer r.Close()

	for xSector := 0; xSector < 32; xSector++ {
		for zSector := 0; zSector < 32; zSector++ {

			if !r.ExistSector(xSector, zSector) {
				continue
			}

			var c save.Chunk
			data, err := r.ReadSector(xSector, zSector)
			if err != nil {
				panic(err)
			}
			err = c.Load(data)
			if err != nil {
				panic(err)
			}

			x := int(xSector * 16)
			z := int(zSector * 16)

			colored := make([]bool, 32*32)
			done := false

			sort.Slice(c.Sections, func(i, j int) bool {
				return int8(c.Sections[i].Y) > int8(c.Sections[j].Y)
			})

			for _, section := range c.Sections {

				//fmt.Printf("section.Y: %v\n", section.Y)
				statePalette := section.BlockStates.Palette
				stateIdPalette := make([]block.StateID, len(statePalette))
				for i, v := range statePalette {
					b, found := block.FromID[v.Name]
					if !found {
						log.Printf("Cannot find block %v", v.Name)
						continue
					}
					stateIdPalette[i] = block.ToStateID[b]
				}

				states := level.NewStatesPaletteContainerWithData(16*16*16, section.BlockStates.Data, stateIdPalette)

				for cY := 15; cY >= 0; cY-- {
					done = true
					for i := 16*16 - 1; i >= 0; i-- {
						if colored[i] {
							continue
						}
						state := states.Get(cY*16*16 + i)
						blockState := block.StateList[state]

						if isAir(blockState) {
							done = false
							continue
						}

						xChunk := i % 16
						xBlock := x + xChunk

						//yBlock := int(section.Y)*16 + cY

						zChunk := i / 16
						zBlock := z + zChunk

						colored[i] = true
						bColor := schema[blockState.ID()]

						/*
							if xChunk == 0 && x != 0 || zChunk == 0 && z != 0 {
								bColor = color.RGBA{R: 255, B: 0, G: 0, A: 255}
							}
						*/

						img.Set(xBlock, zBlock, bColor)

					}
					if done {
						break
					}
				}
				if done {
					break
				}
			}
		}
	}

	pngFile, _ := os.Create(fmt.Sprintf("out/r.%d.%d.png", xMca, zMca))
	png.Encode(pngFile, img)
}

func isAir(s block.Block) bool {
	switch s.(type) {
	case block.Air, block.CaveAir, block.VoidAir:
		return true
	default:
		return false
	}
}
