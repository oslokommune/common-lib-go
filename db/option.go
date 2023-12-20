package db

// Option -.
type Option func(*PostgresConnection)

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *PostgresConnection) {
		c.maxPoolSize = size
	}
}
