package main

type SelectQuery struct {
	tableName string
	columns   []string
	where     string
	orderBy   string
	limit     int
}

type Option func(*SelectQuery)

func WithLimit(n int) Option {
	return func(s *SelectQuery) {
		s.limit = n
	}
}

func WithWhere(where string) Option {
	return func(s *SelectQuery) {
		s.where = where
	}
}

func WithOrderBy(orderBy string) Option {
	return func(s *SelectQuery) {
		s.orderBy = orderBy
	}
}

func WithColumns(cols ...string) Option {
	return func(s *SelectQuery) {
		s.columns = cols
	}
}

func NewSelectQuery(table string, opts ...Option) *SelectQuery {
	query := &SelectQuery{
		tableName: table,
		columns:   []string{"*"},
		where:     "",
		orderBy:   "",
		limit:     0,
	}

	for _, opt := range opts {
		opt(query)
	}
	return query
}
