package airlinetracks

import (
	"fmt"
	"bytes"
	"strings"
)

type Meta struct {
	First_point_id int
	Flt string
	Src string
	Dst string
}


type Track struct {
	Aid string //Hex AID code for plane
	Fin string //Tailfin for plane
	Fid string //One of Airline's flight id
	Points []Point
	Meta []Meta
}

type Tracks []*Track
func (t Tracks) Len() int { return len(t)}
func (t Tracks) Swap(i,j int) { t[i],t[j] = t[j], t[i] }

type ByAID struct { Tracks }
func (t ByAID) Less(i,j int) bool {return t.Tracks[i].Aid < t.Tracks[j].Aid }


func (t *Track) GetWKT() string {
	var buffer bytes.Buffer
	buffer.WriteString("LINESTRING (")
	is_first := true
	for _,p :=range t.Points {
		if(!is_first){
			buffer.WriteString(", ")
		} else {
			is_first = false
		}
		buffer.WriteString(p.ToString())
	}
	buffer.WriteString(")")
	return buffer.String()
}


func (t *Track) ToString() string{
	return "AID: "+t.Aid+" FIN: "+t.Fin+" FID: "+t.Fid+
		fmt.Sprintf(" NumPoints: %v",len(t.Points))
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
	t.Aid = strings.TrimSpace(hdr_args[0])
	t.Fin = strings.TrimSpace(hdr_args[1])
	t.Fid = strings.TrimSpace(hdr_args[2])
	if t.Aid == "" { t.Aid = "-" }
	if t.Fin == "" { t.Fin = "-" }
	if t.Fid == "" { t.Fid = "-" }

	//Check meta filter params, bail if no hit
	hit := t.MetaFilterMatch(filter)
	if(!hit){
		return nil,nil
	}


	//Pull point data out of wkt string
	s:=strings.TrimSpace(wkt)
	s=strings.TrimSuffix(s,")")
	subs:=strings.Split(s,",")
	t.Points = make([]Point, len(subs),len(subs))
	for i,point_string := range subs {
		t.Points[i].ParsePoint(point_string)
	}
	
	//Examine Points
	hit = t.PointFilterMatch(filter)
	if !hit {
		return nil,nil
	}

	return t,nil
}
