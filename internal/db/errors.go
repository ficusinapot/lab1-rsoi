package db

import "github.com/joomcode/errorx"

var (
	ErrorNamespace = errorx.NewNamespace("db").ApplyModifiers(errorx.TypeModifierOmitStackTrace)

	ConnectionError = ErrorNamespace.NewType("connection_error")
)
