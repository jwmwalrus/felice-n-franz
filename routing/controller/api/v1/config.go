package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jwmwalrus/felice-franz/base"
)

func getGlobalConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"port": base.Conf.Port})
}

func getEnvConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"config": "TBD"})
}

func getEnvsList(c *gin.Context) {
	envs := base.Conf.GetEnvsList()
	c.JSON(http.StatusOK, gin.H{"envs": envs})
}

func getEnvTopics(c *gin.Context) {
	envName := c.Param("envname")
	topics := base.Conf.GetEnvTopics(envName)
	c.JSON(http.StatusOK, gin.H{"topics": topics})
}

func getEnvGroups(c *gin.Context) {
	envName := c.Param("envname")
	groups := base.Conf.GetEnvGroups(envName)
	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

func getEnvGroupTopics(c *gin.Context) {
	envName := c.Param("envname")
	groupName := c.Param("groupname")
	topics := base.Conf.GetEnvGroupTopics(envName, groupName)
	c.JSON(http.StatusOK, gin.H{"topics": topics})
}
