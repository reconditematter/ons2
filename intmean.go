// Copyright (c) 2019-2020 Leonid Kneller. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file for full license information.

package ons2

import (
	"github.com/reconditematter/mym"
	"math"
)

// IntMean -- computes the intrinsic sample mean of `points`
// on the unit sphere S².
//
// This computation starts with the extrinsic mean (see `ExtMean`) as
// an initial approximation. The following iterations refine the solution
// by minimizing the sum of squared distances from the current iteration
// mean to all other points.
//
// This function returns the computed intrinsic mean and the number
// of iterations required to converge to the solution.
func IntMean(points []Point) (Point, int) {
	n := len(points)
	xs, ys := make([]float64, n), make([]float64, n)
	cen := ExtMean(points)
	for iter := 1; iter <= 1000; iter++ {
		for i, p := range points {
			xs[i], ys[i] = azeq(cen, p)
		}
		xmean := mym.AccuSum(n, func(i int) float64 { return xs[i] }) / float64(n)
		ymean := mym.AccuSum(n, func(i int) float64 { return ys[i] }) / float64(n)
		newcen := azeqinv(cen, xmean, ymean)
		if cen.Sep(newcen) <= mym.Epsilon {
			return newcen, iter
		}
		cen = newcen
	}
	return cen, 1000
}

// azeq -- azimuthal equidistant projection
func azeq(cen, p Point) (x, y float64) {
	sinc := func(x float64) float64 {
		if math.Abs(x) < mym.Epsilon {
			return 1
		}
		return math.Sin(x) / x
	}
	//
	c := cen.Sep(p)
	kp := 1 / sinc(c)
	//
	lat0, lon0 := cen.Geo()
	lat, lon := p.Geo()
	sin0, cos0 := mym.SinCosD(lat0)
	sin, cos := mym.SinCosD(lat)
	//
	del := lon - lon0
	if del > 180 {
		del -= 360
	} else if del < -180 {
		del += 360
	}
	sind, cosd := mym.SinCosD(del)
	//
	x = kp * cos * sind
	y = kp * (cos0*sin - sin0*cos*cosd)
	//
	return
}

// azeqinv -- inverse azimuthal equidistant projection
func azeqinv(cen Point, x, y float64) Point {
	lat0, lon0 := cen.Geo()
	sin0, cos0 := mym.SinCosD(lat0)
	//
	c := math.Hypot(x, y)
	if c < mym.Epsilon {
		return cen
	}
	sinc, cosc := math.Sincos(c)
	//
	lat := math.Asin(cosc*sin0+y*sinc*cos0/c) * (180 / math.Pi)
	lon := lon0 + math.Atan2(x*sinc, c*cos0*cosc-y*sin0*sinc)*(180/math.Pi)
	if lon > 180 {
		lon -= 360
	} else if lon < -180 {
		lon += 360
	}
	//
	return Geo(lat, lon)
}
