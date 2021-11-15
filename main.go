package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"log"
	"os"
)

/*
Algorithm

Use algorithm (ORB, SIFT, SURF) to extract key points and descriptors

For each image to test:
1. Resize to (1280, 720) --why?
2. Find score overlay
	2a. Use algo to generate key points and descriptors
	2b. Use matcher (flann or brute force) to determing matches
	2c. Lowe's ratio test
	2d. At least 9 matches for SURF, or image under test probs doesn't have scoreboard
	2e. Use the key points in the template can create a transform to the keypoints in the image estimateAffinePartial2D
3. Check if match is under review
4. Find time remaining
5. Determine match state ('pre_match', 'auto', 'teleop', or 'post_match')
6. Get Match Key Name
	6a. _getImgCropThresh => crop and threshold?
	6b. _parseRawMatchName => using tesseract
	6c. _getMatchKey => regex extract
 */

func run(file string) error {

	algo := gocv.NewORB() //nFeatures=1000

	overlay := gocv.IMRead("seasons/score_overlay_2020.png", gocv.IMReadColor)
	kp1, des1 := algo.DetectAndCompute(overlay, gocv.NewMat())

	matcher := gocv.NewBFMatcher() //cv2.NORM_HAMMING, crossCheck=True,

	img := gocv.IMRead("seasons/2020_InfiniteRecharge/frame0010-before.jpg", gocv.IMReadColor)
	kp2, des2 := algo.DetectAndCompute(img, gocv.NewMat())

	matches := matcher.KnnMatch(des1, des2, 4)

	//log.Printf("Matches %#v", matches)

	// Store all the good matches as per Lowe's ratio test
	goodMatches := make([]gocv.DMatch, 0)
	for _, submatches := range matches {
		if submatches[0].Distance < 0.7 * submatches[1].Distance {
			goodMatches = append(goodMatches,submatches[0] )
		}
	}
	const MinGoodMatches = 9
	//log.Printf("GoodMatches: %d", len(goodMatches))
	if len(goodMatches) < MinGoodMatches {
		return fmt.Errorf("only found %d matches, need at least %d", len(goodMatches), MinGoodMatches)
	}

	kp1vec := kp2Point2f(kp1)
	kp2vec := kp2Point2f(kp2)

	t := gocv.EstimateAffinePartial2D(kp1vec, kp2vec)

	log.Printf("T: %#v", t)

	//
	//scale := t.GetDoubleAt(0,0)
	//tx := t.GetDoubleAt(0,2)
	//ty := t.GetDoubleAt(1,2)
	//
	//log.Printf("Scale: %f Tx: %f Ty: %f", scale, tx, ty)

	return nil
}

func kp2Point2f(kp []gocv.KeyPoint) gocv.Point2fVector {
	kp2f := make([]gocv.Point2f, 0)
	for _, kp := range kp {
		kp2f = append(kp2f, gocv.Point2f{
			X: float32(kp.X),
			Y: float32(kp.Y),
		})
	}
	return gocv.NewPoint2fVectorFromPoints(kp2f)
}

func processVideoFile(videoFile string, frameCallback func(frame gocv.Mat, curPosMsec float64), debug bool) error {
	video, err := gocv.VideoCaptureFile(videoFile)
	if err != nil {
		return fmt.Errorf("unable to open %s - %w", videoFile, err)
	}
	defer video.Close()

	window := gocv.NewWindow("Window")
	defer window.Close()

	// storage locations for frames
	img := gocv.NewMat()
	defer img.Close()

	fps := video.Get(gocv.VideoCaptureFPS)

	for {
		if ok := video.Read(&img); !ok {
			return fmt.Errorf("device closed: %v\n", videoFile)
		}
		if img.Empty() {
			continue
		}
		//curFrame := video.Get(gocv.VideoCapturePosFrames)
		curPosMsec := video.Get(gocv.VideoCapturePosMsec)

		frameCallback(img, curPosMsec)

		if debug {
			window.IMShow(img)

			window.WaitKey(1000 / 60) // wait for a keypress before starting playing
		}

		video.Grab(int(fps * 30)) // skip ahead 30 seconds
	}
}

func main() {
	file := "seasons/2020_InfiniteRecharge/Semifinal 2 - 2020 Week 0--Z3Isqo7esc.mp4"

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
