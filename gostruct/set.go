package gostruct

type Empty struct{}
type Set struct {
	inner map[interface{}]Empty
}

func NewSet() *Set {
	return &Set{map[interface{}]Empty{}}
}

func (s *Set) Insert(key interface{}) {
	s.inner[key] = Empty{}
}

func (s *Set) Del(key interface{}) {
	delete(s.inner, key)
}

func (s *Set) Len(key interface{}) int {
	return len(s.inner)
}

func (s *Set) Clear(key interface{}) {
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