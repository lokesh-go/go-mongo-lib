package mongodb

// Config contains all properties required for creating a connection
type Config struct {
	Hosts       []string    // Database server hosts
	AuthEnabled bool        // Enables auth to required user & password to establish connection
	User        string      // Db Username for authentication
	Password    string      // Db password for authentication
	AuthSource  string      // The name of database to use for authentication
	TLSEnabled  bool        // TLS to encrypt all of mongodb's network traffic
	Database    string      // Db name
	Connection  *Connection // More client options
}

// Sets more client options
type Connection struct {
	MinPoolSize            uint64 // The minimum number of connections allowed in the driver's connection pool to each server. (default is 0)
	MaxPoolSize            uint64 // The maximum number of connections allowed in the driver's connection pool to each server. (default is 100)
	MaxConnecting          uint64 // The maximum number of connections a connection pool may establish simultaneously. (default is 2) (not recommended greater than 100)
	MaxConnIdleTime        int    // In milliseconds, The maximum amount of time that a connection will remain idle in a connection pool before it is removed from the pool and closed. (default is 0, meaning a connection can remain unused indefinitely)
	ServerSelectionTimeout int    // In milliseconds, How long the driver will wait to find an available, suitable server to execute an operation. (default is 30 seconds)
	SocketTimeout          int    // In milliseconds, How long the driver will wait for a socket read or write to return before returning a network error. (default is 0, means no timeout is used and socket operations can block indefinitely)
	Timeout                int    // In milliseconds, Amount of time that a single operation run on this client can execute before returning an error. (default value is nil, meaning operations do not inherit a timeout from the client)
}
