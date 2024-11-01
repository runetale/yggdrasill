package namespace

import (
	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/engine/action/shell"
	"github.com/runetale/notch/engine/action/tasklet"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/task"
	"github.com/runetale/notch/types"
)

// managed all namespace actions
type Namespace struct {
	Name        string
	Description string
	Stroages    []*storage.Storage
	Action      action.Action
}

func NewNamespace(ns types.NamespaceType, functions []task.Function) *Namespace {
	var ac action.Action
	switch ns {
	case types.SHELL:
		ac = shell.NewShell()
	case types.CUSTOM:
		ac = tasklet.NewTasklet()
	}

	return &Namespace{
		Name:        ac.Name(),
		Description: ac.Description(),
		Stroages:    nil,
		Action:      ac,
	}
}
