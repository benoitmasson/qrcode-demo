package detect

import "image"

// ValidateSquare returns whether the 4 given points roughly correspond to a square in the image.
// Points should be given in the following order:
// - bottom-left, bottom-right, top-right, top-left
// - or bottom-right, then top-right, top-left, bottom-left
// (as returned by gocv.NewQRCodeDetector().Detect function)
// Parameters width and height correspond to the global image size, to check for out-of-image points.
func ValidateSquare(points []image.Point, width, height int) bool {
	if len(points) != 4 {
		return false
	}

	// invalid coordinates
	if points[0].X < 0 || points[0].X >= width || points[0].Y < 0 || points[0].Y >= height ||
		points[1].X < 0 || points[1].X >= width || points[1].Y < 0 || points[1].Y >= height ||
		points[2].X < 0 || points[2].X >= width || points[2].Y < 0 || points[2].Y >= height ||
		points[3].X < 0 || points[3].X >= width || points[3].Y < 0 || points[3].Y >= height {
		return false
	}

	tolerance := height / 12
	// detected points are ordered as follows: bottom-left first, then bottom-right, then top-right, then top-left
	// make sure they form something which looks like an unrotated square
	if !nearlyEquals(points[0].Y, points[1].Y, tolerance) ||
		!nearlyEquals(points[2].Y, points[3].Y, tolerance) ||
		!nearlyEquals(points[0].X, points[3].X, tolerance) ||
		!nearlyEquals(points[1].X, points[2].X, tolerance) ||
		!nearlyEquals(points[1].X-points[0].X, points[3].Y-points[0].Y, tolerance) { // not a square
		// not a horizontal square, check points are ordered as follows: bottom-right, then top-right, top-left, bottom-left
		if !nearlyEquals(points[0].X, points[1].X, tolerance) ||
			!nearlyEquals(points[2].X, points[3].X, tolerance) ||
			!nearlyEquals(points[0].Y, points[3].Y, tolerance) ||
			!nearlyEquals(points[1].Y, points[2].Y, tolerance) ||
			!nearlyEquals(points[1].Y-points[0].Y, points[0].X-points[3].X, tolerance) { // not a square
			// not a horizontal nor vertical square
			return false
		}
	}

	return true
}
