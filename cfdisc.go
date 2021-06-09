// Copyright (c) 2019-2021 Leonid Kneller. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file for full license information.

package ons2

import (
	"math"
)

// Discrepancy -- computes the generalized discrepancy of `points` to quantify
// a criterion of equidistributed point sets.
//
// Reference: J.Cui and W.Freeden, Equidistribution on the Sphere, SIAM J. Sci. Comput., 18(2), 595–609 (1997).
//
// DOI: https://doi.org/10.1137/S1064827595281344
func Discrepancy(points []Point) float64 {
	const sqrtpi = 1.77245385090551602729816748334114518279754945612239
	n := len(points)
	ln := func(i, j int) float64 {
		u := points[i].C()
		v := points[j].C()
		dot := u[0]*v[0] + u[1]*v[1] + u[2]*v[2]
		if dot > 1 {
			dot = 1
		} else if dot < -1 {
			dot = -1
		}
		return math.Log(1 + math.Sqrt((1-dot)/2))
	}
	//
	D := 0.0
	for i := range points {
		for j := range points {
			D += 1 - 2*ln(i, j)
		}
	}
	return (1 / (2 * sqrtpi * float64(n))) * math.Sqrt(D)
}
