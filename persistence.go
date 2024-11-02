package main

import "log"

func handlePersistence(flags flags) {
	switch flags.persistence {
	case "snapshot":
		log.Println("Persistence mode: snapshot")
		if err := loadLastSnapshot(); err != nil {
			log.Fatal("Snapshot load failed: ", err)
		}
		stop := SetInterval(flags.snapshotInterval, createSnapshot)
		defer func() {
			stop <- true
		}()
	case "aof":
		log.Println("Persistence mode: aof")
		if err := loadFromAOF(); err != nil {
			log.Fatal("Error while executing aof logs: ", err.Error())
		}
	default:
		log.Println("Persistence mode: inmemory")
	}
}
