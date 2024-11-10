package main

// command handler gets called on initial COMMAND command sent from redis-cli like client
func addCOMMANDHandler() {
	commandHandler := func(args []string) string {
		return "+OK\r\n"
	}
	handlers["COMMAND"] = commandHandler
}

func addBasicHandlers() {
	addCOMMANDHandler()
}
