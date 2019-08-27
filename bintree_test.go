package bintree

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"testing"
)

type UserData struct {
	name  string
	email string
}

func generateSlice(n int) []int {
	s := make([]int, 0, n)
	for i := 0; i < n; i++ {
		s = append(s, rand.Intn(100))
	}
	return s
}

func createTree(s []int) (*node, []int) {
	var keys []int
	var err error
	var tr *node
	var user UserData

	for i := 0; i < len(s); i++ {
		user.name = "User " + strconv.Itoa(s[i])
		user.email = "user" + strconv.Itoa(s[i]) + "@semrush.com"
		tr, err = tr.Insert(s[i], user)
		if err == nil {
			keys = append(keys, s[i])
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
}

func TestGet(t *testing.T) {
	s := generateSlice(100)
	tr, exp := createTree(s)
	key := exp[rand.Intn(len(exp))]
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
	key := exp[rand.Intn(len(exp))]
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

func TestDeleteRoot(t *testing.T) {
	s := generateSlice(1)
	tr, exp := createTree(s)
	key := exp[rand.Intn(len(exp))]
	tr, err := tr.Delete(key)
	if err != nil {
		t.Errorf(err.Error())
	}
	if tr != nil {
		t.Errorf("The tree is not empty")
	}
}

func TestDeleteRightLeaf(t *testing.T) {
	s := []int{1, 2}
	tr, _ := createTree(s)
	tr, err := tr.Delete(2)
	if err != nil {
		t.Errorf(err.Error())
	}
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
	tr, err := tr.Delete(1)
	if err != nil {
		t.Errorf(err.Error())
	}
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
	tr, err := tr.Delete(3)
	if err != nil {
		t.Errorf(err.Error())
	}
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
	tr, err := tr.Delete(2)
	if err != nil {
		t.Errorf(err.Error())
	}
	exp := []int{1, 3, 4}
	act := tr.List()
	if !reflect.DeepEqual(act, exp) {
		fmt.Printf("Left elements: %v, expected: %v\n", act, exp)
		t.Errorf("Error deleting node")
	}
}

func TestDeleteRightNodeWithTwoChildren(t *testing.T) {
	s := []int{2, 1, 5, 3, 6, 4}
	tr, _ := createTree(s)
	tr, err := tr.Delete(5)
	if err != nil {
		t.Errorf(err.Error())
	}
	exp := []int{1, 2, 3, 4, 6}
	act := tr.List()
	if !reflect.DeepEqual(act, exp) {
		fmt.Printf("Left elements: %v, expected: %v\n", act, exp)
		t.Errorf("Error deleting node")
	}
}

func TestDeleteLeftNodeWithTwoChildren(t *testing.T) {
	s := []int{5, 2, 1, 4, 3, 6}
	tr, _ := createTree(s)
	tr, err := tr.Delete(2)
	if err != nil {
		t.Errorf(err.Error())
	}
	exp := []int{1, 3, 4, 5, 6}
	act := tr.List()
	if !reflect.DeepEqual(act, exp) {
		fmt.Printf("Left elements: %v, expected: %v\n", act, exp)
		t.Errorf("Error deleting node")
	}
}

func TestMinMax(t *testing.T) {
	s := generateSlice(100)
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
