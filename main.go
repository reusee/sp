package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
)

type Any = interface{}

type State func(args ...Any) State

func bind(s State, args ...Any) State {
	return func(_ ...Any) State {
		return s(args...)
	}
}

var (
	pt = fmt.Printf
)

func init() {
	var seed int64
	binary.Read(crand.Reader, binary.LittleEndian, &seed)
	rand.Seed(seed)
}

func main() {

	// sequence
	var s1, s2, s3 State
	s1 = func(args ...Any) State {
		pt("s1\n")
		return bind(s2, args...)
	}
	s2 = func(args ...Any) State {
		pt("s2\n")
		return bind(s3, args...)
	}
	s3 = func(args ...Any) State {
		pt("s3 %+v\n", args)
		return nil
	}
	for s := bind(s1, 1, 2, 3); s != nil; s = s() {
		pt("-- yield --\n")
	}
	fmt.Printf("\n")

	// loop
	s1 = func(args ...Any) State {
		pt("s1 begin of loop\n")
		i := 0
		s2 = func(args ...Any) State {
			if i == 10 {
				return s3
			} else {
				i++
			}
			pt("s2 %d\n", i)
			return s2
		}
		return s2
	}
	s3 = func(args ...Any) State {
		pt("s3 end of loop\n")
		return nil
	}
	for s := s1; s != nil; s = s() {
		pt("-- yield --\n")
	}
	fmt.Printf("\n")

	// continuation
	var p State
	p = func(args ...Any) State {
		fmt.Printf(args[1].(string), args[2:]...)
		if cont, ok := args[0].(State); ok {
			return cont
		}
		return nil
	}
	s1 = func(args ...Any) State {
		pt("s1 begin of loop\n")
		i := 0
		var s2 State // new variable
		s2 = func(args ...Any) State {
			if i == 10 {
				return s3
			} else {
				i++
			}
			return bind(p, s2, "s2 %d\n", i)
		}
		return s2
	}
	s3 = func(args ...Any) State {
		return bind(p, nil, "s3 end of loop\n")
	}
	for s := s1; s != nil; s = s() {
		pt("-- yield --\n")
	}
	fmt.Printf("\n")

	// schedule
	var states []State
	for i := 0; i < 10; i++ {
		states = append(states, s1)
	}
	for len(states) > 0 {
		i := rand.Intn(len(states))
		if s := states[i](); s != nil {
			states = append(states, s)
		}
		states = append(states[:i], states[i+1:]...)
	}

}
