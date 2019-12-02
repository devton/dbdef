package relation

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/devton/dbdef/config"
	"github.com/devton/dbdef/repository"
	"github.com/jackc/pgtype"
	log "github.com/sirupsen/logrus"
)

// Relation holds the information about definitions of all relations inside a database
type Relation struct {
	Dbname      string         `db:"dbname"`
	Schema      string         `db:"schema"`
	Relation    string         `db:"relation"`
	Kind        string         `db:"kind"`
	Description pgtype.Varchar `db:"description"`
	Definition  pgtype.Varchar `db:"definition"`
	Signature   pgtype.Varchar `db:"signature"`
}

// GetDirPath returns the path for relation
func (r *Relation) GetDirPath() string {
	return fmt.Sprintf("%s/schemas/%s/%s/%s", r.Dbname, r.Schema, r.Kind, r.Relation)
}

// WriteDefinitionToFile write relation definition to file inside base path
func (r *Relation) WriteDefinitionFile(path string) {
	file := fmt.Sprintf("%s/definition.sql", path)
	contextLog := log.WithFields(log.Fields{
		"file": file,
	})
	contextLog.Debugf("WriteDefinitionFile():")
	err := ioutil.WriteFile(file, []byte(r.Definition.String), 0755)
	if err != nil {
		contextLog.Fatal(err)
	}
}

// WriteFirstReadme write readme with comment of relation
func (r *Relation) WriteFirstReadme(path string) {
	file := fmt.Sprintf("%s/readme.md", path)
	contextLog := log.WithFields(log.Fields{
		"file": file,
	})

	if _, err := os.Stat(file); os.IsNotExist(err) {
		contextLog.Info("WriteFirstReadme():")
		err := ioutil.WriteFile(file, []byte(r.Description.String), 0755)
		if err != nil {
			contextLog.Fatal(err)
		}
	}
}

// createDirAll create dir
func createDirAll(path string) {
	contextLog := log.WithFields(log.Fields{
		"dir": path,
	})
	if _, err := os.Stat(path); os.IsNotExist(err) {
		contextLog.Debugf("createDirAll():")
		os.MkdirAll(path, 0755)
	}
}

// Start creates structure from database
func Start(conf *config.Config, repo repository.Repository) {
	createDirAll(conf.BasePath)
	c, err := repo.GetConn()
	if err != nil {
		log.Fatalf("relation.Start(): repo.GetConn() err=%w", err)
	}
	defer c.Conn.Release()

	var sqlToRun string
	if conf.Repository.MajorVersion <= 10 {
		sqlToRun = SQLStructPG10
	} else {
		sqlToRun = SQLStructPG11
	}

	rows, err := c.Conn.Query(context.Background(), sqlToRun, conf.Repository.SchemasFilter)
	if err != nil {
		log.Fatalf("relation.Start(): err=%w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r Relation
		err := rows.Scan(&r.Dbname, &r.Schema, &r.Relation, &r.Kind, &r.Description, &r.Definition, &r.Signature)
		if err != nil {
			log.Fatal(err)
		}
		path := fmt.Sprintf("%s/%s", conf.BasePath, r.GetDirPath())
		createDirAll(path)
		r.WriteDefinitionFile(path)
		r.WriteFirstReadme(path)
	}
}
