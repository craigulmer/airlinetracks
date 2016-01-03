package main

import (
	"flag"
	"fmt"
	atx "github.com/craigulmer/airlinetracks"
	"io/ioutil"
	"sort"
	"strings"
)


func parseDayDirectory(src_dir string, dst_file_name string) {

	discarded := 0
	imported := 0
	var bad_markers Markers
	mm_ok := make(map[string]Markers)
	fin2aid := make(map[string]string)
	fid2aid := make(map[string]map[string]bool)

	files, _ := ioutil.ReadDir(src_dir)
	for _, f := range files {
		tmp_imported, tmp_discarded := parseFile(src_dir+"/"+f.Name(), mm_ok, &bad_markers, fin2aid, fid2aid)
		imported += tmp_imported
		discarded += tmp_discarded

		fmt.Println("Num Ok Planes:", len(mm_ok), 
			"Conflicted Points:",len(bad_markers),
			"Imported Points:",imported,
			"Discarded Points:",discarded)
	}

	//Sort every entry's vals by time
	tracks := make([]*atx.Track, len(mm_ok))
	spot := 0
	for aid, m := range mm_ok {
		//fmt.Println("#\t",id,"\tFLIGHT\t",k,"\t",srcdst[k] )
		sort.Sort(byTime{m})
		tracks[spot] = m.ToATXTrack(aid)
		fmt.Println(tracks[spot].ToString())
		spot++
	}

	fmt.Println("Total Ok:", imported, "Total Discarded:", discarded)

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
