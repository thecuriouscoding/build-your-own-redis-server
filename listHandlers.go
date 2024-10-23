package main

import (
	"fmt"
)

// Node represents an individual node in the doubly linked list
type Node struct {
	data string
	prev *Node
	next *Node
}

type ListValue struct {
	start  *Node
	tail   *Node
	length int
}

var (
	PUSH_ON_LEFT   = "left"
	PUSH_ON_RIGHT  = "right"
	POP_FROM_LEFT  = "left"
	POP_FROM_RIGHT = "right"
)

// it should be adding the elements onto the list on left or right depending upon the push operation being send
func addElementsToList(elements []string, start *Node, tail *Node, pushOn string) (*Node, *Node) {
	if start == nil {
		start = &Node{}
		tail = start
	}
	for _, ele := range elements {
		switch pushOn {
		case PUSH_ON_LEFT:
			if start.data == "" {
				start.data = ele
				tail.data = ele
			} else {
				newListNode := Node{
					data: ele,
				}
				newListNode.next = start
				start.prev = &newListNode
				start = &newListNode
			}
		case PUSH_ON_RIGHT:
			if start.data == "" {
				start.data = ele
				tail.data = ele
			} else {
				newListNode := Node{
					data: ele,
				}
				newListNode.prev = tail
				tail.next = &newListNode
				tail = &newListNode
			}
		}
	}
	return start, tail
}

// it should be removing the elements from the list from left or right depending upon the pop operation being send
func removeElementsFromList(noOfElementsToRemove int, start *Node, tail *Node, popFrom string) (*Node, *Node, []string) {
	elementsRemoved := []string{}
	switch popFrom {
	case POP_FROM_LEFT:
		for start != nil {
			if len(elementsRemoved) == noOfElementsToRemove {
				break
			}
			elementsRemoved = append(elementsRemoved, start.data)
			start = start.next
		}
		if start == nil {
			tail = nil
		}
	case POP_FROM_RIGHT:
		for tail != nil {
			if len(elementsRemoved) == noOfElementsToRemove {
				break
			}
			elementsRemoved = append(elementsRemoved, tail.data)
			tail = tail.prev
		}
		if tail == nil {
			start = nil
		}
	}
	return start, tail, elementsRemoved
}

