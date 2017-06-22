package main

import (
	"fmt"
	"github.com/dgtony/gclient"
)

func main() {
	// create new client
	gc := gclient.NewClient("http://localhost:8080/cache/v1")

	// try to get a key
	testKey1 := "non_exist_key"
	val, ok, err := gc.Get(testKey1, 0)
	if err != nil {
		fmt.Printf("key retrieval error: %s\n", err)
	} else if ok {
		fmt.Printf("value found: %s\n", val)
	} else {
		fmt.Printf("no value could be found for the key %s\n", testKey1)
	}

	// set new string value
	testKey2 := "testkey2"
	_ = gc.Set(testKey2, "testvalue-string", 60, 0)

	// set new dict
	testKey3 := "testkey3"
	_ = gc.Set(testKey3, map[string]int{"a": 1, "b": 2, "c": 3}, 60, 0)

	// set new array
	testKey4 := "testkey4"
	_ = gc.Set(testKey4, []float64{1.01, 2.02, 3.03}, 60, 0)

	// get string value
	val, _, _ = gc.Get(testKey2, 0)
	fmt.Printf("value found: %s\n", val)

	// get value with subkey
	val, _, _ = gc.GetSubKey(testKey3, "b", 0)
	fmt.Printf("subvalue found: %v\n", val)

	// get value with subindex
	// NB: subindexing starts from 1!
	val, _, _ = gc.GetSubIndex(testKey4, 2, 0)
	fmt.Printf("subvalue found: %v\n", val)

	// remove item with key
	_ = gc.Remove(testKey2, 0)

	// get all stored keys
	keys, _ := gc.Keys(1)
	fmt.Printf("all stored keys: %v\n", keys)

	// get stored keys with mask
	mask := "*[23]"
	keys, _ = gc.KeysMask(mask, 1)
	fmt.Printf("stored keys (mask: %s): %v\n", mask, keys)
}
