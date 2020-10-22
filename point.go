// Copyright (c) 2019-2020 Leonid Kneller. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file for full license information.

package ons2

import (
	"github.com/reconditematter/mym"
	"math"
)

// Point -- a point on the unit sphere S².
type Point struct {
	c [3]float64
}

// Geo -- returns a point on the unit sphere S² with given
// geographic coordinates (lat[itude],lon[gitude]).
// This function causes a runtime panic when lat∉[-90,90]
// or lon∉[-180,180].
func Geo(lat, lon float64) Point {
	if !(-90 <= lat && lat <= 90) {
		panic("ons2.Geo: lat not in [-90,90]")
	}
	if !(-180 <= lon && lon <= 180) {
		panic("ons2.Geo: lon not in [-180,180]")
	}
	//
	slat, clat := mym.SinCosD(lat)
	slon, clon := mym.SinCosD(lon)
	x, y := clat*clon, clat*slon
	z := slat
	return Point{[3]float64{x, y, z}}
}

// C -- returns the Cartesian coordinates (x,y,z) of `p`.
func (p Point) C() (xyz [3]float64) {
	xyz = p.c
	if xyz[0] == 0 && xyz[1] == 0 && xyz[2] == 0 {
		xyz[2] = 1
	}
	return
}

// Geo -- returns the geographic coordinates (lat[itude],lon[gitude]) of `p`.
func (p Point) Geo() (lat, lon float64) {
	xyz := p.C()
	r := math.Hypot(xyz[0], xyz[1])
	lat = math.Atan2(xyz[2], r) * (180 / math.Pi)
	lon = math.Atan2(xyz[1], xyz[0]) * (180 / math.Pi)
	return
}

// Sep -- computes the separation angle (geodesic distance) between `p` and `q`.
func (p Point) Sep(q Point) float64 {
	u, v := p.C(), q.C()
	dot := u[0]*v[0] + u[1]*v[1] + u[2]*v[2]
	if dot > 0 {
		x, y, z := u[0]-v[0], u[1]-v[1], u[2]-v[2]
		return 2 * math.Asin(math.Hypot(math.Hypot(x, y), z)/2)
	}
	if dot < 0 {
		x, y, z := u[0]+v[0], u[1]+v[1], u[2]+v[2]
		return math.Pi - 2*math.Asin(math.Hypot(math.Hypot(x, y), z)/2)
	}
	if dot == 0 {
		return math.Pi / 2
	}
	return math.NaN()
}

// Random -- returns a uniform pseudo-random point on the unit sphere S².
// This function is safe for concurrent use by multiple goroutines.
func Random() Point {
	const ε = 1.0 / (1 << 52)
	x, y, z := mym.N01(), mym.N01(), mym.N01()
	r := math.Hypot(math.Hypot(x, y), z)
	for r < ε {
		x, y, z = mym.N01(), mym.N01(), mym.N01()
		r = math.Hypot(math.Hypot(x, y), z)
	}
	return Point{[3]float64{x / r, y / r, z / r}}
}

// Fibonacci -- generates a quasi-uniform sequence of 2n+1 Fibonacci points on the unit sphere S².
//
// Reference: Swinbank R., Purser R.J., Fibonacci grids: A novel approach to global modelling,
// Q.J.R. Meteorol. Soc., vol.132, no.619, pp.1769-1793 (2006).
//
// DOI: https://doi.org/10.1256/qj.05.227
func Fibonacci(n int) []Point {
	if n < 0 {
		n = 0
	}
	n21 := float64(n*2 + 1)
	p := make([]Point, 2*n+1)
	for i := -n; i <= n; i++ {
		fi := float64(i)
		sinφ := 2 * fi / n21
		cosφ := math.Sqrt((1 - sinφ) * (1 + sinφ))
		λ := (2 * math.Pi / math.Phi) * fi
		sinλ, cosλ := math.Sincos(λ)
		//
		x, y, z := cosφ*cosλ, cosφ*sinλ, sinφ
		p[i+n] = Point{[3]float64{x, y, z}}
	}
	return p
}
