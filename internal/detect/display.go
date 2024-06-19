package detect

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

func SetMiniCodeInCorner(img *gocv.Mat, points []image.Point, width, height int) gocv.Mat {
	originVector := gocv.NewPointVectorFromPoints(points)
	defer originVector.Close()
	destinationVector := gocv.NewPointVectorFromPoints([]image.Point{
		{X: 0, Y: 0},
		{X: width - 1, Y: 0},
		{X: width - 1, Y: height - 1},
		{X: 0, Y: height - 1},
	})
	defer destinationVector.Close()
	transform := gocv.GetPerspectiveTransform(originVector, destinationVector)
	defer transform.Close()
	rectangle := (*img).Region(image.Rect(0, 0, width-1, height-1))

	gocv.WarpPerspective(*img, &rectangle, transform, image.Point{X: width - 1, Y: height - 1})
	return rectangle
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
