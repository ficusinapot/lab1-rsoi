package errors

import "github.com/joomcode/errorx"

var commandErrors = errorx.NewNamespace("service_command")

var CommandLineProcessingError = commandErrors.NewType("command_line_processing")
