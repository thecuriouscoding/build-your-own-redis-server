package main

type valueFormat struct {
	ValueType string
	Value     interface{}
}

var (
	handlers  = make(map[string]func(args []string) string)
	dataStore = make(map[string]valueFormat)
)

var (
	VALUE_TYPE_STRING = "string"
)
