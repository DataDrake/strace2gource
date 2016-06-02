package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

var r *regexp.Regexp
var t0 int64 = -1
func ParseLine(line string) {
	match := r.FindStringSubmatch(line)

	if len(match) == 5 {
		process := match[1]
		operation := 'M'
		filename := ""
		found := true
		switch match[3] {
		case "open","openat","mkdir","mkdirat":
			operation = 'A'
			filename = strings.Split(match[4],"\"")[1]
		case "read","lseek","fstat","getdents":
			operation = 'A'
			if strings.Contains(match[4],"<") && strings.Contains(match[4],">"){
				filename = strings.Split(match[4], "<")[1]
				filename = strings.Split(filename, ">")[0]
			}
		case "close","fcntl","write":
			if strings.Contains(match[4],"<") && strings.Contains(match[4],">"){
				filename = strings.Split(match[4], "<")[1]
				filename = strings.Split(filename, ">")[0]
			}
		default:
			found = false
		}

		if found {
			timestamp := strings.Join([]string{"Jan 12",match[2]}," ")
			unix,_ := time.Parse(time.StampMicro,timestamp)
			if t0 == -1 {
				t0 = unix.UnixNano()
			}
			fmt.Printf("%d|%s|%c|%s\n", (unix.UnixNano() - t0) / 1000, process, operation, filename)
		}
	}
}

func ReadGzFile(fi *os.File) error {

	fz, err := gzip.NewReader(fi)
	if err != nil {
		return err
	}
	defer fz.Close()

	buf := bufio.NewReaderSize(fz,10000)

	s,err := buf.ReadString('\n')
	for err == nil {
		ParseLine(s)
		s,err = buf.ReadString('\n')
	}
	return err
}



func main() {
	r, _= regexp.Compile("^(\\d+) (\\S+.\\d+) (\\w+)\\((.*)\\)")

	if len(os.Args) > 1 {
		fi, err := os.Open(os.Args[1])
		if err != nil {
			return
		}
		defer fi.Close()
		ReadGzFile(fi)
	} else {
		ReadGzFile(os.Stdin)
	}
}

