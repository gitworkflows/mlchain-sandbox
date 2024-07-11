package service

import (
	"time"

	"github.com/mlchain/mlchain-sandbox/internal/core/runner/python"
	runner_types "github.com/mlchain/mlchain-sandbox/internal/core/runner/types"
	"github.com/mlchain/mlchain-sandbox/internal/static"
	"github.com/mlchain/mlchain-sandbox/internal/types"
)

type RunCodeResponse struct {
	Stderr string `json:"error"`
	Stdout string `json:"stdout"`
}

func RunPython3Code(code string, preload string, options *runner_types.RunnerOptions) *types.MlchainSandboxResponse {
	if err := checkOptions(options); err != nil {
		return types.ErrorResponse(-400, err.Error())
	}

	timeout := time.Duration(
		static.GetMlchainSandboxGlobalConfigurations().WorkerTimeout * int(time.Second),
	)

	runner := python.PythonRunner{}
	stdout, stderr, done, err := runner.Run(
		code, timeout, nil, preload, options,
	)
	if err != nil {
		return types.ErrorResponse(-500, err.Error())
	}

	stdout_str := ""
	stderr_str := ""

	defer close(done)
	defer close(stdout)
	defer close(stderr)

	for {
		select {
		case <-done:
			return types.SuccessResponse(&RunCodeResponse{
				Stdout: stdout_str,
				Stderr: stderr_str,
			})
		case out := <-stdout:
			stdout_str += string(out)
		case err := <-stderr:
			stderr_str += string(err)
		}
	}
}

type ListDependenciesResponse struct {
	Dependencies []runner_types.Dependency `json:"dependencies"`
}

func ListPython3Dependencies() *types.MlchainSandboxResponse {
	return types.SuccessResponse(&ListDependenciesResponse{
		Dependencies: python.ListDependencies(),
	})
}

type UpdateDependenciesResponse struct{}

func UpdateDependencies() *types.MlchainSandboxResponse {
	err := python.PreparePythonDependenciesEnv()
	if err != nil {
		return types.ErrorResponse(-500, err.Error())
	}

	return types.SuccessResponse(&UpdateDependenciesResponse{})
}
