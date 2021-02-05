// Copyright (c) 2019-2021 Leonid Kneller. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file for full license information.

package ons2

import (
	"github.com/reconditematter/mym"
	"math"
)

// IntMedian -- computes the intrinsic sample median of `points`
// on the unit sphere S².
//
// This computation starts with the extrinsic median (see `ExtMedian`) as
// an initial approximation. The following iterations refine the solution
// by minimizing the sum of distances from the current iteration
// median to all given points.
//
// This function returns the computed intrinsic median and the number
// of iterations required to converge to the solution.
func IntMedian(points []Point) (Point, int) {
	n := len(points)
	if n == 1 {
		return points[0], 1
	}
	xyz := make([][3]float64, n)
	//
	var cen Point
	if n <= 10000 {
		i1, _ := Medoids(points)
		cen = points[i1]
	} else {

		cen = ExtMedian(points)
	}
	//
	for iter := 1; iter <= 1000; iter++ {
		for i, p := range points {
			xyz[i][0], xyz[i][1] = azeq(cen, p)
			// xyz[i][2] is zero
		}
		// call `Vmedian3` to compute the geometric median in the plane (z=0)
		xyzmedian := mym.Vmedian3(xyz)
		newcen := azeqinv(cen, xyzmedian[0], xyzmedian[1])
		if cen.Sep(newcen) <= mym.Epsilon*math.Pi*2 {
			return newcen, iter
		}
		cen = newcen
	}
	return cen, 1000
}
