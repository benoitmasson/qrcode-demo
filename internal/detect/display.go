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
	gocv.Line(img, points[0], points[1], color, width)
	gocv.Line(img, points[1], points[2], color, width)
	gocv.Line(img, points[2], points[3], color, width)
	gocv.Line(img, points[3], points[0], color, width)
}
