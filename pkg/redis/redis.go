package rediswrap

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/go-redis/redis/v8"
	"github.com/gobuffalo/envy"
	"gopkg.in/yaml.v2"
)

// Client is a redis client to interact with redis server
var (
	Client          *redis.Client
	_configFilePath string
)

func init() {
	var err error
	env := envy.Get("GO_ENV", "development")
	Client, err = Connect(env)
	if err != nil {
		log.Fatal(err)
	}
}

// ConnectionDetails is to define redis config struct
type ConnectionDetails struct {
	Host       string `yaml:"host"` // localhost:6379
	Port       int    `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	DB         int    `yaml:"db"`
	PoolSize   int    `yaml:"pool"`
	MaxRetries int    `yaml:"max_retries"`
}

var connections = map[string]*ConnectionDetails{}

// GetConnections is to get redis connections config
func GetConnections() map[string]*ConnectionDetails {
	return connections
}

// Connect accepts an environment variable to connect with corresponding configuration
func Connect(env string) (*redis.Client, error) {
	// LoadConfigFile first
	if err := LoadConfigFile(); err != nil {
		return nil, err
	}

	ctx := context.Background()

	config := connections[env]

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// findConfigPath determines the config file path where should be read from
func findConfigPath() string {
	path := envy.Get("REDIS_CONFIG_PATH", "")
	if path == "" {
		path = filepath.Join(filepath.Dir(filepath.Dir(os.Args[0])), "redis.yml")
	}
	return path
}

// LoadConfigFile return the redis connections config
func LoadConfigFile() error {
	envy.Load()
	path := findConfigPath()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return unmarshalConfig(f)
}

func tmplProcess(content []byte) ([]byte, error) {
	tmpl := template.New("_redis_config_transformer")
	tmpl.Funcs(map[string]interface{}{
		"envOr": func(s1, s2 string) string {
			return envy.Get(s1, s2)
		},
		"env": func(s1 string) string {
			return envy.Get(s1, "")
		},
	})

	t, err := tmpl.Parse(string(content))
	if err != nil {
		return nil, err
	}

	var bb bytes.Buffer
	err = t.Execute(&bb, nil)
	if err != nil {
		return nil, err
	}

	return bb.Bytes(), nil
}

// unmarshalConfig unmarshal the file into the connection structure
func unmarshalConfig(r io.Reader) error {

	content, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	b, err := tmplProcess(content)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, &connections)
	return err
}
