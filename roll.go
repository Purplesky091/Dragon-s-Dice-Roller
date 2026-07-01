package main

import "strconv"

type Roll struct {
	value   int
	dropped bool
}

func (roll Roll) String() string {
	return strconv.Itoa(roll.value)
}
