package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var aofLogsEligibleCmds = map[string]bool{
	"SET":    true,
	"INCR":   true,
	"DECR":   true,
	"EXPIRE": true,
	"LPUSH":  true,
	"LPOP":   true,
	"RPOP":   true,
}

func loadAOFFile() error {
	fileDataBytes, err := os.ReadFile("aofLogs.aof")
	if err != nil {
		return err
	}
	fileData := string(fileDataBytes)
	commandsToExecute := strings.Split(fileData, "\n")
	for _, commandToExecute := range commandsToExecute {
		command := strings.Split(commandToExecute, " ")[0]
		args := strings.Split(commandToExecute, " ")[1:]
		executeCommand(command, args, true)
	}
	return nil

}

func loadFromAOF() error {
	if _, err := os.Stat("aofLogs.aof"); os.IsNotExist(err) {
		return nil
	} else {
		return loadAOFFile()
	}
}

func addToAOFLogs(command string, args []string) {
	if _, ok := aofLogsEligibleCmds[command]; !ok {
		return
	}
	argsJoin := strings.Join(args, " ")
	aofData := fmt.Sprintf("%s %s", command, argsJoin)
	if _, err := os.Stat("aofLogs.aof"); os.IsNotExist(err) {
		// err := os.WriteFile("aofLogs.aof", []byte(aofData), os.ModeAppend)
		err := os.WriteFile("aofLogs.aof", []byte(aofData), 0644)
		if err != nil {
			log.Println("Error while writing to AOF logs file: ", err.Error())
			return
		}
	} else {
		file, err := os.OpenFile("aofLogs.aof", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Error while opening AOF logs file:", err.Error())
			return
		}
		defer file.Close()
		_, err = file.WriteString(fmt.Sprintf("\n%s", aofData))
		if err != nil {
			log.Println("Error while appending to AOF logs file: ", err.Error())
			return
		}
	}
	log.Println("AOF Logs file appended")
}
