package main

var (
	handlers = make(map[string]func(args []string) string)
)
