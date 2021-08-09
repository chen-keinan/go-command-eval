package eval

import "strings"

//specEval hold command and evaluation expression
type specEval struct {
	cmd      []string
	evalExpr string
}

//CommandParams calculate command params map params inorder to inject prev. command result into next command
// accept list of commands return location and result
func CommandParams(commands []string) map[int][]string {
	commandParams := make(map[int][]string)
	for index, command := range commands {
		findIndex(command, "#", index, commandParams)
	}
	return nil
}

// find all params in command to be replace with output
func findIndex(s, c string, commandIndex int, locations map[int][]string) {
	b := strings.Index(s, c)
	if b == -1 {
		return
	}
	if locations[commandIndex] == nil {
		locations[commandIndex] = make([]string, 0)
	}
	locations[commandIndex] = append(locations[commandIndex], s[b+1:b+2])
	findIndex(s[b+2:], c, commandIndex, locations)
}
