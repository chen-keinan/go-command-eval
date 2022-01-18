package main

import (
	"fmt"
	"github.com/chen-keinan/go-command-eval/eval"
)

const (
	policy = `package test

policy_eval = {"kind":kind,"match":allow_policy} {
	input.kind == "PodList"
  	kind = "aaa"
    allow_policy = true
  }`
)

func main() {
	evlResult := eval.NewEvalCmd().EvalKubeAPIPolicy([]string{"namespaces | jq .items |.[]|.metadata.name", "namespaces/${0}/pods"}, "[${1} MATCH allow_mtls_permissive_mode.policy QUERY test.policy_eval RETURN match,kind]", policy)
	fmt.Println(evlResult)
}
