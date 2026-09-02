package pagination

type Pagination struct {
	Limit  uint64
	Offset uint64
}

func (p *Pagination) GetLimit() uint64 {
	if p == nil {
		return 0
	}

	return p.Limit
}

func (p *Pagination) GetOffset() uint64 {
	if p == nil {
		return 0
	}

	return p.Offset
}
