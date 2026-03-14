package query

import (
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/sortmode"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/sortorder"
)

// MissingFirst places documents with missing field values at the beginning of the sort order.
const MissingFirst = "_first"

// MissingLast places documents with missing field values at the end of the sort order.
const MissingLast = "_last"

// SortBuilder constructs Elasticsearch sort specifications using method chaining.
type SortBuilder struct {
	sorts []types.SortCombinations
}

// NewSort returns a new empty SortBuilder.
func NewSort() *SortBuilder {
	return &SortBuilder{}
}

// Field adds a sort on the given field with the specified order.
func (s *SortBuilder) Field(name string, order sortorder.SortOrder) *SortBuilder {
	s.sorts = append(s.sorts, types.SortOptions{
		SortOptions: map[string]types.FieldSort{
			name: {Order: &order},
		},
	})
	return s
}

// FieldWithMissing adds a sort on the given field with a missing value position.
// Use MissingFirst or MissingLast to control where documents with missing field values appear.
func (s *SortBuilder) FieldWithMissing(name string, order sortorder.SortOrder, missing string) *SortBuilder {
	var m types.Missing = missing
	s.sorts = append(s.sorts, types.SortOptions{
		SortOptions: map[string]types.FieldSort{
			name: {Order: &order, Missing: m},
		},
	})
	return s
}

// FieldNested adds a nested field sort with a sort mode.
func (s *SortBuilder) FieldNested(name string, order sortorder.SortOrder, path string, mode sortmode.SortMode) *SortBuilder {
	s.sorts = append(s.sorts, types.SortOptions{
		SortOptions: map[string]types.FieldSort{
			name: {
				Order:  &order,
				Mode:   &mode,
				Nested: &types.NestedSortValue{Path: path},
			},
		},
	})
	return s
}

// FieldCustom adds a sort on the given field with a fully customized FieldSort.
func (s *SortBuilder) FieldCustom(name string, fs types.FieldSort) *SortBuilder {
	s.sorts = append(s.sorts, types.SortOptions{
		SortOptions: map[string]types.FieldSort{
			name: fs,
		},
	})
	return s
}

// ScoreDesc adds a sort by relevance score in descending order.
func (s *SortBuilder) ScoreDesc() *SortBuilder {
	order := sortorder.Desc
	s.sorts = append(s.sorts, types.SortOptions{
		Score_: &types.ScoreSort{Order: &order},
	})
	return s
}

// ScoreAsc adds a sort by relevance score in ascending order.
func (s *SortBuilder) ScoreAsc() *SortBuilder {
	order := sortorder.Asc
	s.sorts = append(s.sorts, types.SortOptions{
		Score_: &types.ScoreSort{Order: &order},
	})
	return s
}

// DocAsc adds a sort by index order in ascending order.
func (s *SortBuilder) DocAsc() *SortBuilder {
	order := sortorder.Asc
	s.sorts = append(s.sorts, types.SortOptions{
		Doc_: &types.ScoreSort{Order: &order},
	})
	return s
}

// DocDesc adds a sort by index order in descending order.
func (s *SortBuilder) DocDesc() *SortBuilder {
	order := sortorder.Desc
	s.sorts = append(s.sorts, types.SortOptions{
		Doc_: &types.ScoreSort{Order: &order},
	})
	return s
}

// Build returns the constructed sort slice ready for use in a search request.
func (s *SortBuilder) Build() []types.SortCombinations {
	return s.sorts
}
