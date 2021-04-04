package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jwmwalrus/felice-franz/base"
	"github.com/jwmwalrus/felice-franz/kafkian"
)

type ConsumersReqJSON struct {
	Env       string   `json:"env"`
	TopicKeys []string `json:"topicKeys"`
}

func consumersEndPoints(r *gin.RouterGroup) {
	ce := r.Group("/consumers")
	{
		ce.POST("/", postConsumers)
	}
}

func postConsumers(c *gin.Context) {
	b := &ConsumersReqJSON{}
	c.BindJSON(b)

	env := base.Conf.GetEnvConfig(b.Env)
	topics, err := env.FindTopicValues(b.TopicKeys)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := kafkian.CreateConsumers(env, topics); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
