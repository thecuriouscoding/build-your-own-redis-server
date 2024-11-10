package main

import "fmt"

// command handler gets called on initial COMMAND command sent from redis-cli like client
func addCOMMANDHandler() {
	commandHandler := func(args []string) string {
		return "+OK\r\n"
	}
	handlers["COMMAND"] = commandHandler
}

func addSETHandler() {
	// SET [key] [value]
	setHandler := func(args []string) string {
		key := args[0]
		value := args[1]
		dataStore[key] = valueFormat{
			ValueType: VALUE_TYPE_STRING,
			Value:     value,
		}
		return "+OK\r\n"
	}
	handlers["SET"] = setHandler
}

func addGETHandler() {
	// GET [key]
	getHandler := func(args []string) string {
		key := args[0]
		if val, ok := dataStore[key]; !ok {
			return "$-1\r\n"
		} else {
			return fmt.Sprintf("$%d\r\n%s\r\n", len(val.Value.(string)), val.Value.(string))
		}
	}
	handlers["GET"] = getHandler
}

func addDELHandler() {
	// DEL [keys...]
	delHandler := func(args []string) string {
		var totalDeleted = 0
		for _, key := range args {
			if _, ok := dataStore[key]; !ok {
				continue
			} else {
				totalDeleted++
				delete(dataStore, key)
			}
		}
		return fmt.Sprintf(":%d\r\n", totalDeleted)
	}
	handlers["DEL"] = delHandler
}

func addBasicHandlers() {
	addCOMMANDHandler()
	addSETHandler()
	addGETHandler()
	addDELHandler()
}
