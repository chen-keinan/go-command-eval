package eval

import (
	"github.com/chen-keinan/go-command-eval/utils"
	"go.uber.org/zap"
)

//CmdEvaluator interface expose one method to evaluate command with evalExpr
type CmdEvaluator interface {
	EvalCommand(commands []string, evalExpr string) CmdEvalResult
}

type commandEvaluate struct {
}

//New instansiate new command evaluator
func New() CmdEvaluator {
	return &commandEvaluate{}
}

//EvalCommand eval command with eval expr
// accept command and evalExpr
// return eval command result
func (cv commandEvaluate) EvalCommand(commands []string, evalExpr string) CmdEvalResult {
	commandParams := CommandParams(commands)
	zlog, err := zap.NewProduction()
	if err != nil {
		return CmdEvalResult{}
	}
	cmdExec := cmd{commandParams: commandParams, commandExec: commands, evalExpr: evalExpr, command: NewShellExec(), cmdExprBuilder: utils.UpdateCmdExprParam, log: zlog}
	val, err := cv.evalCommand(commands, cmdExec)
	if err != nil {
		return CmdEvalResult{Match: false, Error: err}
	}
	return CmdEvalResult{Match: val == 0, Error: err}
}

func (cv commandEvaluate) evalCommand(commands []string, cmdExec cmd) (int, error) {
	cmdTotalRes := make([]string, 0)
	for index := range commands {
		res := cmdExec.execCommand(index, cmdTotalRes, make([]IndexValue, 0))
		cmdTotalRes = append(cmdTotalRes, res)
	}
	// evaluate command result with expression
	return cmdExec.evalExpression(cmdTotalRes, len(cmdTotalRes), make([]string, 0), 0)
}

//CmdEvalResult command result object
type CmdEvalResult struct {
	Match       bool
	CmdEvalExpr string
	Error       error
}
