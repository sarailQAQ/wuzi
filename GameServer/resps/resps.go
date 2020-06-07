package resps

import "github.com/gin-gonic/gin"

func Ok(c *gin.Context){
	c.JSON(200,gin.H{"code":0})
}

func OkWithData(c *gin.Context,data interface{}){
	c.JSON(200,gin.H{"code":0,"data":data})
}

func FormError(c *gin.Context){
	c.JSON(200,gin.H{"code":1,"message":"request form error"})
}

func Error(c *gin.Context,code int,err error){
	c.JSON(200,gin.H{"coed":code,"mesage":err.Error()})
}
