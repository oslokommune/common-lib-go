package db

//DbConf holds the database configuration information.
type DbConf struct {
	Username string
	Password string
	Database string
	Host     string
	Port     int
}

//NewDbConf creates a new DbConf with the prerequisite information.
func NewDbConf(username string, password string, host string, port int, database string) *DbConf {
	return &DbConf{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
	}
}

//Mapped adds the host and port mapped by a container.
func (c *DbConf) UpdateHostAndPort(mappedHost string, mappedPort int) *DbConf {
	c.Host = mappedHost
	c.Port = mappedPort
	return c
}
