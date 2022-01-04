package tool

import (
	"math/rand"
	"sort"
)

// Choice is a generic wrapper that can be used to add weights for any object
type Choice struct {
	Item   interface{}
	Weight int
}

// A Chooser caches many possible Choices in a structure designed to improve
// performance on repeated calls for weighted random selection.
type Chooser struct {
	data   []*Choice
	totals []int
	max    int
}

func (c *Chooser) Data() []*Choice {
	return c.data
}

// Add 增加权重随机项
func (c *Chooser) Add(item interface{}, weight int) {
	c.data = append(c.data, &Choice{Item: item, Weight: weight})
}

// Search 使用二分法查找
func (c *Chooser) Search() int {
	c.totals = make([]int, len(c.data))
	c.max = 0

	for i, ch := range c.data {
		c.max += ch.Weight
		c.totals[i] = c.max
	}

	r := rand.Intn(c.max) + 1         // 使用最大值获取随机数，避免超过范围，随机生成的数需要排除0，故加1
	i := sort.SearchInts(c.totals, r) // 使用二分法，找到对应的下标，如果没有则为大于该数的+1 下标，可能为len(a)即数组长度。

	return i
}

// Pick 返回单个权重随机项
func (c *Chooser) Pick() interface{} {
	// 对元素进行递增排序。
	sort.Slice(c.data, func(i, j int) bool {
		return c.data[i].Weight < c.data[j].Weight
	})

	index := c.Search()

	return c.data[index].Item
}

// NewChooser initializes a new Chooser consisting of the possible Choices.
func NewChooser(cs ...*Choice) *Chooser {
	return &Chooser{data: cs}
}

// MultiChooserPick 返回多个权重随机项
func MultiChooserPick(c *Chooser, count int) []interface{} {
	var res []interface{}

	for i := 0; i < count; i++ {
		index := c.Search()
		res = append(res, c.data[index].Item)

		data := c.Data()
		c = NewChooser(append(data[:index], data[index+1:]...)...)
	}

	return res
}
