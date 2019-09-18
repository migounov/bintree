package bintree

import "errors"

type node struct {
	key         int
	left, right *node
	data        interface{}
	height		int
}

type Updater func(data interface{}) (interface{}, error)

func (t *node) getHeight() int {
	if t == nil {
		return 0
	} else {
		return t.height
	}
}

func (t *node) calcHeight() *node {
	hl := t.left.getHeight()
	hr := t.right.getHeight()
	if hl > hr {
		t.height = hl + 1
	} else {
		t.height = hr + 1
	}
	return t
}

func (t *node) getBalanceFactor() int {
	return t.right.getHeight() - t.left.getHeight()
}

func (t *node) balance() *node {
	t.calcHeight()
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
	t.calcHeight()
	z.calcHeight()
	return z
}

func (t *node) rotateRight() *node {
	z := t.left
	t.left = z.right
	z.right = t
	t.calcHeight()
	z.calcHeight()
	return z
}

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
		r, p := n.findMin(n.right)
		relinkParent(r, p, r.right)
		r.left = n.left
		r.right = n.right
		if parent != nil {
			relinkParent(n, parent, r)
		} else {
			return r.balance(), nil
		}
	case n.left != nil:
		relinkParent(n, parent, n.left)
	case n.right != nil:
		relinkParent(n, parent, n.right)
	}
	return t.balance(), nil
}

func (t *node) findMin(s *node) (*node, *node) {
	var n *node
	parent := t
	for n = s; n != nil && n.left != nil; {
		parent = n
		n = n.left
	}
	return n, parent
}

func (t *node) Min() int {
	min, _ := t.findMin(t.left)
	return min.key
}

func (t *node) findMax(s *node) (*node, *node) {
	var n *node
	parent := t
	for n = t; n != nil && n.right != nil; {
		parent = n
		n = n.right
	}
	return n, parent
}

func (t *node) Max() int {
	max, _ := t.findMax(t.right)
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

func (t *node) ListBalanceFactors() []int {
	var keys []int
	fn := func(t *node) {
		keys = append(keys, t.getBalanceFactor())
	}
	t.traverse(fn)
	return keys
}
