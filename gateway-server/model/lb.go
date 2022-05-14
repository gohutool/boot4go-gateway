package model

import (
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : lb.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/7 23:38
* 修改历史 : 1. [2022/5/7 23:38] 创建文件 by LongYong
*/

type LoadBalance interface {
	Target() (*Target, error)
}

var ErrNoPointer = errors.New("no endpoints available")

func NewRandom(s []*Target, seed int64) LoadBalance {
	return &random{
		s: s,
		r: rand.New(rand.NewSource(seed)),
		l: len(s),
	}
}

type random struct {
	s []*Target
	r *rand.Rand
	l int
}

func (r *random) Target() (*Target, error) {
	if len(r.s) <= 0 {
		return nil, ErrNoPointer
	}
	return r.s[r.r.Intn(r.l)], nil
}

func NewRoundRobin(s []*Target) LoadBalance {
	return &roundRobin{
		s: s,
		c: new(uint64),
		l: len(s),
	}
}

type roundRobin struct {
	s []*Target
	c *uint64
	l int
}

func (rr *roundRobin) Target() (*Target, error) {
	if rr.l <= 0 {
		return nil, ErrNoPointer
	}
	old := atomic.AddUint64(rr.c, 1) - 1
	idx := old % uint64(len(rr.s))
	return rr.s[idx], nil
}

func NewWeightRoundRobin(s []*Target) LoadBalance {
	for _, target := range s {
		if target.Weight == 0 {
			target.Weight = 1
		}
	}
	ss := &safeTargets{
		mtx:     &sync.Mutex{},
		targets: s,
	}
	return &weightRoundRobin{
		s: ss,
		c: 0,
		l: len(s),
	}
}

type safeTargets struct {
	mtx     *sync.Mutex
	targets []*Target
}

type weightRoundRobin struct {
	s *safeTargets
	c uint64
	l int
}

func (wrr *weightRoundRobin) getNextPointerIndex() int {
	index := -1
	var total int8 = 0
	wrr.s.mtx.Lock()
	defer wrr.s.mtx.Unlock()

	for i := 0; i < wrr.l; i++ {
		wrr.s.targets[i].CurrentWeight += wrr.s.targets[i].Weight

		total += wrr.s.targets[i].Weight
		if index == -1 || wrr.s.targets[index].CurrentWeight < wrr.s.targets[i].CurrentWeight {
			index = i
		}
	}
	wrr.s.targets[index].CurrentWeight -= total
	return index
}

func (wrr *weightRoundRobin) Target() (*Target, error) {
	index := wrr.getNextPointerIndex()
	return wrr.s.targets[index], nil
}
