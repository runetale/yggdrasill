package storage

import (
	"sort"
	"strconv"
	"time"

	"github.com/runetale/notch/events"
	"github.com/runetale/notch/types"
)

type Entry struct {
	Time     time.Time
	Complete bool // for Completion storage
	Data     string
}

func NewEntry(data string) *Entry {
	return &Entry{
		Time:     time.Now(),
		Complete: false,
		Data:     data,
	}
}

const CURRENT_TAG = "current"
const PREVIOUS_TAG = "previous"
const STARTED_AT_TAG = "started_at"

type Storage struct {
	name        string
	storageType types.StorageType
	entry       map[string]*Entry

	OnEventCallback func(event events.DisplayEvent)
}

func NewStorage(name string, storageType types.StorageType,
	OnEventCallback func(event events.DisplayEvent),
) *Storage {
	entry := make(map[string]*Entry, 0)

	return &Storage{
		name:            name,
		storageType:     storageType,
		entry:           entry,
		OnEventCallback: OnEventCallback,
	}
}

func (s *Storage) SortedEntries() {
	entries := make([]*Entry, 0, len(s.entry))
	for _, entry := range s.entry {
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Time.UnixNano() < entries[j].Time.UnixNano()
	})
}

func (s *Storage) GetName() string {
	return s.name
}

func (s *Storage) GetEntries() []*Entry {
	values := []*Entry{}
	for _, entry := range s.entry {
		values = append(values, entry)
	}
	return values
}

func (s *Storage) GetEntryList() map[string]*Entry {
	return s.entry
}

func (s *Storage) GetEntry(key string) (*Entry, bool) {
	e, found := s.entry[key]
	return e, found
}

func (s *Storage) IsEmpty() bool {
	return len(s.entry) == 0
}

func (s *Storage) GetStorageType() types.StorageType {
	return s.storageType
}

func (s *Storage) GetStartedAt() time.Time {
	return s.entry[STARTED_AT_TAG].Time
}

func (s *Storage) OnEvent(event events.DisplayEvent) {
	s.OnEventCallback(event)
}

func (s *Storage) AddData(key, data string) {
	s.entry[key] = NewEntry(data)
	s.OnEvent(events.NewStorageUpdateEvent(s.name, s.storageType, key, nil, &data))
}

func (s *Storage) AddTagged(key, data string) {
	s.entry[key] = NewEntry(data)
	s.OnEvent(events.NewStorageUpdateEvent(s.name, s.storageType, key, nil, &data))
}

func (s *Storage) DelTagged(key string) {
	if old, exists := s.entry[key]; exists {
		delete(s.entry, key)
		s.OnEvent(events.NewStorageUpdateEvent(s.name, s.storageType, key, &old.Data, nil))
	}
}

func (s *Storage) GetTagged(key string) string {
	inner := s.entry[key]
	return inner.Data
}

// for planning tasks
func (s *Storage) AddCompletion(data string) {
	var keys []string
	for key := range s.entry {
		keys = append(keys, key)
	}

	lastKey := keys[len(keys)-1]
	lastValue := s.entry[lastKey]
	lastValue.Data = data
	s.OnEvent(events.NewStorageUpdateEvent(s.name, s.storageType, lastKey, nil, &data))
}

func (s *Storage) DelCompletion(pos int) {
	tag := strconv.Itoa(pos)
	if old, exists := s.entry[tag]; exists {
		delete(s.entry, tag)
		s.OnEvent(events.NewStorageUpdateEvent(s.name, s.storageType, tag, &old.Data, nil))
	}
}

func (s *Storage) SetComplete(pos int) bool {
	tag := strconv.Itoa(pos)
	if entry, exists := s.entry[tag]; exists {
		prev := entry.Complete
		entry.Complete = true

		prevState := "incomplete"
		if prev {
			prevState = "complete"
		}
		comp := "complete"
		s.OnEvent(events.NewStorageUpdateEvent(s.name, s.storageType, tag, &prevState, &comp))
		return prev
	}

	return false
}

func (s *Storage) SetInComplete(pos int) bool {
	tag := strconv.Itoa(pos)
	if entry, exists := s.entry[tag]; exists {
		prev := entry.Complete
		entry.Complete = false

		prevState := "incomplete"
		if prev {
			prevState = "complete"
		}
		comp := "incomplete"
		s.OnEvent(events.NewStorageUpdateEvent(s.name, s.storageType, tag, &prevState, &comp))
		return prev
	}

	return false
}

func (s *Storage) SetCurrent(data string) {
	if s.storageType != types.CURRENTPREVIOUS {
		panic("storage type must be CurrentPrevious")
	}

	oldCurrent, exists := s.entry[CURRENT_TAG]
	s.entry[CURRENT_TAG] = NewEntry(data)

	prev := ""
	if exists {
		prev = oldCurrent.Data
		s.entry[PREVIOUS_TAG] = oldCurrent
	}
	s.OnEvent(events.NewStorageUpdateEvent(s.name, s.storageType, CURRENT_TAG, &prev, &data))
}

func (s *Storage) Clear() {
	s.entry = make(map[string]*Entry, 0)
	s.OnEvent(events.NewStorageUpdateEvent(s.name, s.storageType, "", nil, nil))
}
