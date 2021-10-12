package utils

import (
	"log"
	"math"
	"strconv"
)

// DefaultAllocate ...
const DefaultAllocate = 0

// StrToInt casts string to int
func StrToInt(source string) int {
	res, err := strconv.ParseInt(source, 10, 32)
	if err != nil {
		log.Printf("error in StrToInt: %v", err)
		return DefaultAllocate
	}
	return int(res)
}

// StrToInt32 casts string to int32
func StrToInt32(source string) int32 {
	res, err := strconv.ParseInt(source, 10, 32)
	if err != nil {
		log.Printf("error in StrToInt32: %v", err)
		return DefaultAllocate
	}
	return int32(res)
}

// StrToInt64 casts string to int64
func StrToInt64(source string) int64 {
	res, err := strconv.ParseInt(source, 10, 64)
	if err != nil {
		log.Printf("error in StrToInt64: %v", err)
		return DefaultAllocate
	}
	return res
}

// StrToUint casts string to uint
func StrToUint(source string) uint {
	res, err := strconv.ParseUint(source, 10, 32)
	if err != nil {
		log.Printf("error in StrToUint: %v", err)
		return DefaultAllocate
	}

	if res < 0 && res > math.MaxUint32 {
		log.Printf("error in StrToUint: %v", err)
		return DefaultAllocate
	}

	return uint(res)
}
