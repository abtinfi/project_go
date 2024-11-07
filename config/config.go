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

	cluster := gocql.NewCluster("cassandra") // Use service name "cassandra"
	cluster.Port = 9042
	cluster.Keyspace = "system" // Initial connection to the 'system' keyspace
	cluster.Consistency = gocql.Quorum
	cluster.ConnectTimeout = 10 * time.Second
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 5}

	for {
		Session, err = cluster.CreateSession()
		if err != nil {
			log.Printf("Unable to connect to Cassandra: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	log.Println("Connected to Cassandra")

	// Create keyspace if not exists
	keyspace := "user_db" // Replace with your actual keyspace
	query := fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`, keyspace)
	if err := Session.Query(query).Exec(); err != nil {
		log.Fatalf("Failed to create keyspace: %v", err)
	}

	// Connect to the keyspace
	cluster.Keyspace = keyspace
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatalf("Unable to connect to Cassandra keyspace '%s': %v", keyspace, err)
	}
	log.Printf("Connected to keyspace '%s'", keyspace)
}

func CloseSession() {
	if Session != nil {
		Session.Close()
	}
}
