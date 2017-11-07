package main

import "github.com/gin-gonic/gin"

func PanicIf(c *gin.Context, err error) {
	if err != nil {
		c.Abort()
		panic(err)
	}
}

func ErrorReply(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{
		"msg": msg,
	})
	c.Abort()
}

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return e.Message
}