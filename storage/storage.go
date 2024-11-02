package storage

import (
	"sort"
	"strconv"
	"time"

	"github.com/runetale/notch/engine/events"
	"github.com/runetale/notch/types"
)

type Entry struct {
	time     time.Duration
	complete bool // for Completion storage
	data     string
}

func NewEntry(data string) *Entry {
	return &Entry{
		time:     time.Duration(time.Now().UnixNano()),
		complete: false,
		data:     data,
	}
}

const CURRENT_TAG = "current"
const PREVIOUS_TAG = "previous"
const STARTED_AT_TAG = "started_at"

type Storage struct {
	name        string
	storageType types.StorageType
	entry       map[string]*Entry

	OnEventCallback func(event events.Event)
}

func NewStorage(name string, storageType types.StorageType,
	OnEventCallback func(event events.Event),
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
		return entries[i].time < entries[j].time
	})
}

func (s *Storage) GetName() string {
	return s.name
}

func (s *Storage) GetType() types.StorageType {
	return s.storageType
}

func (s *Storage) GetStartedAt() time.Duration {
	return s.entry[STARTED_AT_TAG].time
}

func (s *Storage) OnEvent(event string) {
	s.OnEventCallback(events.Event(event))
}

func (s *Storage) AddData(key, data string) {
	s.entry[key] = NewEntry(data)
	s.OnEvent("add data")
}

func (s *Storage) AddTagged(key, data string) {
	s.entry[key] = NewEntry(data)
	s.OnEvent("add tagged")
}

func (s *Storage) DelTagged(key string) {
	s.entry[key] = nil
	s.OnEvent("del tagged")
}

func (s *Storage) GetTagged(key string) string {
	inner := s.entry[key]
	s.OnEvent("get tagged")
	return inner.data
}

// for planning tasks
func (s *Storage) AddCompletion(data string) {
	var keys []string
	for key := range s.entry {
		keys = append(keys, key)
	}

	lastKey := keys[len(keys)-1]
	lastValue := s.entry[lastKey]
	lastValue.data = data
	s.OnEvent("get completion")
}

func (s *Storage) DelCompletion(pos int) {
	tag := strconv.Itoa(pos)
	s.entry[tag] = nil
	s.OnEvent("delete completion")
}

func (s *Storage) SetComplete(pos int) bool {
	tag := strconv.Itoa(pos)
	s.entry[tag].complete = true
	s.OnEvent("set complete")
	return true
}

func (s *Storage) SetInComplete(pos int) bool {
	tag := strconv.Itoa(pos)
	s.entry[tag].complete = false
	s.OnEvent("set incomplete")
	return true
}

func (s *Storage) SetCurrent(data string) {
	s.entry[CURRENT_TAG] = NewEntry(data)
	s.OnEvent("set current")
}

func (s *Storage) Clear(data string) {
	s.entry = make(map[string]*Entry, 0)
	s.OnEvent("clear storage")
}
