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
### one shell command with single multiple result compared with string eval value

Executing shell command
```
commands:=[]string{"ls /etc/hosts | awk -F \" \" '{print $1}' |awk 'FNR <= 1'"}
```
Assigning match eval expr , ${0} is the result from 1st shell command to match 
```
evalExpr:="'${0}' == '/etc/hosts'"
```
Full code example
```
cmdEval:= New()
cmdEvalResult:=cmdEval.EvalCommand(commands,evalExpr)
if cmdEvalResult.Match {
    fmt.Print("commmand result match eval expression")
}
```
