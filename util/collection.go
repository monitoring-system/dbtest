package util

type Set struct {
	m map[interface{}]struct{}
}

func (s *Set) Contains(item interface{}) bool {
	_, ok := s.m[item]
	return ok
}

var Exists = struct{}{}

func NewSet(items ...interface{}) *Set {
	s := &Set{}
	s.m = make(map[interface{}]struct{})
	s.Put(items...)
	return s
}

func (s *Set) Put(items ...interface{}) error {
	for _, item := range items {
		s.m[item] = Exists
	}
	return nil
}

func (s *Set) Size() int {
	return len(s.m)
}

func (s *Set) ToSlice() []interface{} {
	var data []interface{}
	for key, _ := range s.m {
		data = append(data, key)
	}
	return data
}
