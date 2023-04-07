package panda

import "math"

func Min(nums ...int) int {

	num := math.MaxInt

	for _, v := range nums {

		if v < num {

			num = v
		}
	}

	return num
}

func Max(nums ...int) int {

	num := math.MinInt

	for _, v := range nums {

		if num < v {

			num = v
		}
	}

	return num
}

func Avg(nums ...int) int {

	num := 0
	k := 0

	for _, v := range nums {

		if k == 1 {

			num += v
			num /= 2
		}
		num = v
		k = 1
	}

	return num
}
