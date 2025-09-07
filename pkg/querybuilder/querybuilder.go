package querybuilder

// QueryBuilder helps build Elasticsearch DSL queries in a fluent style.
type QueryBuilder struct {
	query map[string]any
}

// New creates a new QueryBuilder instance.
func New() *QueryBuilder {
	return &QueryBuilder{
		query: make(map[string]any),
	}
}

// Match adds a match query for a given field.
func (qb *QueryBuilder) Match(field string, value any) *QueryBuilder {
	qb.query["match"] = map[string]any{
		field: value,
	}
	return qb
}

// Term adds a term query for a given field.
func (qb *QueryBuilder) Term(field string, value any) *QueryBuilder {
	qb.query["term"] = map[string]any{
		field: value,
	}
	return qb
}

// Range adds a range query for a given field.
func (qb *QueryBuilder) Range(field string, opts []map[string]any) *QueryBuilder {
	qb.query["range"] = map[string]any{
		field: opts,
	}
	return qb
}

// Bool creates a bool query with must/should/must_not/filters.
func (qb *QueryBuilder) Bool(must, should, mustNot, filter []map[string]any) *QueryBuilder {
	boolQuery := map[string]any{}

	if len(must) > 0 {
		boolQuery["must"] = must
	}

	if len(should) > 0 {
		boolQuery["should"] = should
	}

	if len(mustNot) > 0 {
		boolQuery["mustNot"] = mustNot
	}

	if len(filter) > 0 {
		boolQuery["filter"] = filter
	}
	qb.query["bool"] = boolQuery

	return qb
}

// Build finalizes the query  and returns the DSL object.
func (qb *QueryBuilder) Build() map[string]any {
	return qb.query
}
