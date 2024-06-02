package sort

func swap(arr []int, a int, b int) {
	temp := arr[a]
	arr[a] = arr[b]
	arr[b] = temp
}

func QuickSort(arr []int, len int) {
	sort(arr, 0, len)
}

func sort(arr []int, lo int, hi int) {
	if lo >= hi-1 {
		return
	}

	boundary := partition(arr, lo, hi)

	sort(arr, lo, boundary)
	sort(arr, boundary+1, hi)
}

func partition(arr []int, lo int, hi int) int {
	pivot := arr[hi-1]
	boundary := lo - 1
	cur := lo
	for cur < hi {
		if arr[cur] <= pivot {
			boundary++

			swap(arr, cur, boundary)
		}

		cur++
	}

	return boundary
}
