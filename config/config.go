// config/config.go
package config

import (
	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var Session *gocql.Session

func InitCassandra() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_HOST"))
	cluster.Keyspace = os.Getenv("CASSANDRA_KEYSPACE")
	cluster.Consistency = gocql.Quorum

	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal("Unable to connect to Cassandra:", err)
	}
}

func CloseSession() {
	Session.Close()
}
