package main

// Values:
//  0 : registration (A9EC06)
//  1 : lat
//  2 : lon
//  3 : Track (degrees)
//  4 : altitude (ft)
//  5 : speed (kt?)
//  6 : ? "0000"
//  7 : Radar station
//  8 : aircraft type (B753)
//  9 : registration (N73860)
// 10 : Last heard time
// 11 : Source airport (IAD)
// 12 : Dest airport (SFO)
// 13 : Airline code (UA1431)
// 14 : ? 0
// 15 : ? 0 Vertical speed?
// 16 : Airlinecode (UAL1431)
// 17 : Estimated arrival time

import (
	"encoding/json"
	//"fmt"
	"os"
)

func parseFile(src_file string, mm_ok map[string]Markers, bad_markers *Markers, 
	fin2aid map[string]string, fid2aid map[string]map[string]bool) (int, int) {

	imported := 0
	discarded := 0
	
	fi, err := os.Open(src_file)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	dec := json.NewDecoder(fi)
	for {
		var entries map[string]interface{}
		if err := dec.Decode(&entries); err != nil {
			break
		} // panic(err) }

		for ref := range entries {
			//Ref is json id
			if ref == "full_count" { continue }
			if ref == "version"    { continue }

			//fmt.Println(f.Name()+" "+ref);
			vals := entries[ref].([]interface{})

			aid := vals[0].(string)  //aid
			fin := vals[9].(string)  //tail fin number
			fid := vals[13].(string) //airline code
			new_marker := marker{
				ts:  int64(vals[10].(float64)),
				lat: vals[1].(float64),
				lon: vals[2].(float64),
				alt: vals[4].(float64),
				src: vals[11].(string),
				dst: vals[12].(string),
				fid: fid,               //airline code
				fin: fin}               //registration

			if aid!="" {				
				if fin!="" {
					fin2aid[fin] = aid
				}
				if fid!="" {
					tmp_aids, ok:=fid2aid[fid]
					if !ok {
						tmp_aids = make(map[string]bool)
						fid2aid[fid]=tmp_aids
					} else {
						//fmt.Println("Multiple aids for fid. Fid:",fid,"new aids:",aid, "Num Aids:", len(tmp_aids)+1)
					}
					tmp_aids[aid]=true
				}

				marker_array := mm_ok[aid] //ref
				mm_ok[aid] = append(marker_array, &new_marker)
				imported++

			} else {
				if fin=="" {
					//No hope for planes without fins
					//fmt.Println("Bad: ",new_marker.ToString())
					discarded++
				} else {
					*bad_markers = append(*bad_markers, &new_marker)
				}
			}
		}
	}
	return imported, discarded
}

