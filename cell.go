// Copyright (c) 2019-2021 Leonid Kneller. All rights reserved.
// Licensed under the MIT license.
// See the LICENSE file for full license information.

package ons2

import (
	"github.com/reconditematter/mym"
	"math"
)

// CellRnd1x1 -- generates a uniform pseudo-random sequence on `n` points
// in the geographic grid cell [lat,lat+1]x[lon,lon+1].
// This function causes a runtime panic when lat∉{-90,...,89} or lon∉{-180,...,179}.
func CellRnd1x1(lat, lon, n int) []Point {
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
