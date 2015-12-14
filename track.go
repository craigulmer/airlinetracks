package airlinetracks

import (
	"fmt"
	"strings"
)

//type meta_labels struct {
//	first_track_id int
//	flt string
//	src string
//	dst string
//}

type Track struct {
	aid string //Hex AID code for plane
	fin string //Tailfin for plane
	fid string //Airline's flight id
	points []point
	//meta []meta_labels
}




func (t *Track) ToString() string{
	return "AID: "+t.aid+" FIN: "+t.fin+" FID: "+t.fid+
		fmt.Sprintf(" NumPoints: %v",len(t.points))
}
func ParseTrackLine(line string, filter TrackFilter) (*Track,error) {
	t := new(Track)
	//fmt.Println("Parse is :"+line)

	//Split line based on LINESTRING
	split_line := strings.Split(line, "LINESTRING (")
	if len(split_line)!=2 {
		return nil, nil
	}
	prefix := split_line[0]
	wkt    := split_line[1]

	//Parse middle chunk, separated by |
	s_split := strings.Split(prefix,"\t")
	hdr_args := strings.Split(s_split[1],"|")
	t.aid = strings.TrimSpace(hdr_args[0])
	t.fin = strings.TrimSpace(hdr_args[1])
	t.fid = strings.TrimSpace(hdr_args[2])
	if t.aid == "" { t.aid = "-" }
	if t.fin == "" { t.fin = "-" }
	if t.fid == "" { t.fid = "-" }

	//Check meta filter params, bail if no hit
	hit := t.MetaFilterMatch(filter)
	if(!hit){
		return nil,nil
	}


	//Pull point data out of wkt string
	s:=strings.TrimSpace(wkt)
	s=strings.TrimSuffix(s,")")
	subs:=strings.Split(s,",")
	t.points = make([]point, len(subs),len(subs))
	for i,point_string := range subs {
		t.points[i].ParsePoint(point_string)
	}
	
	//Examine points
	hit = t.PointFilterMatch(filter)
	if !hit {
		return nil,nil
	}

	return t,nil
}