// lpush handler should create a list if no such list is present in the pretext that no other format key is already present. It should push the elements to the left of list
func addLPUSHHandler() {
	lpushHandler := func(args []string) string {
		toContinue, err := handlerValidations(argumentLengthCheck{
			toCheck:   true,
			argLength: 2,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, argumentLengthCheck{
			toCheck:   false,
			argLength: 0,
			errorCode: ERR_NO_ERROR,
		}, args, COMMAND_LPUSH)
		if !toContinue {
			return err
		}
		key := args[0]
		if isExpired := checkKeyExpiration(key); isExpired {
			deleteExpiredKey(key)
		}
		if dsVal, ok := dataStore[key]; !ok {
			elements := args[1:]
			start := &Node{}
			tail := start
			start, tail = addElementsToList(elements, start, tail, PUSH_ON_LEFT)
			dataStore[key] = valueFormat{
				valueType: VALUE_TYPE_LIST,
				value: ListValue{
					start:  start,
					tail:   tail,
					length: len(elements),
				},
			}
			return fmt.Sprintf(":%d\r\n", len(elements))
		} else {
			if dsVal.valueType != VALUE_TYPE_LIST {
				return getErrorMessage(ERR_WRONG_TYPE_OPERATION, COMMAND_LPUSH)
			}
			elements := args[1:]
			totalElements := len(elements)
			start := dsVal.value.(ListValue).start
			tail := dsVal.value.(ListValue).tail
			start, tail = addElementsToList(elements, start, tail, PUSH_ON_LEFT)
			dataStore[key] = valueFormat{
				valueType: VALUE_TYPE_LIST,
				value: ListValue{
					start:  start,
					tail:   tail,
					length: dsVal.value.(ListValue).length + totalElements,
				},
			}
			return fmt.Sprintf(":%d\r\n", totalElements)
		}
	}
	handlers["LPUSH"] = lpushHandler
}

// lrange handler should return with the elements of list if key is valid and is of list type
func addLRANGEHandler() {
	lrangeHandler := func(args []string) string {
		toContinue, err := handlerValidations(argumentLengthCheck{
			toCheck:   true,
			argLength: 3,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, argumentLengthCheck{
			toCheck:   true,
			argLength: 3,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, args, COMMAND_LRANGE)
		if !toContinue {
			return err
		}
		key := args[0]
		startIndexStr := args[1]
		endIndexStr := args[2]
		if isExpired := checkKeyExpiration(key); isExpired {
			deleteExpiredKey(key)
		}
		startIndex, convertErr := readIntValue(startIndexStr)
		if convertErr != nil {
			return getErrorMessage(ERR_VALUE_NOT_INTEGER, COMMAND_LRANGE)
		}
		endIndex, convertErr := readIntValue(endIndexStr)
		if convertErr != nil {
			return getErrorMessage(ERR_VALUE_NOT_INTEGER, COMMAND_LRANGE)
		}

		if dsVal, ok := dataStore[key]; !ok {
			return "*0\r\n"
		} else {
			if dsVal.valueType != VALUE_TYPE_LIST {
				return getErrorMessage(ERR_WRONG_TYPE_OPERATION, COMMAND_LRANGE)
			}
			if startIndex >= 0 && endIndex >= 0 {
				if endIndex < startIndex {
					return "*0\r\n"
				}
				if startIndex <= dsVal.value.(ListValue).length && endIndex >= startIndex {
					tempStart := 0
					responseElements := []string{}
					listStart := dsVal.value.(ListValue).start
					for listStart != nil {
						if tempStart >= startIndex && tempStart <= endIndex {
							responseElements = append(responseElements, listStart.data)
						}
						listStart = listStart.next
						tempStart++
					}
					response := fmt.Sprintf("*%d\r\n", len(responseElements))
					for _, ele := range responseElements {
						response = response + fmt.Sprintf("$%d\r\n%s\r\n", len(ele), ele)
					}
					return response
				}
			}
			return "*0\r\n"
		}
	}
	handlers["LRANGE"] = lrangeHandler
}

// rpush handler should create a list if no such list is present in the pretext that no other format key is already present. It should push the elements to the right of list
func addRPUSHHandler() {
	rpushHandler := func(args []string) string {
		toContinue, err := handlerValidations(argumentLengthCheck{
			toCheck:   true,
			argLength: 2,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, argumentLengthCheck{
			toCheck:   false,
			argLength: 0,
			errorCode: ERR_NO_ERROR,
		}, args, COMMAND_RPUSH)
		if !toContinue {
			return err
		}
		key := args[0]
		if isExpired := checkKeyExpiration(key); isExpired {
			deleteExpiredKey(key)
		}
		if dsVal, ok := dataStore[key]; !ok {
			elements := args[1:]
			start := &Node{}
			tail := start
			start, tail = addElementsToList(elements, start, tail, PUSH_ON_RIGHT)
			dataStore[key] = valueFormat{
				valueType: VALUE_TYPE_LIST,
				value: ListValue{
					start:  start,
					tail:   tail,
					length: len(elements),
				},
			}
			return fmt.Sprintf(":%d\r\n", len(elements))
		} else {
			if dsVal.valueType != VALUE_TYPE_LIST {
				return getErrorMessage(ERR_WRONG_TYPE_OPERATION, COMMAND_LPUSH)
			}
			elements := args[1:]
			totalElements := len(elements)
			start := dsVal.value.(ListValue).start
			tail := dsVal.value.(ListValue).tail
			start, tail = addElementsToList(elements, start, tail, PUSH_ON_RIGHT)
			dataStore[key] = valueFormat{
				valueType: VALUE_TYPE_LIST,
				value: ListValue{
					start:  start,
					tail:   tail,
					length: dsVal.value.(ListValue).length + totalElements,
				},
			}
			return fmt.Sprintf(":%d\r\n", totalElements)
		}
	}
	handlers["RPUSH"] = rpushHandler
}

// lpop handler should be removing the passed number of elements to be removed from the list from left
func addLPOPHandler() {
	lpopHandler := func(args []string) string {
		toContinue, err := handlerValidations(argumentLengthCheck{
			toCheck:   true,
			argLength: 1,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, argumentLengthCheck{
			toCheck:   true,
			argLength: 2,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, args, COMMAND_LPOP)
		if !toContinue {
			return err
		}
		key := args[0]
		if isExpired := checkKeyExpiration(key); isExpired {
			deleteExpiredKey(key)
		}
		if dsVal, ok := dataStore[key]; !ok {
			return "$-1\r\n"
		} else {
			if dsVal.valueType != VALUE_TYPE_LIST {
				return getErrorMessage(ERR_WRONG_TYPE_OPERATION, COMMAND_LPOP)
			}
			count := 1
			if len(args) > 1 {
				_count, err := readIntValue(args[1])
				if err != nil {
					return getErrorMessage(ERR_VALUE_OUT_OF_RANGE_MUST_BE_POSITIVE, COMMAND_LPOP)
				}
				if _count == 0 {
					return "*0\r\n"
				}
				count = _count
			}
			start := dsVal.value.(ListValue).start
			tail := dsVal.value.(ListValue).tail
			var elementsRemoved []string
			start, tail, elementsRemoved = removeElementsFromList(count, start, tail, POP_FROM_LEFT)
			dataStore[key] = valueFormat{
				valueType: VALUE_TYPE_LIST,
				value: ListValue{
					start:  start,
					tail:   tail,
					length: dsVal.value.(ListValue).length - len(elementsRemoved),
				},
			}
			response := fmt.Sprintf("*%d\r\n", len(elementsRemoved))
			for _, ele := range elementsRemoved {
				response = response + fmt.Sprintf("$%d\r\n%s\r\n", len(ele), ele)
			}
			return response
		}
	}
	handlers["LPOP"] = lpopHandler
}

// rpop handler should be removing the passed number of elements to be removed from the list from right
func addRPOPHandler() {
	rpopHandler := func(args []string) string {
		toContinue, err := handlerValidations(argumentLengthCheck{
			toCheck:   true,
			argLength: 1,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, argumentLengthCheck{
			toCheck:   true,
			argLength: 2,
			errorCode: ERR_WRONG_NUMBER_OF_ARGUMENTS,
		}, args, COMMAND_RPOP)
		if !toContinue {
			return err
		}
		key := args[0]
		if isExpired := checkKeyExpiration(key); isExpired {
			deleteExpiredKey(key)
		}
		if dsVal, ok := dataStore[key]; !ok {
			return "$-1\r\n"
		} else {
			if dsVal.valueType != VALUE_TYPE_LIST {
				return getErrorMessage(ERR_WRONG_TYPE_OPERATION, COMMAND_RPOP)
			}
			count := 1
			if len(args) > 1 {
				_count, err := readIntValue(args[1])
				if err != nil {
					return getErrorMessage(ERR_VALUE_OUT_OF_RANGE_MUST_BE_POSITIVE, COMMAND_RPOP)
				}
				if _count == 0 {
					return "*0\r\n"
				}
				count = _count
			}
			start := dsVal.value.(ListValue).start
			tail := dsVal.value.(ListValue).tail
			var elementsRemoved []string
			start, tail, elementsRemoved = removeElementsFromList(count, start, tail, POP_FROM_RIGHT)
			dataStore[key] = valueFormat{
				valueType: VALUE_TYPE_LIST,
				value: ListValue{
					start:  start,
					tail:   tail,
					length: dsVal.value.(ListValue).length - len(elementsRemoved),
				},
			}
			response := fmt.Sprintf("*%d\r\n", len(elementsRemoved))
			for _, ele := range elementsRemoved {
				response = response + fmt.Sprintf("$%d\r\n%s\r\n", len(ele), ele)
			}
			return response
		}
	}
	handlers["RPOP"] = rpopHandler
}

func addListHandlers() {
	addLPUSHHandler()
	addLRANGEHandler()
	addRPUSHHandler()
	addLPOPHandler()
	addRPOPHandler()
}
