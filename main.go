package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/kevindoveton/cognito-local/cognito"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gocql/gocql"
)

var cassandraAddress = "localhost"

const cassandraKeyspace = "cognito"

func main() {
	// Set up keyspace and tables
	var session = setupCassandra()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var requestType = r.Header.Get("X-Amz-Target")

		fmt.Println(requestType)

		if requestType == "AWSCognitoIdentityProviderService.CreateUserPool" {
			var	up cognito.UserPool
			json.NewDecoder(r.Body).Decode(&up)
			cognito.CreateUserPool(session, up)
		} else if requestType == "AWSCognitoIdentityProviderService.DeleteUserPool" {
			var	up cognito.UserPool
			json.NewDecoder(r.Body).Decode(&up)
			cognito.DeleteUserPool(session, up)
		} else if requestType == "AWSCognitoIdentityProviderService.AdminCreateUser" {
			var	u cognito.User
			json.NewDecoder(r.Body).Decode(&u)
			cognito.CreateUser(session, u)
	} else if requestType == "AWSCognitoIdentityProviderService.AdminDeleteUser" {
		var	u cognito.User
		json.NewDecoder(r.Body).Decode(&u)
		cognito.DeleteUser(session, u)
	}


	w.Write([]byte("Success"))
	})


	http.ListenAndServe(":3000", r)
}

func setupCassandra() *gocql.Session {
	setupCassandraKeyspace()

	// Create a new session to use through app
	cluster := gocql.NewCluster(cassandraAddress)
	cluster.Consistency = gocql.Quorum
	cluster.Keyspace = cassandraKeyspace
	session, _ := cluster.CreateSession()

	// create needed tables
	setupCassandraTables(session)

	return session
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

func setupCassandraTables(session *gocql.Session) {
	if err := session.Query(`CREATE TABLE IF NOT EXISTS user_pool
  (
    pool_id TEXT,
		pool_name TEXT,
    PRIMARY KEY (pool_name)
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
