package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/itchyny/gojq"
)

//RunJqQuery accept input and jq arg and return data
func RunJqQuery(jqParam string, input []byte) (string, error) {
	if len(jqParam) == 0 {
		return string(input), nil
	}
	var buffer bytes.Buffer
	query, err := gojq.Parse(jqParam)
	if err != nil {
		return "", err
	}
	var obj interface{}
	err = json.Unmarshal(input, &obj)
	if err != nil {
		return "", err
	}
	iter := query.Run(obj)
	var counter int
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return "", err
		}
		if counter == 0 {
			buffer.WriteString(v.(string))
		} else {
			buffer.WriteString(fmt.Sprintf("\n%s", v.(string)))
		}
		counter++
	}
	return buffer.String(), nil
}
