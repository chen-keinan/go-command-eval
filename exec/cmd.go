package exec

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/chen-keinan/go-command-eval/common"
	"github.com/chen-keinan/go-command-eval/utils"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type Cmd struct {
	command        Executor
	log            *zap.Logger
	commandParams  map[int][]string
	commandEval    []string
	evalExpr       string
	cmdExprBuilder utils.CmdExprBuilder
}

func New() *Cmd {
	return &Cmd{}
}

func (c *Cmd) addDummyCommandResponse(expr string, index int, n string) string {
	if n == "[^\"]\\S*'\n" || n == "" || n == common.EmptyValue {
		spExpr := utils.SeparateExpr(expr)
		for _, expr := range spExpr {
			if expr.Type == common.SingleValue {
				if !strings.Contains(expr.Expr, fmt.Sprintf("'$%d'", index)) {
					if strings.Contains(expr.Expr, fmt.Sprintf("$%d", index)) {
						return common.NotValidNumber
					}
				}
			}
		}
		return common.EmptyValue
	}
	return n
}

//IndexValue hold command index and result
type IndexValue struct {
	index int
	value string
}

func (c *Cmd) execCommand(index int, prevResult []string, newRes []IndexValue) string {
	cmd := c.commandEval[index]
	paramArr, ok := c.commandParams[index]
	if ok {
		for _, param := range paramArr {
			paramNum, err := strconv.Atoi(param)
			if err != nil {
				c.log.Info(fmt.Sprintf("failed to convert param for command %s", cmd))
				continue
			}
			if paramNum < len(prevResult) {
				n := c.addDummyCommandResponse(c.evalExpr, index, prevResult[paramNum])
				newRes = append(newRes, IndexValue{index: paramNum, value: n})
			}
		}
		commandRes := c.execCmdWithParams(newRes, len(newRes), make([]IndexValue, 0), cmd, make([]string, 0))
		sb := strings.Builder{}
		for _, cr := range commandRes {
			sb.WriteString(utils.AddNewLineToNonEmptyStr(cr))
		}
		return sb.String()
	}
	result, _ := c.command.Exec(cmd)
	if result.Stderr != "" {
		c.log.Info(fmt.Sprintf("Failed to execute command %s\n %s", result.Stderr, cmd))
	}
	return c.addDummyCommandResponse(c.evalExpr, index, result.Stdout)
}

func (c *Cmd) execCmdWithParams(arr []IndexValue, index int, prevResHolder []IndexValue, currCommand string, resArr []string) []string {
	if len(arr) == 0 {
		return c.execShellCmd(prevResHolder, resArr, currCommand, c.command)
	}
	sArr := strings.Split(utils.RemoveNewLineSuffix(arr[0].value), "\n")
	for _, a := range sArr {
		prevResHolder = append(prevResHolder, IndexValue{index: arr[0].index, value: a})
		resArr = c.execCmdWithParams(arr[1:index], index-1, prevResHolder, currCommand, resArr)
		prevResHolder = prevResHolder[:len(prevResHolder)-1]
	}
	return resArr
}

func (c *Cmd) execShellCmd(prevResHolder []IndexValue, resArr []string, currCommand string, se Executor) []string {
	for _, param := range prevResHolder {
		if param.value == common.EmptyValue || param.value == common.NotValidNumber || param.value == "" {
			resArr = append(resArr, param.value)
			break
		}
		cmd := strings.ReplaceAll(currCommand, fmt.Sprintf("#%d", param.index), param.value)
		result, _ := se.Exec(cmd)
		if result.Stderr != "" {
			c.log.Info(fmt.Sprintf("Failed to execute command %s", result.Stderr))
		}
		if len(strings.TrimSpace(result.Stdout)) == 0 {
			result.Stdout = common.EmptyValue
		}
		resArr = append(resArr, result.Stdout)
	}
	return resArr
}

//evalExpression expression eval as cartesian product
func (c *Cmd) evalExpression(cmd string,
	commandRes []string, commResSize int, permutationArr []string, testFailure int) int {
	if len(commandRes) == 0 {
		return c.evalCommand(cmd, permutationArr, testFailure)
	}
	outputs := strings.Split(utils.RemoveNewLineSuffix(commandRes[0]), "\n")
	for _, o := range outputs {
		permutationArr = append(permutationArr, o)
		testFailure = c.evalExpression(cmd, commandRes[1:commResSize], commResSize-1, permutationArr, testFailure)
		permutationArr = permutationArr[:len(permutationArr)-1]
	}
	return testFailure
}

func (c *Cmd) evalCommand(cmd string, permutationArr []string, testExec int) int {
	// build command expression with params
	expr := c.cmdExprBuilder(permutationArr, c.evalExpr)
	testExec++
	// eval command expression
	testSucceeded, err := evalCommandExpr(strings.ReplaceAll(expr, common.EmptyValue, ""))
	if err != nil {
		c.log.Info(fmt.Sprintf("failed to evaluate command expr %s for cmd  %s : err %s", expr, cmd, err.Error()))
	}
	return testExec - testSucceeded
}

func evalCommandExpr(expr string) (int, error) {
	expression, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return 0, err
	}
	result, err := expression.Evaluate(nil)
	if err != nil {
		return 0, err
	}
	b, ok := result.(bool)
	if ok && b {
		return 1, nil
	}
	return 0, nil
}
