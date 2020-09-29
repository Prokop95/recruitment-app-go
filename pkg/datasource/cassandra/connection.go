package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/pawel_prokop/recruitment-project-go/config"
	"log"
)

func InitCassandra(config *config.Config) *gocql.Session {
	cluster := gocql.NewCluster(config.Cassandra.Host)
	cluster.Consistency = gocql.LocalOne
	var err error
	var session *gocql.Session
	session, err = cluster.CreateSession()
	if err != nil {
		log.Println(err)
		return nil
	}

	createKeyspace(session)
	createTables(session)

	return session
}

func createKeyspace(session *gocql.Session) {
	var err error

	err = session.Query("CREATE KEYSPACE IF NOT EXISTS recruitment WITH REPLICATION = {'class' : 'NetworkTopologyStrategy', 'datacenter1' : 1};").Exec()
	if err != nil {
		log.Println(err)
		return
	}
}

func createTables(session *gocql.Session) {

	err := session.Query("CREATE TABLE IF NOT EXISTS recruitment.message (id uuid, email text, magic_number int, title text, content text, PRIMARY KEY (email, id));").Exec()
	if err != nil {
		log.Println(err)
		return
	}

	err = session.Query("CREATE TABLE IF NOT EXISTS recruitment.message_by_magic_number (id uuid, email text, magic_number int, title text, content text, PRIMARY KEY (magic_number, id));").Exec()
	if err != nil {
		log.Println(err)
		return
	}
}
