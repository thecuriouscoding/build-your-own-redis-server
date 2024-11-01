package main

import (
	"flag"
	"time"
)

type flags struct {
	persistence      string
	snapshotInterval time.Duration
}

func getFlags() flags {
	flg := &flags{}
	flag.StringVar(&flg.persistence, "persistence", "inmemory", "What kind of persistence should be used: inmemory, snapshot?")
	flag.DurationVar(&flg.snapshotInterval, "snapshot-interval", 5*time.Second, "Snapshot interval time")
	flag.Parse()
	return *flg
}
