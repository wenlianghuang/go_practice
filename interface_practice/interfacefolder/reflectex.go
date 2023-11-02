package interfacefolder

import (
	"errors"
	"fmt"
	"reflect"
)

type Savings interface {
	Deposit(amount float64) error
	Withdraw(amount float64) error
}

type Account struct {
	id      string
	name    string
	balance float64
}

func (ac *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("必須存入正數")
	}
	ac.balance += amount
	return nil
}

func (ac *Account) Withdraw(amount float64) error {
	if amount > ac.balance {
		return errors.New("餘額不足")
	}
	ac.balance -= amount
	return nil
}

func Reflectex() {
	var savings Savings = &Account{"X123", "Justin Lin", 1000}
	t := reflect.TypeOf(savings)

	for i, n := 0, t.NumMethod(); i < n; i++ {
		f := t.Method(i)
		fmt.Println(f.Name, f.Type)
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fmt.Println(t.Kind())
	fmt.Println(t.String())
	for i, n := 0, t.NumField(); i < n; i++ {
		f := t.Field(i)
		fmt.Println(f.Name, f.Type)
	}
}
