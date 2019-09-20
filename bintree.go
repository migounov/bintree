package bintree

import "errors"

type node struct {
	key         int
	left, right *node
	data        interface{}
	height      int
}

type Updater func(data interface{}) (interface{}, error)

func (t *node) getHeight() int {
	if t == nil {
		return 0
	} else {
		return t.height
	}
}

func (t *node) setHeight() {
	hl := t.left.getHeight()
	hr := t.right.getHeight()
	if hl > hr {
		t.height = hl + 1
	} else {
		t.height = hr + 1
	}
}

func (t *node) getBalanceFactor() int {
	return t.right.getHeight() - t.left.getHeight()
}

func (t *node) balance() *node {
	t.setHeight()
	if t.getBalanceFactor() == 2 {
		if t.right.getBalanceFactor() < 0 {
			t.right = t.right.rotateRight()
		}
		return t.rotateLeft()
	} else if t.getBalanceFactor() == -2 {
		if t.left.getBalanceFactor() > 0 {
			t.left = t.left.rotateLeft()
		}
		return t.rotateRight()
	}
	return t
}

func (t *node) rotateLeft() *node {
	z := t.right
	t.right = z.left
	z.left = t
	t.setHeight()
	z.setHeight()
	return z
}

func (t *node) rotateRight() *node {
	z := t.left
	t.left = z.right
	z.right = t
	t.setHeight()
	z.setHeight()
	return z
}

func (t *node) search(key int) *node {
	for n := t; n != nil; {
		switch {
		case key == n.key:
			return n
		case key < n.key:
			n = n.left
		case key > n.key:
			n = n.right
		}
	}
	return nil
}

func (t *node) Insert(key int, data interface{}) (*node, error) {
	var err error
	if t == nil {
		return &node{key, nil, nil, data, 1}, nil
	}
	if key == t.key {
		return t, errors.New("duplicate key")
	} else if key < t.key {
		t.left, err = t.left.Insert(key, data)
	} else if key > t.key {
		t.right, err = t.right.Insert(key, data)
	}
	return t.balance(), err
}

func (t *node) Get(key int) (interface{}, error) {
	n := t.search(key)
	if n == nil {
		return nil, errors.New("key not found")
	}
	return n.data, nil
}

func (t *node) Update(key int, u Updater) error {
	n := t.search(key)
	if n == nil {
		return errors.New("key not found")
	}
	updatedData, err := u(n.data)
	if err == nil {
		n.data = updatedData
	}
	return err
}

func (t *node) Delete(key int) *node {
	switch {
	case key < t.key:
		t.left = t.left.Delete(key)
	case key > t.key:
		t.right = t.right.Delete(key)
	case key == t.key:
		l := t.left
		r := t.right
		if r == nil {
			return l
		}
		min := r.findMin()
		min.right = r.deleteMin()
		min.left = l
		return min.balance()
	}
	return t.balance()
}

func (t *node) findMin() *node {
	if t.left != nil {
		return t.left.findMin()
	} else {
		return t
	}
}

func (t *node) deleteMin() *node {
	if t.left == nil {
		return t.right
	}
	t.left = t.left.deleteMin()
	return t.balance()
}

func (t *node) Min() int {
	return t.findMin().key
}

func (t *node) findMax() *node {
	if t.right != nil {
		return t.right.findMax()
	} else {
		return t
	}
}

func (t *node) Max() int {
	return t.findMax().key
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

func (t *node) ListBalanceFactors() []int {
	var keys []int
	fn := func(t *node) {
		keys = append(keys, t.getBalanceFactor())
	}
	t.traverse(fn)
	return keys
}
