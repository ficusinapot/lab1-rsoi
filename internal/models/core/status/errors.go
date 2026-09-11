package status

import "github.com/joomcode/errorx"

var statusErrors = errorx.NewNamespace("status")

var ErrConnectionNotAvailable = statusErrors.NewType("connection_not_available")
