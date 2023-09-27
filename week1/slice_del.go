package week1

import "fmt"

// []int 删除特定下标元素.
func intSliceDel(slice []int, index int) ([]int, error) {
	length := len(slice)

	// 异常参数检测.
	if index < 0 || index >= length {
		return nil, fmt.Errorf("%d out of range", index)
	}

	if index == 0 {
		return slice[1:], nil
	}

	if index == length-1 {
		return slice[:index], nil
	}

	// 预估新 slice 大小为 length - 1
	// 这种情况下实现缩容，其他两种情况不做缩容.
	newSlice := make([]int, length-1)

	copy(newSlice, slice[:index])
	copy(newSlice[index:], slice[index+1:])

	return newSlice, nil
}
