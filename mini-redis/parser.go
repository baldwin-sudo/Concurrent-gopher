package main

import (
	"fmt"
	"strings"
)

type Request struct {
	Command CommandType
	Args    []string
}
type RequestError struct {
	Message string
}

func (reqErr *RequestError) Error() string {
	return fmt.Sprint("Error  :", reqErr.Message)
}

type CommandType int

const (
	CMD_UNKNOWN CommandType = iota
	CMD_GET
	CMD_SET
	CMD_DEL
	CMD_HELP
	CMD_QUIT
)

func ParseCommand(input string) (*Request, error) {
	input = strings.TrimSpace(input)
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return nil, &RequestError{Message: "Invalid command !"}
	}
	command := strings.ToUpper(fields[0])
	switch command {
	case "GET":
		if len(fields) != 2 {
			return nil, &RequestError{Message: " invalid format please use GET <key> "}
		}

		return &Request{Command: CMD_GET, Args: []string{fields[1]}}, nil
	case "SET":
		if len(fields) != 3 {
			return nil, &RequestError{Message: " invalid format please use SET <key> <value> "}
		}
		return &Request{Command: CMD_SET, Args: fields[1:]}, nil

	case "DEL":
		if len(fields) != 2 {
			return nil, &RequestError{Message: " invalid format please use DEL <key>  "}
		}
		return &Request{Command: CMD_DEL, Args: fields[1:]}, nil

	case "HELP":
		if len(fields) != 1 {
			return nil, &RequestError{Message: " invalid format please use HELP  "}
		}
		return &Request{Command: CMD_HELP, Args: nil}, nil

	case "QUIT":
		if len(fields) != 1 {
			return nil, &RequestError{Message: " invalid format please use QUIT  "}
		}
		return &Request{Command: CMD_QUIT, Args: nil}, nil

	default:
		return &Request{Command: CMD_UNKNOWN, Args: nil}, nil
	}

}
