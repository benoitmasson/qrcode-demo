package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"

	"github.com/benoitmasson/qrcode-demo/internal/decode"
	"github.com/benoitmasson/qrcode-demo/internal/detect"
	"github.com/benoitmasson/qrcode-demo/internal/extract"
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

const (
	miniCodeWidth  = 200
	miniCodeHeight = 200
)

// scanCode extracts the QR-code from the given image, then decodes it.
// If successful, returns a new image with miniature QR-code in the top-left corner and the message.
// Otherwise, returns the original image.
func scanCode(img, imgWithMiniCode *gocv.Mat, points *gocv.Mat, width, height int) (gocv.Mat, bool, string) {
	qrcodeDetector := gocv.NewQRCodeDetector()
	found := qrcodeDetector.Detect(*img, points) // false positives
	if !found {
		return *img, false, ""
	}

	imagePoints := newImagePointsFromPoints(points)

	found = detect.ValidateSquare(imagePoints, width, height)
	if !found {
		return *img, false, ""
	}
	// fmt.Println("Points form a square, proceed")

	img.CopyTo(imgWithMiniCode)
	miniCode := detect.SetMiniCodeInCorner(imgWithMiniCode, imagePoints, miniCodeWidth, miniCodeHeight)
	detect.EnhanceImage(&miniCode)

	dots, ok := detect.GetDots(miniCode)
	miniCode.Close()
	if !ok {
		return *img, false, ""
	}
	fmt.Println("Dots scanned successfully, proceed")

	bits, version, errorCorrectionLevel, err := extractBits(dots)
	if err != nil {
		fmt.Printf("Dots do not form a valid QR-code: %v\n", err)
		return *img, false, ""
	}
	fmt.Println("Bits extracted successfully, proceed")

	message, err := decodeMessage(bits, version, errorCorrectionLevel)
	if err != nil {
		fmt.Printf("QR-code cannot be decoded: %v\n", err)
		return *img, false, ""
	}

	printQRCode(dots)
	detect.OutlineQRCode(imgWithMiniCode, imagePoints, color.RGBA{255, 0, 0, 255}, 5)

	return *imgWithMiniCode, true, message
}

// extractBits follows explanations from https://typefully.com/DanHollick/qr-codes-T7tLlNi
// to extract the QR-code bits from the 2D dots grid.
func extractBits(dots detect.QRCode) ([]bool, uint, decode.ErrorCorrectionLevel, error) {
	if len(dots) < 17 {
		return nil, 0, 0, errors.New("dots array too small")
	}

	version, err := extract.Version(dots)
	if err != nil {
		return nil, 0, 0, err
	}
	maskID, errorCorrectionLevel, err := extract.Format(dots)
	if err != nil {
		return nil, 0, 0, err
	}
	fmt.Printf("Mask ID is %d / Error correction level is %s\n", maskID, errorCorrectionLevel.String())

	bits := extract.ReadBits(dots, maskID)

	return bits, version, errorCorrectionLevel, nil
}

// decodeMessages performs error correction on the bits read, then decodes the message.
// In case error correction fails, the uncorrected message is returned (if possible).
func decodeMessage(bits []bool, version uint, errorCorrectionLevel decode.ErrorCorrectionLevel) (string, error) {
	bitsCorrected, err := decode.Correct(bits, version, errorCorrectionLevel)
	if err != nil {
		return "", err
	}

	mode, bits := decode.GetMode(bitsCorrected)
	length, contents, err := decode.GetContentLength(bits, version, mode, errorCorrectionLevel)
	if err != nil {
		return "", err
	}
	fmt.Printf("Mode is %s / Content length is %d bytes\n", mode.String(), length)

	message, err := decode.Message(mode, length, contents)
	if err != nil {
		return "", err
	}

	return message, nil
}
