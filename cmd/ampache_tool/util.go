package main

import (
	"errors"
	"fmt"
)

func parseArgs(args []string) (host, action string, input map[string]string, err error) {
	input = make(map[string]string)

	var name string
	var needValue bool = false
	for _, s := range args {
		// Capture value part of "--flag value"
		if needValue {
			needValue = false
			name = ""
			input[name] = s
			continue
		} else if s[0] == '-' {
			numMinuses := 1
			if s[1] == '-' {
				numMinuses++
			}

			name = s[numMinuses:]
			if len(name) == 0 || name[0] == '-' || name[0] == '=' {
				err = fmt.Errorf("bad flag syntax: %s", s)
				return
			}

			// Capture value part of "--flag=value"
			needValue = true
			value := ""
			for i := 1; i < len(name); i++ {
				if name[i] == '=' {
					value = name[i+1:]
					needValue = false
					name = name[0:i]
					break
				}
			}
			if !needValue {
				input[name] = value
			}
		} else {
			if host == "" {
				host = s
				continue
			}
			if action != "" {
				err = fmt.Errorf("too many commands; first is '%s'", action)
				return
			}

			action = s
		}
	}
	if needValue {
		err = fmt.Errorf("missing value for '%s'", args[len(args)-1])
		return
	}
	if host == "" {
		err = errors.New("missing host")
		return
	}

	if action == "" {
		err = errors.New("missing action")
		return
	}

	return
}
