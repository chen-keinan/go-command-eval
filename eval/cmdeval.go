package eval

import (
	"github.com/chen-keinan/go-command-eval/utils"
	"github.com/chen-keinan/go-opa-validate/validator"
	"go.uber.org/zap"
	"strings"
)

//CmdEvaluator interface expose one method to evaluate command with evalExpr
type CmdEvaluator interface {
	EvalCommand(commands []string, evalExpr string) CmdEvalResult
	EvalCommandPolicy(commands []string, policy string, propertyEval string, commNum int) CmdEvalResult
}

type commandEvaluate struct {
}

//NewEvalCmd instantiate new command evaluator
func NewEvalCmd() CmdEvaluator {
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

//EvalCommandPolicy eval command with opa policy
// accept command and policy and property to eval
// return eval command result
func (cv commandEvaluate) EvalCommandPolicy(commands []string, policy string, propertyEval string, commNum int) CmdEvalResult {
	commandParams := CommandParams(commands)
	zlog, err := zap.NewProduction()
	if err != nil {
		return CmdEvalResult{}
	}
	cmdExec := cmd{commandParams: commandParams, commandExec: commands, command: NewShellExec(), cmdExprBuilder: utils.UpdateCmdExprParam, log: zlog}
	val, err := cv.evalPolicy(commands, cmdExec, policy, propertyEval, commNum)
	if err != nil {
		return CmdEvalResult{Match: false, Error: err}
	}
	return CmdEvalResult{Match: val == 0, Error: err}
}

func (cv commandEvaluate) evalPolicy(commands []string, cmdExec cmd, policy string, propertyEval string, compareComm int) (int, error) {
	resMap := make(map[int][]string)
	cmdTotalRes := make([]string, 0)
	var commNum = 0
	for index := range commands {
		res := cmdExec.execCommand(index, cmdTotalRes, make([]IndexValue, 0))
		sb := strings.Builder{}
		for _, s := range res {
			sb.WriteString(s)
		}
		resMap[commNum] = res
		cmdTotalRes = append(cmdTotalRes, sb.String())
		commNum++
	}
	policyEvalResults := make([]*validator.ValidateResult, 0)
	if val, ok := resMap[compareComm]; ok {
		for _, cmdRes := range val {
			res, err := validator.NewPolicyEval().EvaluatePolicy([]string{propertyEval}, policy, cmdRes)
			policyEvalResults = append(policyEvalResults, res...)
			if err != nil {
				return 0, err
			}
		}
		for _, per := range policyEvalResults {
			if !per.Value {
				return 1, nil
			}
		}
	}
	return 0, nil
}

func (cv commandEvaluate) evalCommand(commands []string, cmdExec cmd) (int, error) {
	cmdTotalRes := make([]string, 0)
	for index := range commands {
		res := cmdExec.execCommand(index, cmdTotalRes, make([]IndexValue, 0))
		sb := strings.Builder{}
		for _, s := range res {
			sb.WriteString(s)
		}
		cmdTotalRes = append(cmdTotalRes, sb.String())
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
