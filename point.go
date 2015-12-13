package airlinetracks

import (
	"strings"
	"strconv"
	"math"
)

type point struct {
	x float64
	y float64
	z float64
	t int64
}
func (p *point) ParsePoint(s string) {
	sa:=strings.Fields(s)
	p.x,_ =strconv.ParseFloat(sa[0],64)
	p.y,_ =strconv.ParseFloat(sa[1],64)
	p.z,_ =strconv.ParseFloat(sa[2],64)
	p.t,_ =strconv.ParseInt(sa[3],0,32)
}

func ParsePoint(s string) point{
	var p point
	p.ParsePoint(s)
	return p
}

func GetDistanceMiles(p1 point, p2 point) (miles float64) {
	degree_to_rad := math.Pi / 180.0
	d_lon := (p2.x - p1.x) * degree_to_rad
	d_lat := (p2.y - p1.y) * degree_to_rad
	a := math.Pow(math.Sin(d_lat/2), 2) +
	   math.Cos(p1.y*degree_to_rad)*
		math.Cos(p2.y*degree_to_rad)*
		math.Pow(math.Sin(d_lon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	miles = 3956 * c
	return miles
}

func (p1 point) GetDistanceMiles(p2 point) (miles float64) {
	return GetDistanceMiles(p1,p2)
}
