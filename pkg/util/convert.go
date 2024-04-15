package util

import (
	"math/big"
	"reflect"
)

func InterToArray(obj interface{}) []uint {
	if reflect.TypeOf(obj).Kind() == reflect.Slice {
		s := reflect.ValueOf(obj)
		arrays := make([]uint, 0)
		for i := 0; i < s.Len(); i++ {
			ele := s.Index(i)
			e := ele.Interface().(*big.Int)
			arrays = append(arrays, uint(e.Uint64()))
		}
		return arrays
	}
	return make([]uint, 0)
}
