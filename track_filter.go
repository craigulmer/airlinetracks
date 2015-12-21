package airlinetracks

import (
//	"fmt"
	"strings"
)

type TrackFilterVariable int
const (
	NONE TrackFilterVariable = iota 
	FID 
	FIN
	AID
	TRACK
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
	lleft       Point
	uright      Point
}
func NewFilterByCarrier(carrier string) TrackFilter{
	return TrackFilter{Search_var:FID, Search_op:BEGINS_WITH, Search_term:carrier }
}
func NewFilterByRegion(region string) TrackFilter{
	lleft,uright := GetNamedRegionBoundingBox(region)
	return TrackFilter{TRACK, CONTAINS, "", lleft, uright}
}

func (t *Track) MetaFilterMatch(filter TrackFilter) bool {
	if filter.Search_op == NOP || filter.Search_var == NONE {
		return true
	}

	var var_val string
	switch filter.Search_var {
	case FID: var_val = t.Fid
	case FIN: var_val = t.Fin
	case AID: var_val = t.Aid
	case TRACK: return true
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

	return hit
}


func (t *Track) PointFilterMatch(filter TrackFilter) bool{
	if filter.Search_var != TRACK {
		return true
	}
	return t.AnyPointsInsideRegion2D(filter.lleft, filter.uright)

}
