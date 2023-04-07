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

	window := gocv.NewWindow("Capture Window")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	points := gocv.NewMat()
	qrcode := gocv.NewMat()

	first := true
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
			// fmt.Println(img.Channels(), img.Size(), img.Type(), img.Total(), img.Rows(), img.Cols())
			first = false
		}

		qrcodeDetector := gocv.NewQRCodeDetector()
		// found := qrcodeDetector.Detect(img, &points) // false positives
		// => https://docs.opencv.org/4.x/d9/d61/tutorial_py_morphological_ops.html to clean image (open)
		// => gocv.GetPerspectiveTransform() (https://www.projectpro.io/recipes/what-are-warpaffine-and-warpperspective-opencv) puis WarpPerspective() (https://docs.opencv.org/4.x/da/d6e/tutorial_py_geometric_transformations.html) to transform (rotate/scale/translate)
		code := qrcodeDetector.DetectAndDecode(img, &points, &qrcode)
		if code != "" {
			// https://docs.opencv.org/4.x/de/dc3/classcv_1_1QRCodeDetector.html#a64373f7d877d27473f64fe04bb57d22b
			// fmt.Println()
			// fmt.Println(code, qrcode.Channels(), qrcode.Size(), qrcode.Type(), qrcode.Rows(), qrcode.Cols())
			printQRCode(qrcode)
			fmt.Println("=>", code, "<=")

			r, c := points.Rows(), points.Cols()
			// fmt.Println(points.Channels(), points.Size(), points.Type(), points.Total(), r, c)
			imagePoints := make([]image.Point, 0, r*c)
			for i := 0; i < r; i++ {
				for j := 0; j < c; j++ {
					vec := points.GetVecfAt(i, j)
					x, y := vec[0], vec[1]
					fmt.Printf("Point %d: (%f, %f)\n", i*r+j, x, y)
					imagePoints = append(imagePoints, image.Point{
						X: int(x),
						Y: int(y),
					})
				}
			}
			fmt.Println()

			gocv.Line(&img, imagePoints[0], imagePoints[1], color.RGBA{255, 0, 0, 255}, 5)
			gocv.Line(&img, imagePoints[1], imagePoints[2], color.RGBA{255, 0, 0, 255}, 5)
			gocv.Line(&img, imagePoints[2], imagePoints[3], color.RGBA{255, 0, 0, 255}, 5)
			gocv.Line(&img, imagePoints[3], imagePoints[0], color.RGBA{255, 0, 0, 255}, 5)

			points = gocv.NewMat()
			qrcode = gocv.NewMat()
		}

		window.IMShow(img)
		if window.WaitKey(1) == 27 {
			break
		}

		if code != "" {
			time.Sleep(5 * time.Second)
			// break
		}
	}
}

func printQRCode(qrcode gocv.Mat) {
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
