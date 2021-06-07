// Copyright (c) 2019-2021 Leonid Kneller. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file for full license information.

package ons2

import (
	"github.com/reconditematter/mym"
	"math"
)

// CellRnd1x1 -- generates a uniform pseudo-random sequence of `n` points
// in the geographic grid cell [lat,lat+1]x[lon,lon+1].
// This function causes a runtime panic when lat∉{-90,...,89} or lon∉{-180,...,179}.
func CellRnd1x1(lat, lon, n int) []Point {
	if !(-90 <= lat && lat < 90) {
		panic("ons2.CellRnd1x1: lat not in [-90,89]")
	}
	if !(-180 <= lon && lon < 180) {
		panic("ons2.CellRnd1x1: lon not in [-180,179]")
	}
	ps := make([]Point, n)
	latmin, latmax := float64(lat), float64(lat+1)
	lonmin, lonmax := float64(lon), float64(lon+1)
	lat0 := latmin + 0.5
	lon0 := lonmin + 0.5
	// use cylindrical equal-area projection to find the cell boundary in the plane
	xmin, ymin := cyleq(lat0, lon0, latmin, lonmin)
	xmax, ymax := cyleq(lat0, lon0, latmax, lonmax)
	dx, dy := xmax-xmin, ymax-ymin
	//
	for i := range ps {
		x := mym.U01()*dx + xmin
		y := mym.U01()*dy + ymin
		// use inverse cylindrical equal-are projection to get geographic coordinates
		lat, lon := cyleqinv(lat0, lon0, x, y)
		ps[i] = Geo(lat, lon)
	}
	return ps
}

// CellFib1x1 -- generates a quasi-uniform sequence of _approximately_ `n` Fibonacci
// points in the geographic grid cell [lat,lat+1]x[lon,lon+1].
// This function causes a runtime panic when lat∉{-90,...,89} or lon∉{-180,...,179}.
//
// Reference: Swinbank R., Purser R.J., Fibonacci grids: A novel approach to global modelling,
// Q.J.R. Meteorol. Soc., vol.132, no.619, pp.1769-1793 (2006).
//
// DOI: https://doi.org/10.1256/qj.05.227
func CellFib1x1(lat, lon, n int) []Point {
	if !(-90 <= lat && lat < 90) {
		panic("ons2.CellFib1x1: lat not in [-90,89]")
	}
	if !(-180 <= lon && lon < 180) {
		panic("ons2.CellFib1x1: lon not in [-180,179]")
	}
	// compute 1x1 cell area
	const S = 4 * math.Pi
	cellA := (math.Pi / 180) * math.Abs(math.Sin(float64(lat+1)*(math.Pi/180))-math.Sin(float64(lat)*(math.Pi/180)))
	nn := int(math.Ceil(float64(n) * S / cellA))
	m := nn / 2
	m21 := float64(m*2 + 1)
	p := make([]Point, 0, n)
	//
	for i := -m; i <= m; i++ {
		fi := float64(i)
		sinφ := 2 * fi / m21
		cosφ := math.Sqrt((1 - sinφ) * (1 + sinφ))
		φdeg := math.Atan2(sinφ, cosφ) * (180 / math.Pi)
		if !(float64(lat) < φdeg && φdeg < float64(lat+1)) {
			continue
		}
		λ := (2 * math.Pi / math.Phi) * fi
		sinλ, cosλ := math.Sincos(λ)
		λdeg := math.Atan2(sinλ, cosλ) * (180 / math.Pi)
		if !(float64(lon) < λdeg && λdeg < float64(lon+1)) {
			continue
		}
		//
		x, y, z := cosφ*cosλ, cosφ*sinλ, sinφ
		p = append(p, Point{[3]float64{x, y, z}})
	}
	return p
}

// cylindrical equal-area projection
func cyleq(lat0, lon0, lat, lon float64) (x, y float64) {
	cos := math.Cos(lat0 * (math.Pi / 180))
	x = (math.Pi / 180) * (lon - lon0) * cos
	y = math.Sin(lat*(math.Pi/180)) / cos
	return
}

// cylindrical equal-area projection (inverse)
func cyleqinv(lat0, lon0, x, y float64) (lat, lon float64) {
	cos := math.Cos(lat0 * (math.Pi / 180))
	lat = math.Asin(y*cos) * (180 / math.Pi)
	lon = (x/cos)*(180/math.Pi) + lon0
	return
}
