[![Go Report Card](https://goreportcard.com/badge/github.com/chen-keinan/go-simple-config)](https://goreportcard.com/report/github.com/chen-keinan/go-simple-config)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/chen-keinan/go-simple-config/blob/master/LICENSE)
<img src="./pkg/img/coverage_badge.png" alt="test coverage badge">
[![Gitter](https://badges.gitter.im/beacon-sec/community.svg)](https://gitter.im/beacon-sec/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

<br><img src="./pkg/img/cmd_eval.png" width="300" alt="cmd_eval logo"><br>
# go-command-eval

Go-command-eval is an open source lib who evaluate shell command results against eval expr.

* [Installation](#installation)
* [Usage](#usage)
* [Contribution](#Contribution)


## Installation

```
go get github.com/chen-keinan/go-command-eval
```

## Usage
### one shell command with single result evaluated against eval expression

create shell command which return one result
```
commands:=[]string{"ls /etc/hosts | awk -F \" \" '{print $1}' |awk 'FNR <= 1'"}
```
evaluate command result with eval expression ( ${0} is the result from 1st shell command) 
```
evalExpr := "'${0}' == '/etc/hosts'"
```

### two shell commands with single result each evaluated with eval expression

create two shell commands with one result for each
```
commands := []string{"ls /etc/hosts | awk -F \" \" '{print $1}' |awk 'FNR <= 1'",
                    "ls /etc/group | awk -F \" \" '{print $1}' |awk 'FNR <= 1'"}
```
evaluate each command result with eval expression
```
evalExpr := "'${0}' == '/etc/hosts'; && '${1}' == '/etc/group';"
```

### shell command return two results evaluated with IN Clause eval expression

create shell command with return two results
```
commands := []string{"ls /etc | awk -F \" \" '{print $1}' |awk 'FNR <= 2'"}
```
evaluate command result with IN Clause eval expression
```
evalExpr := "'${0}' IN ('afpovertcp.cfg','aliases')"
```

### shell command result passed as an arg to the following shell command; both results are evaluated against eval expression

create tow shell commands 1st command result passed as an arg to the following shell command
```
commands := []string{"ls /etc/hosts | awk -F " " '{print $1}' |awk 'FNR <= 1'",
                    "stat -f %A" ${0}}
```
both results are evaluated against eval expression 1st result is evaluated as string 
2nd result is evaluated as an Integer
```
evalExpr := "'${0}' == '/etc/hosts'; && ${1} <= 766"
```

Full code example
```
commands := []string{"ls /etc/hosts | awk -F \" \" '{print $1}' |awk 'FNR <= 1'",
		"ls /etc/group | awk -F \" \" '{print $1}' |awk 'FNR <= 1'"}
		
evalExpr := "'${0}' == '/etc/hosts'; && '${1}' == '/etc/group';"

cmdEval := NewEvalCmd()
cmdEvalResult := cmdEval.EvalCommand(commands, evalExpr)
if cmdEvalResult.Match {
    fmt.Print("commmand result match eval expression")
}
```


## Contribution
code contribution is welcome !! , contribution with passing tests and linter is more than welcome :)