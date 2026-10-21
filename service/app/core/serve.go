package core

type Serve interface {
	Start()
	Stop()
}

type ServeGroup struct {
	serves []Serve
}

// Add 添加实现了 Serve 接口的结构体。
func (g *ServeGroup) Add(s Serve) {
	if s == nil {
		return
	}
	g.serves = append(g.serves, s)
}

func (g *ServeGroup) Start() {
	for _, s := range g.serves {
		s.Start()
	}
}

func (g *ServeGroup) Stop() {
	for _, s := range g.serves {
		s.Stop()
	}
}
