package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/devton/dbdef/config"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

type Relation struct {
	Schema     string         `db:"schema"`
	Relation   string         `db:"relation"`
	Kind       string         `db:"kind"`
	Definition pgtype.Varchar `db:"definition"`
	Signature  pgtype.Varchar `db:"signature"`
}

func createDirAll(path string) {
	contextLog := log.WithFields(log.Fields{
		"path": path,
	})
	if _, err := os.Stat(path); os.IsNotExist(err) {
		contextLog.Info("createDirAll():")
		os.MkdirAll(path, 0755)
	}
}

func (r *Relation) WriteDefinitionToFile(path string) {
	file := fmt.Sprintf("%s/%s.sql", path, r.Relation)
	contextLog := log.WithFields(log.Fields{
		"file": file,
	})
	contextLog.Info("writeFile():")
	ioutil.WriteFile(path, []byte(r.Definition.String), 0755)
}

func main() {
	conf := config.New()
	config.Load(conf)
	createDirAll(conf.BasePath)
	conn, err := pgx.Connect(context.Background(), conf.GetRepositoryURL())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), structureQuery)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var r Relation
		err := rows.Scan(&r.Schema, &r.Relation, &r.Kind, &r.Definition, &r.Signature)
		if err != nil {
			log.Fatal(err)
		}

		path := fmt.Sprintf("%s/%s/%s", conf.BasePath, r.Schema, r.Kind)
		createDirAll(path)
		r.WriteDefinitionToFile(path)
	}

}
