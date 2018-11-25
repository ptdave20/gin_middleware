package mongo

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

const ginMongoKey = "_mongo"

//Config is used to initialize the middleware
type Config struct {
	Host        []string `json:"host"`
	UseAuth     bool     `json:"useAuth"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Database    string   `json:"database"`
	Clone       bool     `json:"clone"` // Are we going to create a fresh instance or clone an existing one
	FailOnIssue bool     `json:"failOnIssue"`
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
	var dInfo mgo.DialInfo
	dInfo.Addrs = config.Host
	dInfo.Database = config.Database
	if config.UseAuth {
		dInfo.Username = config.Username
		dInfo.Password = config.Password
	}
	if config.Clone {
		session, err := mgo.DialWithInfo(&dInfo)
		if err != nil && config.FailOnIssue {
			log.Fatalln(err)
		}
		return func(c *gin.Context) {
			cSession := session.Clone()

			c.Set(ginMongoKey, cSession.DB(dInfo.Database))
			defer cSession.Close()
			c.Next()
		}
	}

	return func(c *gin.Context) {
		session, err := mgo.DialWithInfo(&dInfo)
		if err != nil && config.FailOnIssue {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		defer session.Close()
		c.Set(ginMongoKey, session.DB(dInfo.Database))
		c.Next()
	}
}
