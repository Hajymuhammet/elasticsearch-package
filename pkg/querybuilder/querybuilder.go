// querybuilder.go
package querybuilder

// QueryBuilder helps build Elasticsearch DSL queries in a fluent style.
type QueryBuilder struct {
	must    []map[string]any
	should  []map[string]any
	mustNot []map[string]any
	filter  []map[string]any
}

// New creates a new QueryBuilder instance.
func New() *QueryBuilder {
	return &QueryBuilder{}
}

// MultiMatch adds a multi_match query across multiple fields (goes into must by default).
func (qb *QueryBuilder) MultiMatch(fields []string, value any) *QueryBuilder {
	q := map[string]any{
		"multi_match": map[string]any{
			"query":  value,
			"fields": fields,
		},
	}
	qb.must = append(qb.must, q)
	return qb
}

// Match adds a match query for a given field (goes into must).
func (qb *QueryBuilder) Match(field string, value any) *QueryBuilder {
	q := map[string]any{
		"match": map[string]any{
			field: map[string]any{"query": value},
		},
	}
	qb.must = append(qb.must, q)
	return qb
}

// Term adds a term query for a given field (goes into filter).
func (qb *QueryBuilder) Term(field string, value any) *QueryBuilder {
	q := map[string]any{
		"term": map[string]any{
			field: value,
		},
	}
	qb.filter = append(qb.filter, q)
	return qb
}

// Range adds a range query for a given field (goes into filter).
func (qb *QueryBuilder) Range(field string, opts map[string]any) *QueryBuilder {
	q := map[string]any{
		"range": map[string]any{
			field: opts,
		},
	}
	qb.filter = append(qb.filter, q)
	return qb
}

// MustNot adds a must_not query (negation).
func (qb *QueryBuilder) MustNot(q map[string]any) *QueryBuilder {
	qb.mustNot = append(qb.mustNot, q)
	return qb
}

// Should adds a should query (OR logic).
func (qb *QueryBuilder) Should(q map[string]any) *QueryBuilder {
	qb.should = append(qb.should, q)
	return qb
}

// Build finalizes the query and returns the DSL object.
func (qb *QueryBuilder) Build() map[string]any {
	boolQuery := make(map[string]any)

	if len(qb.must) > 0 {
		boolQuery["must"] = qb.must
	}
	if len(qb.should) > 0 {
		boolQuery["should"] = qb.should
	}
	if len(qb.mustNot) > 0 {
		boolQuery["must_not"] = qb.mustNot
	}
	if len(qb.filter) > 0 {
		boolQuery["filter"] = qb.filter
	}

	return map[string]any{
		"query": map[string]any{
			"bool": boolQuery,
		},
	}
}
