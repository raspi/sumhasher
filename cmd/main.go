package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"github.com/raspi/sumhasher"
	"io"
	"os"
)

func main() {
	h := sumhasher.New()

	flag.Parse()

	if flag.NArg() == 0 {
		os.Exit(0)
	}

	fname := flag.Arg(flag.NArg() - 1)

	var (
		f   io.Reader
		err error
	)

	switch fname {
	case `-`:
		f = bufio.NewReader(os.Stdin)
	default:
		tmpf, err := os.Open(fname)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, `%v`, err)
			os.Exit(1)
		}
		defer tmpf.Close()
		f = tmpf
	}

	buffer := make([]byte, h.BlockSize())

	for {
		rb, err := f.Read(buffer)
		if err != nil {

			if errors.Is(err, io.EOF) {
				break
			}

			_, _ = fmt.Fprintf(os.Stderr, `%v`, err)
			os.Exit(1)
		}

		if rb == 0 {
			break
		}

		_, err = h.Write(buffer[:rb])
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, `%v`, err)
			os.Exit(1)
		}
	}

	res := uint64(0)
	result := h.Sum(nil)
	err = binary.Read(bytes.NewReader(result), binary.BigEndian, &res)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, `%v`, err)
		os.Exit(1)
	}

	fmt.Printf(`%*x`, h.Size(), res)
	fmt.Println()
}
