package main

import (
	"image"

	"gocv.io/x/gocv"
)

func newImagePointsFromPoints(points *gocv.Mat) []image.Point {
	r, c := points.Rows(), points.Cols()
	// fmt.Println(points.Channels(), points.Size(), points.Type(), points.Total(), r, c)

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
