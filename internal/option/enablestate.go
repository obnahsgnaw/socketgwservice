package option

type EnableState int32

const (
	Enabled  EnableState = 1
	Default  EnableState = -1
	Disabled EnableState = 0
)

func (e EnableState) String() string {
	if e == Enabled {
		return "已启用"
	}
	if e == Disabled {
		return "已禁用"
	}

	return "默认"
}

func (e EnableState) Is(v int32) bool {
	return int32(e) == v
}

func (e EnableState) Value() int32 {
	return int32(e)
}

