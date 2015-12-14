package main

import (
	"fmt"
	"flag"
	atx "github.com/craigulmer/airlinetracks"
)


func readFile(fname string, carrier string, region string, verbose bool){

	var trdr *atx.TrackReader
	var err error
	var op_name string
	if carrier!="" && region!="" {
		panic("Currently, can only filter on ONE thing")
	}

	switch {
	case region!="":  
		trdr, err = atx.OpenWithFilter(fname, atx.NewFilterByRegion(region));
		op_name="Region: "+region
	case carrier!="": 
		trdr, err = atx.OpenWithFilter(fname, atx.NewFilterByCarrier(carrier)); 
		op_name="Carrier: "+carrier
	default:          
		trdr, err = atx.Open(fname)
	}
	if err!=nil {
		fmt.Println("Error opening "+fname)
		return
	}

	//Walk through each track until nothing left
	for {
		track,_ := trdr.GetNext()
		if track==nil { break }
		if verbose {
			fmt.Println("Track: "+track.ToString())		
		}
	}

	//Stats are maintained in the reader
	fmt.Println("ValidTracks:",trdr.ValidTracks,"TotalTracks:",trdr.TotalTracks, op_name)
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
