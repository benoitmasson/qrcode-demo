package detect

import (
	"image"

	"gocv.io/x/gocv"
)

// EnhanceImage improves image contrast and sharpness
func EnhanceImage(img *gocv.Mat) {
	// "open" to clean image: https://docs.opencv.org/4.x/d9/d61/tutorial_py_morphological_ops.html
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Point{X: 3, Y: 3})
	defer kernel.Close()
	gocv.MorphologyEx(*img, img, gocv.MorphOpen, kernel)
	// increase contrast
	gocv.AddWeighted(*img, 2, *img, 0, -50, img)
}
