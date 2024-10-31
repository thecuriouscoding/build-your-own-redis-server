package main

import "flag"

type flags struct {
	persistence string
}

func getFlags() flags {
	flg := &flags{}
	flag.StringVar(&flg.persistence, "persistence", "inmemory", "What kind of persistence should be used: inmemory, snapshot?")
	flag.Parse()
	return *flg
}
