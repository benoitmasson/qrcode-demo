package main

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

type QRCode [][]bool

func newImagePointsFromPoints(points *gocv.Mat) []image.Point {
	r, c := points.Rows(), points.Cols()

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

func outlineQRCode(img *gocv.Mat, points []image.Point, color color.RGBA, width int) {
	for i := 1; i < len(points); i++ {
		gocv.Line(img, points[i-1], points[i], color, width)
	}
	if len(points) > 0 {
		// close outline
		gocv.Line(img, points[len(points)-1], points[0], color, width)
	}
}
