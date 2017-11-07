package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type EndpointHandler struct {
	usecaseHandler Usecase
}

func (h *EndpointHandler) Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.defaultSignup(c)
	}
}

func (h *EndpointHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.defaultLogin(c)
	}
}

func (h *EndpointHandler) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.defaultPost(c)
	}
}

func (h *EndpointHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.defaultGet(c)
	}
}

func (h *EndpointHandler) Put() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.defaultPut(c)
	}
}

func (h *EndpointHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.defaultDelete(c)
	}
}

func (h *EndpointHandler) GetOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.defaultGetOne(c)
	}
}

func (h *EndpointHandler) PutOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.defaultPutOne(c)
	}
}

func (h *EndpointHandler) DeleteOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		h.defaultDeleteOne(c)
	}
}

func (h *EndpointHandler) defaultSignup(c *gin.Context) {

	var err error

	var email string
	var password string
	var name string
	var age uint
	var number int
	var date time.Time

	if _, ok := c.GetPostForm("email"); !ok {
		ErrorReply(c, http.StatusBadRequest, "Parameter email missing")
		return
	}
	if _, ok := c.GetPostForm("password"); !ok {
		ErrorReply(c, http.StatusBadRequest, "Parameter password missing")
		return
	}
	if _, ok := c.GetPostForm("name"); !ok {
		ErrorReply(c, http.StatusBadRequest, "Parameter name missing")
		return
	}
	if _, ok := c.GetPostForm("age"); !ok {
		ErrorReply(c, http.StatusBadRequest, "Parameter age missing")
		return
	}
	if _, ok := c.GetPostForm("number"); !ok {
		ErrorReply(c, http.StatusBadRequest, "Parameter number missing")
		return
	}
	if _, ok := c.GetPostForm("date"); !ok {
		ErrorReply(c, http.StatusBadRequest, "Parameter date missing")
		return
	}

	if param, ok := c.GetPostForm("email"); ok {
		email = param
	}
	if param, ok := c.GetPostForm("password"); ok {
		password = param
	}
	if param, ok := c.GetPostForm("name"); ok {
		name = param
	}
	if param, ok := c.GetPostForm("age"); ok {
		age, err = ParseAgeFromString(param)
		if err != nil {
			ErrorReply(c, http.StatusBadRequest, "Invalid value for age")
			return
		}
	}
	if param, ok := c.GetPostForm("number"); ok {
		number, err = ParseNumberFromString(param)
		if err != nil {
			ErrorReply(c, http.StatusBadRequest, "Invalid value for number")
			return
		}
	}
	if param, ok := c.GetPostForm("date"); ok {
		date, err = ParseDateFromString(param)
		if err != nil {
			ErrorReply(c, http.StatusBadRequest, "Invalid value for date")
			return
		}
	}

	model, err := h.usecaseHandler.Create(email, password, name, age, number, date)
	if err != nil {
		if v, ok := err.(Error); ok {
			c.JSON(v.Code, gin.H{"msg": v.Message})
		} else {
			panic(err)
		}
	}

	c.JSON(http.StatusCreated, gin.H{"model": model})
}

