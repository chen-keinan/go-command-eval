package eval

type CmdEvaluator interface {
}

type commandEvaluate struct {
}

func NewCmdEval() CmdEvaluator {
	return &commandEvaluate{}
}

func New() CmdEvaluator {
	return &commandEvaluate{}
}

func (cv commandEvaluate) EvalCommand(commands []string, evalExpr string) bool {
	commandParams := CommandParams(commands)
	cmdExec := cmd{commandParams: commandParams, commandExec: commands, evalExpr: evalExpr}
	cv.evalCommand(commands, cmdExec)
}

func (cv commandEvaluate) evalCommand(commands []string, cmdExec cmd) {
	cmdTotalRes := make([]string, 0)
	for index := range commands {
		res := cmdExec.execCommand(index, cmdTotalRes, make([]IndexValue, 0))
		cmdTotalRes = append(cmdTotalRes, res)
	}
}
