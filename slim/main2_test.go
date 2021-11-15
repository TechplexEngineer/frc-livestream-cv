package main

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
	"gocv.io/x/gocv"
)

//func Test_GetImgCropThresh(t *testing.T) {
//	overlayFile := "score_overlay_2021_1280.png"
//
//	overlay := gocv.IMRead(overlayFile, gocv.IMReadColor)
//	if overlay.Empty() {
//		t.Fatalf("unable to load %s", overlayFile)
//	}
//
//	//tl := image.Pt(0, 0)
//	//br := image.Pt(200, 100)
//
//	win := gocv.NewWindow("Review")
//	win.IMShow(morphRes)
//	win.WaitKey(-1)
//
//	//GetImgCropThresh(overlay, image.Pt(0,0), image.Pt(100,100))
//
//}

func Test(t *testing.T) {
	is := is.New(t)
	//m := gocv.NewMatWithSizesWithScalar([]int{2, 2}, gocv.MatTypeCV8S, gocv.NewScalar(1, 1, 1, 1))
	//log.Printf("%#v", m)
	//
	//WhiteLow := gocv.Ones(1, 3, gocv.MatTypeCV8U) // np.array([120, 120, 120])
	//WhiteLow.MultiplyFloat(120)
	//WhiteHigh := gocv.Ones(1, 3, gocv.MatTypeCV8U) // np.array([255, 255, 255])
	//WhiteHigh.MultiplyFloat(255)
	//
	//BlackLow := gocv.Zeros(1, 3, gocv.MatTypeCV8U) // np.array([0, 0, 0])
	//BlackHigh := gocv.Ones(1, 3, gocv.MatTypeCV8U) // np.array([135, 135, 155])
	//BlackHigh.MultiplyFloat(135)

	mats := make(map[string]gocv.Mat, 0)
	WhiteLow, err := gocv.NewMatFromBytes(1, 3, gocv.MatTypeCV8S, []byte{120, 120, 120})
	is.NoErr(err)
	mats["WhiteLow"] = WhiteLow

	WhiteHigh, err := gocv.NewMatFromBytes(1, 3, gocv.MatTypeCV8S, []byte{255, 255, 255})
	is.NoErr(err)
	mats["WhiteHigh"] = WhiteHigh

	BlackLow, err := gocv.NewMatFromBytes(1, 3, gocv.MatTypeCV8S, []byte{0, 0, 0})
	is.NoErr(err)
	mats["BlackLow"] = BlackLow

	BlackHigh, err := gocv.NewMatFromBytes(1, 3, gocv.MatTypeCV8S, []byte{135, 135, 155})
	is.NoErr(err)
	mats["BlackHigh"] = BlackHigh

	test, err := gocv.NewMatFromBytes(1, 3, gocv.MatTypeCV8S, []byte{1, 2, 3})
	is.NoErr(err)
	mats["test"] = test

	for k, img := range mats {
		fmt.Println(k)
		fmt.Printf("[\n")
		for r := 0; r < img.Rows(); r++ {
			fmt.Printf("\t[")
			for c := 0; c < img.Cols(); c++ {
				fmt.Printf("%v, ", img.GetVecbAt(r, c))
			}
			fmt.Printf("]\n")
		}
		fmt.Printf("]\n\n")
	}

}
