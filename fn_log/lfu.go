package fn_log

import (
	"math/rand"
	"time"
)

var RandMax int64 = 0x7fffffff

var LfuInitVal int64 = 5

/*# +--------+------------+------------+------------+------------+------------+
# | factor | 100 hits   | 1000 hits  | 100K hits  | 1M hits    | 10M hits   |
# +--------+------------+------------+------------+------------+------------+
# | 0      | 104        | 255        | 255        | 255        | 255        |
# +--------+------------+------------+------------+------------+------------+
# | 1      | 18         | 49         | 255        | 255        | 255        |
# +--------+------------+------------+------------+------------+------------+
# | 10     | 10         | 18         | 142        | 255        | 255        |
# +--------+------------+------------+------------+------------+------------+
# | 100    | 8          | 11         | 49         | 143        | 255        |
# +--------+------------+------------+------------+------------+------------+
*/
var LfuLogFactor = 10

var LfuDecayTime int64 = 1

// counter 默认值
var Adapter CacheAdapter

type CacheAdapter interface {
	Get(key string, data interface{}) *LfuNode
	Save(key string, score int64, data *LfuNode)
	Size(key string) int64
	DeleteSuffix() //删除末尾的
}

type Lfu interface {
	Lookup(key string, data interface{})
	List(key string) []LfuNode
}

func NewLfu(capacity int64) Lfu {
	return &lfu{
		capacity: capacity,
	}
}

type lfu struct {
	capacity int64
}

func (l lfu) Lookup(key string, data interface{}) {
	//获取 这个数据
	node := Adapter.Get(key, data)
	if node == nil {
		node = &LfuNode{
			lru:  (lfuGetTimeInMinutes() << 8) | LfuInitVal,
			data: data,
		}
		Adapter.Save(key, node.lru, node)
	} else {
		oldScore := node.lru
		updateLFU(node)
		if node.lru != oldScore {
			Adapter.Save(key, node.lru, node)
		}
		if Adapter.Size(key) > l.capacity {
			Adapter.DeleteSuffix()
		}
	}
}

func (l lfu) List(key string) []LfuNode {
	return []LfuNode{}
}

type LfuNode struct {
	lru  int64
	data interface{}
}

func lfuDecrAndReturn(l *LfuNode) int64 {
	ldt := l.lru >> 8
	counter := l.lru & 255
	var numPeriods int64 = 0
	if LfuDecayTime > 0 {
		numPeriods = lfuTimeElapsed(ldt) / LfuDecayTime
	}

	if numPeriods > 0 {
		if numPeriods > counter {
			counter = 0
		} else {
			counter = counter - numPeriods
		}
	}
	return counter
}

func lfuGetTimeInMinutes() int64 {
	return time.Now().Unix() & 65535
}

func lfuTimeElapsed(ldt int64) int64 {
	now := lfuGetTimeInMinutes()
	if now >= ldt {
		return now - ldt
	}
	return 65535 - ldt + now
}

func lfuLogIncr(counter int64) int64 {
	if counter == 255 {
		return 255
	}
	r := float64(rand.Int63()) / float64(RandMax)
	baseVal := float64(counter) - float64(LfuInitVal)
	if baseVal < 0 {
		baseVal = 0
	}

	p := 1.0 / (baseVal*float64(LfuLogFactor) + 1)
	if r < p {
		counter++
	}
	return counter
}

func updateLFU(l *LfuNode) {
	counter := lfuDecrAndReturn(l)
	counter = lfuLogIncr(counter)
	l.lru = (lfuGetTimeInMinutes() << 8) | counter
}
