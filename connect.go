package mongodb

import (
	"context"
	"crypto/tls"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

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
	ReplicaSetName           string // Replica set name of the cluster, the cluster will be treated as a replica set and the driver will automatically discover all servers in the set, starting with the nodes specified through ApplyURI or SetHosts. All nodes in the replica set must have the same replica set name, or they will not be considered as part of the client. (default empty)
	MinPoolSize              uint64 // The minimum number of connections allowed in the driver's connection pool to each server. (default is 0)
	MaxPoolSize              uint64 // The maximum number of connections allowed in the driver's connection pool to each server. (default is 100)
	MaxConnecting            uint64 // The maximum number of connections a connection pool may establish simultaneously. (default is 2) (not recommended greater than 100)
	MaxConnIdleTime          int    // In milliseconds, The maximum amount of time that a connection will remain idle in a connection pool before it is removed from the pool and closed. (default is 0, meaning a connection can remain unused indefinitely)
	ServerSelectionTimeout   int    // In milliseconds, How long the driver will wait to find an available, suitable server to execute an operation. (default is 30 seconds)
	SocketTimeout            int    // In milliseconds, How long the driver will wait for a socket read or write to return before returning a network error. (default is 0, means no timeout is used and socket operations can block indefinitely)
	Timeout                  int    // In milliseconds, Amount of time that a single operation run on this client can execute before returning an error. (default value is nil, meaning operations do not inherit a timeout from the client)
	RetryReads               bool   // Supported read operations should be retried once on certain error, such as network errors. (default is true)
	RetryWrites              bool   // Supported write operations should be retried once on certain error, such as network errors. (default is true)
	ReadConcernWithMajority  bool   // Majority specifies that the query should return the instance's most recent data acknowledged as having been written to a majority of members in the replica set.
	ReadSecondaryPreferred   bool   // In most situations, operation read from secondary members but if no secondary members are available, operations read from the primary on sharded clusters.
	WriteConcernWithMajority bool   // Majority of nodes must acknowledge write operations before the operation returns.
	WriteConcernTimeout      int    // In milliseconds, How long write operations should wait for the correct number of nodes to acknowledge the operation.
}

// Intialise and return new mongodb client connection
func New(config *Config) (dbClient *Client, err error) {
	// Connects
	dbClient, err = config.connect()
	if err != nil {
		return nil, err
	}

	// Returns
	return dbClient, nil
}

// Connect
func (c *Config) connect() (dbClient *Client, err error) {
	// Assigns hosts
	mongoConnOptions := &options.ClientOptions{
		Hosts: c.Hosts,
	}

	// Checks TLS
	if c.TLSEnabled {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		mongoConnOptions.TLSConfig = tlsConfig
	}

	// Checks auth
	if c.AuthEnabled {
		creds := &options.Credential{
			Username:   c.User,
			Password:   c.Password,
			AuthSource: c.AuthSource,
		}
		mongoConnOptions.Auth = creds
	}

	// Checks & apply connection config
	if c.Connection != nil {
		// Sets min pool size
		if c.Connection.MinPoolSize != 0 {
			mongoConnOptions.SetMinPoolSize(c.Connection.MinPoolSize)
		}

		// Sets max pool size
		if c.Connection.MaxPoolSize != 0 {
			mongoConnOptions.SetMaxPoolSize(c.Connection.MaxPoolSize)
		}

		// Sets max connecting
		if c.Connection.MaxConnecting != 0 {
			mongoConnOptions.SetMaxConnecting(c.Connection.MaxConnecting)
		}

		// Sets max connection idle time
		if c.Connection.MaxConnIdleTime != 0 {
			mongoConnOptions.SetMaxConnIdleTime(time.Duration(c.Connection.MaxConnIdleTime) * time.Millisecond)
		}

		// Sets server selection timeout
		if c.Connection.ServerSelectionTimeout != 0 {
			mongoConnOptions.SetServerSelectionTimeout(time.Duration(c.Connection.ServerSelectionTimeout) * time.Millisecond)
		}

		// Sets socket timeout
		if c.Connection.SocketTimeout != 0 {
			mongoConnOptions.SetSocketTimeout(time.Duration(c.Connection.SocketTimeout) * time.Millisecond)
		}

		// Sets timeout
		if c.Connection.Timeout != 0 {
			mongoConnOptions.SetTimeout(time.Duration(c.Connection.Timeout) * time.Millisecond)
		}

		// Sets read concern majority
		if c.Connection.ReadConcernWithMajority {
			majorityLevel := readconcern.Majority().GetLevel()
			rc := readconcern.New(readconcern.Level(majorityLevel))
			mongoConnOptions.SetReadConcern(rc)
		}

		// Sets read secondary preferred
		if c.Connection.ReadSecondaryPreferred {
			rp := readpref.SecondaryPreferred()
			mongoConnOptions.SetReadPreference(rp)
		}

		// Sets write concern with majority
		if c.Connection.WriteConcernWithMajority {
			wc := writeconcern.New(writeconcern.WMajority())

			// Sets write conecern timeout
			if c.Connection.WriteConcernTimeout != 0 {
				wc.WithOptions(writeconcern.WTimeout(time.Duration(c.Connection.WriteConcernTimeout) * time.Millisecond))
			}
			mongoConnOptions.SetWriteConcern(wc)
		}
	}

	// Gets new mongodb client
	client, err := mongo.NewClient(mongoConnOptions)
	if err != nil {
		return nil, err
	}

	// Connect client with timeout
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, errors.New("client connection failed " + err.Error())
	}

	// Ping
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, errors.New("client ping failed ->" + err.Error())
	}

	// Sets client connection
	dbClient = &Client{}
	dbClient.database = client.Database(c.Database)

	// Returns
	return dbClient, nil
}
