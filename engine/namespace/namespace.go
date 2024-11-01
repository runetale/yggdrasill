package namespace

import (
	"github.com/runetale/notch/engine/namespace/shell"
	"github.com/runetale/notch/types"
)

// managed all namespace actions
type Namespace struct {
	Shell *shell.Shell
}

func NewNamespace(ns types.NamespaceType) *Namespace {
	var sh *shell.Shell
	switch ns {
	case types.SHELL:
		sh = shell.NewShell()
	}

	return &Namespace{
		Shell: sh,
	}
}
