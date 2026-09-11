//go:build tools

// Package tools tracks generator dependencies that are not imported by runtime code.
package tools

import (
	_ "entgo.io/ent/cmd/ent"
)
