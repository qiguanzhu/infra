/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2024/7/13 -- 14:53
 @Author  : bishop ❤️ MONEY
 @Description: xrand.go
*/

package xrand

import (
	"math/rand"
	"time"
)

var RanD random

type random struct {
	rd *rand.Rand
}

var rd *rand.Rand

const (
	PROBABILITY int32 = 10000
)

func init() {
	rd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (r *random) GetRandomNumber(n int32) int32 {
	return r.getRandomNumber(n)
}

func (r *random) getRandomNumber(n int32) int32 {
	if n == 0 {
		n = PROBABILITY
	}
	return rd.Int31n(n)
}

func (r *random) CheckProbabilityJackpot(prob int32) bool {
	number := r.getRandomNumber(PROBABILITY)
	if number < prob {
		return true
	}
	return false
}

func (r *random) calRandomNumberBetween(min, max int32) int32 {
	// 认为只要调用数量至少为1
	if max == 0 {
		return 1
	}
	number := r.getRandomNumber(max - min)
	cal := min + number
	if cal != 0 {
		return cal
	}
	return 1
}

func Wave(base int64, lowRate int64, upRate int64) int64 {
	rate := rand.Intn(int((upRate-lowRate)+1)) + int(lowRate)

	return base + int64(float64(base)*float64(rate)/100)
}

func RandBetween(start int64, end int64) int64 {
	if start > end {
		return 0
	}
	return int64(rand.Intn(int(end-start)+1)) + start
}

// choose n from m eg: we have m=10 and n=5, then we will get [1,3,5,7,8]
func ChooseMN(m, n int) []int {
	chosen := make([]int, 0, n)
	if m <= n {
		for i := 0; i < m; i++ {
			chosen = append(chosen, i)
		}
		return chosen
	}
	choosed := map[int]bool{}
	for {
		if len(chosen) >= n {
			break
		}
		idx := rand.Intn(m)
		if !choosed[idx] {
			chosen = append(chosen, idx)
			choosed[idx] = true
		}
	}
	return chosen
}
