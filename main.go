package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gocql/gocql"
)

var cassandraAddress = "localhost"

const cassandraKeyspace = "cognito"

func main() {
	// Set up keyspace and tables
	setupCassandra()

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	http.ListenAndServe(":3000", r)
}

func setupCassandra() {
	setupCassandraKeyspace()
	setupCassandraTables()
}

func setupCassandraKeyspace() {
	cluster := gocql.NewCluster(cassandraAddress)
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	if err := session.Query(fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = { 
    'class' : 'SimpleStrategy', 
    'replication_factor' : 1 
  }`, cassandraKeyspace)).Exec(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created Keyspace " + cassandraKeyspace)
}

func setupCassandraTables() {
	cluster := gocql.NewCluster(cassandraAddress)
	cluster.Keyspace = cassandraKeyspace
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	if err := session.Query(`CREATE TABLE IF NOT EXISTS user_pool
  (
    pool_id TEXT,
    PRIMARY KEY (pool_id)
  );`).Exec(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created User Pool Table")

	if err := session.Query(`CREATE TABLE IF NOT EXISTS user
  (
    username TEXT,
    pool_id TEXT,
    email TEXT,
    phone TEXT,
    password TEXT,
    salt TEXT,
    PRIMARY KEY (username, pool_id)
  );`).Exec(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created User Table")

	if err := session.Query(`CREATE TABLE IF NOT EXISTS user_attribute
  (
    username TEXT,
    attribute TEXT,
    value TEXT,
    PRIMARY KEY (username)
  );`).Exec(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created User Attribute Table")

}
