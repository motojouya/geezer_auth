package utility_test

import (
	"github.com/motojouya/geezer_auth/internal/utility"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter(t *testing.T) {
	var list = []string{"this", "test", "item"}
	var predicate = func(item string) bool {
		var chars = []rune(item)
		if len(chars) > 0 && chars[0] == 't' {
			return true
		}
		return false
	}

	var tList = utility.Filter(list, predicate)

	assert.Equal(t, 2, len(tList))
	assert.Equal(t, "this", tList[0])
	assert.Equal(t, "test", tList[1])

	t.Logf("filtered list: %v", tList)
}

func TestMap(t *testing.T) {
	var list = []string{"this", "test", "item"}
	var mapper = func(item string) string {
		return item + "_mapped"
	}

	var tList = utility.Map(list, mapper)

	assert.Equal(t, 3, len(tList))
	assert.Equal(t, "this_mapped", tList[0])
	assert.Equal(t, "test_mapped", tList[1])
	assert.Equal(t, "item_mapped", tList[2])

	t.Logf("mapped list: %v", tList)
}

func TestFold(t *testing.T) {
	var list = []string{"this", "test", "item"}
	var folder = func(accumulator string, item string) (string, error) {
		return accumulator + "_" + item, nil
	}

	var result, err = utility.Fold(list, folder, "first")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, "first_this_test_item", result)

	t.Logf("reduced result: %s", result)
}

func TestFoldError(t *testing.T) {
	var list = []string{"this", "test", "item"}
	var folder = func(accumulator string, item string) (string, error) {
		return "", errors.New("test error")
	}

	var result, err = utility.Fold(list, folder, "first")
	if err == nil {
		t.Fatal("expected error, but got nil")
	}
}

func TestReduce(t *testing.T) {
	var list = []uint{1, 2, 3}
	var reducer = func(accumulator uint, item uint) uint {
		return accumulator + item
	}

	var result = utility.Reduce(list, reducer)

	assert.Equal(t, 6, result)

	t.Logf("reduced result: %d", result)
}

func TestSome(t *testing.T) {
	var list1 = []string{"this", "test", "item"}
	var list2 = []string{"this", "tast", "item"}
	var predicate = func(item string) bool {
		return item == "test"
	}

	var contains1 = utility.Some(list1, predicate)
	var contains2 = utility.Some(list2, predicate)

	assert.True(t, contains1)
	assert.False(t, contains2)
}

func TestEvery(t *testing.T) {
	var list1 = []string{"this", "test", "item"}
	var list2 = []string{"this", "test", "temi"}
	var predicate = func(item string) bool {
		var chars = []rune(item)
		return chars[0] == 't'
	}

	var allMatch1 = utility.Every(list1, predicate)
	var allMatch2 = utility.Every(list2, predicate)

	assert.True(t, allMatch1)
	assert.False(t, allMatch2)
}

func TestFind(t *testing.T) {
	var list = []string{"this", "test", "item"}
	var predicate = func(item string) bool {
		var chars = []rune(item)
		return chars == 't'
	}

	var foundItem = utility.Find(list, predicate)

	assert.Equal(t, "this", foundItem)

	t.Logf("found item: %s", foundItem)
}

func TestFindLast(t *testing.T) {
	var list = []string{"this", "test", "item"}
	var predicate = func(item string) bool {
		var chars = []rune(item)
		return chars == 't'
	}

	var foundItem = utility.FindLast(list, predicate)

	assert.Equal(t, "test", foundItem)

	t.Logf("found item: %s", foundItem)
}

func TestKeys(t *testing.T) {
	var m = map[string]int{
		"this": 1,
		"test": 2,
		"item": 3,
	}

	var keys = utility.Keys(m)

	assert.Equal(t, 3, len(keys))
	assert.Equal(t, keys[0], "this")
	assert.Equal(t, keys[1], "test")
	assert.Equal(t, keys[2], "item")

	t.Logf("keys: %v", keys)
}

func TestValues(t *testing.T) {
	var m = map[string]int{
		"this": 1,
		"test": 2,
		"item": 3,
	}

	var values = utility.Values(m)

	assert.Equal(t, 3, len(values))
	assert.Equal(t, values[0], 1)
	assert.Equal(t, values[1], 2)
	assert.Equal(t, values[2], 3)

	t.Logf("values: %v", values)
}

func TestEntries(t *testing.T) {
	var m = map[string]int{
		"this": 1,
		"test": 2,
		"item": 3,
	}

	var entries = utility.Entries(m)

	assert.Equal(t, 3, len(entries))
	assert.Equal(t, entries[0].Key, "this")
	assert.Equal(t, entries[0].Value, 1)
	assert.Equal(t, entries[1].Key, "test")
	assert.Equal(t, entries[1].Value, 2)
	assert.Equal(t, entries[2].Key, "item")
	assert.Equal(t, entries[2].Value, 3)

	t.Logf("entries: %v", entries)
}

func TestToMap(t *testing.T) {
	var list = []string{"this", "test", "item"}
	var getKey = func(item string) string {
		var chars = []rune(item)
		return string(chars[:2])
	}

	var result = utility.ToMap(list, getKey)

	assert.Equal(t, 3, len(result))
	assert.Equal(t, result["th"], "this")
	assert.Equal(t, result["te"], "test")
	assert.Equal(t, result["it"], "item")

	t.Logf("map: %v", result)
}

