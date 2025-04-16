package option

type EndType uint32

const (
	Backend  EndType = 0
	Frontend EndType = 1
)

func (e EndType) String() string {
	if e == Backend {
		return "后台"
	}

	return "前台"
}

func (e EndType) Is(v uint32) bool {
	return uint32(e) == v
}

func (e EndType) Value() uint32 {
	return uint32(e)
}
