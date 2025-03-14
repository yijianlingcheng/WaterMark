package images

var MAX_WORKER_NUM = 3

// Workerlimit 协程数量控制
type Workerlimit struct {
	n int
	c chan struct{}
}

// New 初始化
//
//	@param n
//	@return *Workerlimit
func NewWorkerlimit(n int) *Workerlimit {
	return &Workerlimit{
		n: n,
		c: make(chan struct{}, n),
	}
}

// worker 运行
//
//	@param f
func (g *Workerlimit) Worker(f func()) {
	g.c <- struct{}{}
	go func() {
		f()
		<-g.c
	}()
}
