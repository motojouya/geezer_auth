package utility

import (
	"reflect"
	"errors"
)

func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func Map[T any, U any](slice []T, mapper func(T) U) []U {
	var result []U
	for _, item := range slice {
		result = append(result, mapper(item))
	}
	return result
}

func Every[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if !predicate(item) {
			return false
		}
	}
	return true
}

func Some[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}

// 使い方は、`internal/authorization/authorization#GetPriorityRolePermission`を参照。叙述関数は高階関数になる
func Find[T any](slice []T, predicate func(T) bool) *T {
	for _, item := range slice {
		if predicate(item) {
			return &item
		}
	}
	return nil
}

func FindLast[T any](slice []T, predicate func(T) bool) *T {
	for i := len(slice) - 1; i >= 0; i-- {
		if predicate(slice[i]) {
			return &slice[i]
		}
	}
	return nil
}

func Keys[T comparable, V any](m map[T]V) []T {
	var keys []T
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func Values[T any, V any](m map[T]V) []V {
	var values []V
	for _, value := range m {
		values = append(values, value)
	}
	return values
}

// TODO Findあるから不要な気がしてきた。ただ長いlistを繰り返し探索するときにはmapに変換しておくと高速化できる
func ToMap[T any, K comparable](slice []T, getKey func(T) K) map[K]T {
	result := make(map[K]T)
	for _, item := range slice {
		var key = getKey(item)
		result[key] = item
	}
	return result
}

func Entries[T comparable, V any](m map[T]V) []struct {
	Key   T
	Value V
} {
	var entries []struct {
		Key   T
		Value V
	}
	for key, value := range m {
		entries = append(entries, struct {
			Key   T
			Value V
		}{Key: key, Value: value})
	}
	return entries
}

/*
 * branch側のリストの要素に対して、leaves側の要素のリストを紐づける
 * 
 * ORMapperで用いられるRelationの機能の代替するためのもの
 * branchは、[]leavesの要素を持っているが、DBでqueryを投げる際には、branch,leafで別々に投げたい。
 * 別々に投げた後に紐づけを行うための関数
 */
func Relate[B any, L any](property string, branches []*P, leaves []*L, predicate func(B, L) bool) []P {

	if len(branches) == 0 {
		return branches
	}

	var workings = make([]*L, 0, len(leaves))

	for _, branch := range branches {

		var e = reflect.ValueOf(&branch).Elem()
		var list = e.FieldByName(property)
		if !list.IsValid() || list.Kind() != reflect.Slice {
			panic("property must be a valid slice field")
		}

		for i, leaf := range workings {
			if relate(branch, leaf) {
				var item = []L{leaf}
				list.Set(reflect.AppendSlice(list, reflect.ValueOf(item)))

				workings = append(workings[:i], workings[i+1:]...)
			}
		}
	}

	return branches
}
