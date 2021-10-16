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
	val, err := cv.evalPolicy(commands, cmdExec, evalExpr, policy, pep.EvalParamNum, []string{pep.PolicyQueryParam}, pep.ReturnKeys)
	if err != nil {
		return CmdEvalResult{Match: false, Error: err}
	}
	return CmdEvalResult{Match: val.EvalExpResult == 0, Error: err, PolicyResult: val.PolicyResult}
}

func (cv commandEvaluate) evalPolicy(commands []string, cmdExec cmd, evalExpr string, policy string, compareComm int, propertyEval []string, ReturnFields []string) (*FinalResult, error) {
	resMap, cmdTotalRes := cv.ExecCommands(commands, cmdExec, evalExpr)
	policyEvalResults := make([]utils.PolicyResult, 0)
	var policyRes int
	if val, ok := resMap[compareComm]; ok {
		var policyResult utils.PolicyResult
		for _, cmdRes := range val {
			res, err := validator.NewPolicyEval().EvaluatePolicy(propertyEval, policy, cmdRes)
			if err != nil {
				return nil, err
			}
			if len(res) > 0 {
				policyResult = utils.MatchPolicy(res[0].ExpressionValue[0].Value, ReturnFields)
			} else {
				policyResult = utils.PolicyResult{ReturnValues: map[string]string{propertyEval[0]: "false"}}
			}
			policyEvalResults = append(policyEvalResults, policyResult)
		}
		for _, per := range policyEvalResults {
			if returnVal, ok := per.ReturnValues[propertyEval[0]]; ok {
				val, err := strconv.ParseBool(returnVal)
				if err != nil {
					continue
				}
				if !val {
					policyRes = 1
					break
				}
			}
		}
	}
	match := policyRes == 0
	policyExpr := utils.GetPolicyExpr(evalExpr)
	if len(policyExpr) == len(evalExpr) {
		return &FinalResult{EvalExpResult: policyRes, PolicyResult: policyEvalResults}, nil
	}
	neweEvalExpr := strings.Replace(evalExpr, policyExpr, fmt.Sprintf("'true' == '%s'", strconv.FormatBool(match)), -1)
	evalExpResult, err := cmdExec.evalExpression(cmdTotalRes, len(cmdTotalRes), make([]string, 0), 0, neweEvalExpr)
	return &FinalResult{EvalExpResult: evalExpResult, PolicyResult: policyEvalResults}, err

}

//ExecCommands execute shell commands and encapsulate it results
func (cv commandEvaluate) ExecCommands(commands []string, cmdExec cmd, evalExpr string) (map[int][]string, []string) {
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
	return resMap, cmdTotalRes
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
	Match        bool
	CmdEvalExpr  string
	PolicyResult []utils.PolicyResult
	Error        error
}

//FinalResult  eval result object
type FinalResult struct {
	EvalExpResult int
	PolicyResult  []utils.PolicyResult
}
