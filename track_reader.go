package airlinetracks

import (
	"io"
	"os"
	"compress/gzip"
	"strings"
	"bufio"
)



type TrackReader struct {	
	fi io.Reader
	Rdr *bufio.Reader
	filter TrackFilter
	TotalTracks int
	ValidTracks int
}




func OpenWithFilter(filename string, filter TrackFilter) (*TrackReader, error) {

	var err error
	tr :=new(TrackReader)
	tr.filter = filter
	
	tr.fi, err = os.Open(filename)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(filename,".gz") {
		gzipReader,err := gzip.NewReader(tr.fi)
		if err!=nil{
			return nil, err
		}
		tr.Rdr = bufio.NewReader(gzipReader)
		//fmt.Println("Open gzip")

	} else {
		tr.Rdr = bufio.NewReader(tr.fi)
		//fmt.Println("Open regular")
	}

	return tr, nil
}
func Open(filename string) (*TrackReader, error) {
	return OpenWithFilter(filename, TrackFilter{})
}


func (ts *TrackReader) Close(){
	//ts.rd.Close()
	//ts.fi.Close()
}

func (ts *TrackReader) GetNext() (*Track, error){

	for {
		ln,err := ts.Rdr.ReadString('\n')
		//fmt.Println(ln)
		if err!=nil {
			return nil,err
		}

		ts.TotalTracks++
		t,err :=ParseTrackLine(ln, ts.filter)
		if(t!=nil){
			ts.ValidTracks++
			return t,nil
		}
	}
	return nil,nil


}
