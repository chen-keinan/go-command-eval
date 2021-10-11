module github.com/chen-keinan/go-command-eval

go 1.16

replace github.com/chen-keinan/go-opa-validate => github.com/chen-keinan/go-opa-validate v0.0.6

require (
	github.com/Knetic/govaluate v3.0.0+incompatible
	github.com/chen-keinan/go-opa-validate v0.0.0-00010101000000-000000000000
	github.com/golang/mock v1.6.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.19.1
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac // indirect
)
