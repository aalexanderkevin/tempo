package container

import (
	"database/sql"

	"tempo/config"
	"tempo/repository"
)

type Container struct {
	db     *sql.DB
	config config.Config

	// repo
	userRepo repository.User
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) Db() *sql.DB {
	return c.db
}

func (c *Container) SetDb(db *sql.DB) {
	c.db = db
}

func (c *Container) Config() config.Config {
	return c.config
}

func (c *Container) SetConfig(config config.Config) {
	c.config = config
}

func (c *Container) UserRepo() repository.User {
	return c.userRepo
}

func (c *Container) SetUserRepo(userRepo repository.User) {
	c.userRepo = userRepo
}
