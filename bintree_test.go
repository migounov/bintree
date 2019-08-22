package bintree

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestInsert(t *testing.T) {
	var exp []int
	var err error

	key := rand.Intn(100)
	user := "User " + strconv.Itoa(key)

	tr := &node{
		key:  key,
		data: user,
	}

	exp = append(exp, key)

	for i := 1; i <= 100; i++ {
		key := rand.Intn(100)
		user := "User " + strconv.Itoa(key)
		tr, err = tr.Insert(key, user)
		if err == nil {
			exp = append(exp, key)
		}
	}

	sort.Ints(exp)
	act := tr.list()
	if !reflect.DeepEqual(act, exp) {
		fmt.Printf("Tree elements:       %v\n", act)
		fmt.Printf("Original elements:   %v\n", exp)
		t.Errorf("Tree insert is messed up")
	}
}
