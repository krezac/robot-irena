package main

import (
	"bufio"
	"fmt"
	"log"

	"github.com/tarm/serial"
)

func main() {
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	//n, err := s.Write([]byte("test"))
	//if err != nil {
	//        log.Fatal(err)
	//}
	for {
		reader := bufio.NewReader(s)
		reply, err := reader.ReadBytes('\x0a')
		if err != nil {
			panic(err)
		}
		fmt.Printf("%q\n", reply)

		//       buf := make([]byte, 128)
		//       n, err := s.Read(buf)
		//       if err != nil {
		//               log.Fatal(err)
		//       }
		//       log.Printf("%q", buf[:n])
	}
}
