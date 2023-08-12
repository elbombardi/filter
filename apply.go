// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package filter contains utility functions for filtering slices through the
// distributed application of a filter function.
//
// The package is an experiment to see how easy it is to write such things
// in Go. It is easy, but for loops are just as easy and more efficient.
//
// You should not use this package.
package filter // import "robpike.io/filter"

import (
	"reflect"
)

// Apply takes a slice of type []T and a function of type func(T) T. (If the
// input conditions are not satisfied, Apply panics.) It returns a newly
// allocated slice where each element is the result of calling the function on
// successive elements of the slice.
func Apply[T any, R any](slice []T, function func(T) R) []R {
	return apply(slice, function, false)
}

// ApplyInPlace is like Apply, but overwrites the slice rather than returning a
// newly allocated slice.
func ApplyInPlace[T any, R any](slice []T, function func(T) R) {
	apply(slice, function, true)
}

// Choose takes a slice of type []T and a function of type func(T) bool. (If
// the input conditions are not satisfied, Choose panics.) It returns a newly
// allocated slice containing only those elements of the input slice that
// satisfy the function.
func Choose[T any](slice []T, function func(T) bool) []T {
	out, _ := chooseOrDrop(slice, function, false, true)
	return out
}

// Drop takes a slice of type []T and a function of type func(T) bool. (If the
// input conditions are not satisfied, Drop panics.) It returns a newly
// allocated slice containing only those elements of the input slice that do
// not satisfy the function, that is, it removes elements that satisfy the
// function.
func Drop[T any](slice []T, function func(T) bool) []T {
	out, _ := chooseOrDrop(slice, function, false, false)
	return out
}

// ChooseInPlace is like Choose, but overwrites the slice rather than returning
// a newly allocated slice. Since ChooseInPlace must modify the header of the
// slice to set the new length, it takes as argument a pointer to a slice
// rather than a slice.
func ChooseInPlace[T any](pointerToSlice *[]T, function func(T) bool) {
	chooseOrDropInPlace(pointerToSlice, function, true)
}

// DropInPlace is like Drop, but overwrites the slice rather than returning a
// newly allocated slice. Since DropInPlace must modify the header of the slice
// to set the new length, it takes as argument a pointer to a slice rather than
// a slice.
func DropInPlace[T any](pointerToSlice *[]T, function func(T) bool) {
	chooseOrDropInPlace(pointerToSlice, function, false)
}

func apply[T any, R any](slice []T, function func(T) R, inPlace bool) []R {
	var out []R
	intype := reflect.TypeOf(slice)
	outtype := reflect.TypeOf(out)
	if inPlace && intype == outtype {
		out = reflect.ValueOf(slice).Interface().([]R)
	} else {
		out = make([]R, len(slice))
	}
	for i, s := range slice {
		out[i] = function(s)
	}
	return out
}

func chooseOrDropInPlace[T any](slice *[]T, function func(T) bool, truth bool) {
	inp := reflect.ValueOf(slice)
	if inp.Kind() != reflect.Ptr {
		panic("choose/drop: not pointer to slice")
	}
	_, n := chooseOrDrop(*slice, function, true, truth)
	inp.Elem().SetLen(n)
}

var boolType = reflect.ValueOf(true).Type()

func chooseOrDrop[T any](slice []T, function func(T) bool, inPlace, truth bool) ([]T, int) {
	var r []T
	if inPlace {
		r = slice[:0]
	}
	for _, s := range slice {
		if function(s) == truth {
			r = append(r, s)
		}
	}
	return r, len(r)
}
