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

type Relation struct {
	Schema     string         `db:"schema"`
	Relation   string         `db:"relation"`
	Kind       string         `db:"kind"`
	Definition pgtype.Varchar `db:"definition"`
	Signature  pgtype.Varchar `db:"signature"`
}

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

func New(basePath string, repo repository.Repository) {
	c, err := repo.GetConn()
	if err != nil {
		log.Fatalf("relation.New(): repo.GetConn() err=%w", err)
	}
	defer c.Conn.Release()
	rows, err := c.Conn.Query(context.Background(), SQLStruct)
	if err != nil {
		log.Fatalf("relation.New(): err=%w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r Relation
		err := rows.Scan(&r.Schema, &r.Relation, &r.Kind, &r.Definition, &r.Signature)
		if err != nil {
			log.Fatal(err)
		}

		path := fmt.Sprintf("%s/%s/%s", basePath, r.Schema, r.Kind)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.MkdirAll(path, 0755)
		}
		r.WriteDefinitionToFile(path)
	}
}