func (h *EndpointHandler) defaultLogin(c *gin.Context) {

	var err error

	var email string
	var password string

	if _, ok := c.GetPostForm("email"); !ok {
		ErrorReply(c, http.StatusBadRequest, "Parameter email missing")
		return
	}
	if _, ok := c.GetPostForm("password"); !ok {
		ErrorReply(c, http.StatusBadRequest, "Parameter password missing")
		return
	}

	if param, ok := c.GetPostForm("email"); ok {
		email = param
	}
	if param, ok := c.GetPostForm("password"); ok {
		password = param
	}

	token, user, err := h.usecaseHandler.Login(email, password)
	if err != nil {
		if v, ok := err.(Error); ok {
			c.JSON(v.Code, gin.H{"msg": v.Message})
		} else {
			panic(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

func (h *EndpointHandler) defaultPost(c *gin.Context) {

	var err error

	var email string
	var password string
	var name string
	var age uint
	var number int
	var date time.Time

	if _, ok := c.GetPostForm("email"); !ok {
		ErrorReply(c, http.StatusBadRequest, "Parameter email missing")
		return
	}
	if _, ok := c.GetPostForm("password"); !ok {
		ErrorReply(c, http.StatusBadRequest, "Parameter password missing")
		return
	}
	if _, ok := c.GetPostForm("name"); !ok {
		ErrorReply(c, http.StatusBadRequest, "Parameter name missing")
		return
	}
	if _, ok := c.GetPostForm("age"); !ok {
		ErrorReply(c, http.StatusBadRequest, "Parameter age missing")
		return
	}
	if _, ok := c.GetPostForm("number"); !ok {
		ErrorReply(c, http.StatusBadRequest, "Parameter number missing")
		return
	}
	if _, ok := c.GetPostForm("date"); !ok {
		ErrorReply(c, http.StatusBadRequest, "Parameter date missing")
		return
	}

	if param, ok := c.GetPostForm("email"); ok {
		email = param
	}
	if param, ok := c.GetPostForm("password"); ok {
		password = param
	}
	if param, ok := c.GetPostForm("name"); ok {
		name = param
	}
	if param, ok := c.GetPostForm("age"); ok {
		age, err = ParseAgeFromString(param)
		if err != nil {
			ErrorReply(c, http.StatusBadRequest, "Invalid value for age")
			return
		}
	}
	if param, ok := c.GetPostForm("number"); ok {
		number, err = ParseNumberFromString(param)
		if err != nil {
			ErrorReply(c, http.StatusBadRequest, "Invalid value for number")
			return
		}
	}
	if param, ok := c.GetPostForm("date"); ok {
		date, err = ParseDateFromString(param)
		if err != nil {
			ErrorReply(c, http.StatusBadRequest, "Invalid value for date")
			return
		}
	}

	model, err := h.usecaseHandler.Create(email, password, name, age, number, date)
	if err != nil {
		if v, ok := err.(Error); ok {
			c.JSON(v.Code, gin.H{"msg": v.Message})
		} else {
			panic(err)
		}
	}

	c.JSON(http.StatusCreated, gin.H{"model": model})
}

func (h *EndpointHandler) defaultGet(c *gin.Context) {

	filter := c.MustGet("filter").([]map[string]string)
	order := c.MustGet("order").(map[string]string)
	offset := c.MustGet("offset").(int)
	limit := c.MustGet("limit").(int)

	models, err := h.usecaseHandler.Find(filter, order, offset, limit)
	if err != nil {
		if v, ok := err.(Error); ok {
			c.JSON(v.Code, gin.H{"msg": v.Message})
		} else {
			panic(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"models": models})
}

func (h *EndpointHandler) defaultPut(c *gin.Context) {

	updates := map[string]interface{}{}

	if param, ok := c.GetPostForm("name"); ok {
		updates["Name"] = param
	}
	if param, ok := c.GetPostForm("age"); ok {
		age, err := ParseAgeFromString(param)
		if err != nil {
			ErrorReply(c, http.StatusBadRequest, "Invalid value for age")
			return
		}
		updates["Age"] = age
	}
	if param, ok := c.GetPostForm("number"); ok {
		number, err := ParseNumberFromString(param)
		if err != nil {
			ErrorReply(c, http.StatusBadRequest, "Invalid value for number")
			return
		}
		updates["Number"] = number
	}
	if param, ok := c.GetPostForm("param"); ok {
		date, err := ParseDateFromString(param)
		if err != nil {
			ErrorReply(c, http.StatusBadRequest, "Invalid value for date")
			return
		}
		updates["Date"] = date
	}

	filter := c.MustGet("filter").([]map[string]string)

	err := h.usecaseHandler.Update(updates, filter)
	if err != nil {
		if v, ok := err.(Error); ok {
			c.JSON(v.Code, gin.H{"msg": v.Message})
		} else {
			panic(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *EndpointHandler) defaultDelete(c *gin.Context) {

	filter := c.MustGet("filter").([]map[string]string)

	err := h.usecaseHandler.Delete(filter)
	if err != nil {
		if v, ok := err.(Error); ok {
			c.JSON(v.Code, gin.H{"msg": v.Message})
		} else {
			panic(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *EndpointHandler) defaultGetOne(c *gin.Context) {
	model := c.MustGet("one")
	c.JSON(http.StatusOK, gin.H{"model": model})
}

func (h *EndpointHandler) defaultPutOne(c *gin.Context) {
	model := c.MustGet("one").(*Model)

	params := struct {
		name   string
		age    string
		number string
		date   string
	}{
		c.PostForm("name"),
		c.PostForm("age"),
		c.PostForm("number"),
		c.PostForm("date"),
	}

	updates := map[string]interface{}{}

	if params.name != "" {
		updates["Name"] = params.name
	}
	if params.age != "" {
		age, err := ParseAgeFromString(params.age)
		if err != nil {
			ErrorReply(c, http.StatusBadRequest, "Invalid value for age")
			return
		}
		updates["Age"] = age
	}
	if params.number != "" {
		number, err := ParseNumberFromString(params.number)
		if err != nil {
			ErrorReply(c, http.StatusBadRequest, "Invalid value for number")
			return
		}
		updates["Number"] = number
	}
	if params.date != "" {
		date, err := ParseDateFromString(params.date)
		if err != nil {
			ErrorReply(c, http.StatusBadRequest, "Invalid value for date")
			return
		}
		updates["Date"] = date
	}

	model, err := h.usecaseHandler.UpdateOne(model, updates)
	if err != nil {
		if v, ok := err.(Error); ok {
			c.JSON(v.Code, gin.H{"msg": v.Message})
		} else {
			panic(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"model": model})
}

func (h *EndpointHandler) defaultDeleteOne(c *gin.Context) {
	model := c.MustGet("one").(*Model)

	err := h.usecaseHandler.DeleteOne(model)
	if err != nil {
		if v, ok := err.(Error); ok {
			c.JSON(v.Code, gin.H{"msg": v.Message})
		} else {
			panic(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}
