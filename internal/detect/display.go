package detect

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

func SetMiniCodeInCorner(img *gocv.Mat, points []image.Point, width, height int) gocv.Mat {
	// TODO (1.2): perform a perspective transform + translation to show detected QR-code in top-left corner
	return gocv.NewMat()
}

func OutlineQRCode(img *gocv.Mat, points []image.Point, color color.RGBA, width int) {
	for i := 1; i < len(points); i++ {
		gocv.Line(img, points[i-1], points[i], color, width)
	}
	if len(points) > 0 {
		// close outline
		gocv.Line(img, points[len(points)-1], points[0], color, width)
	}
}
