package main

import (
	"fmt"
	"bytes"
	atx "github.com/craigulmer/airlinetracks"
	
)

type marker struct {
	ts  int64    //10
	lat float64  //1
	lon float64  //2
   alt float64  //4
	src string   //11
	dst string   //12
	fid string   //13 airline code
	fin string   //9  registratioin
}
type Markers []*marker

func (m Markers) Len() int { return len(m) }
func (m Markers) Swap(i,j int) { m[i],m[j] = m[j],m[i] }

type byTime struct{ Markers }
func (m byTime) Less(i,j int) bool { return m.Markers[i].ts < m.Markers[j].ts }

func (m1 marker) ToString() string {
	return fmt.Sprintf("%d\t%f\t%f\t%f\t%s\t%s", m1.ts, m1.lat, m1.lon, m1.alt, m1.fin, m1.fid)
}


func (m1 *marker) MetaIsSame(m2 *marker) bool {
	return m1.src==m2.src &&
		m1.dst==m2.dst &&
		m1.fid==m2.fid	
}

func (m Markers) ToATXMeta() []atx.Meta {
	//fmt.Println("marker num points: ",len(m))
	meta := make([]atx.Meta,0,len(m))
	for i,o := range m {
		if i>0 && m[i-1].MetaIsSame(o) { continue }
		mnew:= atx.Meta{First_point_id:i, Flt:o.fid, Src:o.src, Dst:o.dst}
		meta = append(meta,mnew)
	}
	return meta
}
func (m Markers) ToATXPoints() []atx.Point {
	points := make([]atx.Point, len(m), len(m))
	for i,o := range m {
		points[i] = atx.Point{X: o.lon, Y: o.lat, Z: o.alt, T: o.ts}
	}
	return points
}
func (m Markers) ToATXTrack(aid string) *atx.Track {
	t := atx.Track{Aid:aid, Fin: m[0].fin, Fid: m[0].fid, 
		Points: m.ToATXPoints(),
		Meta: m.ToATXMeta() }
	return &t
}


func (m Markers) ToString(id int) {
	for _,o := range m {
		fmt.Printf("%d\t%d\t%f\t%f\t%d\n", id, o.ts, o.lat, o.lon, o.alt);
	}
}
func (m Markers) ToWKT() string {
	var buffer bytes.Buffer
	buffer.WriteString("LINESTRING (");
	is_first := true;
	for _,o := range m {
		if(!is_first){
			buffer.WriteString(", ");
		} else {
			is_first = false;
		}
		s := fmt.Sprintf("%f %f %d %d", o.lon, o.lat, o.alt, o.ts);
		buffer.WriteString(s);
		//buffer.WriteString(o.lon+" "+o.lat+" "+o.alt+" ",o.ts);
		//fmt.Printf("%d\t%d\t%f\t%f\t%d\n", id, o.ts, o.lat, o.lon, o.alt);
	}
	buffer.WriteString(")");
	return buffer.String();
}



