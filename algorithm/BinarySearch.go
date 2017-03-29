package algorithm

import (
	"errors"
)

//compare == 0  值i
//compare < 0 值i 太小
//compare > 0 值i 太大
type BinarySearchCompare func(i int) (compare int, e error)

//使用 二分查找 尋找 [low,high] 中 滿足 compare 的 值
func BinarySearch(low, high int, compareFunc BinarySearchCompare) (int, error) {
	if low > high {
		return 0, errors.New("not found")
	}
	mid := low + (high-low)/2

	compare, e := compareFunc(mid)
	if e != nil {
		return 0, e
	}
	if compare == 0 {
		return mid, nil
	} else if compare < 0 {
		return BinarySearch(mid+1, high, compareFunc)
	}
	return BinarySearch(low, mid-1, compareFunc)
}
