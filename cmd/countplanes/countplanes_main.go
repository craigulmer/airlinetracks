package main

import (
	"fmt"
	"flag"
	atx "github.com/craigulmer/airlinetracks"
)


func readFile(fname string, carrier string, region string, verbose bool){

	var trdr *atx.TrackReader
	var err error

	if carrier!="" && region!="" {
		panic("Currently, can only filter on ONE thing")
	}


	if carrier=="" {

		if region=="" {
			trdr, err = atx.Open(fname)
		} else {
			trdr, err = atx.OpenWithFilter(fname, atx.NewFilterByRegion(region))
		}
	} else {
		trdr, err = atx.OpenWithFilter(fname, atx.NewFilterByCarrier(carrier))
	}
	if err!=nil {
		fmt.Println("Error opening "+fname)
		return
	}

	for {
		track,_ := trdr.GetNext()
		if track==nil { break }
		if verbose {
			fmt.Println("Track: "+track.ToString())		
		}
	}

	//Stats are maintained in the reader
	fmt.Println("ValidTracks:",trdr.ValidTracks,"TotalTracks:",trdr.TotalTracks, "CarrierFilter: '"+carrier+"'")
	trdr.Close()
}

func main() {

	var file_src = flag.String("file", "2014-10-28.txt.gz", "Input File")
	var carrier = flag.String("carrier", "", "Filter to a specific carrier")
	var region = flag.String("region","", "Filter to flights that crossed a particular region")
	var verbose = flag.Bool("verbose", false, "Verbose")
	flag.Parse()

	readFile(*file_src, *carrier,  *region, *verbose,)

}
