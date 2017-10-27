package main

import (
	"sync"
)

type set struct {
	lock     sync.Mutex
	innerMap map[int]uint8
}

func (s *set) Length() int {
	s.lock.Lock()
	defer s.lock.Unlock()

	return len(s.innerMap)
}

func (s *set) Add(i int) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, hasIt := s.innerMap[i]; hasIt {
		return false
	}

	s.innerMap[i] = 1
	return true
}

func (s *set) Pop(i int) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, hasIt := s.innerMap[i]; hasIt {
		return false
	}

	delete(s.innerMap, i)
	return true
}

func (s *set) PopOne() (int, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for k := range s.innerMap {
		delete(s.innerMap, k)
		return k, true
	}
	return 0, false
}
