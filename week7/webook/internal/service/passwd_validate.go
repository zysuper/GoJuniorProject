package service

import (
	"crypto/md5"
	"encoding/hex"
	bcrypt "golang.org/x/crypto/bcrypt"
	"log"
)

// PasswordValidateService 密码校验接口抽象.
type PasswordValidateService interface {
	// Hash 获取密码 Hash 值.
	Hash(password string) ([]byte, error)
	// ComparePassword 比较密码
	ComparePassword(hashedPasswd, passwd string) error
}

type DefaultPvs struct {
}

func (d *DefaultPvs) Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (d *DefaultPvs) ComparePassword(hashedPasswd, passwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPasswd), []byte(passwd))
}

type Md5Pvs struct {
}

func (m *Md5Pvs) Hash(password string) ([]byte, error) {
	ret := md5.Sum([]byte(password))
	var r = make([]byte, 32)
	hex.Encode(r, ret[:])
	return r, nil
}

func (m *Md5Pvs) ComparePassword(hashedPasswd, passwd string) error {
	r, _ := m.Hash(passwd)
	log.Println(string(r), passwd)
	if hashedPasswd != string(r) {
		return ErrInvalidUserOrPassword
	}
	return nil
}

func NewPasswordValidator() PasswordValidateService {
	return &DefaultPvs{}
}
