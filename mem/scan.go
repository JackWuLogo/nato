package mem

import "math"

type PkFkData struct {
	Pk string
	Fk string
}

// 分页结果
type ScanResult struct {
	Cur   int64       // 当前页数
	Total int64       // 总页数
	Limit int64       // 每页显示数
	Count int64       // 总数据量
	Rows  interface{} // 数据列表
}

// 获取单主键列表数据
func (sr *ScanResult) GetRows() []string {
	if val, ok := sr.Rows.([]string); ok {
		return val
	}
	return nil
}

// 获取含有外键的主键列表
func (sr *ScanResult) GetFkRows() []*PkFkData {
	if val, ok := sr.Rows.([]*PkFkData); ok {
		return val
	}
	return nil
}

func NewScanResult(cur int64, count int64, limit int64) *ScanResult {
	// 计算最大页数
	total := int64(math.Ceil(float64(count) / float64(limit)))

	if cur <= 0 {
		cur = 1
	}
	if cur > total {
		cur = total
	}

	return &ScanResult{
		Cur:   cur,
		Total: total,
		Limit: limit,
		Count: count,
	}
}
