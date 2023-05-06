package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"time"

	"gocv.io/x/gocv"

	"github.com/benoitmasson/qrcode-demo/detect"
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

		img, found := scanCode(&img, &points, width, height)

		window.IMShow(img)
		if window.WaitKey(1) == 27 {
			break
		}

		if found {
			time.Sleep(3 * time.Second)
		}
	}
}

func scanCode(img *gocv.Mat, points *gocv.Mat, width, height int) (gocv.Mat, bool) {
	newImg := img.Clone()
	qrcodeDetector := gocv.NewQRCodeDetector()
	found := qrcodeDetector.Detect(newImg, points) // false positives
	if !found {
		return *img, false
	}

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

	found = detect.ValidateSquare(imagePoints, width, height)
	if !found {
		return *img, false
	}
	fmt.Println("Points form a square, proceed")

	miniCode := setMiniCodeInCorner(&newImg, imagePoints, miniCodeWidth, miniCodeHeight)
	detect.EnhanceImage(&miniCode)

	dots, ok := detect.GetDots(miniCode)
	miniCode.Close()
	if !ok {
		return *img, false
	}
	fmt.Println("Dots scanned successfully, proceed")

	printQRCode(dots)
	outlineQRCode(&newImg, imagePoints, color.RGBA{255, 0, 0, 255}, 5)

	newImg.CopyTo(img)
	newImg.Close()

	return *img, true
}
