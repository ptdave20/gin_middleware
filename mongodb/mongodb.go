package mongo

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

const ginMongoKey = "_mongo"

//Config is used to initialize the middleware
type Config struct {
	Host     string `json:"host"`
	UseAuth  bool   `json:"useAuth"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Clone    bool   `json:"clone"`
}

// GetMongo takes a context and returns the database or nil if it doesn't exist
func GetMongo(c *gin.Context) *mgo.Database {
	v, exists := c.Get(ginMongoKey)
	if exists {
		return v.(*mgo.Database)
	}
	return nil
}

// InjectMongo takes a configuration and injects MGO into the your context
func InjectMongo(config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
