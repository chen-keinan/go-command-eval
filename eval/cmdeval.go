package eval

type CmdEvaluator interface {
	EvalCommand(commands []string, evalExpr string) CmdEvalResult
}

type commandEvaluate struct {
}

func NewCmdEval() CmdEvaluator {
	return &commandEvaluate{}
}

func New() CmdEvaluator {
	return &commandEvaluate{}
}

func (cv commandEvaluate) EvalCommand(commands []string, evalExpr string) CmdEvalResult {
	commandParams := CommandParams(commands)
	cmdExec := cmd{commandParams: commandParams, commandExec: commands, evalExpr: evalExpr ,command: NewShellExec()}
	cv.evalCommand(commands, cmdExec)
	return CmdEvalResult{}
}

func (cv commandEvaluate) evalCommand(commands []string, cmdExec cmd) {
	cmdTotalRes := make([]string, 0)
	for index := range commands {
		res := cmdExec.execCommand(index, cmdTotalRes, make([]IndexValue, 0))
		cmdTotalRes = append(cmdTotalRes, res)
	}

	// evaluate command result with expression
	NumFailedTest := cmd.evalExpression(cmdTotalRes, len(cmdTotalRes), make([]string, 0), 0)
}

type CmdEvalResult struct {
	Match       bool
	CmdEvalExpr string
	Error       error
}
