package bintree

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"testing"
	"time"
)

type UserData struct {
	name  string
	email string
}

func randInt(n int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(n)
}

func generateSlice(n int) []int {
	s := make([]int, 0, n)
	for i := 0; i < n; i++ {
		s = append(s, randInt(100))
	}
	return s
}

func removeFromSlice(s []int, value int) []int {
	for i, v := range s {
		if v == value {
			s[i] = s[len(s)-1]
			s = s[:len(s)-1]
		}
	}
	return s
}

func createTree(s []int) (*node, []int) {
	var keys []int
	var err error
	var tr *node
	var user UserData

	for _, v := range s {
		user.name = "User " + strconv.Itoa(v)
		user.email = "user" + strconv.Itoa(v) + "@bintree.com"
		tr, err = tr.Insert(v, user)
		if err == nil {
			keys = append(keys, v)
		}
	}
	return tr, keys
}

func MinMax(s []int) (int, int, error) {
	if len(s) == 0 {
		return 0, 0, errors.New("the slice is empty")
	}
	min := s[0]
	max := s[0]
	for _, value := range s {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max, nil
}

func TestInsert(t *testing.T) {
	s := generateSlice(100)
	tr, exp := createTree(s)
	sort.Ints(exp)
	act := tr.List()
	if !reflect.DeepEqual(act, exp) {
		fmt.Printf("Tree elements:       %v\n", act)
		fmt.Printf("Original elements:   %v\n", exp)
		t.Errorf("Tree insert is messed up")
	}
	bf := tr.ListBalanceFactors()
	if min, max, _ := MinMax(bf); min <= -2 || max >= 2 {
		t.Errorf("Tree is not balanced")
	}
}

func TestGet(t *testing.T) {
	s := generateSlice(100)
	tr, exp := createTree(s)
	key := exp[randInt(len(exp))]
	n, err := tr.Get(key)
	if err != nil {
		t.Errorf(err.Error())
	}
	expectedName := "User " + strconv.Itoa(key)
	if n.(UserData).name != expectedName {
		fmt.Printf("Found user: %v, expected: %v\n", n.(UserData).name, expectedName)
		t.Errorf("Error finding user")
	}
}

func TestGetInvalidKey(t *testing.T) {
	s := generateSlice(100)
	tr, _ := createTree(s)
	_, err := tr.Get(-1)
	if err == nil {
		t.Errorf("Found nonexistent user")
	}
}

func TestUpdate(t *testing.T) {
	var u Updater
	var updatedUser UserData

	s := generateSlice(100)
	tr, exp := createTree(s)
	key := exp[randInt(len(exp))]
	updatedUser.name = "updated name"
	updatedUser.email = "updated email"
	u = func(data interface{}) (interface{}, error) {
		d := data.(UserData)
		d.name = updatedUser.name
		d.email = updatedUser.email
		return d, nil
	}
	err := tr.Update(key, u)
	if err != nil {
		t.Errorf(err.Error())
	}
	updatedData, err := tr.Get(key)
	if updatedData.(UserData).name != updatedUser.name {
		fmt.Printf("Updated name: %v, expected: %v\n", updatedData.(UserData).name, updatedUser.name)
		t.Errorf("Error updating user name")
	}
	if updatedData.(UserData).email != updatedUser.email {
		fmt.Printf("Updated name: %v, expected: %v\n", updatedData.(UserData).email, updatedUser.email)
		t.Errorf("Error updating user email")
	}
}

func TestDeleteLoneRoot(t *testing.T) {
	s := []int{1}
	tr, _ := createTree(s)
	tr = tr.Delete(1)
	if tr != nil {
		t.Errorf("The tree is not empty")
	}
}

func TestDeleteRoot(t *testing.T) {
	s := []int{2, 1, 3}
	tr, _ := createTree(s)
	tr = tr.Delete(2)
	exp := []int{1, 3}
	act := tr.List()
	if !reflect.DeepEqual(act, exp) {
		fmt.Printf("Left elements: %v, expected: %v\n", act, exp)
		t.Errorf("Error deleting leaf")
	}
}

func TestDeleteRightLeaf(t *testing.T) {
	s := []int{1, 2}
	tr, _ := createTree(s)
	tr = tr.Delete(2)
	exp := []int{1}
	act := tr.List()
	if !reflect.DeepEqual(act, exp) {
		fmt.Printf("Left elements: %v, expected: %v\n", act, exp)
		t.Errorf("Error deleting leaf")
	}
}

func TestDeleteLeftLeaf(t *testing.T) {
	s := []int{2, 1}
	tr, _ := createTree(s)
	tr = tr.Delete(1)
	exp := []int{2}
	act := tr.List()
	if !reflect.DeepEqual(act, exp) {
		fmt.Printf("Left elements: %v, expected: %v\n", act, exp)
		t.Errorf("Error deleting leaf")
	}
}

func TestDeleteRightNodeWithOneChild(t *testing.T) {
	s := []int{2, 1, 3, 4}
	tr, _ := createTree(s)
	tr = tr.Delete(3)
	exp := []int{1, 2, 4}
	act := tr.List()
	if !reflect.DeepEqual(act, exp) {
		fmt.Printf("Left elements: %v, expected: %v\n", act, exp)
		t.Errorf("Error deleting node")
	}
}

func TestDeleteLeftNodeWithOneChild(t *testing.T) {
	s := []int{3, 2, 1, 4}
	tr, _ := createTree(s)
	tr = tr.Delete(2)
	exp := []int{1, 3, 4}
	act := tr.List()
	if !reflect.DeepEqual(act, exp) {
		fmt.Printf("Left elements: %v, expected: %v\n", act, exp)
		t.Errorf("Error deleting node")
	}
}

func TestDeleteRightNodeWithTwoChildren(t *testing.T) {
	s := []int{3, 2, 7, 1, 8, 5, 6, 4}
	tr, _ := createTree(s)
	tr = tr.Delete(7)
	exp := []int{1, 2, 3, 4, 5, 6, 8}
	act := tr.List()
	if !reflect.DeepEqual(act, exp) {
		fmt.Printf("Left elements: %v, expected: %v\n", act, exp)
		t.Errorf("Error deleting node")
	}
}

func TestDeleteLeftNodeWithTwoChildren(t *testing.T) {
	s := []int{6, 7, 4, 8, 5, 2, 1, 3}
	tr, _ := createTree(s)
	tr = tr.Delete(4)
	exp := []int{1, 2, 3, 5, 6, 7, 8}
	act := tr.List()
	if !reflect.DeepEqual(act, exp) {
		fmt.Printf("Left elements: %v, expected: %v\n", act, exp)
		t.Errorf("Error deleting node")
	}
}

func TestDeleteRandomNode(t *testing.T) {
	s := generateSlice(100)
	tr, exp := createTree(s)
	key := exp[randInt(len(exp))]
	exp = removeFromSlice(exp, key)
	sort.Ints(exp)
	tr = tr.Delete(key)
	act := tr.List()
	if !reflect.DeepEqual(act, exp) {
		fmt.Printf("Tree elements:       %v\n", act)
		fmt.Printf("Original elements:   %v\n", exp)
		t.Errorf("Tree deletion is messed up")
	}
	bf := tr.ListBalanceFactors()
	if min, max, _ := MinMax(bf); min <= -2 || max >= 2 {
		t.Errorf("Tree is not balanced")
	}
}

func TestMinMax(t *testing.T) {
	s := generateSlice(10)
	tr, exp := createTree(s)
	minAct := tr.Min()
	maxAct := tr.Max()
	minExp, maxExp, err := MinMax(exp)
	if err != nil {
		t.Errorf(err.Error())
	}
	if minAct != minExp {
		fmt.Printf("Tree minimum: %v, expected: %v\n", minAct, minExp)
		t.Errorf("Tree Min function needs more work")
	}
	if maxAct != maxExp {
		fmt.Printf("Tree maximum: %v, expected: %v\n", maxAct, maxExp)
		t.Errorf("Tree Min function needs more work")
	}
}
