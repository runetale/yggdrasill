package namespace

import (
	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/engine/action/shell"
	"github.com/runetale/notch/engine/action/tasklet"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/task"
	"github.com/runetale/notch/types"
)

type StorageDescriptor struct {
	name        string
	storageType types.StorageType
	predefined  *map[string]string
}

func NewStorageDescriptor(name string, storagetype types.StorageType, predefined *map[string]string) *StorageDescriptor {
	return &StorageDescriptor{
		name:        name,
		storageType: storagetype,
		predefined:  predefined,
	}
}

// managed all namespace actions
type Namespace struct {
	name        string
	description string
	stroages    []*storage.Storage
	action      action.Action
	// description of storages
	storageDescriptor *StorageDescriptor
}

// get namespace by types.Namespacetype
func NewNamespace(ns types.NamespaceType, functions []task.Function,
) *Namespace {
	var ac action.Action
	switch ns {
	case types.SHELL:
		ac = shell.NewShell()
	case types.CUSTOM:
		ac = tasklet.NewTasklet()
	case types.HTTP:
		predefined := map[string]string{}
		predefined["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
		predefined["Accept-Encoding"] = "deflate"
		// TODO: NewHTTP need for some pre header value
		// ac = tasklet.NewHTTP(predefined)
	}

	sd := NewStorageDescriptor(ac.Name(), ac.StorageType(), ac.Predefined())

	return &Namespace{
		name:              ac.Name(),
		description:       ac.Description(),
		stroages:          nil,
		action:            ac,
		storageDescriptor: sd,
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

func (n *Namespace) GetStorageType() types.StorageType {
	return n.storageDescriptor.storageType
}
