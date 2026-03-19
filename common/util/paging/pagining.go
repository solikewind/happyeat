package paging

type PageParam struct {
	Current  int
	PageSize int
}

func NewPageParam(current, pageSize int64) PageParam {
	curr := int(current)
	size := int(pageSize)

	if curr <= 0 {
		curr = 1
	}
	if size <= 0 {
		size = 10
	} else if size > 100 { // 增加安全上限
		size = 100
	}

	return PageParam{
		Current:  curr,
		PageSize: size,
	}
}

// Offset 计算数据库偏移量
func (p PageParam) Offset() int {
	return (p.Current - 1) * p.PageSize
}
