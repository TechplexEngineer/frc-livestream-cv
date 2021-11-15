package main

import (
	"bytes"
	"fmt"
	"github.com/otiai10/gosseract/v2"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"gocv.io/x/gocv"
)

func run() error {
	var nfeatures = 1000                    // default is 500
	var scaleFactor float32 = 1.2           // default 1.2
	var nlevels = 8                         // default 8
	var edgeThreshold = 31                  // default 32
	var firstLevel = 0                      // default 0
	var WtaK = 2                            // default 2
	var scoreType = gocv.ORBScoreTypeHarris // default ORBScoreTypeHarris
	var patchSize = 31                      // default 31
	var fastThreshold = 20                  // default 20

	algo := gocv.NewORBWithParams(nfeatures, scaleFactor, nlevels, edgeThreshold, firstLevel, WtaK, scoreType, patchSize, fastThreshold)

	overlayFile := "score_overlay_2021_1280.png"
	frameFile := "frame-00570.jpg"

	matcher := gocv.NewBFMatcher()

	overlay := gocv.IMRead(overlayFile, gocv.IMReadColor)
	if overlay.Empty() {
		return fmt.Errorf("unable to load %s", overlayFile)
	}

	overlayW := overlay.Cols()
	//overlayH := overlay.Rows()

	img := gocv.IMRead(frameFile, gocv.IMReadColor)
	if img.Empty() {
		return fmt.Errorf("unable to load %s", frameFile)
	}

	newWidth := 1280.0
	newHeight := newWidth * (float64(img.Rows()) / float64(img.Cols()))

	gocv.Resize(img, &img, image.Pt(int(newWidth), int(newHeight)), 0, 0, gocv.InterpolationDefault)

	emptyMask := gocv.NewMat()
	kp1, des1 := algo.DetectAndCompute(overlay, emptyMask)
	kp2, des2 := algo.DetectAndCompute(img, emptyMask)

	matches := matcher.KnnMatch(des1, des2, 2)

	// Store all the good matches as per Lowe's ratio test
	goodMatches := make([]gocv.DMatch, 0)
	for _, submatches := range matches {
		if submatches[0].Distance < 0.75*submatches[1].Distance {
			goodMatches = append(goodMatches, submatches[0])
		}
	}
	const MinGoodMatches = 7
	if len(goodMatches) < MinGoodMatches {
		return fmt.Errorf("only found %d matches, need at least %d", len(goodMatches), MinGoodMatches)
	}

	// use goodMatches to build list of "good" keypoints
	kp1f := make([]gocv.Point2f, 0)
	for _, gm := range goodMatches {
		kp1f = append(kp1f, gocv.Point2f{
			X: float32(kp1[gm.QueryIdx].X),
			Y: float32(kp1[gm.QueryIdx].Y),
		})
	}
	kp1vec := gocv.NewPoint2fVectorFromPoints(kp1f)

	kp2f := make([]gocv.Point2f, 0)
	for _, gm := range goodMatches {
		kp2f = append(kp2f, gocv.Point2f{
			X: float32(kp2[gm.TrainIdx].X),
			Y: float32(kp2[gm.TrainIdx].Y),
		})
	}
	kp2vec := gocv.NewPoint2fVectorFromPoints(kp2f)

	log.Printf("kp1:%d kp2:%d", kp1vec.Size(), kp2vec.Size())

	data := gocv.EstimateAffinePartial2D(kp1vec, kp2vec)

	t := Transform{
		Scale: data.GetDoubleAt(0, 0),
		Tx:    data.GetDoubleAt(0, 2),
		Ty:    data.GetDoubleAt(1, 2),
	}

	log.Printf("Scale: %f Tx: %f Ty: %f", t.Scale, t.Tx, t.Ty)

	// outline whole image
	//drawBox(&img, image.Pt(0,0), image.Pt(w,h), color.RGBA{
	//	R: 255,
	//	G: 0,
	//	B: 0,
	//	A: 0,
	//})

	// Outline scoreboard in blue
	//drawBox(&img, t.TransformPt(image.Pt(0, 0)), t.TransformPt(image.Pt(overlayW, overlayH)), color.RGBA{
	//	R: 0,
	//	G: 0,
	//	B: 255,
	//	A: 0,
	//})

	// show review point 1
	//reviewPt1 := image.Pt(624, 93)
	//drawBox(&img, t.TransformPt(reviewPt1), t.TransformPt(image.Pt(reviewPt1.X+2, reviewPt1.Y+2)), color.RGBA{
	//	R: 255,
	//	G: 0,
	//	B: 255,
	//	A: 0,
	//})

	// show review point 2
	//reviewPt2 := image.Pt(1279-624, 93)
	//drawBox(&img, t.TransformPt(reviewPt2), t.TransformPt(image.Pt(reviewPt2.X+2, reviewPt2.Y+2)), color.RGBA{
	//	R: 255,
	//	G: 0,
	//	B: 255,
	//	A: 0,
	//})

	// outline time in green

	horizCenter := overlayW / 2
	timeTL := image.Pt(horizCenter-25, 57)
	timeBR := image.Pt(horizCenter+25, 85)

	//drawBox(&img, t.TransformPt(timeTL), t.TransformPt(timeBR), color.RGBA{
	//	R: 0,
	//	G: 255,
	//	B: 0,
	//	A: 0,
	//})

	// show mode point 1 (In the time bar)
	//modePt1 := image.Pt(520, 70)
	//drawBox(&img, t.TransformPt(modePt1), t.TransformPt(image.Pt(modePt1.X+2, modePt1.Y+2)), color.RGBA{
	//	R: 0,
	//	G: 255,
	//	B: 255,
	//	A: 0,
	//})

	// show mode point 2 (In the time bar)
	//modePt2 := image.Pt(581, 70)
	//drawBox(&img, t.TransformPt(modePt2), t.TransformPt(image.Pt(modePt2.X+2, modePt2.Y+2)), color.RGBA{
	//	R: 0,
	//	G: 255,
	//	B: 255,
	//	A: 0,
	//})

	matchNameTl := image.Pt(220, 6)
	matchNameBR := image.Pt(570, 43)
	//drawBox(&img, t.TransformPt(matchNameTl), t.TransformPt(matchNameBR), color.RGBA{
	//	R: 0,
	//	G: 0,
	//	B: 255,
	//	A: 0,
	//})

	thresh := GetImgCropThresh(img, t.TransformPt(timeTL), t.TransformPt(timeBR))
	ocrText, err := OCRImage(thresh)
	if err != nil {
		return err
	}
	log.Printf("Match Time: %s", ocrText)

	matchNameImg := GetImgCropThresh(img, t.TransformPt(matchNameTl), t.TransformPt(matchNameBR))
	matchNameRaw, err := OCRImage(matchNameImg)
	if err != nil {
		return err
	}
	log.Printf("Match Name: %s", matchNameRaw)

	win := gocv.NewWindow("Review")
	win.IMShow(img)
	win.WaitKey(-1)

	return nil
}

