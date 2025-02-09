package testing

import "testing"

type Record[A any, E any] struct {
	Name    string
	Subject func(t *testing.T) (A, error)
	Expect  func(t *testing.T) (E, error)
	Assert  func(t *testing.T, r *Record[A, E]) bool
}

type Table[A, E any] struct {
	records []Record[A, E]
}

func NewTable[A, E any](records []Record[A, E]) *Table[A, E] {
	return &Table[A, E]{
		records: records,
	}
}

func (table *Table[A, E]) Run(test *testing.T) {
	for _, record := range table.records {
		test.Run(record.Name, func(t *testing.T) {
			if !record.Assert(t, &record) {
				t.Errorf("Test %s failed", record.Name)
			}
		})
	}
}
