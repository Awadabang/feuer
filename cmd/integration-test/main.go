package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Awadabang/feuer/client"
)

const maxN = 10000000
const maxBufferSize = 1024 * 1024

func main() {
	s := client.NewSimple([]string{"localhost"})

	want, err := send(s)
	if err != nil {
		log.Fatalf("Send error: %v", err)
	}

	got, err := receive(s)
	if err != nil {
		log.Fatalf("Receive error: %v", err)
	}

	if want != got {
		log.Fatalf("The expected sum %d is not equal to the actual sum %d", want, got)
	}

	log.Printf("The test passed!")
}

func send(s *client.Simple) (sum int64, err error) {
	var b bytes.Buffer

	for i := 0; i <= maxN; i++ {
		sum += int64(i)

		fmt.Fprintf(&b, "%d\n", i)
		if b.Len() >= maxBufferSize {
			if err := s.Send(b.Bytes()); err != nil {
				return 0, err
			}

			b.Reset()
		}
	}

	if b.Len() != 0 {
		if err := s.Send(b.Bytes()); err != nil {
			return 0, err
		}
	}

	return 0, nil
}

func receive(s *client.Simple) (sum int64, err error) {
	buf := make([]byte, maxBufferSize)

	for {
		res, err := s.Receive(buf)
		if err == io.EOF {
			return sum, nil
		} else if err != nil {
			return 0, err
		}

		ints := strings.Split(res, "\n")

	}

	return 0, nil
}
