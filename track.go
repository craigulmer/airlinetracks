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

type TrackFilterVariable int
const (
	NONE TrackFilterVariable = iota 
	FID 
	FIN
	AID
	//TRACK
)
type TrackFilterOp int
const (
	NOP TrackFilterOp = iota
	BEGINS_WITH 
	EQUALS
	ENDS_WITH
	CONTAINS
)

type TrackFilter struct {
	Search_var  TrackFilterVariable
	Search_op   TrackFilterOp
	Search_term string
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

	if filter.Search_var != NONE && filter.Search_op != NOP {
		var var_val string
		switch filter.Search_var {
		case FID: var_val = t.fid
		case FIN: var_val = t.fin
		case AID: var_val = t.aid
		default:
			panic("Unknown type in filter search")
		}

		//fmt.Println("Compare "+var_val+" to "+filter.Search_term)
		var hit bool
		switch filter.Search_op {
		case BEGINS_WITH: hit = strings.HasPrefix(var_val, filter.Search_term)
		case EQUALS:      hit = var_val==filter.Search_term
		case ENDS_WITH:   hit = strings.HasSuffix(var_val, filter.Search_term)
		case CONTAINS:    hit = strings.Contains(var_val, filter.Search_term)
		default:
			panic("Unknown type in filter op")
		}

		if(!hit){
			return nil,nil
		}
	}


	//Pull point data out of wkt string
	s:=strings.TrimSpace(wkt)
	s=strings.TrimSuffix(s,")")
	subs:=strings.Split(s,",")
	t.points = make([]point, len(subs),len(subs))
	for i,point_string := range subs {
		t.points[i].ParsePoint(point_string)
	}
	

	return t,nil
}
