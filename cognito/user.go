package cognito

import (
	"github.com/gocql/gocql"
	"log"
)

type User struct {
	Id string
	Username string
	Email string
	Phone string
	UserPoolId string
}

func CreateUser(session *gocql.Session, u User) {
	if err := session.Query(`
	INSERT INTO user (pool_id, username) 
	VALUES (?, ?)
	`, u.UserPoolId, u.Username).Exec(); err != nil {
		log.Fatal("Error", err)
	}
}

func DeleteUser(session *gocql.Session, u User) {
	if err := session.Query(`
	DELETE FROM user
	WHERE pool_id = ? and username = ?
	`, u.UserPoolId, u.Username).Exec(); err != nil {
		log.Fatal("Error", err)
	}
}