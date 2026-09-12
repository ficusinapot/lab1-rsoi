package persons

import "github.com/joomcode/errorx"

var personErrors = errorx.NewNamespace("persons")

var ErrPersonNotFound = personErrors.NewType("not_found")
