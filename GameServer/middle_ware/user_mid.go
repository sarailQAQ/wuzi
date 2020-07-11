package middle_ware

import (
	"Server/jwt"
	"Server/resps"
	"errors"
	"github.com/gin-gonic/gin"
)

type UserClaim struct {
	Id uint
	Username string
}

func LoginStatus(c *gin.Context)  {
	auth:= c.GetHeader("token")
	//fmt.Println(auth)
	if len(auth) < 7 {
		resps.Error(c,2,errors.New("Illegal jwt"))
		c.Abort()
		return
	}
	//token := auth[7:]
	uid,user,err := jwt.Check(auth)
	if err != nil {
		resps.Error(c,2,err)
		c.Abort()
		return
	}

	c.Set("user",UserClaim{
		Id:       uid,
		Username: user,
	})

	c.Next()
	return
}
