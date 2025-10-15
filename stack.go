package main

import (
	"fmt"
	"slices"
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

func (s *Stack) Pop(n int) ([]float64, error) {
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
	result, err := s.Pop(n)
	if err != nil {
		return nil, err
	}
	slices.Reverse(result)
	return result, nil
}

func (s *Stack) Swap() error {
	topTwo, err := s.Pop(2)
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
