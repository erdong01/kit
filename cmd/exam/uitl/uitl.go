package uitl

import (
	"math/rand"
	"time"
)

func Interset(nums1 []string, nums2 []string) []string {
	m := make(map[string]int)
	for _, v := range nums1 {
		m[v]++
	}
	for _, v := range nums2 {
		times, _ := m[v]
		if times == 1 {
			nums1 = append(nums2, v)
		}
	}
	return nums1
}

func Difference(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := Interset(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}
	for _, value := range slice2 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}

//生成count个[start,end)结束的不重复的随机数
func GenerateRandomNumber(start int, end int, count int) []int {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}

	//存放结果的slice
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn((end - start)) + start

		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}
