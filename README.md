# Demo QR-Code Reader

[![MIT license](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## Disclaimer

This project was inspired by [Dan Hollick's article](https://typefully.com/DanHollick/qr-codes-T7tLlNi), which gives a nice overview of how QR-codes work. It is an educational _demo_ used to explain how QR-codes work, it is not meant to be used in real-life applications.

If you need an efficient and complete QR-code reader written in Go, you may use the powerful [makiuchi-d/gozxing](https://github.com/makiuchi-d/gozxing) library.

Or if you really want to rely on OpenCV, then use [QRCodeDetector.DetectAndDecode](https://pkg.go.dev/gocv.io/x/gocv#QRCodeDetector.DetectAndDecode) method from [gocv.io](https://gocv.io/) bindings.

## Presentation

The [presentation slides](./slides/QR%20Codes.pdf) can be found in the [`slides`](./slides/) folder.

## Instructions

To be able to compile and run this tool, [OpenCV 4](https://opencv.org/) library is required.

Install it using [official instructions](https://docs.opencv.org/4.x/df/d65/tutorial_table_of_content_introduction.html) or using [Homebrew](https://formulae.brew.sh/formula/opencv).

Then, compile the application with `go build` or run with `go run .`  
The first build/run may be slow, because OpenCV files need to be compiled beforehand.
The application can be stopped by pressing `Esc` key at any time.

A different capture device may be chosen with parameter `--device-id` (default to `0`).

## Explanations

### 1. Code detection

The application captures frames from the webcam, then tries to detect QR-codes in each frame.

Image capture and processing is performed thanks to [OpenCV 4](https://opencv.org/) Go bindings from [gocv.io](https://pkg.go.dev/gocv.io/x/gocv).

1. For each frame, try to detect a QR-code with [QRCodeDetector.Detect](https://pkg.go.dev/gocv.io/x/gocv#QRCodeDetector.Detect).

   This function returns a set of 4 points delimiting a QR-code candidate in the image.

1. Eliminate false positives by keeping only coordinates forming a square (with some tolerance).

1. Project and display the detected QR-code in the top-left corner, using [WarpPerspective](https://pkg.go.dev/gocv.io/x/gocv#WarpPerspective) function

   Enhance QR-code image contrast (with [AddWeighted](https://pkg.go.dev/gocv.io/x/gocv#AddWeighted)) and "[open](https://docs.opencv.org/4.x/d9/d61/tutorial_py_morphological_ops.html)" image to remove noise (with [GetStructuringElement](https://pkg.go.dev/gocv.io/x/gocv#GetStructuringElement))

1. Compute QR-code dots width in pixels, then scan the image pixel to construct the dot matrix, and display it on the console.

When all steps are successful, the QR-code is highlighted in the image, and the video freezes for a few seconds to show the result.

### 2. Extracting contents

Once the QR-code dots have been detected, the code contents bits are extracted from it.

1. Get metadata from the QR-code: version information (the "size" of the code), the mask ID to apply to the dots and the error correction level used on the contents.

   For the last two (the code "format"), error correction is used on the selected dots to make sure the value found is correct. This error correction implements the Reed-Solomon algorithm, as explained on [this page](https://www.thonky.com/qr-code-tutorial/format-version-information).

2. Read the contents bits in the correct order, starting from the bottom-right, 2 columns at a time from right to left, alternating upwards and downwards and avoiding reserved areas.

   See [explanations and picture](https://www.thonky.com/qr-code-tutorial/module-placement-matrix#step-6-place-the-data-bits) for an illustration.

   Note that the current implementation supports only 1 alignement pattern, and thus works only with QR-codes version 6 and below.

The returned bits contain metadata (content type and length), and the contents with error correction data.

### 3. Decoding message

Finally, decode the message from the bits contents.

1. First of all, use the [Reed-Solomon algorithm](https://en.m.wikiversity.org/wiki/Reed%E2%80%93Solomon_codes_for_coders) to perform error correction on the full contents bits, to correct any errors that may have been introduced during the scanning process. Library http://github.com/colin-davis/reedSolomon is used for this purpose.

   Note that interleaved blocks are not supported yes, hence the higher version-error correction level pairs will not be decoded.

2. Then, read data from the contents bits: first, metadata (character mode and message length), then the message itself. See [this page](https://www.thonky.com/qr-code-tutorial/data-encoding) for more details on how data is encoded.

   Note that Kanji mode is not supported at all, and ECI (unicode) may produce strange results.

If the QR-code is successfully decoded, the message is revealed in the console.

## References

The following explanations have been used to implement this project. Many thanks to their authors.

### Big picture

- Dan Hollick's blog: https://typefully.com/DanHollick/qr-codes-T7tLlNi

### Encoding & Decoding

- Thonky's guide: https://www.thonky.com/qr-code-tutorial
- Linux Magazine n°194: https://connect.ed-diamond.com/GNU-Linux-Magazine/glmf-194/decoder-un-code-qr (French)

### Error correction

- Linux Magazine n°198: https://connect.ed-diamond.com/GNU-Linux-Magazine/glmf-198/reparer-un-code-qr (French)
- Wikiversity: https://en.m.wikiversity.org/wiki/Reed%E2%80%93Solomon_codes_for_coders
