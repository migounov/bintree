package bintree

import "errors"

type node struct {
	key         int
	left, right *node
	data        interface{}
}

type Updater func(data interface{}) (interface{}, error)

func (t *node) search(key int) (*node, *node) {
	var parent *node
	for n := t; n != nil; {
		switch {
		case key == n.key:
			return n, parent
		case key < n.key:
			parent = n
			n = n.left
		case key > n.key:
			parent = n
			n = n.right
		}
	}
	return nil, nil
}

func relinkParent(n *node, parent *node, newChild *node) *node {
	switch {
	case n.key < parent.key:
		parent.left = newChild
	case n.key > parent.key:
		parent.right = newChild
	}
	return parent
}

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

func (t *node) Get(key int) (interface{}, error) {
	n, _ := t.search(key)
	if n == nil {
		return nil, errors.New("key not found")
	}
	return n.data, nil
}

func (t *node) Update(key int, u Updater) error {
	n, _ := t.search(key)
	if n == nil {
		return errors.New("key not found")
	}
	updatedData, err := u(n.data)
	if err == nil {
		n.data = updatedData
	}
	return err
}

func (t *node) Delete(key int) (*node, error) {
	n, parent := t.search(key)

	if n == nil {
		return nil, errors.New("key not found")
	}

	switch {
	case n.left == nil && n.right == nil:
		if parent == nil {
			return nil, nil
		}
		relinkParent(n, parent, nil)
	case n.left != nil && n.right != nil:
		switch {
		case n.key > parent.key:
			r, p := n.left.findMax()
			relinkParent(r, p, r.left)
			r.left = n.left
			r.right = n.right
			relinkParent(n, parent, r)
		case n.key < parent.key:
			r, p := n.right.findMin()
			relinkParent(r, p, r.right)
			r.left = n.left
			r.right = n.right
			relinkParent(n, parent, r)
		}
	case n.left != nil:
		relinkParent(n, parent, n.left)
	case n.right != nil:
		relinkParent(n, parent, n.right)
	}
	return t, nil
}

func (t *node) findMin() (*node, *node) {
	var parent *node
	for n := t; n != nil; {
		if n.left == nil {
			return n, parent
		}
		parent = n
		n = n.left
	}
	return nil, nil
}

func (t *node) Min() int {
	min, _ := t.findMin()
	return min.key
}

func (t *node) findMax() (*node, *node) {
	var parent *node
	for n := t; n != nil; {
		if n.right == nil {
			return n, parent
		}
		parent = n
		n = n.right
	}
	return nil, nil
}

func (t *node) Max() int {
	max, _ := t.findMax()
	return max.key
}

func (t *node) traverse(fn func(t *node)) {
	if t != nil {
		t.left.traverse(fn)
		fn(t)
		t.right.traverse(fn)
	}
}

func (t *node) List() []int {
	var keys []int
	fn := func(t *node) {
		keys = append(keys, t.key)
	}
	t.traverse(fn)
	return keys
}
