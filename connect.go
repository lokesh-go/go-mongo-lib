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
	ReplicaSetName          string // Replica set name of the cluster, the cluster will be treated as a replica set and the driver will automatically discover all servers in the set, starting with the nodes specified through ApplyURI or SetHosts. All nodes in the replica set must have the same replica set name, or they will not be considered as part of the client. (default empty)
	MinPoolSize             uint64 // The minimum number of connections allowed in the driver's connection pool to each server. (default is 0)
	MaxPoolSize             uint64 // The maximum number of connections allowed in the driver's connection pool to each server. (default is 100)
	MaxConnecting           uint64 // The maximum number of connections a connection pool may establish simultaneously. (default is 2) (not recommended greater than 100)
	MaxConnIdleTime         int    // In milliseconds, The maximum amount of time that a connection will remain idle in a connection pool before it is removed from the pool and closed. (default is 0, meaning a connection can remain unused indefinitely)
	ServerSelectionTimeout  int    // In milliseconds, How long the driver will wait to find an available, suitable server to execute an operation. (default is 30 seconds)
	SocketTimeout           int    // In milliseconds, How long the driver will wait for a socket read or write to return before returning a network error. (default is 0, means no timeout is used and socket operations can block indefinitely)
	Timeout                 int    // In milliseconds, Amount of time that a single operation run on this client can execute before returning an error. (default value is nil, meaning operations do not inherit a timeout from the client)
	RetryReads              bool   // Supported read operations should be retried once on certain error, such as network errors. (default is true)
	RetryWrites             bool   // Supported write operations should be retried once on certain error, such as network errors. (default is true)
	ReadConcernWithMajority bool   // Majority specifies that the query should return the instance's most recent data acknowledged as having been written to a majority of members in the replica set.
	ReadSecondaryPreferred  bool   // In most situations, operation read from secondary members but if no secondary members are available, operations read from the primary on sharded clusters.
}
