package eval

import (
	"fmt"
	"github.com/chen-keinan/go-command-eval/utils"
	"github.com/chen-keinan/go-opa-validate/validator"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

//CmdEvaluator interface expose one method to evaluate command with evalExpr
type CmdEvaluator interface {
	EvalCommand(commands []string, evalExpr string) CmdEvalResult
	EvalCommandPolicy(commands []string, evalExpr string, policy string) CmdEvalResult
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
	cmdExec := cmd{commandParams: commandParams, commandExec: commands, command: NewShellExec(), cmdExprBuilder: utils.UpdateCmdExprParam, log: zlog}
	val, err := cv.evalCommand(commands, cmdExec, evalExpr)
	if err != nil {
		return CmdEvalResult{Match: false, Error: err}
	}
	return CmdEvalResult{Match: val == 0, Error: err}
}

//EvalCommandPolicy eval command with opa policy
// accept command and policy and property to eval
// return eval command result
func (cv commandEvaluate) EvalCommandPolicy(commands []string, evalExpr string, policy string) CmdEvalResult {
	commandParams := CommandParams(commands)
	zlog, err := zap.NewProduction()
	if err != nil {
		return CmdEvalResult{}
	}
	cmdExec := cmd{commandParams: commandParams, commandExec: commands, command: NewShellExec(), cmdExprBuilder: utils.UpdateCmdExprParam, log: zlog}
	pep, err := utils.ReadPolicyExpr(evalExpr)
	if err != nil {
		return CmdEvalResult{Match: false}
	}
	val, err := cv.evalPolicy(commands, cmdExec, evalExpr, policy, pep.EvalParamNum, []string{pep.PolicyQueryParam}...)
	if err != nil {
		return CmdEvalResult{Match: false, Error: err}
	}
	return CmdEvalResult{Match: val == 0, Error: err}
}

func (cv commandEvaluate) evalPolicy(commands []string, cmdExec cmd, evalExpr string, policy string, compareComm int, propertyEval ...string) (int, error) {
	resMap := make(map[int][]string)
	cmdTotalRes := make([]string, 0)
	var commNum = 0
	for index := range commands {
		res := cmdExec.execCommand(index, cmdTotalRes, make([]IndexValue, 0), evalExpr)
		sb := strings.Builder{}
		for _, s := range res {
			sb.WriteString(s)
		}
		resMap[commNum] = res
		cmdTotalRes = append(cmdTotalRes, sb.String())
		commNum++
	}
	policyEvalResults := make([]*validator.ValidateResult, 0)
	var policyRes int
	if val, ok := resMap[compareComm]; ok {
		for _, cmdRes := range val {
			res, err := validator.NewPolicyEval().EvaluatePolicy(propertyEval, policy, cmdRes)
			policyEvalResults = append(policyEvalResults, res...)
			if err != nil {
				return 0, err
			}
		}
		for _, per := range policyEvalResults {
			if !per.Value {
				policyRes = 1
				break
			}
		}
	}
	match := policyRes == 0
	PolicyExpr := utils.GetPolicyExpr(evalExpr)
	if len(PolicyExpr) == len(evalExpr) {
		return policyRes, nil
	}
	neweEvalExpr := strings.Replace(evalExpr, PolicyExpr, fmt.Sprintf("'true' == '%s'", strconv.FormatBool(match)), -1)
	return cmdExec.evalExpression(cmdTotalRes, len(cmdTotalRes), make([]string, 0), 0, neweEvalExpr)
}

func (cv commandEvaluate) evalCommand(commands []string, cmdExec cmd, evalExpr string) (int, error) {
	cmdTotalRes := make([]string, 0)
	for index := range commands {
		res := cmdExec.execCommand(index, cmdTotalRes, make([]IndexValue, 0), evalExpr)
		sb := strings.Builder{}
		for _, s := range res {
			sb.WriteString(s)
		}
		cmdTotalRes = append(cmdTotalRes, sb.String())
	}
	// evaluate command result with expression
	return cmdExec.evalExpression(cmdTotalRes, len(cmdTotalRes), make([]string, 0), 0, evalExpr)
}

//CmdEvalResult command result object
type CmdEvalResult struct {
	Match       bool
	CmdEvalExpr string
	Error       error
}
