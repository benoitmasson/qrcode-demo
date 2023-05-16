# Demo QR-Code Reader

[![MIT license](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## Disclaimer

This project is an educational _demo_ used to explain how QR-codes work, it is not meant to be used in real-life applications.

If you need an efficient and complete QR-code reader written in Go, you may use the powerful [makiuchi-d/gozxing](https://github.com/makiuchi-d/gozxing) library.

Or if you really want to rely on OpenCV, then use [QRCodeDetector.DetectAndDecode](https://pkg.go.dev/gocv.io/x/gocv#QRCodeDetector.DetectAndDecode) method from [gocv.io](https://gocv.io/) bindings.

## Instructions

To be able to compile and run this tool, [OpenCV 4](https://opencv.org/) library is required.

Install it using [official instructions](https://docs.opencv.org/4.x/df/d65/tutorial_table_of_content_introduction.html) or using [Homebrew](https://formulae.brew.sh/formula/opencv).

Then, compile the application with `go build` or run with `go run .`  
The first build/run may be slow, because OpenCV files need to be compiled beforehand.
The application can be stopped by pressing `Esc` key at any time.

A different capture device may be chosen with parameter `--device-id` (default to `0`).

## Description

### Detection

The application captures frames from the webcam, then tries to detect QR-codes in each frame.

Image capture and processing is performed thanks to [OpenCV 4](https://opencv.org/) Go bindings from [gocv.io](https://pkg.go.dev/gocv.io/x/gocv).

1. For each frame, try to detect a QR-code with [QRCodeDetector.Detect](https://pkg.go.dev/gocv.io/x/gocv#QRCodeDetector.Detect).

   This function returns a set of 4 points delimiting a QR-code candidate in the image.

1. Eliminate false positives by keeping only coordinates forming a square (with some tolerance).

1. Project and display the detected QR-code in the top-left corner, using [WarpPerspective](https://pkg.go.dev/gocv.io/x/gocv#WarpPerspective) function

   Enhance QR-code image contrast (with [AddWeighted](https://pkg.go.dev/gocv.io/x/gocv#AddWeighted)) and "[open](https://docs.opencv.org/4.x/d9/d61/tutorial_py_morphological_ops.html)" image to remove noise (with [GetStructuringElement](https://pkg.go.dev/gocv.io/x/gocv#GetStructuringElement))

1. Compute QR-code dots width in pixels, then scan the image pixel to construct the dot matrix, and display it on the console.

When all steps are successful, the QR-code is highlighted in the image, and the video freezes for a few seconds to show the result.

### Decoding

<!-- TODO -->

https://typefully.com/DanHollick/qr-codes-T7tLlNi
