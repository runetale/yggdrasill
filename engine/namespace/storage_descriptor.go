package namespace

import "github.com/runetale/yggdrasill/types"

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
