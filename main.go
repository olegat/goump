package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func repeat(n int, value string) []string {
	result := make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = value
	}
	return result
}

func makefmt(nBytesPerLines int, offsetWidth int) string {
	offsetFmt := "%0" + strconv.Itoa(offsetWidth) + "x"
	dumpFmt := strings.Join(repeat(nBytesPerLines, "%02x"), " ")
	return offsetFmt + "  " + dumpFmt + "\n"
}

func Fdump(path string, in *os.File, ierr error, out *os.File, oerr error) {
	const ncol = 16 // number of bytes to display per line
	f := makefmt(ncol, 4)
	off := 0
	bytes := [ncol]byte{0}
	fmt.Fprintf(out, "%v:\n", path)
	for ierr != io.EOF {
		var n int
		n, ierr = in.Read(bytes[:])
		//fmt.Fprintf(out, f, off, bytes[:]...)
		fmt.Fprintf(out, f, off, bytes[0], bytes[1], bytes[2], bytes[3], bytes[4], bytes[5], bytes[6], bytes[7], bytes[8], bytes[9], bytes[10], bytes[11], bytes[12], bytes[13], bytes[14], bytes[15])
		off += n
	}
}

func Dump(path string) {
	file, err := os.Open(path)
	Fdump(path, file, err, os.Stdout, nil)
}

func main() {
	if len(os.Args) == 1 {
		Fdump("<stdin>", os.Stdin, nil, os.Stdout, nil)
	}

	for _, path := range os.Args[1:] {
		Dump(path)
	}
}
