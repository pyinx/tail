package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
)

var (
	filename *string
	line     *int64
)
var f *os.File

func init() {
	filename = flag.String("filename", "php_err.log", "tail filename")
	line = flag.Int64("line", 10, "tail lines")
}

func tail(offset int64, n int64) ([]byte, error) {
	buf := make([]byte, offset)
	if n == 1 {
		f.Seek(-offset, 2)
	} else {
		// f.Seek(-offset*(n-1), 2)
		f.Seek(-offset, 1)
	}
	_, err := f.Read(buf)
	return buf, err
}

func Tail(filename string, l int64) {
	var err error
	f, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var (
		n         int64 //loop read index
		filesize  int64
		output    [][]byte //final output
		data      []byte   // tmp output
		lines     [][]byte // split data to slice
		count     int64    //read lines
		endTag    bool
		buferSize = int64(100) * l // read bufer size
		fi        os.FileInfo      //read file info
	)
	fi, _ = f.Stat()
	filesize = fi.Size()
	for n = 1; ; n++ {
		if (n * buferSize) <= filesize {
			data, _ = tail(buferSize, n)
		} else {
			data, _ = tail(filesize, n)
			endTag = true
		}
		lines = bytes.Split(data, []byte("\n"))
		output = append(lines, output...)
		count = int64(len(output))
		if count > l {
			endTag = true
			output = output[count-l-1 : count-1]
		}
		if endTag {
			for _, text := range output {
				fmt.Println(string(text))
			}
			// for _, _ = range output {
			// }
			return
		}
	}

}

func main() {
	flag.Parse()
	Tail(*filename, *line)
}
