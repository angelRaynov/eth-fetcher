package authentication

import "github.com/gin-gonic/gin"

type Authenticator interface {
	Authenticate(c *gin.Context)
}
