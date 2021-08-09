package eval

import (
	"fmt"
	"testing"
)

func TestEvalCommand(t *testing.T) {
	cmd := []string{"ls /etc/hosts | awk -F \" \" '{print $1}' |awk 'FNR <= 1'"}
	evalExpr := "'$0' == '/etc/hosts'"
	res := New().EvalCommand(cmd, evalExpr)
	if res.Match {
		fmt.Print("OK")
	}

}
