// config/config.go
package config

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

var Session *gocql.Session

func InitCassandra() {
	var err error

	// Connect initially to the 'system' keyspace to ensure other operations can proceed
	cluster := gocql.NewCluster("cassandra")
	cluster.Port = 9042
	cluster.Keyspace = "system" // Connect to 'system' keyspace first
	cluster.Consistency = gocql.Quorum
	cluster.ConnectTimeout = 10 * time.Second
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 5}

	// Retry loop to connect to Cassandra
	for {
		Session, err = cluster.CreateSession()
		if err != nil {
			log.Printf("Unable to connect to Cassandra: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	log.Println("Connected to Cassandra 'system' keyspace")

	// Create keyspace if it does not exist
	keyspace := "user_db" // Replace with your keyspace name
	query := fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`, keyspace)
	if err := Session.Query(query).Exec(); err != nil {
		log.Fatalf("Failed to create keyspace '%s': %v", keyspace, err)
	}

	// Close initial session and reconnect to the newly created keyspace
	Session.Close()
	cluster.Keyspace = keyspace
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatalf("Unable to connect to Cassandra keyspace '%s': %v", keyspace, err)
	}
	log.Printf("Connected to keyspace '%s'", keyspace)

	// Create table if not exists
	tableQuery := `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			name TEXT,
			email TEXT
		)
	`
	if err := Session.Query(tableQuery).Exec(); err != nil {
		log.Fatalf("Failed to create 'users' table: %v", err)
	}
	log.Println("Table 'users' ensured")
}

func CloseSession() {
	if Session != nil {
		Session.Close()
	}
}
