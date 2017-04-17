package algorithm

import (
	"errors"
)

//使用 二分查找 尋找 [low,high] 中 滿足 compare 的 值
//compare == 0  值i
//compare < 0 值i 太小
//compare > 0 值i 太大
func BinarySearch(low, high int,
	compareFunc func(i int) (compare int, e error)) (int, error) {
	for low <= high {
		mid := low + (high-low)/2
		compare, e := compareFunc(mid)
		if e != nil {
			return 0, e
		}

		if compare == 0 {
			return mid, nil
		} else if compare < 0 {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return 0, errors.New("not found")
}
