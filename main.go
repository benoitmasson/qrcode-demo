package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
)

func main() {
	// parse args
	var deviceID int
	flag.IntVar(&deviceID, "device-id", 0, "Webcam device ID, for image capture")
	flag.Parse()

	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening video capture device %d: %v\n", deviceID, err)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow("QR-code decoder")
	defer window.Close()

	// pre-allocate matrices once and for all
	img, imgWithMiniCode := gocv.NewMat(), gocv.NewMat()
	defer img.Close()
	defer imgWithMiniCode.Close()
	points := gocv.NewMat()
	defer points.Close()

	first := true
	var width, height, fps int
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
			fps = int(math.Round(webcam.Get(gocv.VideoCaptureFPS)))
			fmt.Printf("[%s] %dx%d, %dfps\n", img.Type(), width, height, fps)
			first = false
		}

		img, found, message := scanCode(&img, &imgWithMiniCode, &points, width, height)

		window.IMShow(img)
		if window.WaitKey(1) == 27 {
			break
		}

		if found {
			fmt.Printf("QR-code message is: '\033[1m%s\033[0m'\n", message)
			fmt.Println()
			webcam.Grab(3 * fps) // drop frames and sleep for 3 seconds
		}
	}
}

// scanCode extracts the QR-code from the given image, then decodes it.
// If successful, returns a new image with miniature QR-code in the top-left corner and the message.
// Otherwise, returns the original image.
func scanCode(img, imgWithMiniCode *gocv.Mat, points *gocv.Mat, width, height int) (gocv.Mat, bool, string) {
	var (
		imagePoints []image.Point
		message     string
		found       bool
	)

	// TODO: decode QR-code

	if found {
		imagePoints = newImagePointsFromPoints(points)
		outlineQRCode(img, imagePoints, color.RGBA{255, 0, 0, 255}, 5)
	}

	return *img, found, message
}

func newImagePointsFromPoints(points *gocv.Mat) []image.Point {
	r, c := points.Rows(), points.Cols()
	// fmt.Println(points.Channels(), points.Size(), points.Type(), points.Total(), r, c)

	imagePoints := make([]image.Point, 0, r*c)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			vec := points.GetVecfAt(i, j)
			x, y := vec[0], vec[1]

			imagePoints = append(imagePoints, image.Point{
				X: int(x),
				Y: int(y),
			})
		}
	}
	return imagePoints
}
