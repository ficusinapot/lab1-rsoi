package manager

import "github.com/joomcode/errorx"

var (
	ErrorNamespace = errorx.NewNamespace("manager").ApplyModifiers(errorx.TypeModifierOmitStackTrace)

	InvalidStateError       = ErrorNamespace.NewType("invalid_state")
	PrepareError            = ErrorNamespace.NewType("prepare")
	RunError                = ErrorNamespace.NewType("run")
	StopCommunicationsError = ErrorNamespace.NewType("stop_communications")
	ShutdownError           = ErrorNamespace.NewType("shutdown")
)
