package service

type ColumnsBuffer struct {
	columns []string
}

func ColumnBuffer(col ...string) *ColumnsBuffer {
	return &ColumnsBuffer{
		columns: col,
	}
}

func (c *ColumnsBuffer) Add(col ...string) {
	c.columns = append(c.columns, col...)
}

func (c *ColumnsBuffer) Val() []string {
	return c.columns
}
