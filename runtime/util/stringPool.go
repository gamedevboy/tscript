package util

import (
    "sync"
)

type StringPool interface {
    Insert(value string) *string
}

var _ StringPool = &stringPool{}

var globalPool = sync.Map{}

type stringPool struct {
    pool sync.Map
}

func (s *stringPool) Insert(value string) *string {
    ret, _ := s.pool.LoadOrStore(value, &value)
    return ret.(*string)
}

func NewStringPool() StringPool {
    return &stringPool{}
}
