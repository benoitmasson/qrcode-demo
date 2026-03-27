package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"math"

	"gocv.io/x/gocv"
)

func main() {
	// parse args
	var deviceID int
	flag.IntVar(&deviceID, "device-id", 0, "Webcam device ID, for image capture")
	flag.Parse()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		slog.Error(fmt.Sprintf("Error opening video capture device %d: %v", deviceID, err))
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
	slog.Info(fmt.Sprintf("Start reading device: %v", deviceID))
	for {
		if ok := webcam.Read(&img); !ok {
			slog.Error(fmt.Sprintf("Device closed: %v", deviceID))
			return
		}
		if img.Empty() {
			continue
		}

		if first {
			width = img.Cols()
			height = img.Rows()
			fps = int(math.Round(webcam.Get(gocv.VideoCaptureFPS)))
			slog.Info(fmt.Sprintf("[%s] %dx%d, %dfps", img.Type(), width, height, fps))
			first = false
		}

		img, found, message := scanCode(&img, &imgWithMiniCode, &points, width, height)

		window.IMShow(img)
		if window.WaitKey(1) == 27 {
			break
		}

		if found {
			slog.Warn(fmt.Sprintf("QR-code message is: '\033[1m%s\033[0m'", message))
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

	qrcodeDetector := gocv.NewQRCodeDetector()
	message = qrcodeDetector.DetectAndDecode(*img, points, imgWithMiniCode)
	found = message != ""

	if found {
		imagePoints = newImagePointsFromPoints(points)
		outlineQRCode(img, imagePoints, color.RGBA{255, 0, 0, 255}, 5)
	}

	return *img, found, message
}
