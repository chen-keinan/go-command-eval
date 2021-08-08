package eval

type CmdEval interface {
}

type commandEvaluate struct {
}

func NewCmdEval() CmdEval {
	return &commandEvaluate{}
}

func (cv commandEvaluate) runAuditTest(commands []string, evalExpr string) {

}
