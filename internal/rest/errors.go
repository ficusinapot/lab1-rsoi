package rest

import "github.com/joomcode/errorx"

var ErrorNamespace = errorx.NewNamespace("rest").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
