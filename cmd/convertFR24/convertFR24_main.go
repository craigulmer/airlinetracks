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
	"fmt"
	//"bytes"
	"sort"
	"os"
	"flag"
	"encoding/json"
	"io/ioutil"
	atx "github.com/craigulmer/airlinetracks"
	
)

func parseFile(src_file string, mm map[string]Markers) {

	fi, err := os.Open(src_file);
	if err != nil { panic(err) }
	defer func(){
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
		
	dec := json.NewDecoder(fi);
	for {
		var entries map[string]interface{} 
		if err:= dec.Decode(&entries); err != nil { break } // panic(err) }
			
		for ref := range entries {
			if ref == "full_count" { continue }
			if ref == "version"    { continue }

			//fmt.Println(f.Name()+" "+ref);
			vals:=entries[ref].([]interface{})

			new_marker := marker{
				ts:  int64(vals[10].(float64)), 
				lat: vals[1].(float64),
				lon: vals[2].(float64),
				alt: vals[4].(float64),
				src: vals[11].(string),
				dst: vals[12].(string),
				fid: vals[13].(string), //airline code
				fin: vals[9].(string) }  //registration

			flabel := vals[0].(string); //aid
			marker_array := mm[ flabel ]; //ref
			mm[ flabel ] = append(marker_array, &new_marker);
		}
	}
}

func parseDayDirectory(src_dir string, dst_file_name string){

	mm:=make(map[string]Markers);
	
	files,_ := ioutil.ReadDir(src_dir)
	for _,f := range files {
		parseFile(src_dir+"/"+f.Name(), mm)
		fmt.Println("mm size: ",len(mm))
	}

	
	//Sort every entry's vals by time
	tracks := make([]*atx.Track, len(mm))
	spot:=0
	for aid, m := range mm {
		//fmt.Println("#\t",id,"\tFLIGHT\t",k,"\t",srcdst[k] )
		sort.Sort(byTime{m})
		tracks[spot]=m.ToATXTrack(aid)
		fmt.Println(tracks[spot].ToString())
		spot++
	}
	
	sort.Sort(atx.ByAID{tracks})
	for _,v := range tracks {
		fmt.Println(v.ToString())

	}


}



func main(){

	
	var idir = flag.String("input", "data/2014-03-31", "Input Directory")
	var odir = flag.String("output", "stdout", "Output Directory")
	flag.Parse()
	fmt.Println("Input Directory: "+*idir+" Output Directory: "+*odir)
	
	parseDayDirectory(*idir, *odir)
}
