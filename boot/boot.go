package boot

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v9"
	"lihood/conf"
	"lihood/g"
	"os"
	"time"
)

// Setup set up your application here
func Setup() error {
	// initialize configuration
	if err := setupConfig(); err != nil {
		return err
	}
	// initialize database
	if err := setupDB(); err != nil {
		return err
	}
	// initialize redis
	if err := setupRedis(); err != nil {
		return err
	}
	return nil
}

// MustSetup set up your application here
func MustSetup() {
	if err := Setup(); err != nil {
		panic(err)
	}
}

// setupConfig initialize configuration
func setupConfig() error {
	file, err := os.Open("config.yaml")
	if err != nil {
		return err
	}
	defer file.Close()
	return conf.Encode(file)
}

// initialize database
func setupDB() error {
	var err error
	g.DB, err = sql.Open("mysql", conf.Instance.DB.DSN())
	if err != nil {
		return err
	}
	if err = g.DB.Ping(); err != nil {
		return err
	}
	g.DB.SetConnMaxLifetime(time.Hour * 2)
	g.DB.SetMaxIdleConns(10)
	g.DB.SetMaxOpenConns(50)
	return nil
}

// initialize redis
func setupRedis() error {
	g.Redis = redis.NewClient(&redis.Options{
		Addr:     conf.Instance.Redis.Addr(),
		Network:  "tcp",
		Password: conf.Instance.Redis.Password,
		DB:       conf.Instance.Redis.DB,
	})
	return g.Redis.Ping(context.Background()).Err()
}
