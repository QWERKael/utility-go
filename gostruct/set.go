package gostruct

type Empty struct{}
type Set struct {
	inner map[interface{}]Empty
}

func NewEmptySet() *Set {
	return &Set{map[interface{}]Empty{}}
}

func NewSetFromStringList(values []string) *Set {
	s := NewEmptySet()
	for _, value := range values {
		s.Insert(value)
	}
	return s
}

func (s *Set) Insert(key interface{}) {
	s.inner[key] = Empty{}
}

func (s *Set) Del(key interface{}) {
	delete(s.inner, key)
}

func (s *Set) Len() int {
	return len(s.inner)
}

func (s *Set) Clear() {
	s.inner = make(map[interface{}]Empty)
}

func (s *Set) List() []interface{} {
	list := make([]interface{}, 0)
	for k := range s.inner {
		list = append(list, k)
	}
	return list
}

func (s *Set) Exists(key interface{}) bool {
	for k := range s.inner {
		if k == key {
			return true
		}
	}
	return false
}

func (s *Set) Equal(t *Set) bool {
	if s.Len() != t.Len() {
		return false
	}
	for k := range s.inner {
		if !t.Exists(k) {
			return false
		}
	}
	return true
}
