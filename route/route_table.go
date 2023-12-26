package route

type RouteTable map[uint32] func(data []byte)

func (t RouteTable) Regist(id uint32, fn func(data []byte)) {
	t[id] = fn
}
