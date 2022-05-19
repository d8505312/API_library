package lib

import "sync"

// goroutine pool
type GoroutinePool struct {
	c  chan struct{}
	wg *sync.WaitGroup
}

// 採用有緩衝channel實現,當channel滿的時候阻塞
func NewGoroutinePool(maxSize int) *GoroutinePool {
	if maxSize <= 0 {
		panic("max size too small")
	}
	return &GoroutinePool{
		c:  make(chan struct{}, maxSize),
		wg: new(sync.WaitGroup),
	}
}

// add
func (g *GoroutinePool) Add(delta int) {
	g.wg.Add(delta)
	for i := 0; i < delta; i++ {
		g.c <- struct{}{}
	}

}

// done
func (g *GoroutinePool) Done() {
	<-g.c
	g.wg.Done()
}

// wait
func (g *GoroutinePool) Wait() {
	g.wg.Wait()
}
