// Copyright (c) 2019-2021 Leonid Kneller. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file for full license information.

package ons2

import (
	"github.com/reconditematter/mym"
	"math"
)

// Medoids -- finds two medoids of `points` on the unit sphere S².
// The first medoid point `i1` minimizes the sum of distances from points[i1] to all given points.
// The second medoid point `i2` minimizes the sum of squared distances from points[i2] to all given points.
func Medoids(points []Point) (i1, i2 int) {
	n := len(points)
	if n == 0 {
		i1, i2 = -1, -1
		return
	}
	//
	D := mym.NewSym0(n)
	for i, pi := range points {
		for j := i + 1; j < n; j++ {
			pj := points[j]
			D.Set(i, j, pi.Sep(pj))
		}
	}
	//
	min1, min2 := math.Inf(1), math.Inf(1)
	for k := 0; k < n; k++ {
		sum1 := mym.AccuSum(n, func(i int) float64 { return D.Get(k, i) })
		sum2 := mym.AccuDot(n, func(i int) float64 { return D.Get(k, i) }, func(i int) float64 { return D.Get(k, i) })
		if sum1 < min1 {
			i1 = k
			min1 = sum1
		}
		if sum2 < min2 {
			i2 = k
			min2 = sum2
		}
	}
	//
	return
}
