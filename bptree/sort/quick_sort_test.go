package sort

import "testing"

func TestQuickSort(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	QuickSort(arr, len(arr))

	for i, n := range arr {
		if n != i+1 {
			t.Errorf("test failed! expected %d got %d", i+1, n)
		}
	}
}

func TestQuickSortReverse(t *testing.T) {
	arr := []int{4, 3, 2, 1}
	QuickSort(arr, len(arr))

	for i, n := range arr {
		if n != i+1 {
			t.Errorf("test failed! expected %d got %d", i+1, n)
		}
	}
}

func TestQuickSortFullLenArr(t *testing.T) {
	arr := []int{1, 3, 2, 4, 5, 10, 9, 9, 8}
	QuickSort(arr, len(arr))
	validAns := []int{1, 2, 3, 4, 5, 8, 9, 9, 10}
	for i, n := range arr {
		if n != validAns[i] {
			t.Errorf("test failed! expected %d got %d", validAns[i], n)
		}
	}

	arr = []int{10, 9, 9, 8, 1, 3, 2, 4, 5}
	QuickSort(arr, len(arr))
	for i, n := range arr {
		if n != validAns[i] {
			t.Errorf("test failed! expected %d got %d", validAns[i], n)
		}
	}
}

func TestQuickSortCapLenArr(t *testing.T) {
	arr := []int{1, 3, 2, 4, 5, 10, 9, 9, 8, 0, 0, 0, 0}
	QuickSort(arr, len(arr)-4)

	validAns := []int{1, 2, 3, 4, 5, 8, 9, 9, 10, 0, 0, 0, 0}
	for i, n := range arr {
		if n != validAns[i] {
			t.Errorf("test failed! expected %d got %d", validAns[i], n)
		}
	}
}