func OCRImage(thresh gocv.Mat) (string, error) {
	ocr := gosseract.NewClient()
	defer ocr.Close()

	img, err := thresh.ToImage()
	if err != nil {
		return "", err
	}
	buf := bytes.Buffer{}
	err = png.Encode(&buf, img)
	if err != nil {
		return "", err
	}

	err = ocr.SetImageFromBytes(buf.Bytes())
	if err != nil {
		return "", err
	}

	text, err := ocr.Text()
	if err != nil {
		return "", err
	}

	return text, nil
}

const OCRHeight = 64.0 // Do all OCR at this size

// prepare a crop of the image for OCR
func GetImgCropThresh(img gocv.Mat, tl, br image.Point) gocv.Mat {

	// crop
	cropped := img.Region(image.Rectangle{
		Min: tl,
		Max: br,
	})

	// scale to 64 pixels tall for ocr
	scale := OCRHeight / float64(cropped.Rows())

	resized := gocv.NewMat()
	gocv.Resize(cropped, &resized, image.Point{}, scale, scale, gocv.InterpolationDefault)

	// threshold the image
	bw := gocv.NewMat()
	gocv.CvtColor(resized, &bw, gocv.ColorBGRToGray)

	thresh := gocv.NewMat()
	gocv.Threshold(bw, &thresh, 120, 255, gocv.ThresholdBinary)

	return thresh

	//return bw
	//WhiteLow := gocv.Ones(1, 3, gocv.MatTypeCV8U) // np.array([120, 120, 120])
	//WhiteLow.MultiplyFloat(120)
	//WhiteHigh := gocv.Ones(1, 3, gocv.MatTypeCV8U) // np.array([255, 255, 255])
	//WhiteHigh.MultiplyFloat(255)
	//
	//BlackLow := gocv.Zeros(1, 3, gocv.MatTypeCV8U) // np.array([0, 0, 0])
	//BlackHigh := gocv.Ones(1, 3, gocv.MatTypeCV8U) // np.array([135, 135, 155])
	//BlackHigh.MultiplyFloat(135)

	//WhiteLow, err := gocv.NewMatFromBytes(1, 3, gocv.MatTypeCV8S, []byte{120, 120, 120})
	//if err != nil {
	//	panic(err)
	//}
	//WhiteHigh, err := gocv.NewMatFromBytes(1, 3, gocv.MatTypeCV8S, []byte{255, 255, 255})
	//if err != nil {
	//	panic(err)
	//}
	//BlackLow, err := gocv.NewMatFromBytes(1, 3, gocv.MatTypeCV8S, []byte{0, 0, 0})
	//if err != nil {
	//	panic(err)
	//}
	//BlackHigh, err := gocv.NewMatFromBytes(1, 3, gocv.MatTypeCV8S, []byte{135, 135, 155})
	//if err != nil {
	//	panic(err)
	//}
	//
	//lb := BlackLow
	//ub := BlackHigh
	//if darkBg {
	//	lb = WhiteLow
	//	ub = WhiteHigh
	//}
	//
	//kernel := gocv.Ones(3, 3, gocv.MatTypeCV8U)
	//
	////cv2.morphologyEx(
	////	cv2.inRange(resized, self._BLACK_LOW, self._BLACK_HIGH),
	////	cv2.MORPH_OPEN,
	////	self._morph_kernel)
	//inRange := gocv.NewMat()
	//gocv.InRange(resized, lb, ub, &inRange)
	//morphRes := gocv.NewMat()
	//gocv.MorphologyEx(inRange, &morphRes, gocv.MorphOpen, kernel)
	//
	//return morphRes
}

func drawBox(img *gocv.Mat, tl image.Point, br image.Point, clr color.RGBA) {
	pts := make([]image.Point, 0)

	pts = append(pts, tl)
	pts = append(pts, image.Pt(tl.X, br.Y))
	pts = append(pts, br)
	pts = append(pts, image.Pt(br.X, tl.Y))

	v := gocv.NewPointsVectorFromPoints([][]image.Point{pts})

	thickness := 2
	gocv.Polylines(img, v, true, clr, thickness)
}

type Transform struct {
	Scale float64
	Tx    float64
	Ty    float64
}

func (t Transform) TransformPt(pt image.Point) image.Point {
	x := float64(pt.X)*t.Scale + t.Tx
	y := float64(pt.Y)*t.Scale + t.Ty

	return image.Pt(int(x), int(y))
}

func main() {
	if err := run(); err != nil {
		log.Printf("Error: %s", err)
		os.Exit(1)
	}
}
