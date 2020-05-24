package main

import (
	"fmt"
	"github.com/vearne/gtimer"
	"log"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	PRODUCER_COUNT = 10
	CONSUMER_COUNT = 10
	TARGET_COUNT   = 1000000
)

var ops int64 = 0

func main() {
	st := gtimer.NewSuperTimer(CONSUMER_COUNT)
	t1 := time.Now()
	for i := 0; i < PRODUCER_COUNT; i++ {
		go push(st, "worker"+strconv.Itoa(i))
	}

	time.Sleep(100 * time.Millisecond)

	for {
		v := atomic.LoadInt64(&ops)
		if v >= TARGET_COUNT {
			st.Stop()
			break
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
	t2 := time.Now()
	log.Printf("cost:%v\n", t2.Sub(t1))
}

func DefaultAction(t time.Time, value string) {
	// fmt.Printf("trigger_time:%v, value:%v\n", t, value)
	atomic.AddInt64(&ops, 1)
}

func push(timer *gtimer.SuperTimer, name string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 1000000; i++ {
		now := time.Now()
		t := now.Add(time.Millisecond * time.Duration(r.Int63n(300)))
		value := fmt.Sprintf("%v:value:%v", name, strconv.Itoa(i))
		// create a delayed task
		item := gtimer.NewDelayedItemFunc(t, value, DefaultAction)
		timer.Add(item)
	}
}