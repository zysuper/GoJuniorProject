package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type BaseEntity struct {
	Id         int64
	CreateTime time.Time
	UpdateTime time.Time
}

type User struct {
	BaseEntity
	name string
}

func NewUserByName(name string) User {
	return User{
		name: name,
	}
}

func NewUserById(id int64) User {
	return User{
		BaseEntity: BaseEntity{
			Id: id,
		},
	}
}

// 插入 T 到数据库
func Insert[T BaseEntity](t T) {

}

//func TestInsert(t *testing.T) {
//	Insert[User](User{})
//}

type Stream[T any] struct {
}

//func (s *Stream[T]) Map[E any](func(t T) E) *Stream[E] {
//
//}

func (s *Stream[T]) Filter() {

}

//type Orm interface {
//	Select[T any]() (*T, error)
//}
//
//func TestUseOrm(t *testing.T) {
//	var o Orm
//	user, err := o.Select[User]()
//	order, err := o.Select[Order]()
//}

type Selector[T any] struct {
}

func (s *Selector[T]) Get() (*T, error) {
	return new(T), nil
}

func TestUseSelector(t *testing.T) {
	s := &Selector[User]{}
	user, err := s.Get()
	assert.NoError(t, err)
	t.Log(user)
}

//func NewA() {
//
//}

//public class A {
//	public A() {
//
//	}
//}
