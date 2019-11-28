package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

type Config struct {
	ConnectionUrl string `toml:"connection_url"`
	BasePath      string `toml:"base_path"`
}

type Relation struct {
	Schema     string         `db:"schema"`
	Relation   string         `db:"relation"`
	Kind       string         `db:"kind"`
	Definition pgtype.Varchar `db:"definition"`
	Signature  pgtype.Varchar `db:"signature"`
}

func main() {
	var conf Config
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(conf.BasePath); os.IsNotExist(err) {
		os.Mkdir(conf.BasePath, 0755)
	}

	conn, err := pgx.Connect(context.Background(), conf.ConnectionUrl)
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
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.MkdirAll(path, 0755)
		}
		file := fmt.Sprintf("%s/%s.sql", path, r.Relation)
		ioutil.WriteFile(file, []byte(r.Definition.String), 0755)
	}

}
