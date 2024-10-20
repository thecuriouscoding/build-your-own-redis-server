package main

import (
	"fmt"
	"strconv"
	"time"
)

type valueFormat struct {
	valueType string
	value     interface{}
}

var (
	dataStore      = make(map[string]valueFormat)
	handlers       = make(map[string]func(args []string) string)
	keyExpirations = make(map[string]time.Time)
)

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

func addSETHandler() {
	setHandler := func(args []string) string {
		if len(args) < 2 {
			return "-ERR wrong number of arguments for 'set' command\r\n"
		}
		if len(args) > 2 {
			return "-ERR syntax error\r\n"
		}
		key := args[0]
		value := args[1]
		dataStore[key] = valueFormat{
			valueType: "string",
			value:     value,
		}
		return "+OK\r\n"
	}
	handlers["SET"] = setHandler
}

func addGETHandler() {
	getHandler := func(args []string) string {
		if len(args) < 1 {
			return "-ERR wrong number of arguments for 'get' command\r\n"
		}
		key := args[0]
		if val, ok := dataStore[key]; !ok {
			return "$-1\r\n"
		} else {
			if val.valueType != "string" {
				return "$-1\r\n"
			}
			isExpired := checkKeyExpiration(key)
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
		if len(args) < 1 {
			return "-ERR wrong number of arguments for 'expire' command\r\n"
		}
		key := args[0]
		value := args[1]
		if _, ok := dataStore[key]; !ok {
			return ":0\r\n"
		} else {
			expirationTime, err := strconv.Atoi(value)
			if err != nil {
				return ":0\r\n"
			}
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
		if len(args) < 1 {
			return "-ERR wrong number of arguments for 'del' command\r\n"
		}
		for _, key := range args {
			if _, ok := dataStore[key]; !ok {
				return ":0\r\n"
			} else {
				isExpired := checkKeyExpiration(key)
				if isExpired {
					delete(dataStore, key)
					delete(keyExpirations, key)
					return ":0\r\n"
				}
				delete(dataStore, key)
				delete(keyExpirations, key)
				return ":1\r\n"
			}
		}
		return ":1\r\n"
	}
	handlers["DEL"] = deleteHandler
}

// ttl handler should check if the key is present in data store and should return whats the remaining seconds for its expiry
func addTTLHandler() {
	ttlHandler := func(args []string) string {
		if len(args) != 1 {
			return "-ERR wrong number of arguments for 'ttl' command\r\n"
		}
		key := args[0]
		if _, ok := dataStore[key]; !ok {
			return ":-2\r\n"
		} else {
			if expireVal, ok := keyExpirations[key]; !ok {
				return ":-1\r\n"
			} else {
				remainingSeconds := time.Until(expireVal).Seconds()
				if int(remainingSeconds) <= 0 {
					delete(dataStore, key)
					delete(keyExpirations, key)
					return ":-2\r\n"
				}
				return fmt.Sprintf(":%d\r\n", int(remainingSeconds))
			}
		}
	}
	handlers["TTL"] = ttlHandler
}

// incr handler should increase the value of key if the value is in number format, if there is no key, a key should be created
func addINCRHandler() {
	incrHandler := func(args []string) string {
		if len(args) != 1 {
			return "-ERR wrong number of arguments for 'incr' command\r\n"
		}
		key := args[0]
		if dsVal, ok := dataStore[key]; !ok {
			handlers["SET"]([]string{key, "1"})
			return ":1\r\n"
		} else {
			if dsVal.valueType != "string" {
				return "-ERR WRONGTYPE Operation against a key holding the wrong kind of value\r\n"
			}
			if _, ok := keyExpirations[key]; !ok {
				if intVal, err := strconv.Atoi(dsVal.value.(string)); err != nil {
					return "-ERR value is not an integer or out of range\r\n"
				} else {
					newVal := intVal + 1
					dataStore[key] = valueFormat{
						valueType: "string",
						value:     fmt.Sprintf("%d", newVal),
					}
					return fmt.Sprintf(":%d\r\n", newVal)
				}
			} else {
				isKeyExpired := checkKeyExpiration(key)
				if isKeyExpired {
					delete(dataStore, key)
					delete(keyExpirations, key)
					handlers["SET"]([]string{key, "1"})
					return ":1\r\n"
				} else {
					if intVal, err := strconv.Atoi(dsVal.value.(string)); err != nil {
						return "-ERR value is not an integer or out of range\r\n"
					} else {
						newVal := intVal + 1
						dataStore[key] = valueFormat{
							valueType: "string",
							value:     fmt.Sprintf("%d", newVal),
						}
						return fmt.Sprintf(":%d\r\n", newVal)
					}
				}
			}
		}
	}
	handlers["INCR"] = incrHandler
}

// decr handler should decrease the value of key if the value is in number format, if there is no key, a key should be created
func addDECRHandler() {
	decrHandler := func(args []string) string {
		if len(args) != 1 {
			return "-ERR wrong number of arguments for 'decr' command\r\n"
		}
		key := args[0]
		if dsVal, ok := dataStore[key]; !ok {
			handlers["SET"]([]string{key, "-1"})
			return ":-1\r\n"
		} else {
			if dsVal.valueType != "string" {
				return "-ERR WRONGTYPE Operation against a key holding the wrong kind of value\r\n"
			}
			if _, ok := keyExpirations[key]; !ok {
				if intVal, err := strconv.Atoi(dsVal.value.(string)); err != nil {
					return "-ERR value is not an integer or out of range\r\n"
				} else {
					newVal := intVal - 1
					dataStore[key] = valueFormat{
						valueType: "string",
						value:     fmt.Sprintf("%d", newVal),
					}
					return fmt.Sprintf(":%d\r\n", newVal)
				}
			} else {
				isKeyExpired := checkKeyExpiration(key)
				if isKeyExpired {
					delete(dataStore, key)
					delete(keyExpirations, key)
					handlers["SET"]([]string{key, "-1"})
					return ":-1\r\n"
				} else {
					if intVal, err := strconv.Atoi(dsVal.value.(string)); err != nil {
						return "-ERR value is not an integer or out of range\r\n"
					} else {
						newVal := intVal - 1
						dataStore[key] = valueFormat{
							valueType: "string",
							value:     fmt.Sprintf("%d", newVal),
						}
						return fmt.Sprintf(":%d\r\n", newVal)
					}
				}
			}
		}
	}
	handlers["DECR"] = decrHandler
}

func addHandlers() {
	addCOMMANDHandler()
	addSETHandler()
	addGETHandler()
	addEXPIREHandler()
	addDELHandler()
	addTTLHandler()
	addINCRHandler()
	addDECRHandler()
}
