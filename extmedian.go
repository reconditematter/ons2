// Copyright (c) 2019-2021 Leonid Kneller. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file for full license information.

package ons2

import (
	"github.com/reconditematter/mym"
	"math"
	"math/rand"
	"sort"
	"time"
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

// ExtMedianBoot -- computes the extrinsic geometric median of `points`
// on the unit sphere S² (see `ExtMedian`); also computes two
// confidence cones using the non-parametric bootstrap quantile method.
// When seed≠0, it is used as the seed for a pseudo-random number generator;
// otherwise, a seed based on a Unix time is used.
//
//      c95 -- ½ of the vertex angle (degrees) for the 95% confidence cone
//      c99 -- ½ of the vertex angle (degrees) for the 99% confidence cone
func ExtMedianBoot(points []Point, seed int64) (median Point, c95, c99 float64) {
	const (
		B   = 5000
		B95 = B * 95 / 100
		B99 = B * 99 / 100
	)
	//
	src := mym.MT19937()
	if seed == 0 {
		src.Seed(time.Now().UnixNano())
	} else {
		src.Seed(seed)
	}
	rng := rand.New(src)
	//
	// compute the exrinsic median
	median = ExtMedian(points)
	//
	// compute `B` extrinsic medians using bootstrap method
	n := len(points)
	tpoints := make([]Point, n)
	bseps := make([]float64, B)
	for b := 0; b < B; b++ {
		for k := range tpoints {
			tpoints[k] = points[rng.Intn(n)]
		}
		tmedian := ExtMedian(tpoints)
		bseps[b] = median.Sep(tmedian)
	}
	//
	// compute the 95% and 99% confidence cones
	sort.Float64s(bseps)
	c95 = bseps[B95-1] * (180 / math.Pi)
	c99 = bseps[B99-1] * (180 / math.Pi)
	//
	return
}
