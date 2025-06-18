package utility

import (
	"slices"
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

func Flatten[T any](slice [][]T) []T {
	var result []T
	for _, subSlice := range slice {
		result = append(result, subSlice...)
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
func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
	for _, item := range slice {
		if predicate(item) {
			return item, true
		}
	}

	var zero T
	return zero, false
}

func FindLast[T any](slice []T, predicate func(T) bool) (T, bool) {
	for i := len(slice) - 1; i >= 0; i-- {
		if predicate(slice[i]) {
			return slice[i], true
		}
	}

	var zero T
	return zero, false
}

/*
 * たとえば、Order(注文)に商品(Item)を追加する処理とかを想定する
 *
 * こんな感じ
 * func OrderAddItem(order Order, item Item) (Order, error) { ... }
 * var order, _ = Fold(items, order, OrderAddItem)
 */
func Fold[T any, U any](slice []T, initial U, folder func(U, T) (U, error)) (U, error) {
	if len(slice) == 0 {
		return initial, nil
	}

	var result = initial
	for _, item := range slice {
		var working, err = folder(result, item)
		result = working
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

/*
 * Foldと違って、初期値はゼロ値であることを前提としている。より単純な型に対して使うことを想定している
 *
 * たとえばこんなかんじ
 * var itemSubTotals = []int{100, 200, 300}
 * func Sum(a, b int) int { return a + b }
 * var total = Reduce(itemSubTotals, Sum)
 */
func Reduce[T any](slice []T, reducer func(T, T) T) T {
	var initial T

	if len(slice) == 0 {
		return initial
	}

	result := initial
	for _, item := range slice {
		result = reducer(result, item)
	}
	return result
}

func Keys[T comparable, V any](m map[T]V) []T {
	var keys []T
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func Values[T comparable, V any](m map[T]V) []V {
	var values []V
	for _, value := range m {
		values = append(values, value)
	}
	return values
}

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
func Relate[B any, L any](branches []B, leaves []L, relateIfCase func(B, L) (B, bool)) []B {

	if len(branches) == 0 {
		return branches
	}

	var related = make([]B, 0, len(branches))

	var rest = slices.Clone(leaves)

	for _, branch := range branches {
		var workingBranch = branch
		var matchedIndexes = []int{}

		for index, leaf := range rest {
			var result, hasRelated = relateIfCase(workingBranch, leaf)
			if hasRelated {
				workingBranch = result
				matchedIndexes = append(matchedIndexes, index)
			}
		}

		/*
		 * indexの削除は、後ろから行うことで、削除によるindexのずれを防ぐ
		 * 同一のloopで実現しようとすると、その場合は削除のために逆順でloopを回すので、sort orderが逆転してしまう
		 * そのため、ここでは、削除のために逆順でloopを回る
		 */
		for i := len(matchedIndexes) - 1; i >= 0; i-- {
			var indexValue = matchedIndexes[i]
			rest = slices.Delete(rest, indexValue, indexValue+1)
		}

		related = append(related, workingBranch)
	}

	return related
}

func Intersect[V any, H any](verticals []V, horizontals []H, predicate func(V, H) bool) ([]V, []H, []V, []H) {

	var verticalMatched []V = []V{}
	var horizontalMatched []H = []H{}

	var verticalUnmatched []V = slices.Clone(verticals)
	var horizontalUnmatched []H = slices.Clone(horizontals)

	var vIndex = len(verticalUnmatched)
	for {
		if vIndex == 0 {
			break
		}

		var vertical = verticalUnmatched[vIndex-1]
		for hIndex, horizontal := range horizontalUnmatched {
			if predicate(vertical, horizontal) {
				verticalMatched = append(verticalMatched, vertical)
				horizontalMatched = append(horizontalMatched, horizontal)
				verticalUnmatched = slices.Delete(verticalUnmatched, vIndex-1, vIndex)
				horizontalUnmatched = slices.Delete(horizontalUnmatched, hIndex, hIndex+1)
				break
			}
		}

		vIndex -= 1
	}
	return verticalMatched, horizontalMatched, verticalUnmatched, horizontalUnmatched
}

func Group[T any](slice []T, predicate func(T, T) bool) [][]T {
	grouped := make([][]T, 0)

	var workings = slices.Clone(slice)
	var index = len(workings)
	for {
		if index == 0 {
			break
		}

		var item = workings[index]
		workings = slices.Delete(workings, index, index+1)
		var groupedItems = []T{item}

		for i := len(workings) - 1; i >= 0; i-- {
			var compare = workings[i]
			if predicate(item, compare) {
				groupedItems = append(groupedItems, compare)
				workings = slices.Delete(workings, i, i+1)
			}
		}

		grouped = append(grouped, groupedItems)
		index = len(workings)
	}

	return grouped
}

func Duplicate[T any](slice []T, predicate func(T, T) bool) []T {
	var grouped = Group(slice, predicate)
	var duplicates = Filter(grouped, func(group []T) bool {
		return len(group) > 1
	})
	return Flatten(duplicates)
}
