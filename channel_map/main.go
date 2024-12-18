package main

type OpType uint8

const (
	Put OpType = iota
	Get
	Delete
	Len
)

type Request[K comparable, V any] struct {
	OpType  OpType
	Key     K
	Value   V
	Result  chan any
	Succeed bool
}

type SafeMap[K comparable, V any] struct {
	m     map[K]V
	reqCh chan *Request[K, V]
}

func (sm *SafeMap[K, V]) Run() {
	for req := range sm.reqCh {
		switch req.OpType {
		case Get:
			value, exists := sm.m[req.Key]
			req.Succeed = exists
			req.Result <- value
		case Put:
			sm.m[req.Key] = req.Value
			req.Succeed = true
			req.Result <- nil
		case Delete:
			_, exists := sm.m[req.Key]
			delete(sm.m, req.Key)
			req.Succeed = exists
			req.Result <- nil
		case Len:
			req.Succeed = true
			req.Result <- len(sm.m)
		}
	}
}

func (sm *SafeMap[K, V]) Get(key K) (V, bool) {
	req := &Request[K, V]{
		OpType: Get,
		Key:    key,
		Result: make(chan any),
	}
	sm.reqCh <- req
	res := <-req.Result
	if v, ok := res.(V); ok {
		return v, req.Succeed
	}
	panic("type assertion failed")
}

func (sm *SafeMap[K, V]) Set(key K, value V) {
	req := &Request[K, V]{
		OpType: Put,
		Key:    key,
		Value:  value,
		Result: make(chan any),
	}
	sm.reqCh <- req
	<-req.Result
}

func (sm *SafeMap[K, V]) Delete(key K) {
	req := &Request[K, V]{
		OpType: Delete,
		Key:    key,
		Result: make(chan any),
	}
	sm.reqCh <- req
	<-req.Result
}

func (sm *SafeMap[K, V]) Len() int {
	req := &Request[K, V]{
		OpType: Len,
		Result: make(chan any),
	}
	sm.reqCh <- req
	res := <-req.Result
	v, ok := res.(int)
	if !ok {
		panic(ok)
	}
	return v
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	sm := SafeMap[K, V]{
		m:     make(map[K]V),
		reqCh: make(chan *Request[K, V]),
	}
	go sm.Run()
	return &sm
}

//func main() {
//	sm := NewSafeMap[string, int]()
//	sm.Set("apple", 10)
//	sm.Set("banana", 20)
//
//	length := sm.Len()
//	fmt.Println(length)
//
//	value, exists := sm.Get("apple")
//	if exists {
//		fmt.Printf("key = apple, value = %d", value)
//	}
//
//}
