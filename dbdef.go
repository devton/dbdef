package main

import (
	"github.com/devton/dbdef/config"
	"github.com/devton/dbdef/relation"
	"github.com/devton/dbdef/repository"
	log "github.com/sirupsen/logrus"
)

func main() {
	conf := config.New()
	config.Load(conf)

	repo := repository.New(conf.GetRepositoryURL())
	defer repo.Close()
	log.Debugf("main(): repo=%+v", repo)

	relation.Start(conf.BasePath, repo)
}
