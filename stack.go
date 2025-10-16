package main

import (
	"fmt"
	"slices"
	"strings"
)

type Stack struct {
	storage []float64
}

func NewStack() *Stack {
	return &Stack{}
}

func (s *Stack) Push(value float64) {
	s.storage = append(s.storage, value)
}

func (s *Stack) Pop() (float64, error) {
	elems, err := s.PopN(1)
	if err != nil {
		return 0.0, err
	}
	return elems[0], err
}

func (s *Stack) PopN(n int) ([]float64, error) {
	if len(s.storage) < n {
		return nil, fmt.Errorf("insufficient stack")
	}
	result := []float64{}
	for range n {
		index := len(s.storage) - 1
		result = append(result, s.storage[index])
		s.storage = s.storage[:index]
	}
	return result, nil
}

func (s *Stack) PopR(n int) ([]float64, error) {
	result, err := s.PopN(n)
	if err != nil {
		return nil, err
	}
	slices.Reverse(result)
	return result, nil
}

func (s *Stack) Swap() error {
	topTwo, err := s.PopN(2)
	if err != nil {
		return err
	}
	for _, n := range topTwo {
		s.Push(n)
	}
	return nil
}

func (s *Stack) Clear() {
	s.storage = s.storage[:0]
}

func (s *Stack) Len() int {
	return len(s.storage)
}

func (s *Stack) Empty() bool {
	return s.Len() <= 0
}

func (s *Stack) String() string {
	stackSize := s.Len()
	var b strings.Builder
	fmt.Fprintf(&b, "[ ")
	for i, n := range s.storage {
		fmt.Fprintf(&b, "%g", n)
		if i < stackSize-1 {
			fmt.Fprintf(&b, "  ")
		}
	}
	fmt.Fprintf(&b, " ]")
	return b.String()
}

func (s *Stack) Sort() {
	slices.Sort(s.storage)
}

func (s *Stack) Copy() []float64 {
	result := make([]float64, s.Len())
	copy(result, s.storage)
	return result
}
