package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Dumper struct {
	offset      int
	offsetWidth int
	nColumns    int
}

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

func (this *Dumper) Fdump(path string, in *os.File, ierr error, out *os.File, oerr error) {
	f := makefmt(this.nColumns, this.offsetWidth)
	bytes := make([]byte, this.nColumns)
	args := make([]any, this.nColumns+1)
	fmt.Fprintf(out, "%v:\n", path)

	var n int
	n, ierr = in.Read(bytes[:])
	for ierr != io.EOF {
		for i, e := range bytes {
			args[i+1] = e
		}

		args[0] = this.offset
		this.offset += n

		if n == this.nColumns {
			fmt.Fprintf(out, f, args...)
		} else {
			g := makefmt(n, this.offsetWidth)
			fmt.Fprintf(out, g, args[:n+1]...)
		}
		n, ierr = in.Read(bytes[:])
	}
}

func (this *Dumper) Dump(path string) {
	file, err := os.Open(path)
	this.Fdump(path, file, err, os.Stdout, nil)
}

func main() {
	dumper := Dumper{
		offset:      0,
		offsetWidth: 4,
		nColumns:    16,
	}

	if len(os.Args) == 1 {
		dumper.Fdump("<stdin>", os.Stdin, nil, os.Stdout, nil)
	}

	for _, path := range os.Args[1:] {
		dumper.Dump(path)
	}
}
