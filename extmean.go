package ons2

import (
	"github.com/reconditematter/mym"
	"math"
)

// ExtMean -- computes the extrinsic sample mean of `points`
// on the unit sphere S².
//
// This computation involves two steps: (1) Computation of the mean
// of the data points seen as vectors in the Euclidean space E³.
// (2) A shortest distance projection back to the sphere S².
func ExtMean(points []Point) Point {
	n := len(points)
	sumx := mym.AccuSum(n, func(i int) float64 { return points[i].C()[0] })
	sumy := mym.AccuSum(n, func(i int) float64 { return points[i].C()[1] })
	sumz := mym.AccuSum(n, func(i int) float64 { return points[i].C()[2] })
	r := math.Hypot(math.Hypot(sumx, sumy), sumz)
	return Point{[3]float64{sumx / r, sumy / r, sumz / r}}
}
