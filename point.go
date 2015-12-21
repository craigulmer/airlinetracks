package airlinetracks

import (
	"fmt"
	"strings"
	"strconv"
	"math"
)


type Point struct {
	X float64
	Y float64
	Z float64
	T int64
}

func (p *Point) ParsePoint(s string) {
	sa:=strings.Fields(s)
	p.X,_ =strconv.ParseFloat(sa[0],64)
	p.Y,_ =strconv.ParseFloat(sa[1],64)
	p.Z,_ =strconv.ParseFloat(sa[2],64)
	p.T,_ =strconv.ParseInt(sa[3],0,32)
}

func (p *Point) ToString() string {
	return fmt.Sprintf("%f %f %f %d", p.X, p.Y, p.Z, p.T)
}

func ParsePoint(s string) Point{
	var p Point
	p.ParsePoint(s)
	return p
}

func GetDistanceMiles(p1 Point, p2 Point) (miles float64) {
	degree_to_rad := math.Pi / 180.0
	d_lon := (p2.X - p1.X) * degree_to_rad
	d_lat := (p2.Y - p1.Y) * degree_to_rad
	a := math.Pow(math.Sin(d_lat/2), 2) +
	   math.Cos(p1.Y*degree_to_rad)*
		math.Cos(p2.Y*degree_to_rad)*
		math.Pow(math.Sin(d_lon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	miles = 3956 * c
	return miles
}

func (p1 Point) GetDistanceMiles(p2 Point) (miles float64) {
	return GetDistanceMiles(p1,p2)
}
