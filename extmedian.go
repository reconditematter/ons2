// Copyright (c) 2019-2020 Leonid Kneller. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file for full license information.

package ons2

import (
	"github.com/reconditematter/mym"
)

// ExtMedian -- computes the extrinsic geometric median of `points`
// on the unit sphere S².
//
// This computation involves two steps: (1) Computation of the geometric
// median of the data points seen as vectors in the Euclidean space E³.
// (2) A shortest distance projection back to the sphere S².
func ExtMedian(points []Point) Point {
	n := len(points)
	u := make([][3]float64, n)
	//
	for i, p := range points {
		u[i] = p.C()
	}
	//
	return Point{mym.Vhat3(mym.Vmedian3(u))}
}
