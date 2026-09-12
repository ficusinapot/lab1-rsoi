package core

import "github.com/joomcode/errorx"

var ErrorNamespace = errorx.NewNamespace("core").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