func TestFlatten(t *testing.T) {
	var nestedList = [][]string{
		{"this", "test"},
		{"item", "example"},
	}

	var flattened = utility.Flatten(nestedList)

	assert.Equal(t, 4, len(flattened))
	assert.Equal(t, flattened[0], "this")
	assert.Equal(t, flattened[1], "test")
	assert.Equal(t, flattened[2], "item")
	assert.Equal(t, flattened[3], "example")

	t.Logf("flattened list: %v", flattened)
}

type ItemRequest struct {
	ID       string
	Quantity uint
}

type Item struct {
	ID       string
	OrderID  string
	Name     string
	Quantity uint
}

type Order struct {
	ID       string
	Customer string
	Items    []Item
}

func TestRelated(t *testing.T) {
	var orderList = []Order{
		{ID: "1", Customer: "Alice", Items: []Item{}},
		{ID: "2", Customer: "Bob", Items: []Item{}},
	}
	var itemList = []Item{
		{ID: "a1", OrderID: "1", Name: "Apple", Quantity: 2},
		{ID: "b1", OrderID: "1", Name: "Banana", Quantity: 3},
		{ID: "c1", OrderID: "2", Name: "Carrot", Quantity: 5},
	}
	var predicate = func(order Order, item Item) bool {
		return order.ID == item.OrderID
	}

	var related = utility.Related("Items", orderList, itemList, predicate)

	assert.Equal(t, 2, len(related))
	assert.Equal(t, related[0].Customer, "Alice")
	assert.Equal(t, 2, len(related[0].Items))
	assert.Equal(t, "Apple", related[0].Items[0].Name)
	assert.Equal(t, "Banana", related[0].Items[1].Name)
	assert.Equal(t, related[1].Customer, "Bob")
	assert.Equal(t, 1, len(related[1].Items))
	assert.Equal(t, "Carrot", related[1].Items[1].Name)

	t.Logf("related items: %v", related)
}

func TestIntersect(t *testing.T) {
	var list1 = []ItemRequest{
		{ID: "a1", Quantity: 2},
		{ID: "b1", Quantity: 3},
		{ID: "d1", Quantity: 4},
	}
	var list2 = []Item{
		{ID: "a1", OrderID: "1", Name: "Apple", Quantity: 1},
		{ID: "b1", OrderID: "1", Name: "Banana", Quantity: 1},
		{ID: "c1", OrderID: "2", Name: "Carrot", Quantity: 1},
	}
	var predicate = func(itemRequest ItemRequest, item Item) bool {
		return itemRequest.ID == item.ID
	}

	var varticalMatched, horizontalMatched, varticalUnMatched, horizontalUnMatched = utility.Intersect(list1, list2, predicate)

	assert.Equal(t, 2, len(varticalMatched))
	assert.Equal(t, "Apple", horizontalMatched[0].Name)
	assert.Equal(t, "Banana", horizontalMatched[1].Name)
	assert.Equal(t, 2, len(horizontalMatched))
	assert.Equal(t, 2, horizontalMatched[0].Quantity)
	assert.Equal(t, 3, horizontalMatched[1].Quantity)
	assert.Equal(t, 1, len(varticalUnMatched))
	assert.Equal(t, "Carrot", varticalUnMatched[0].Name)
	assert.Equal(t, 1, len(horizontalUnMatched))
	assert.Equal(t, 4, horizontalMatched[0].Quantity)

	t.Logf("varticalMatched: %v", varticalMatched)
	t.Logf("horizontalMatched: %v", horizontalMatched)
	t.Logf("varticalUnMatched: %v", varticalUnMatched)
	t.Logf("horizontalUnMatched: %v", horizontalUnMatched)
}

func TestGroup(t *testing.T) {
	var list = []Item{
		{ID: "a1", OrderID: "1", Name: "Apple", Quantity: 2},
		{ID: "b1", OrderID: "1", Name: "Banana", Quantity: 3},
		{ID: "c1", OrderID: "2", Name: "Carrot", Quantity: 2},
	}

	var grouppper = func(item1 Item, item2 Item) string {
		return item1.Quantity == item2.Quantity
	}

	var grouped = utility.Group(list, grouppper)

	assert.Equal(t, 2, len(grouped))
	assert.Equal(t, 2, len(grouped[0]))
	assert.Equal(t, 1, len(grouped[1]))
	assert.Equal(t, grouped[0][0].Name, "Apple")
	assert.Equal(t, grouped[0][1].Name, "Carrot")
	assert.Equal(t, grouped[1][0].Name, "Banana")

	t.Logf("grouped items: %v", grouped)
}

func TestDuplicates(t *testing.T) {
	var list = []string{"this", "test", "item", "test", "this"}
	var predicate = func(item1 string, item2 string) bool {
		return item1 == item2
	}

	var duplicates = utility.Duplicates(list, predicate)

	assert.Equal(t, 4, len(duplicates))
	assert.Equal(t, duplicates[0], "this")
	assert.Equal(t, duplicates[1], "this")
	assert.Equal(t, duplicates[2], "test")
	assert.Equal(t, duplicates[3], "test")

	t.Logf("duplicates: %v", duplicates)
}
