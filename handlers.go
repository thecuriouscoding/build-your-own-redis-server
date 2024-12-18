package main

import (
	"fmt"
	"strconv"
	"time"
)

// command handler gets called on initial COMMAND command sent from redis-cli like client
func addCOMMANDHandler() {
	commandHandler := func(args []string) string {
		return "+OK\r\n"
	}
	handlers["COMMAND"] = commandHandler
}

// set handler should be able to set a key and value in data store
func addSETHandler() {
	setHandler := func(args []string) string {
		toContinue, err := handlerValidations(argumentLengthCheck{
			toCheck:   true,
			argLength: 2,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, argumentLengthCheck{
			toCheck:   true,
			argLength: 2,
			errorCode: ERR_SYNTAX_ERROR,
		}, args, COMMAND_SET)
		if !toContinue {
			return err
		}
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

// get handler should be returning the value of key if key is of string type
func addGETHandler() {
	getHandler := func(args []string) string {
		toContinue, err := handlerValidations(argumentLengthCheck{
			toCheck:   true,
			argLength: 1,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, argumentLengthCheck{
			toCheck:   false,
			argLength: 0,
			errorCode: ERR_NO_ERROR,
		}, args, COMMAND_GET)
		if !toContinue {
			return err
		}
		key := args[0]
		if val, ok := dataStore[key]; !ok {
			return "$-1\r\n"
		} else {
			if val.ValueType != VALUE_TYPE_STRING {
				return "$-1\r\n"
			}
			isExpired := checkKeyExpiration(key)
			if isExpired {
				deleteExpiredKey(key)
				return "$-1\r\n"
			}
			return fmt.Sprintf("$%d\r\n%s\r\n", len(val.Value.(string)), val.Value)
		}
	}
	handlers["GET"] = getHandler
}

// expire handler should add the passed number of seconds as expire time against the key in data store
func addEXPIREHandler() {
	expireHandler := func(args []string) string {
		toContinue, err := handlerValidations(argumentLengthCheck{
			toCheck:   true,
			argLength: 1,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, argumentLengthCheck{
			toCheck:   true,
			argLength: 2,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, args, COMMAND_EXPIRE)
		if !toContinue {
			return err
		}
		key := args[0]
		value := args[1]
		if _, ok := dataStore[key]; !ok {
			return ":0\r\n"
		} else {
			expirationTime, err := strconv.Atoi(value)
			if err != nil {
				return getErrorMessage(ERR_VALUE_NOT_INTEGER, COMMAND_EXPIRE)
			}
			keyExpirations[key] = time.Now().Add(time.Duration(expirationTime) * time.Second)
		}
		return ":1\r\n"
	}
	handlers["EXPIRE"] = expireHandler
}

// delete handler should be deleting key from data store
func addDELHandler() {
	deleteHandler := func(args []string) string {
		toContinue, err := handlerValidations(argumentLengthCheck{
			toCheck:   true,
			argLength: 1,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, argumentLengthCheck{
			toCheck:   false,
			argLength: 0,
			errorCode: ERR_NO_ERROR,
		}, args, COMMAND_DEL)
		if !toContinue {
			return err
		}
		for _, key := range args {
			if _, ok := dataStore[key]; !ok {
				return ":0\r\n"
			} else {
				isExpired := checkKeyExpiration(key)
				if isExpired {
					deleteExpiredKey(key)
					return ":0\r\n"
				}
				deleteExpiredKey(key)
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
		toContinue, err := handlerValidations(argumentLengthCheck{
			toCheck:   true,
			argLength: 1,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, argumentLengthCheck{
			toCheck:   true,
			argLength: 1,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, args, COMMAND_TTL)
		if !toContinue {
			return err
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
					deleteExpiredKey(key)
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
		toContinue, err := handlerValidations(argumentLengthCheck{
			toCheck:   true,
			argLength: 1,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, argumentLengthCheck{
			toCheck:   true,
			argLength: 1,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, args, COMMAND_INCR)
		if !toContinue {
			return err
		}
		key := args[0]
		if dsVal, ok := dataStore[key]; !ok {
			handlers["SET"]([]string{key, "1"})
			return ":1\r\n"
		} else {
			if dsVal.ValueType != VALUE_TYPE_STRING {
				return getErrorMessage(ERR_WRONG_TYPE_OPERATION, COMMAND_INCR)
			}
			if _, ok := keyExpirations[key]; !ok {
				if intVal, err := strconv.Atoi(dsVal.Value.(string)); err != nil {
					return getErrorMessage(ERR_VALUE_NOT_INTEGER, COMMAND_INCR)
				} else {
					newVal := intVal + 1
					dataStore[key] = valueFormat{
						ValueType: VALUE_TYPE_STRING,
						Value:     fmt.Sprintf("%d", newVal),
					}
					return fmt.Sprintf(":%d\r\n", newVal)
				}
			} else {
				isKeyExpired := checkKeyExpiration(key)
				if isKeyExpired {
					deleteExpiredKey(key)
					handlers["SET"]([]string{key, "1"})
					return ":1\r\n"
				} else {
					if intVal, err := strconv.Atoi(dsVal.Value.(string)); err != nil {
						return getErrorMessage(ERR_VALUE_NOT_INTEGER, COMMAND_INCR)
					} else {
						newVal := intVal + 1
						dataStore[key] = valueFormat{
							ValueType: "string",
							Value:     fmt.Sprintf("%d", newVal),
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
		toContinue, err := handlerValidations(argumentLengthCheck{
			toCheck:   true,
			argLength: 1,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, argumentLengthCheck{
			toCheck:   true,
			argLength: 1,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, args, COMMAND_DECR)
		if !toContinue {
			return err
		}
		key := args[0]
		if dsVal, ok := dataStore[key]; !ok {
			handlers["SET"]([]string{key, "-1"})
			return ":-1\r\n"
		} else {
			if dsVal.ValueType != VALUE_TYPE_STRING {
				return getErrorMessage(ERR_WRONG_TYPE_OPERATION, COMMAND_DECR)
			}
			if _, ok := keyExpirations[key]; !ok {
				if intVal, err := strconv.Atoi(dsVal.Value.(string)); err != nil {
					return getErrorMessage(ERR_WRONG_TYPE_OPERATION, COMMAND_DECR)
				} else {
					newVal := intVal - 1
					dataStore[key] = valueFormat{
						ValueType: VALUE_TYPE_STRING,
						Value:     fmt.Sprintf("%d", newVal),
					}
					return fmt.Sprintf(":%d\r\n", newVal)
				}
			} else {
				isKeyExpired := checkKeyExpiration(key)
				if isKeyExpired {
					deleteExpiredKey(key)
					handlers["SET"]([]string{key, "-1"})
					return ":-1\r\n"
				} else {
					if intVal, err := strconv.Atoi(dsVal.Value.(string)); err != nil {
						return getErrorMessage(ERR_WRONG_TYPE_OPERATION, COMMAND_DECR)
					} else {
						newVal := intVal - 1
						dataStore[key] = valueFormat{
							ValueType: VALUE_TYPE_STRING,
							Value:     fmt.Sprintf("%d", newVal),
						}
						return fmt.Sprintf(":%d\r\n", newVal)
					}
				}
			}
		}
	}
	handlers["DECR"] = decrHandler
}

func addBasicHandlers() {
	addCOMMANDHandler()
	addSETHandler()
	addGETHandler()
	addEXPIREHandler()
	addDELHandler()
	addTTLHandler()
	addINCRHandler()
	addDECRHandler()
}
