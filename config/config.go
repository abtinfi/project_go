// config/config.go
// config/config.go
package config

import (
	"log"
	"os"
	"time"

	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
)

var Session *gocql.Session

func InitCassandra() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_HOST"))
	cluster.Consistency = gocql.Quorum
	cluster.Keyspace = "system" // Temporarily connect to the system keyspace to create user_db if needed

	// Retry mechanism
	for i := 0; i < 10; i++ {
		Session, err = cluster.CreateSession()
		if err == nil {
			break
		}
		log.Printf("Unable to connect to Cassandra: %v. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatal("Unable to connect to Cassandra after retries:", err)
	}

	// Create the keyspace if it doesn't exist
	err = Session.Query(`CREATE KEYSPACE IF NOT EXISTS user_db WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`).Exec()
	if err != nil {
		log.Fatal("Unable to create keyspace:", err)
	}

	// Close session and reconnect to the actual keyspace
	Session.Close()
	cluster.Keyspace = os.Getenv("CASSANDRA_KEYSPACE")
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal("Unable to connect to Cassandra with keyspace:", err)
	}
}

func CloseSession() {
	if Session != nil {
		Session.Close()
	}
}
