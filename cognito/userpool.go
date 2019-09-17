package cognito

import (
	"github.com/gocql/gocql"
	"github.com/lucasjones/reggen"
	"log"
)

type UserPool struct {
	UserPoolId string
	PoolName string
}

func DeleteUserPool(session *gocql.Session, up UserPool) {
	if err := session.Query(`
	DELETE FROM user_pool
	WHERE pool_id = ?
	`, up.UserPoolId).Exec(); err != nil {
		log.Fatal("Error", err)
	}
}

func CreateUserPool(session *gocql.Session, up UserPool) {
	if err := session.Query(`
	INSERT INTO user_pool (pool_id, pool_name) 
	VALUES (?, ?)
	`, createId(), up.PoolName).Exec(); err != nil {
		log.Fatal("Error", err)
	}
}

func createId() string {
	str, err := reggen.Generate("^[\\w-]+_[0-9a-zA-Z]+$", 10)
	if err != nil {
		panic(err)
	}

	return str
}