package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/Awadabang/feuer/client"
)

const maxN = 10000000
const maxBufferSize = 1024 * 1024

func main() {
	s := client.NewSimple([]string{"localhost"})

	var b bytes.Buffer

	for i := 0; i <= maxN; i++ {
		fmt.Fprintf(&b, "%d\n", i)

		if b.Len() >= maxBufferSize {
			if err := s.Send(b.Bytes()); err != nil {
				log.Fatalf("Send error: %v", err)
			}

			b.Reset()
		}
	}
}
