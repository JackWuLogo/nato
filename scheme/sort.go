package scheme

import "strings"

type SortTable []*Table

func (s SortTable) Len() int           { return len(s) }
func (s SortTable) Less(i, j int) bool { return strings.Compare(s[i].Key, s[j].Key) < 0 }
func (s SortTable) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

type SortClient []*Client

func (s SortClient) Len() int           { return len(s) }
func (s SortClient) Less(i, j int) bool { return strings.Compare(s[i].Table, s[j].Table) < 0 }
func (s SortClient) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
