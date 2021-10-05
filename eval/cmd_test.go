package eval

import (
	"fmt"
	"github.com/chen-keinan/go-command-eval/utils"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
	"strings"
	"testing"
)

func TestCommandParams(t *testing.T) {
	tests := []struct {
		name string
		cmd  []string
		want map[int][]string
	}{
		{name: "two command and one param", cmd: []string{" aaa", "bb ${1}"}, want: map[int][]string{1: {"1"}}},
		{name: "two command and 2 params on 2 commands", cmd: []string{" aaa", "bb ${1}", "cc ${2}"}, want: map[int][]string{1: {"1"}, 2: {"2"}}},
		{name: "two command and 2 params on one command", cmd: []string{" aaa", "bb ${1}", "cc ${1} ${2}"}, want: map[int][]string{1: {"1"}, 2: {"1", "2"}}},
		{name: "two command no params", cmd: []string{" aaa", "bb ", "cc"}, want: map[int][]string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CommandParams(tt.cmd)
			if len(tt.want) != len(got) {
				t.Errorf("CommandParams() = %v, want %v", got, tt.want)
			}
			for key, value := range tt.want {
				if val, ok := got[key]; ok {
					for k, v := range val {
						if v != value[k] {
							t.Errorf("CommandParams() = %v, want %v", got, tt.want)
						}
					}
				} else {
					{
						t.Errorf("CommandParams() = %v, want %v", got, tt.want)
					}
				}
			}
		})
	}
}

func TestEvalExpression(t *testing.T) {

	tests := []struct {
		name        string
		commandRes  []string
		commResSize int
		testFailure int
		evalExpr    string
		want        int
		wantErr     error
	}{
		{name: "one command res and one param good", commandRes: []string{"/etc/hosts"}, commResSize: 1, testFailure: 0, evalExpr: "'${0}' == '/etc/hosts'", want: 0, wantErr: nil},
		{name: "one command res and one param bad", commandRes: []string{"/etc/hosts"}, commResSize: 1, testFailure: 0, evalExpr: "'${0}' == '/etc/hosts1'", want: 1, wantErr: nil},
		{name: "one command res and one param bad", commandRes: []string{"/etc/hosts"}, commResSize: 1, testFailure: 0, evalExpr: "'${0}' == /etc/hosts", want: 0, wantErr: fmt.Errorf("failed to evaluate command expr '/etc/hosts' == /etc/hosts for : err Cannot transition token types from COMPARATOR [==] to MODIFIER [/]")},
		{name: "two command res and one param good", commandRes: []string{"/etc/hosts", "/etc/groups"}, commResSize: 2, testFailure: 0, evalExpr: "'${0}' == /etc/hosts && '${0}' == /etc/groups", want: 0, wantErr: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmdEval := cmd{cmdExprBuilder: utils.UpdateCmdExprParam, evalExpr: tt.evalExpr}
			got, err := cmdEval.evalExpression(tt.commandRes, tt.commResSize, make([]string, 0), tt.testFailure)
			if tt.want != got && err.Error() != tt.wantErr.Error() {
				t.Errorf("evalExpression() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecCommand(t *testing.T) {
	tests := []struct {
		name          string
		index         int
		prevResult    []string
		commandParams map[int][]string
		commandExec   []string
		newRes        []IndexValue
		want          string
		wantErr       error
	}{
		{name: "one command res and one param good", index: 0, prevResult: []string{}, commandParams: nil, commandExec: []string{"aaa"}, newRes: []IndexValue{}, want: "bbb", wantErr: fmt.Errorf("")},
		{name: "two command res and one param bad", index: 0, prevResult: []string{}, commandParams: nil, commandExec: []string{"aaa"}, newRes: []IndexValue{}, want: "EmptyValue", wantErr: fmt.Errorf("failed to exec command")},
	}
	zlog, err := zap.NewProduction()
	if err != nil {
		t.Errorf("failed to instansiate logger")
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			executor := NewMockExecutor(ctrl)
			for _, c := range tt.commandExec {
				executor.EXPECT().Exec(c).Return(&CommandResult{Stdout: tt.want, Stderr: tt.wantErr.Error()}, nil).Times(1)
			}
			cmdEval := cmd{command: executor, commandParams: tt.commandParams, commandExec: tt.commandExec, log: zlog}
			got := cmdEval.execCommand(tt.index, tt.prevResult, tt.newRes)
			sb := strings.Builder{}
			for _, s := range got {
				sb.WriteString(s)
			}
			if tt.want != sb.String() {
				t.Errorf("execCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
