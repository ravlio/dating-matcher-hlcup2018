package main

import (
	"flag"
	"github.com/ravlio/highloadcup2018"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

var host = flag.String("host", "0.0.0.0", "hostname")
var port = flag.String("port", "8181", "port")
var profile = flag.Bool("profile", true, "use profile")
var dpath = flag.String("dpath", "/tmp/data/", "data path")

func main() {
	flag.Parse()
	if *profile {
		go func() {
			log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
		}()
	}

	db := hl.NewDB(&hl.Options{Indexer: hl.OptionsIndexer{Enabled: true, WorkerCount: 4, BufferSize: 1}, Debug: true})
	err := db.Start()
	if err != nil {
		log.Fatal("error starting database")
	}

	err = db.LoadFromJSON(*dpath)
	if err != nil {
		log.Fatalf("error loading json: %s", err.Error())
	}

	srv := &hl.HTTPServer{DB: db, Addr: net.JoinHostPort(*host, *port)}

	err = srv.Run()
	if err != nil {
		log.Fatalf("error running server: %v", err)
	}

}
