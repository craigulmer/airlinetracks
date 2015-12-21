package airlinetracks

import (
	"fmt"
	"os"
	"compress/gzip"
	"strings"
	"bufio"
)

type TrackWriter struct {	
	fo  *os.File
	gz  *gzip.Writer
	wr  *bufio.Writer
	TotalTracks int
	ValidTracks int
}

func OpenWriter(filename string) (*TrackWriter, error){
	var err error
	tw :=new(TrackWriter)

	tw.fo, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0660)
	if err != nil {
		return nil,err
	}
	
	if strings.HasSuffix(filename,".gz") {
		tw.gz = gzip.NewWriter(tw.fo)
		tw.wr = bufio.NewWriter(tw.gz)
	} else {
		tw.wr = bufio.NewWriter(tw.fo)
	}
	return tw,nil
}

func (tw *TrackWriter) Close(){
	tw.wr.Flush()
	tw.gz.Close()
	tw.fo.Close()
}

func (tw *TrackWriter) Append(t *Track){

	s := fmt.Sprintf("%d\t%s|%s|%s\t%s\n",
		len(t.Points),t.Aid,t.Fin,t.Fid,
		t.GetWKT())
	tw.wr.WriteString(s)
}
