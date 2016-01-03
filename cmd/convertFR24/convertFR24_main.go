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
	"flag"
	"fmt"
	atx "github.com/craigulmer/airlinetracks"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

func parseFile(src_file string,
	mm map[string]Markers, fin2aid map[string]string) (int, int, int, int) {

	discarded := 0
	imported := 0
	repaired := 0
	duplicates := 0

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

			aid := vals[0].(string) //aid
			fin := vals[9].(string) //fin

			if aid == "" {

				if fin=="" {
					discarded++
					continue
				}

				prv_aid, found := fin2aid[fin]
				if !found {
					discarded++
					continue
				}
				aid = prv_aid
				repaired++
			} else {
				imported++
			}

			new_marker := marker{
				ts:  int64(vals[10].(float64)),
				lat: vals[1].(float64),
				lon: vals[2].(float64),
				alt: vals[4].(float64),
				src: vals[11].(string),
				dst: vals[12].(string),
				fid: vals[13].(string), //airline code
				fin: fin}               //registration

			fin2aid[fin] = aid
			marker_array := mm[aid] //ref
			num_pts:=len(marker_array)
			if num_pts>0 {
				if marker_array[num_pts-1].ts == new_marker.ts {
					if !new_marker.IsIdentical(marker_array[num_pts-1]) {
						fmt.Println("Duplicates, but not identical: (previous first)",aid,"-->",fin,"\n",
							marker_array[num_pts-1].ToString(),"\n",
							new_marker.ToString())
					} else {
						fmt.Println("Duplicate IS identical")
					}

					duplicates++
					continue
				}
			}
			mm[aid] = append(marker_array, &new_marker)
		}
	}
	return discarded, imported, repaired, duplicates
}

func parseDayDirectory(src_dir string, dst_file_name string) {

	discarded := 0
	imported := 0
	repaired := 0
	duplicates := 0
	mm := make(map[string]Markers)
	fin2aid := make(map[string]string)

	files, _ := ioutil.ReadDir(src_dir)
	for _, f := range files {
		tmp_discarded, tmp_imported, tmp_repaired, tmp_duplicates := parseFile(src_dir+"/"+f.Name(), mm, fin2aid)
		discarded += tmp_discarded
		imported += tmp_imported
		repaired += tmp_repaired
		duplicates += tmp_duplicates
		fmt.Println("mm size: ", len(mm), " Total Ok: ", imported, 
			"Total Discarded: ", discarded, "Total Repaired: ", repaired,
			"Total Duplictes: ", duplicates)
	}

	//Sort every entry's vals by time
	tracks := make([]*atx.Track, len(mm))
	spot := 0
	for aid, m := range mm {
		//fmt.Println("#\t",id,"\tFLIGHT\t",k,"\t",srcdst[k] )
		sort.Sort(byTime{m})
		tracks[spot] = m.ToATXTrack(aid)
		fmt.Println(tracks[spot].ToString())
		spot++
	}

	fmt.Println("Total Ok:", imported, "Total Discarded:", discarded, "Total Duplicates:",duplicates)

	sort.Sort(atx.ByAID{tracks})

	tw, err := atx.OpenWriter(dst_file_name)
	if err != nil {
		fmt.Println("Error opening output file", dst_file_name, err)
		return
	}
	for _, t := range tracks {
		//fmt.Println(v.ToString())
		tw.Append(t)
	}
	tw.Close()

}

func makeOutputFilename(idir string, odir string) string {
	args := strings.Split(idir, "/")
	return odir + "/" + args[len(args)-1] + ".gz"
}

func main() {

	var idir = flag.String("input", "data/2014-03-31", "Input Directory")
	var odir = flag.String("output", "./", "Output Directory")
	flag.Parse()
	fmt.Println("Input Directory: " + *idir + " Output Directory: " + *odir)

	fmt.Println("Val is " + makeOutputFilename(*idir, *odir))
	parseDayDirectory(*idir, makeOutputFilename(*idir, *odir))
}
