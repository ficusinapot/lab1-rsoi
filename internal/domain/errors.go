package domain

import "github.com/joomcode/errorx"

var ErrorNamespace = errorx.NewNamespace("coreifc").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
