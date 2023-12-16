package utils

import (
	"math"
	"math/rand"
)

// distance: calculate the distance between two points
func Calcul_Distance(x1, y1, x2, y2 int) float64 {
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	return math.Sqrt(dx*dx + dy*dy)
}

// min returns the minimum of two integers.
func Intmin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers.
func Intmax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Float32min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func Float32max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func RandomPositionInRectangle(posX, posY, rayon, x_lower_bound, x_upper_bound, y_lower_bound, y_upper_bound int) (int, int) {

	X_minVal := Intmax(x_lower_bound, posX-rayon)
	X_maxVal := Intmin(x_upper_bound, posX+rayon)
	Y_minVal := Intmax(y_lower_bound, posY-rayon)
	Y_maxVal := Intmin(y_upper_bound, posY+rayon)

	X := X_minVal + rand.Intn(X_maxVal-X_minVal+1)
	Y := Y_minVal + rand.Intn(Y_maxVal-Y_minVal+1)

	return X, Y
}
