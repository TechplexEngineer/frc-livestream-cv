package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"log"
	"os"
)

func run(file string) error {
	video, err := gocv.VideoCaptureFile(file)
	if err != nil {
		return fmt.Errorf("unable to open %s - %w", file, err)
	}
	defer video.Close()

	window := gocv.NewWindow("Window")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := video.Read(&img); !ok {
			return fmt.Errorf("device closed: %v\n", file)
		}
		if img.Empty() {
			continue
		}

		window.IMShow(img)

		window.WaitKey(1)// wait for a keypress before starting playing
	}


}

func main() {
	file := "InfiniteRecharge/Semifinal 2 - 2020 Week 0--Z3Isqo7esc.mp4"

	if err := run(file); err != nil {
		log.Printf("Error: %s", err)
		os.Exit(1)
	}
	//webcam, _ := gocv.VideoCaptureDevice(0)
	//window := gocv.NewWindow("Hello")
	//img := gocv.NewMat()
	//
	//for {
	//	webcam.Read(&img)
	//	window.IMShow(img)
	//	window.WaitKey(1)
	//}
}
