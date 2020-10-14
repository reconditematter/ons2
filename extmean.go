package ons2

import (
	"github.com/reconditematter/mym"
	"math"
	"math/rand"
	"sort"
	"time"
)

// ExtMean -- computes the extrinsic sample mean of `points`
// on the unit sphere S².
//
// This computation involves two steps: (1) Computation of the mean
// of the data points seen as vectors in the Euclidean space E³.
// (2) A shortest distance projection back to the sphere S².
func ExtMean(points []Point) Point {
	n := len(points)
	sumx := mym.AccuSum(n, func(i int) float64 { return points[i].C()[0] })
	sumy := mym.AccuSum(n, func(i int) float64 { return points[i].C()[1] })
	sumz := mym.AccuSum(n, func(i int) float64 { return points[i].C()[2] })
	r := math.Hypot(math.Hypot(sumx, sumy), sumz)
	return Point{[3]float64{sumx / r, sumy / r, sumz / r}}
}

// ExtMeanBoot -- computes the extrinsic sample mean of `points`
// on the unit sphere S² (see `ExtMean`); also computes two
// confidence cones using the non-parametric bootstrap quantile method.
// When seed≠0, it is used as the seed for a pseudo-random number generator;
// otherwise, a seed based on a Unix time is used.
//
//	c95 -- ½ of the vertex angle (degrees) for the 95% confidence cone
//	c99 -- ½ of the vertex angle (degrees) for the 99% confidence cone
func ExtMeanBoot(points []Point, seed int64) (mean Point, c95, c99 float64) {
	const (
		B   = 40000
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
	// compute the exrinsic mean
	mean = ExtMean(points)
	//
	// compute `B` extrinsic means using bootstrap method
	n := len(points)
	tpoints := make([]Point, n)
	bmeans := make([]locsep, B)
	for b := 0; b < B; b++ {
		for k := range tpoints {
			tpoints[k] = points[rng.Intn(n)]
		}
		tmean := ExtMean(tpoints)
		bmeans[b] = locsep{loc: tmean, sep: mean.Sep(tmean)}
	}
	//
	// compute the 95% and 99% confidence cones
	sort.Sort(locseps(bmeans))
	c95 = bmeans[B95-1].sep * (180 / math.Pi)
	c99 = bmeans[B99-1].sep * (180 / math.Pi)
	//
	return
}

type locsep struct {
	loc Point
	sep float64
}

type locseps []locsep

func (s locseps) Len() int           { return len(s) }
func (s locseps) Less(i, j int) bool { return s[i].sep < s[j].sep }
func (s locseps) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
