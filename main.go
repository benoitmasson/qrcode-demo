package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strings"
	"time"

	"gocv.io/x/gocv"
)

const (
	miniCodeWidth  = 200
	miniCodeHeight = 200
)

func main() {
	// parse args
	deviceID := "0"
	if len(os.Args) >= 2 {
		deviceID = os.Args[1]
	}

	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow("QR-code decoder")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	points := gocv.NewMat()
	defer points.Close()

	first := true
	var width, height int
	fmt.Printf("Start reading device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		if first {
			width = img.Cols()
			height = img.Rows()
			fmt.Printf("[%s] %dx%d\n", img.Type(), width, height)
			first = false
		}

		qrcodeDetector := gocv.NewQRCodeDetector()
		found := qrcodeDetector.Detect(img, &points) // false positives
		if found {
			r, c := points.Rows(), points.Cols()
			// fmt.Println(points.Channels(), points.Size(), points.Type(), points.Total(), r, c)
			imagePoints := make([]image.Point, 0, r*c)
			for i := 0; i < r; i++ {
				for j := 0; j < c; j++ {
					vec := points.GetVecfAt(i, j)
					x, y := vec[0], vec[1]
					// fmt.Printf("Point %d: (%f, %f)\n", i*r+j, x, y)
					imagePoints = append(imagePoints, image.Point{
						X: int(x),
						Y: int(y),
					})
				}
			}

			found = validateSquare(imagePoints, width, height)
			if found {
				fmt.Println("Points form a square, proceed")
				miniCode := setMiniCodeInCorner(&img, imagePoints, miniCodeWidth, miniCodeHeight)
				enhanceImage(&miniCode)

				firstColumn := getFirstColumnSequences(miniCode)
				fmt.Printf("First column: %v\n", firstColumn)
				if len(firstColumn) >= 3 && len(firstColumn)%2 == 1 && // column starts and ends with 2 different black sequences
					nearlyEquals(firstColumn[0], firstColumn[len(firstColumn)-1], 5) { // column should start and end with similar black blocks
					scale := float64(firstColumn[0]) / 7. // a marker is 7 dots high
					fmt.Printf("Dots are %f pixels wide\n", scale)
					// TODO: get dots
					miniCode.Close()

					fmt.Println("Dots scanned successfully, proceed")
					outlineQRCode(&img, imagePoints, color.RGBA{255, 0, 0, 255}, 5)
				} else {
					found = false
				}
			} // else {
			// 	fmt.Println("Inconsistent points, discard")
			// }

			// fmt.Println()
		}

		window.IMShow(img)
		if window.WaitKey(1) == 27 {
			break
		}

		if found {
			time.Sleep(3 * time.Second)
			// break
		}
	}
}

func validateSquare(points []image.Point, width, height int) bool {
	if len(points) != 4 {
		return false
	}

	// invalid coordinates
	if points[0].X < 0 || points[0].X >= width || points[0].Y < 0 || points[0].Y >= height ||
		points[1].X < 0 || points[1].X >= width || points[1].Y < 0 || points[1].Y >= height ||
		points[2].X < 0 || points[2].X >= width || points[2].Y < 0 || points[2].Y >= height ||
		points[3].X < 0 || points[3].X >= width || points[3].Y < 0 || points[3].Y >= height {
		return false
	}

	tolerance := height / 12
	// detected points are ordered as follows: bottom-left first, then bottom-right, then top-right, then top-left
	// make sure they form something which looks like an unrotated square
	if !nearlyEquals(points[0].Y, points[1].Y, tolerance) ||
		!nearlyEquals(points[2].Y, points[3].Y, tolerance) ||
		!nearlyEquals(points[0].X, points[3].X, tolerance) ||
		!nearlyEquals(points[1].X, points[2].X, tolerance) ||
		!nearlyEquals(points[1].X-points[0].X, points[3].Y-points[0].Y, tolerance) { // not a square
		// not a horizontal square, check points are ordered as follows: bottom-right, then top-right, top-left, bottom-left
		if !nearlyEquals(points[0].X, points[1].X, tolerance) ||
			!nearlyEquals(points[2].X, points[3].X, tolerance) ||
			!nearlyEquals(points[0].Y, points[3].Y, tolerance) ||
			!nearlyEquals(points[1].Y, points[2].Y, tolerance) ||
			!nearlyEquals(points[1].Y-points[0].Y, points[0].X-points[3].X, tolerance) { // not a square
			// not a horizontal nor vertical square
			return false
		}
	}

	return true
}

