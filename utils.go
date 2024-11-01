package main

import (
	"bytes"
	"errors"
	"strconv"
	"time"
)

// Helper function to remove empty lines
func filterNonEmpty(lines [][]byte) [][]byte {
	var result [][]byte
	for _, line := range lines {
		if len(line) > 0 {
			result = append(result, line)
		}
	}
	return result
}

func getCommandAndArgs(data []byte) (string, []string, error) {
	lines := bytes.Split(data, []byte("\r\n"))
	lines = filterNonEmpty(lines)
	if len(lines) <= 2 {
		return "", nil, errors.New("incomplete command")
	}
	command := string(bytes.ToUpper(lines[2]))
	var args []string
	for i := 4; i < len(lines); i += 2 {
		args = append(args, string(lines[i]))
	}
	return command, args, nil
}

// Helper function to determine if the command is complete
func isCompleteRedisCommand(data []byte) bool {
	lines := bytes.Split(data, []byte("\r\n"))
	// Redis protocol starts with an array indicator like "*3"
	if len(lines) > 0 && len(lines[0]) > 0 && lines[0][0] == '*' {
		// Get the total number of elements expected
		numElements, _ := strconv.Atoi(string(lines[0][1:]))
		expectedLines := 2 + 2*numElements
		// Check if the number of lines matches the expected number of elements
		return len(lines) == expectedLines
	}
	return false
}

// function will help in converting a string value into its int equivalent, if it errors, it will return the same
func readIntValue(input string) (int, error) {
	if intVal, err := strconv.Atoi(input); err != nil {
		return 0, err
	} else {
		return intVal, nil
	}
}

// SetInterval will fire the passed function at every <interval> passed
func SetInterval(interval time.Duration, task func()) chan bool {
	stop := make(chan bool)
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				task()
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
	return stop
}
