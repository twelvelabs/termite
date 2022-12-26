package run

//go:generate moq --rm --out=executor_test.go . Executor

// Executor is an interface representing the ability to execute an external command.
type Executor interface {
	// ExitCode returns the exit code of the exited process, or -1
	// if the process hasn't exited or was terminated by a signal.
	ExitCode(cmd *Cmd) int
	// Output runs the command and returns its standard output.
	Output(cmd *Cmd) ([]byte, error)
	// Run starts the specified command and waits for it to complete.
	Run(cmd *Cmd) error
}
