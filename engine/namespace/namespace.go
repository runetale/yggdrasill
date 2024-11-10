package namespace

import (
	"github.com/runetale/notch/engine/action"
	"github.com/runetale/notch/engine/action/goal"
	"github.com/runetale/notch/engine/action/memory"
	"github.com/runetale/notch/engine/action/shell"
	"github.com/runetale/notch/engine/action/tasklet"
	"github.com/runetale/notch/storage"
	"github.com/runetale/notch/task"
	"github.com/runetale/notch/types"
)

type StorageDescriptor struct {
	name        string
	storageType types.StorageType
	predefined  map[string]*string
}

func NewStorageDescriptor(name string, storagetype types.StorageType, predefined map[string]*string) *StorageDescriptor {
	return &StorageDescriptor{
		name:        name,
		storageType: storagetype,
		predefined:  predefined,
	}
}

func (s *StorageDescriptor) Name() string {
	return s.name
}

func (s *StorageDescriptor) Type() types.StorageType {
	return s.storageType
}

func (s *StorageDescriptor) Predefined() map[string]*string {
	return s.predefined
}

func (s *StorageDescriptor) StorageType() types.StorageType {
	return s.storageType
}

// managed all namespace actions
type Namespace struct {
	name        string
	description string
	stroages    []*storage.Storage
	actions     []action.Action
	// description of storages, using memory types
	storageDescriptor []*StorageDescriptor
}

// get namespace by types.Namespacetype
func NewNamespace(ns types.NamespaceType, functions []*task.Function,
) *Namespace {
	var (
		name        string
		description string
	)

	actions := []action.Action{}
	descriptors := []*StorageDescriptor{}

	switch ns {
	case types.SHELL:
		s := shell.NewShell()
		name = "Shell"
		description = s.NamespaceDescription()
		actions = append(actions, s)
	case types.TASKLET:
		t := tasklet.NewTasklet()
		name = "Task"
		description = t.NamespaceDescription()
		actions = append(actions, t)
	case types.GOAL:
		g := goal.NewGoal()
		name = "Goal"
		description = g.NamespaceDescription()
		actions = append(actions, g)
		descriptors = append(descriptors, NewStorageDescriptor("goal", types.CURRENTPREVIOUS, nil))
	case types.MEMORY:
		sm := memory.NewSaveMemroy()
		dm := memory.NewDeleteMemory()
		name = "Memory"
		description = sm.NamespaceDescription()
		actions = append(actions, sm)
		actions = append(actions, dm)
		descriptors = append(descriptors, NewStorageDescriptor("memories", types.TAGGED, nil))
	case types.HTTP:
		predefined := map[string]string{}
		predefined["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
		predefined["Accept-Encoding"] = "deflate"
		// TODO: NewHTTP need for some pre header value
		// ac = tasklet.NewHTTP(predefined)
	default:
		panic("not implemented namespaces")
	}

	return &Namespace{
		name:              name,
		description:       description,
		actions:           actions,
		storageDescriptor: descriptors,
	}
}

func (n *Namespace) Description() string {
	return n.description
}

func (n *Namespace) GetStrorageDescriptor() []*StorageDescriptor {
	return n.storageDescriptor
}

func (n *Namespace) Name() string {
	return n.name
}

func (n *Namespace) Actions() []action.Action {
	return n.actions
}

func (n *Namespace) SetStorage(s *storage.Storage) {
	n.stroages = append(n.stroages, s)
}

// list of actions with action itself
func (n *Namespace) GetActions() []action.Action {
	return n.actions
}
