[![Go Report Card](https://goreportcard.com/badge/github.com/chen-keinan/go-simple-config)](https://goreportcard.com/report/github.com/chen-keinan/go-simple-config)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/chen-keinan/go-simple-config/blob/master/LICENSE)
<img src="./pkg/img/coverage_badge.png" alt="test coverage badge">
[![Gitter](https://badges.gitter.im/beacon-sec/community.svg)](https://gitter.im/beacon-sec/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

# go-command-eval

Go-command-eval is an open source lib who evaluate shell command actual vs. expected results.

* [Installation](#installation)
* [Usage](#usage)

## Installation

```
go get github.com/chen-keinan/go-command-eval
```

## Usage
### one shell command with single / multiple result compared with string eval value

Executing shell command
```
commands:=[]string{"ls /etc/hosts | awk -F \" \" '{print $1}' |awk 'FNR <= 1'"}
```
evaluate command result with eval expr ( ${0} is the result from 1st shell command) 
```
evalExpr:="'${0}' == '/etc/hosts'"
```

### multi shell command with single / multiple result compared with string eval value

Executing two shell commands
```
commands:=[]string{"ls /etc/hosts | awk -F \" \" '{print $1}' |awk 'FNR <= 1'",
                    "ls /etc/group | awk -F \" \" '{print $1}' |awk 'FNR <= 1'"}
```
evaluate both shell command results with eval expr ${0} is the 1st / ${1} is the 2nd shell commands results
```
evalExpr:="'${0}' == '/etc/hosts'; && '${1}' == '/etc/group';"
```

### shell command  with IN Clause eval expr

Executing shell command
```
commands:=[]string{"ls /etc | awk -F \" \" '{print $1}' |awk 'FNR <= 2'"}
```
evaluate command result with IN Clause eval expr
```
evalExpr:="'${0}' IN ('afpovertcp.cfg','aliases')"
```

### shell command result passed as an arg to next shell command 

pass 1st command result as an arg to 2nd command
```
commands:=[]string{"ls /etc/hosts | awk -F " " '{print $1}' |awk 'FNR <= 1'",
                    "stat -f %A" ${0}}
```
eval both 1st and 2nd command result with eval expr
```
evalExpr:="'${0}' == '/etc/hosts'; && ${1} <= 766"
```

Full code example
```
cmdEval:= New()
cmdEvalResult:=cmdEval.EvalCommand(commands,evalExpr)
if cmdEvalResult.Match {
    fmt.Print("commmand result match eval expression")
}
```
