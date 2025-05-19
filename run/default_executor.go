package run

var (
	// DefaultExecutor is the default implementation of Executor used by new Clients.
	// Each of it's methods simply delegate to the underlying [os/exec.Cmd].
	DefaultExecutor Executor = &defaultExecutor{}
)

type defaultExecutor struct {
}

func (e *defaultExecutor) ExitCode(cmd *Cmd) int {
	return cmd.ProcessState.ExitCode()
}

func (e *defaultExecutor) Output(cmd *Cmd) ([]byte, error) {
	return cmd.Cmd.Output()
}

func (e *defaultExecutor) Run(cmd *Cmd) error {
	return cmd.Cmd.Run()
}
