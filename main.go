package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"time"

	"gocv.io/x/gocv"

	"github.com/benoitmasson/qrcode-demo/internal/decode"
	"github.com/benoitmasson/qrcode-demo/internal/detect"
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

		img, found, message := scanCode(&img, &imgWithMiniCode, &points, width, height)

		window.IMShow(img)
		if window.WaitKey(1) == 27 {
			break
		}

		if found {
			fmt.Printf("QR-code message is: '\033[1m%s\033[0m'\n", message)
			fmt.Println()
			time.Sleep(3 * time.Second)
		}
	}
}

// scanCode extracts the QR-code from the given image, then decodes it.
// If successful, returns a new image with miniature QR-code in the top-left corner and the message.
// Otherwise, returns the original image.
func scanCode(img, imgWithMiniCode *gocv.Mat, points *gocv.Mat, width, height int) (gocv.Mat, bool, string) {
	qrcodeDetector := gocv.NewQRCodeDetector()
	found := qrcodeDetector.Detect(*img, points) // false positives
	if !found {
		return *img, false, ""
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
		return *img, false, ""
	}
	fmt.Println("Points form a square, proceed")

	img.CopyTo(imgWithMiniCode)
	miniCode := setMiniCodeInCorner(imgWithMiniCode, imagePoints, miniCodeWidth, miniCodeHeight)
	detect.EnhanceImage(&miniCode)

	dots, ok := detect.GetDots(miniCode)
	miniCode.Close()
	if !ok {
		return *img, false, ""
	}
	fmt.Println("Dots scanned successfully, proceed")

	message, err := decodeMessage(dots)
	if err != nil {
		fmt.Printf("Dots do not form a valid QR-code: %v\n", err)
		return *img, false, ""
	}

	printQRCode(dots)
	outlineQRCode(imgWithMiniCode, imagePoints, color.RGBA{255, 0, 0, 255}, 5)

	return *imgWithMiniCode, true, message
}

// decodeMessage follows explanations from https://typefully.com/DanHollick/qr-codes-T7tLlNi
// to decode the QR-code contents.
func decodeMessage(dots detect.QRCode) (string, error) {
	if len(dots) < 17 {
		return "", errors.New("dots array too small")
	}

	_, err := decode.Version(dots)
	if err != nil {
		return "", err
	}
	maskID, errorCorrectionLevel, err := decode.Format(dots)
	if err != nil {
		return "", err
	}
	fmt.Printf("Mask ID is %d, error correction level is %d\n", maskID, errorCorrectionLevel)

	bits := decode.ReadBits(dots, maskID)
	_ = bits

	return "TODO", nil
}
