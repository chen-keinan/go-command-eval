package eval

import (
	"testing"
)

func TestEvalCommand(t *testing.T) {
	res := New()
	tests := []struct {
		name     string
		cmd      []string
		evalExpr string
		want     bool
	}{
		{name: "single command and evalExpr match", cmd: []string{"ls /etc/hosts | awk -F \" \" '{print $1}' |awk 'FNR <= 1'"}, evalExpr: "'$0' == '/etc/hosts'", want: true},
		{name: "two command and evalExpr match", cmd: []string{"ls /etc/hosts | awk -F \" \" '{print $1}' |awk 'FNR <= 1'", "ls /etc/group | awk -F \" \" '{print $1}' |awk 'FNR <= 1'"}, evalExpr: "'$0' == '/etc/hosts'; && '$1' == '/etc/group';", want: true},
		{name: "two command and evalExpr do one not match", cmd: []string{"ls /etc/hosts | awk -F \" \" '{print $1}' |awk 'FNR <= 1'", "ls /etc/group | awk -F \" \" '{print $1}' |awk 'FNR <= 1'"}, evalExpr: "'$0' == '/etc/hosts'; && '$1' == '/etc/group1';", want: false},
		{name: "two command and evalExpr do both not match", cmd: []string{"ls /etc/hosts | awk -F \" \" '{print $1}' |awk 'FNR <= 1'", "ls /etc/group | awk -F \" \" '{print $1}' |awk 'FNR <= 1'"}, evalExpr: "'$0' == '/etc/hosts1'; && '$1' == '/etc/group1';", want: false},
		{name: "command and evalExpr in clause match", cmd: []string{"ls /etc | awk -F \" \" '{print $1}' |awk 'FNR <= 2'"}, evalExpr: "'$0' IN ('afpovertcp.cfg','aliases')", want: true},
		{name: "command and evalExpr in clause notmatch", cmd: []string{"ls /etc | awk -F \" \" '{print $1}' |awk 'FNR <= 2'"}, evalExpr: "'$0' IN ('afpovertcp.cfg1','aliases')", want: false},
		{name: "single command and evalExpr match", cmd: []string{"ls /etc/hosts | awk -F \" \" '{print $1}' |awk 'FNR <= 1'", "grep local #0| awk -F \"127.0.0.1\" '{print $1}' |awk 'FNR <= 1' | awk -F \" \" '{print $3}' |awk 'FNR <= 1'", "grep local #0 | awk -F \"localhost\" '{print $2}' |awk 'FNR <= 2'|grep tservice | awk -F \" \" '{print $2}'"}, evalExpr: "'$0' == '/etc/hosts'; && '$1' == 'is'; && '$2' == 'tservice';", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := res.EvalCommand(tt.cmd, tt.evalExpr); got.Match != tt.want {
				t.Errorf("EvalCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
