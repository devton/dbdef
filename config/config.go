package config

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	configFile        string
	defaultConfigFile = "./config.toml"
)

// RepositoryConfig a config to repository
type RepositoryConfig struct {
	URL             string
	SchemasFilter   string
	TablesFilter    string
	ViewsFilter     string
	FunctionsFilter string
}

// Config basic config
type Config struct {
	Repository *RepositoryConfig
	Dev        bool
	Trace      bool
	BasePath   string
}

func getEnvConfig(env string) (cfg string) {
	cfg = os.Getenv(env)
	return
}

func getDefaultConfig(file string) (fileConfig string) {
	fileConfig = defaultConfigFile
	if file != "" {
		fileConfig = file
	}

	_, err := os.Stat(fileConfig)
	if err != nil {
		fileConfig = ""
	}

	return
}

func viperCfg() {
	configFile = getDefaultConfig(getEnvConfig("DBDEF_CONFIG"))
	dir, file := filepath.Split(configFile)
	file = strings.TrimSuffix(file, filepath.Ext(file))
	viper.AddConfigPath(dir)
	viper.SetConfigName(file)
	viper.SetConfigType("toml")
	viper.SetDefault("dbdef.dev", true)
	viper.SetDefault("dbdef.trace", false)
	viper.SetDefault("dbdef.base_path", "dbdef_report")
}

// parse Config configs
func parse(cfg *Config) (err error) {
	err = viper.ReadInConfig()
	if err != nil {
		log.Errorf("config.Parse(): error=%w", err)
		return
	}

	cfg.Dev = viper.GetBool("dbdef.dev")
	cfg.Trace = viper.GetBool("dbdef.trace")
	cfg.BasePath = viper.GetString("dbdef.base_path")
	cfg.Repository.URL = viper.GetString("repository.url")
	cfg.Repository.SchemasFilter = viper.GetString("repository.schemas_filter")
	cfg.Repository.TablesFilter = viper.GetString("repository.tables_filter")
	cfg.Repository.ViewsFilter = viper.GetString("repository.views_filter")
	cfg.Repository.FunctionsFilter = viper.GetString("repository.functions_filter")

	return
}

func logConfig(cfg *Config) {
	log.SetReportCaller(false)
	log.SetLevel(log.InfoLevel)

	if cfg.Trace {
		log.SetReportCaller(true)
		log.Debug("init(): trace enabled")
	}

	if cfg.Dev {
		log.SetLevel(log.DebugLevel)
		log.Debug("init(): dev environment")
	}

}

// New initialize the basic config
func New() *Config {
	return &Config{
		Repository: &RepositoryConfig{},
	}
}

// Load configuration
func Load(c *Config) {
	viperCfg()

	if err := parse(c); err != nil {
		log.Fatalf("config.Load(): Parse(ConfigConf): %w", err)
	}

	logConfig(c)

	log.Debug("config.Load(): configuration loaded")
	log.Debugf("config.ConfigConf=%+v\n", c)
}

// GetRepositoryURL get url from config
func (c *Config) GetRepositoryURL() string {
	return c.Repository.URL
}

// GetRepositorySchemasFilter get schemas_filter from config
func (c *Config) GetRepositorySchemasFilter() string {
	return c.Repository.SchemasFilter
}
