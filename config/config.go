// config/config.go
package config

import (
	"github.com/gocql/gocql"
	"log"
	"time"
)

var Session *gocql.Session

func InitCassandra() {
	var err error

	cluster := gocql.NewCluster("cassandra") // Use service name "cassandra"
	cluster.Port = 9042
	cluster.Keyspace = "your_keyspace" // Replace with your keyspace name
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
}

func CloseSession() {
	if Session != nil {
		Session.Close()
	}
}
