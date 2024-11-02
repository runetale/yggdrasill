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
	name        string
	description string
	stroages    []*storage.Storage
	action      action.Action
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
		name:        ac.Name(),
		description: ac.Description(),
		stroages:    nil,
		action:      ac,
	}
}

func (n *Namespace) GetName() string {
	return n.action.Name()
}

func (n *Namespace) GetDescription() string {
	return n.action.Description()
}

func (n *Namespace) GetAction() action.Action {
	return n.action
}

func (n *Namespace) GetStorages() []*storage.Storage {
	return n.stroages
}
