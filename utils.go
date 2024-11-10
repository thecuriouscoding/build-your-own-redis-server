package main

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

// Helper function to remove empty lines
func filterNonEmpty(lines []string) []string {
	var result []string
	for _, line := range lines {
		if len(line) > 0 {
			result = append(result, line)
		}
	}
	return result
}

func getCommandAndArgs(data string) (string, []string, error) {
	lines := strings.Split(data, ("\r\n"))
	lines = filterNonEmpty(lines)
	if len(lines) <= 2 {
		return "", nil, errors.New("incomplete command")
	}
	command := string(strings.ToUpper(lines[2]))
	var args []string
	for i := 4; i < len(lines); i += 2 {
		args = append(args, string(lines[i]))
	}
	return command, args, nil
}

// Helper function to determine if the command is complete
func isCompleteRedisCommand(data []byte) bool {
	stringData := strings.TrimRight(string(data), " ")
	log.Println("stringData: ", stringData)
	lines := strings.Split(stringData, "\r\n")
	if len(lines) > 0 && len(lines[0]) > 0 && lines[0][0] == '*' {
		numElements, _ := strconv.Atoi(string(lines[0][1:]))
		expectedLines := 2 + 2*numElements
		return len(lines) == expectedLines
	}
	return false
}
