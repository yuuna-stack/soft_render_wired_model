package main

import (
	"log"
	"time"

	"github.com/softrender/model"
	"github.com/softrender/tgalib"
)

func abs(val int) int {
	if val >= 0 {
		return val
	} else {
		return val * -1
	}
}

func line(x0 int, y0 int, x1 int, y1 int, img *tgalib.TGAImage, color *tgalib.RGBColor) {
	steep := false
	if abs(x0-x1) < abs(y0-y1) {
		y0, x0 = x0, y0
		y1, x1 = x1, y1
		steep = true
	}
	if x0 > x1 {
		x1, x0 = x0, x1
		y1, y0 = y0, y1
	}
	dx := x1 - x0
	dy := y1 - y0
	derr := abs(dy) * 2
	err := 0
	y := y0
	for x := x0; x <= x1; x++ {
		if steep {
			tgalib.SetRGBColor(img, y, x, color)
		} else {
			tgalib.SetRGBColor(img, x, y, color)
		}
		err += derr
		if err > dx {
			if y1 > y0 {
				y += 1
			} else {
				y -= 1
			}
			err -= dx * 2
		}
	}
}

//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	// defer duration(track("main"))
	// flag.Parse()
	// if *cpuprofile != "" {
	// 	f, err := os.Create(*cpuprofile)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	pprof.StartCPUProfile(f)
	// 	defer pprof.StopCPUProfile()
	// }
	width, height := 800, 800
	white := tgalib.NewRGB(255, 255, 255)
	//red := tgalib.NewRGB(255, 0, 0)
	img := tgalib.NewTGAImage(width, height)
	model, err := model.ReadModel("obj/african_head.obj")
	if err != nil {
		log.Fatal(err.Error())
	}
	for i := 0; i < model.FacesCount(); i++ {
		face := model.GetFace(i)
		for j := 0; j < 3; j++ {
			v0 := model.GetVertex((*face)[j])
			v1 := model.GetVertex((*face)[(j+1)%3])
			x0 := int((v0.GetX() + 1.0) * float32(width) / 2.0)
			y0 := int((v0.GetY() + 1.0) * float32(height) / 2.0)
			x1 := int((v1.GetX() + 1.0) * float32(width) / 2.0)
			y1 := int((v1.GetY() + 1.0) * float32(height) / 2.0)
			line(x0, y0, x1, y1, img, white)
		}
	}
	img.FlipVertically()
	img.WriteTgaFile("output.tga", true)
}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}
