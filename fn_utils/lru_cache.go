package fn_utils

import (
	"errors"
	"github.com/vmihailenco/msgpack/v5"
	"log"
	"sync"
	"time"
)

type (
	tNode struct {
		key      string
		data     interface{}
		dataSize int64
		lastTime int64
		ttl      int64
		lruPre   *tNode
		lruNext  *tNode
	}
	//LruCache  lru ttl cache
	LruCache struct {
		spaceSize int64
		leftSpace int64
		lruHead   *tNode
		lruTail   *tNode
		datas     map[string]*tNode
		mvMux     *sync.Mutex
		rwMux     *sync.RWMutex
	}
)

//Add 增加数据
func (l *LruCache) Add(key string, data interface{}, ttl int64) {
	//startTIme := time.Now().Nanosecond()
	dataSize := getSize(data)
	// free space
	l.rwMux.Lock()
	defer l.rwMux.Unlock()
	if !l.checkSpace(dataSize) {
		if err := l.free(dataSize); err != nil {
			panic(err)
		}
	}

	var node *tNode
	var ok bool
	if node, ok = l.datas[key]; !ok {
		node = &tNode{
			key: key,
		}
		l.datas[key] = node
	} else {
		l.leftSpace = l.leftSpace + node.dataSize
	}

	node.data = data
	node.dataSize = dataSize
	node.ttl = ttl
	if l.lruTail == nil {
		l.lruTail = node
	}

	l.lruMoveToHead(node)
	l.leftSpace = l.leftSpace - dataSize
	//log.Printf("Add cost := %v leftSpace =%d", time.Now().Nanosecond()-startTIme, l.leftSpace)
	return
}

// Get 获取数据
func (l *LruCache) Get(key string) interface{} {
	l.rwMux.RLock()
	defer l.rwMux.RUnlock()
	if data, ok := l.datas[key]; ok {
		if data.ttl-time.Now().UnixMilli() > 0 {
			l.lruMoveToHead(data)
			return data.data
		} else {
			if data.lruPre != nil {
				data.lruPre.lruNext = data.lruNext
			}
			if data.lruNext != nil {
				data.lruNext.lruPre = data.lruPre
			}
			delete(l.datas, data.key)
		}
	}
	return nil
}

func (l *LruCache) checkSpace(dataSize int64) bool {
	if l.leftSpace-dataSize > 0 {
		return true
	}
	return false
}

func (l *LruCache) free(dataSize int64) error {
	for l.freeLru() {
		if l.leftSpace > dataSize {
			return nil
		} else if l.leftSpace == l.spaceSize {
			return errors.New("so big data")
		}
	}
	return nil
}

func (l *LruCache) freeLru() bool {
	tail := l.lruTail
	l.lruTail = tail.lruPre
	l.leftSpace += tail.dataSize
	delete(l.datas, tail.key)
	return true
}

// 将节点移动到头节点
func (l *LruCache) lruMoveToHead(node *tNode) {
	l.mvMux.Lock()
	defer l.mvMux.Unlock()
	if node.lruPre != nil {
		node.lruPre.lruNext = node.lruNext
	}

	if node.lruNext != nil {
		node.lruNext.lruPre = node.lruPre
	}
	if l.lruHead != nil {
		l.lruHead.lruPre = node
	}
	node.lruNext = l.lruHead
	l.lruHead = node
}

// 默认不报错
func getSize(data interface{}) int64 {
	marshal, _ := msgpack.Marshal(data)
	return int64(len(marshal))
}

func (l LruCache) GetLeftSpace() int64 {
	log.Printf("%d", len(l.datas))
	return l.leftSpace
}

// NewLruCache TODO
func NewLruCache(spaceSize int64) *LruCache {
	cache := &LruCache{
		spaceSize: spaceSize,
		leftSpace: spaceSize,
		rwMux:     &sync.RWMutex{},
		mvMux:     &sync.Mutex{},
		datas:     make(map[string]*tNode),
	}
	return cache
}
