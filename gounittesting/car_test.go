package gounittesting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Simple testing what different between Fatal and Error
func TestNew(t *testing.T) {
	c, err := New("Matt", 100)
	if err != nil {
		t.Fatal("got errors:", err)
	}

	if c == nil {
		t.Error("car should be nil")
	}
}

// Simple testing with testify tool
func TestNewWithAssert(t *testing.T) {
	c, err := New("", 100)
	assert.NotNil(t, err)
	assert.Nil(t, c)
}

func TestCar_SetName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		c    *Car
		args args
		want string
	}{
		{
			name: "no input name",
			c: &Car{
				Name:  "foo",
				Price: 100,
			},
			args: args{
				name: "",
			},
			want: "foo",
		},
		{
			name: "input name",
			c: &Car{
				Name:  "foo",
				Price: 100,
			},
			args: args{
				name: "foo",
			},
			want: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.SetName(tt.args.name); got != tt.want {
				t.Errorf("Car.SetName() = %v, want %v", got, tt.want)
			}
		})
	}
}
