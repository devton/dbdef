package relation

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/devton/dbdef/repository"
	"github.com/jackc/pgtype"
	log "github.com/sirupsen/logrus"
)

// Relation holds the information about definitions of all relations inside a database
type Relation struct {
	Schema     string         `db:"schema"`
	Relation   string         `db:"relation"`
	Kind       string         `db:"kind"`
	Definition pgtype.Varchar `db:"definition"`
	Signature  pgtype.Varchar `db:"signature"`
}

// WriteDefinitionToFile write relation definition to file inside base path
func (r *Relation) WriteDefinitionToFile(path string) {
	file := fmt.Sprintf("%s/%s.sql", path, r.Relation)
	contextLog := log.WithFields(log.Fields{
		"file": file,
	})
	contextLog.Info("writeFile():")
	err := ioutil.WriteFile(file, []byte(r.Definition.String), 0755)
	if err != nil {
		contextLog.Fatal(err)
	}
}

// createDirAll create dir
func createDirAll(path string) {
	contextLog := log.WithFields(log.Fields{
		"path": path,
	})
	if _, err := os.Stat(path); os.IsNotExist(err) {
		contextLog.Info("createDirAll():")
		os.MkdirAll(path, 0755)
	}
}

// Start creates structure from database
func Start(basePath string, repo repository.Repository) {
	createDirAll(basePath)
	c, err := repo.GetConn()
	if err != nil {
		log.Fatalf("relation.Start(): repo.GetConn() err=%w", err)
	}
	defer c.Conn.Release()
	rows, err := c.Conn.Query(context.Background(), SQLStruct)
	if err != nil {
		log.Fatalf("relation.Start(): err=%w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r Relation
		err := rows.Scan(&r.Schema, &r.Relation, &r.Kind, &r.Definition, &r.Signature)
		if err != nil {
			log.Fatal(err)
		}
		path := fmt.Sprintf("%s/%s/%s", basePath, r.Schema, r.Kind)
		createDirAll(path)
		r.WriteDefinitionToFile(path)
	}
}
