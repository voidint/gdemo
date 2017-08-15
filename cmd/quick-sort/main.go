package main

import "fmt"

func main() {
	arr := []int{9, 93, -12, 23, 923, -34, 38}
	quickSortAsc(arr, 0, len(arr)-1)

	fmt.Println(arr)
}

// quickSortAsc 整数切片的递增快速排序
func quickSortAsc(arr []int, left, right int) {
	if left < 0 || right <= 0 || left >= right || len(arr) <= 0 {
		return
	}

	i, j, base := left, right, arr[left]

	for i != j {
		for ; arr[j] >= base && i < j; j-- {
		}

		for ; arr[i] <= base && i < j; i++ {
		}

		if i < j {
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	arr[left], arr[i] = arr[i], arr[left]
	quickSortAsc(arr, left, i-1)
	quickSortAsc(arr, i+1, right)
}
