package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/krezac/robot-irena/vectornav"
	"github.com/tarm/serial"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/ajax.html"))
	t.Execute(w, "user") // merge.
}

func main() {

	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	var data vectornav.YMRDataFull

	go func() {
		for {
			reader := bufio.NewReader(s)
			reply, err := reader.ReadBytes('\x0a')
			//fmt.Println("Got data")
			if err != nil {
				panic(err)
			}
			//fmt.Printf("%q\n", reply)
			vectornav.ParseYMR(string(reply), &data)
		}
	}()
	//vectornav.StartJSONServer(8080, &data)

	http.HandleFunc("/imu", vectornav.GetVectornavHTTPHandler(&data.YMRData))
	//http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("."))))
	//http.Handle("/", http.FileServer(http.Dir("/tmp")))
	http.Handle("/static_content/", http.StripPrefix("/static_content/", http.FileServer(http.Dir("static_content"))))
	http.HandleFunc("/", mainHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil)
}
