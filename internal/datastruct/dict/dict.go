package dict

import (
	"sync"
)

type Consumer func(key string, val any)

type Dict interface {
	Get(key string) (val any, exist bool)
	Len() int
	Put(key string, val any) (result int)
	PutIfAbsent(key string, val any) (result int)
	PutIfExists(key string, val any) (result int)
	Remove(key string) (result int)
	ForEach(consumer Consumer)
	Keys() []string
	RandomKeys(limit int) []string
	RandomDistinctKeys(limit int) []string
	Clear()
}

type SyncDict struct {
	*sync.Map
}

func (s *SyncDict) Get(key string) (val any, exist bool) {
	return s.Load(key)
}

func (s *SyncDict) Len() int {
	length := 0
	s.Range(func(key, value any) bool {
		length++
		return true
	})
	return length
}

func (s *SyncDict) Put(key string, val any) (result int) {
	_, ok := s.Get(key)
	s.Store(key, val)
	if ok {
		return 0
	}
	return 1
}

func (s *SyncDict) PutIfAbsent(key string, val any) (result int) {
	_, exist := s.Get(key)
	if exist {
		return 0
	}
	s.Store(key, val)
	return 1
}

func (s *SyncDict) PutIfExists(key string, val any) (result int) {
	_, exist := s.Get(key)
	if !exist {
		return 0
	}
	s.Store(key, val)
	return 1
}

func (s *SyncDict) Remove(key string) (result int) {
	_, exist := s.Get(key)
	s.Delete(key)
	if exist {
		return 1
	}
	return 0
}

func (s *SyncDict) ForEach(consumer Consumer) {
	s.Range(func(key, value any) bool {
		consumer(key.(string), value)
		return true
	})
}

func (s *SyncDict) Keys() []string {
	result := make([]string, s.Len())
	i := 0
	s.Range(func(key, value any) bool {
		result[i] = key.(string)
		i++
		return true
	})
	return result
}

func (s *SyncDict) RandomKeys(limit int) []string {
	result := make([]string, s.Len())
	for i := range limit {
		s.Range(func(key, value any) bool {
			result[i] = key.(string)
			return true
		})
	}
	return result
}

func (s *SyncDict) RandomDistinctKeys(limit int) []string {
	result := make([]string, s.Len())
	i := 0
	s.Range(func(key, value any) bool {
		result[i] = key.(string)
		i++
		if i == limit {
			return false
		}
		return true
	})
	return result
}

func (s *SyncDict) Clear() {
	s.Map = new(sync.Map)
}

func NewSyncDict() Dict {
	return &SyncDict{Map: new(sync.Map)}
}
