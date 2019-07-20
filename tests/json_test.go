package tests

import (
	"github.com/ravlio/highloadcup2018"
	"log"
	"testing"
)

func BenchmarkRequestUnmarshal(b *testing.B) {
	dpath := "/Users/ravlio/Downloads/elim_accounts_261218/"

	db := hl.NewDB(&hl.Options{Indexer: hl.OptionsIndexer{Enabled: false, WorkerCount: 0}, Debug: true})
	err := db.Start()
	if err != nil {
		log.Fatal("error starting database")
	}

	err = db.LoadFromJSON(dpath + "data/")
}