func nearlyEquals(i, j, tolerance int) bool {
	return j-i >= -tolerance && j-i <= tolerance
}

func setMiniCodeInCorner(img *gocv.Mat, points []image.Point, width, height int) gocv.Mat {
	originVector := gocv.NewPointVectorFromPoints(points)
	defer originVector.Close()
	destinationVector := gocv.NewPointVectorFromPoints([]image.Point{
		{X: 0, Y: 0},
		{X: width - 1, Y: 0},
		{X: width - 1, Y: height - 1},
		{X: 0, Y: height - 1},
	})
	defer destinationVector.Close()
	transform := gocv.GetPerspectiveTransform(originVector, destinationVector)
	defer transform.Close()
	rectangle := (*img).Region(image.Rect(0, 0, width-1, height-1))

	gocv.WarpPerspective(*img, &rectangle, transform, image.Point{X: width - 1, Y: height - 1})
	return rectangle
}

func enhanceImage(img *gocv.Mat) {
	// "open" to clean image: https://docs.opencv.org/4.x/d9/d61/tutorial_py_morphological_ops.html
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{X: 3, Y: 3})
	defer kernel.Close()
	gocv.MorphologyEx(*img, img, gocv.MorphOpen, kernel)
	// increase contrast
	gocv.AddWeighted(*img, 1.5, *img, 0, 0, img)
}

const (
	luminosityThreshold         = 500
	luminosityPersistenceOffset = 0 // used to favor sequence color persistence if in doubt
)

func getFirstColumnSequences(img gocv.Mat) []int {
	sequences := make([]int, 0)
	const offset = 2

	isFirstSequence := true
	currentSequenceIsBlack := false
	currentSequenceLength := 0

	for row := offset; row < img.Rows()-offset; row++ {
		vec := img.GetVecbAt(row, offset)
		pixelLuminosity := int(vec[0]) + int(vec[1]) + int(vec[2])

		threshold := luminosityThreshold
		if currentSequenceIsBlack {
			threshold += luminosityPersistenceOffset
		} else {
			threshold -= luminosityPersistenceOffset
		}
		pixelIsBlack := true
		if pixelLuminosity >= threshold {
			pixelIsBlack = false
		}

		// fmt.Printf("[row %d] lum: %v / black: %v / sequence: %v\n", y, pixelLuminosity, pixelIsBlack, currentSequenceLength)
		if pixelIsBlack == currentSequenceIsBlack {
			currentSequenceLength++
		} else {
			if isFirstSequence && !currentSequenceIsBlack {
				isFirstSequence = false
			} else {
				sequences = append(sequences, currentSequenceLength)
			}
			currentSequenceIsBlack = !currentSequenceIsBlack
			currentSequenceLength = 1
		}
	}
	if !isFirstSequence {
		sequences = append(sequences, currentSequenceLength)
	}

	return sequences
}

func outlineQRCode(img *gocv.Mat, points []image.Point, color color.RGBA, width int) {
	gocv.Line(img, points[0], points[1], color, width)
	gocv.Line(img, points[1], points[2], color, width)
	gocv.Line(img, points[2], points[3], color, width)
	gocv.Line(img, points[3], points[0], color, width)
}

func printQRCodeMat(qrcode gocv.Mat) {
	fmt.Println("\033[7;m", strings.Repeat(" ", 2*qrcode.Cols()+3)) // turn on inverse mode, start with blank line

	for i := 0; i < qrcode.Rows(); i++ {
		fmt.Print("  ") // start line with blank characters
		for j := 0; j < qrcode.Cols(); j++ {
			val := qrcode.GetUCharAt(i, j)
			char := " " // blank, displayed white
			if val == 0 {
				char = "█" // filled, displayed black
			}
			fmt.Print(char, char) // double print to achieve 1:1 scale
		}
		fmt.Println("  ") // end line with blank characters
	}

	fmt.Println(strings.Repeat(" ", 2*qrcode.Cols()+3), "\033[0;m") // end with blank line, turn off inverse mode
}
