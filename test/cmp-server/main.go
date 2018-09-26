package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rs/cors"
)

func main() {
	startServer()
}

func startServer() {

	router := InitRoutes()
	server := &http.Server{
		Addr:    "0.0.0.0:8088",
		Handler: router,
	}
	log.Println("Listening...")
	server.ListenAndServe()
	log.Println("not Listening........")

}

func InitRoutes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/demo/data/media", func(w http.ResponseWriter, r *http.Request) {

		log.Print("Listener got the request to give report :")

		absPath, _ := filepath.Abs("demo_data_media.json")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		w.Header().Set("Content-type", "application/json")
		w.Header().Set("Access-Control-Expose-Headers", "File-Name")

		filebytes, err := ioutil.ReadFile(absPath)
		if err != nil {
			fmt.Println(err)

		}
		b := bytes.NewBuffer(filebytes)
		if _, err := b.WriteTo(w); err != nil {
			fmt.Fprintf(w, "%s", err)
		}

		w.Write(b.Bytes())
		return
	})

	mux.HandleFunc("/demo/data/operator", func(w http.ResponseWriter, r *http.Request) {

		log.Print("Listener got the request to give report :")

		absPath, _ := filepath.Abs("demo_data_operator.json")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		w.Header().Set("Content-type", "application/json")
		w.Header().Set("Access-Control-Expose-Headers", "File-Name")

		filebytes, err := ioutil.ReadFile(absPath)
		if err != nil {
			fmt.Println(err)

		}
		b := bytes.NewBuffer(filebytes)
		if _, err := b.WriteTo(w); err != nil {
			fmt.Fprintf(w, "%s", err)
		}

		w.Write(b.Bytes())
		return
	})

	mux.HandleFunc("/demo/media/data/save", func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
			return
		}
		f, err := os.OpenFile("demo_data_media.json", os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		return
	})

	mux.HandleFunc("/demo/operator/data/save", func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
			return
		}
		f, err := os.OpenFile("demo_data_operator.json", os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		return
	})

	router := cors.Default().Handler(mux)

	return router
}
