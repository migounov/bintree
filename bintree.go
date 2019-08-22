package bintree

import "errors"

type node struct {
	key         int
	left, right *node
	data        interface{}
}

type Updater func(data interface{}) (interface{}, error)

func (t *node) Insert(key int, data interface{}) (*node, error) {
	var err error
	if t == nil {
		return &node{key, nil, nil, data}, nil
	}
	if key == t.key {
		return t, errors.New("duplicate key")
	}
	if key < t.key {
		t.left, err = t.left.Insert(key, data)
		return t, err
	}
	t.right, err = t.right.Insert(key, data)
	return t, err
}

func (t *node) Get() (interface{}, error) {
	return nil, nil
}

func (t *node) Update(key int, updater Updater) error {
	return nil
}

func (t *node) Delete(key int) error {
	return nil
}

func (t *node) Min() int {
	return -1
}

func (t *node) Max() int {
	return -1
}

func (t *node) traverse(fn func(t *node)) {
	if t != nil {
		t.left.traverse(fn)
		fn(t)
		t.right.traverse(fn)
	}
}

func (t *node) list() []int {
	var keys []int
	fn := func(t *node) {
		keys = append(keys, t.key)
	}
	t.traverse(fn)
	return keys
}
