package eval

import (
	"bytes"
	"context"
	"github.com/chen-keinan/go-command-eval/utils"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/exec"
	"path"
	"strings"
)

//ShellToUse bash shell
const (
	//jqPrefix expr
	jqPrefix = "| jq"
	// ShellToUse shell command
	ShellToUse = "sh"
)

//Executor defines the interface for shell command executor
//exec.go
//go:generate mockgen -destination=./mock_Executor.go -package=eval . Executor
type Executor interface {
	Exec(command string) (*CommandResult, error)
}

//CommandExec object
type CommandExec struct {
}

//KubeClientExec object
type KubeClientExec struct {
	client *rest.RESTClient
}

//Exec make api call to k8s apiserver
func (k KubeClientExec) Exec(endpoint string) (*CommandResult, error) {
	var jqExpr string
	if strings.Contains(endpoint, jqPrefix) {
		endpointParts := strings.Split(endpoint, jqPrefix)
		if len(endpointParts[1]) > 0 {
			jqExpr = endpointParts[1]
			endpoint = strings.TrimSpace(endpointParts[0])
		}
	}
	var errString string
	kubeAPIRes, err := k.client.Get().AbsPath(path.Join("/api/v1/", endpoint)).DoRaw(context.TODO())
	if err != nil {
		return &CommandResult{Stdout: "", Stderr: errString}, err
	}
	out, err := utils.RunJqQuery(jqExpr, kubeAPIRes)
	if err != nil {
		errString = err.Error()
	}
	return &CommandResult{Stdout: out, Stderr: errString}, err
}

// NewKubeClientExec return new instance of kube client executor
func NewKubeClientExec() Executor {
	var config *rest.Config
	var err error
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		panic(err)
	}
	crdConfig := *config
	//crdConfig.APIPath = "/apis"
	crdConfig.APIPath = "/api/v1"
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	exampleRestClient, err := rest.UnversionedRESTClientFor(&crdConfig)
	if err != nil {
		panic(err)
	}
	return &KubeClientExec{client: exampleRestClient}
}

//NewShellExec return new instance of shell executor
func NewShellExec() Executor {
	return &CommandExec{}
}

//CommandResult return data
type CommandResult struct {
	Stdout string
	Stderr string
}

//Exec execute shell command
// #nosec
func (ce CommandExec) Exec(command string) (*CommandResult, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return &CommandResult{Stdout: stdout.String(), Stderr: stderr.String()}, err
}
