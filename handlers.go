package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type valueFormat struct {
	valueType string
	value     interface{}
}

var dataStore = make(map[string]valueFormat)

type commandHandler func(args []string) string

var handlers = make(map[string]commandHandler)

var keyExpirations = make(map[string]time.Time)

func checkKeyExpiration(key string) bool {
	log.Println("Key expiration value is: ", keyExpirations[key])
	if val, ok := keyExpirations[key]; !ok {
		return false
	} else {
		if time.Now().After(val) {
			return true
		} else {
			return false
		}
	}
}

func addSETHandler() {
	// var setHandler commandHandler
	setHandler := func(args []string) string {
		log.Println("SET handler called with args: ", args)
		if len(args) < 2 {
			return "-ERR wrong number of arguments for 'set' command\r\n"
		}
		key := args[0]
		value := args[1]
		dataStore[key] = valueFormat{
			valueType: "string",
			value:     value,
		}
		log.Println("Data store updated to: ", dataStore)
		return "+OK\r\n"
	}
	handlers["SET"] = setHandler
}

func addGETHandler() {
	getHandler := func(args []string) string {
		log.Println("GET handler called with args: ", args)
		if len(args) < 1 {
			return "-ERR wrong number of arguments for 'get' command\r\n"
		}
		key := args[0]
		log.Println("Data store key value to to args is: ", dataStore[key])
		if val, ok := dataStore[key]; !ok {
			return "$-1\r\n"
		} else {
			if val.valueType != "string" {
				return "$-1\r\n"
			}
			isExpired := checkKeyExpiration(key)
			log.Println("Is key expired: ", isExpired)
			if isExpired {
				delete(dataStore, key)
				delete(keyExpirations, key)
				return "$-1\r\n"
			}
			return fmt.Sprintf("$%d\r\n%s\r\n", len(val.value.(string)), val.value)
		}
	}
	handlers["GET"] = getHandler
}

func addEXPIREHandler() {
	expireHandler := func(args []string) string {
		log.Println("EXPIRE handler called with args: ", args)
		if len(args) < 1 {
			return "-ERR wrong number of arguments for 'expire' command\r\n"
		}
		key := args[0]
		value := args[1]
		if _, ok := dataStore[key]; !ok {
			return ":0\r\n"
		} else {
			// check in redis actual server for invalid value type i.e. not number, what error it gives
			expirationTime, _ := strconv.Atoi(value)
			keyExpirations[key] = time.Now().Add(time.Duration(expirationTime) * time.Second)
		}
		return ":1\r\n"
	}
	handlers["EXPIRE"] = expireHandler
}

func addCOMMANDHandler() {
	commandHandler := func(args []string) string {
		return "+OK\r\n"
	}
	handlers["COMMAND"] = commandHandler
}

func addDELHandler() {
	deleteHandler := func(args []string) string {
		log.Println("DEL handler called with args: ", args)
		if len(args) < 1 {
			return "-ERR wrong number of arguments for 'del' command\r\n"
		}
		key := args[0]
		if _, ok := dataStore[key]; !ok {
			return ":0\r\n"
		} else {
			delete(dataStore, key)
			delete(keyExpirations, key)
			return ":1\r\n"
		}
	}
	handlers["DEL"] = deleteHandler
}

func addHandlers() {
	addCOMMANDHandler()
	addSETHandler()
	addGETHandler()
	addEXPIREHandler()
	addDELHandler()
}
