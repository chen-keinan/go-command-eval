package eval

import (
	"testing"
)

func TestReverseString1(t *testing.T) {
	res := New()
	tests := []struct {
		name     string
		cmd      []string
		evalExpr string
		want     bool
	}{
		{name: "single command and evalExpr", cmd: []string{"ls /etc/hosts | awk -F \" \" '{print $1}' |awk 'FNR <= 1'"}, evalExpr: "'$0' == '/etc/hosts'", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := res.EvalCommand(tt.cmd, tt.evalExpr); got.Match != tt.want {
				t.Errorf("CvssScoreToSeverity() = %v, want %v", got, tt.want)
			}
		})
	}
}
