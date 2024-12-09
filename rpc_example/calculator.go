package main

import "errors"

type Calculator struct{}

func (c *Calculator) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}

func (c *Calculator) Divide(args *Args, reply *float64) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	*reply = float64(args.A) / float64(args.B)
	return nil
}

type Args struct {
	A, B int
}
