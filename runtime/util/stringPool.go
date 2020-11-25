package util

import (
	"sync"
)

type StringPool interface {
	Insert(value string) string
}

var _ StringPool = &stringPool{}

var globalPool = sync.Map{}

type stringPool struct {
	pool map[string]string
}

func (s *stringPool) Insert(value string) string {
	if ret, ok := s.pool[value]; ok {
		return ret
	}

	s.pool[value] = value

	return s.pool[value]
}

func NewStringPool() StringPool {
	return &stringPool{
		pool: make(map[string]string),
	}
}
