// Copyright (c) 2019-2021 Leonid Kneller. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file for full license information.

package ons2

import (
	"math"
)

// Kpois -- computes the value of K function for complete spatial random (Poisson)
// process on the unit sphere S² given an angular separation `θ`. This function
// returns a NaN when θ∉[0,π].
func Kpois(θ float64) float64 {
	if !(0 <= θ && θ <= math.Pi) {
		return math.NaN()
	}
	//
	return 2 * math.Pi * (1 - math.Cos(θ))
}

// Kripley -- computes the estimate of Ripley's K function for a sample of `points`
// on the unit sphere S² given an angular separation `θ`. This function
// returns a NaN when len(points)<2 or θ∉[0,π].
func Kripley(points []Point, θ float64) float64 {
	n := len(points)
	if n < 2 {
		return math.NaN()
	}
	if !(0 <= θ && θ <= math.Pi) {
		return math.NaN()
	}
	//
	k := 0
	for i, pi := range points {
		for j := i + 1; j < n; j++ {
			pj := points[j]
			s := pi.Sep(pj)
			if s <= θ {
				k++
			}
		}
	}
	return 8 * math.Pi * float64(k) / float64(n*(n-1))
}
