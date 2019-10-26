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
