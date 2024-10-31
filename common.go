package main

import (
	"fmt"
	"time"
)

type valueFormat struct {
	ValueType string      `json:"value_type"`
	Value     interface{} `json:"value"`
}

type argumentLengthCheck struct {
	toCheck   bool
	argLength int
	errorCode int
}

var (
	dataStore      = make(map[string]valueFormat)
	handlers       = make(map[string]func(args []string) string)
	keyExpirations = make(map[string]time.Time)
)
var (
	ERR_NO_ERROR                            = 0
	ERR_WRONG_NUMBER_OF_ARGUMENTS           = 1
	ERR_SYNTAX_ERROR                        = 2
	ERR_WRONG_TYPE_OPERATION                = 3
	ERR_VALUE_NOT_INTEGER                   = 4
	ERR_VALUE_OUT_OF_RANGE_MUST_BE_POSITIVE = 5
)

var (
	COMMAND_SET    = "set"
	COMMAND_DEL    = "del"
	COMMAND_GET    = "get"
	COMMAND_EXPIRE = "expire"
	COMMAND_TTL    = "ttl"
	COMMAND_INCR   = "incr"
	COMMAND_DECR   = "decr"
	COMMAND_LPUSH  = "lpush"
	COMMAND_RPUSH  = "rpush"
	COMMAND_LRANGE = "lrange"
	COMMAND_LPOP   = "lpop"
	COMMAND_RPOP   = "rpop"
)

var (
	VALUE_TYPE_STRING = "string"
	VALUE_TYPE_LIST   = "list"
)

func getErrorMessage(code int, command string) string {
	switch code {
	case ERR_SYNTAX_ERROR:
		return "-ERR syntax error\r\n"
	case ERR_WRONG_NUMBER_OF_ARGUMENTS:
		return fmt.Sprintf("-ERR wrong number of arguments for '%s' command\r\n", command)
	case ERR_WRONG_TYPE_OPERATION:
		return "-ERR WRONGTYPE Operation against a key holding the wrong kind of value\r\n"
	case ERR_VALUE_NOT_INTEGER:
		return "-ERR value is not an integer or out of range\r\n"
	case ERR_VALUE_OUT_OF_RANGE_MUST_BE_POSITIVE:
		return "-ERR value is out of range, must be positive\r\n"
	default:
		return ""
	}
}
func checkKeyExpiration(key string) bool {
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

func deleteExpiredKey(key string) {
	delete(dataStore, key)
	delete(keyExpirations, key)
}

func handlerValidations(minArgument argumentLengthCheck, maxArgument argumentLengthCheck, args []string, command string) (bool, string) {
	if minArgument.toCheck {
		if len(args) < minArgument.argLength {
			return false, getErrorMessage(minArgument.errorCode, command)
		}
	}
	if maxArgument.toCheck {
		if len(args) > maxArgument.argLength {
			return false, getErrorMessage(maxArgument.errorCode, command)
		}
	}
	return true, ""
}
