package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"
)

func Authenticate(config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		defaultAuthenticate(c, config)
	}
}

func Filter() gin.HandlerFunc {
	return func(c *gin.Context) {
		defaultFilter(c)
	}
}

func AuthenticatedID() gin.HandlerFunc {
	return func(c *gin.Context) {
		defaultAuthenticatedID(c)
	}
}

func GetID() gin.HandlerFunc {
	return func(c *gin.Context) {
		defaultGetID(c)
	}
}

func FindOne(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		defaultFindOne(c, db)
	}
}

func Order() gin.HandlerFunc {
	return func(c *gin.Context) {
		defaultOrder(c)
	}
}

func Paginate() gin.HandlerFunc {
	return func(c *gin.Context) {
		defaultPaginate(c)
	}
}

func defaultAuthenticate(c *gin.Context, config *Config) {

	authorizationHeader := strings.ToLower(c.Request.Header.Get("Authorization"))
	if strings.HasPrefix(authorizationHeader, "bearer ") == false {
		ErrorReply(c, http.StatusUnauthorized, "")
		return
	}
	tokenString := strings.TrimPrefix(authorizationHeader, "bearer ")

	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JwtSecret), nil
	})
	if err != nil {
		panic(err)
	}

	claims, ok := token.Claims.(*JWTCustomClaims)
	if ok == false || token.Valid == false {
		ErrorReply(c, http.StatusUnauthorized, "")
		return
	}

	//Lib only checks validity of "exp, iat, nbf" and only if the claims are present. So...
	if claims.IssuedAt == 0 {
		ErrorReply(c, http.StatusUnauthorized, "")
		return
	}
	if claims.ExpiresAt == 0 {
		ErrorReply(c, http.StatusUnauthorized, "")
		return
	}
	if claims.Issuer != config.AppName {
		ErrorReply(c, http.StatusUnauthorized, "")
		return
	}
	if claims.Subject == "" {
		ErrorReply(c, http.StatusUnauthorized, "")
		return
	}
	if claims.Audience != config.AppName {
		ErrorReply(c, http.StatusUnauthorized, "")
		return
	}

	id, err := strconv.ParseUint(claims.Subject, 10, 32)
	if err != nil {
		ErrorReply(c, http.StatusUnauthorized, "")
		return
	}
	c.Set("authenticatedID", id)

	c.Next()
}

func defaultFilter(c *gin.Context) {

	var queries []map[string]string

	name := c.Query("name")
	if name != "" {
		queries = append(queries, map[string]string{"name = ?": name})
	}
	age := c.Query("age")
	if age != "" {
		age, operator, valid := getFilterOperator([]byte(age))
		if !valid {
			ErrorReply(c, http.StatusBadRequest, "Invalid operator for age")
			return
		}
		_, err := ParseAgeFromString(age)
		if err != nil {
			ErrorReply(c, http.StatusBadRequest, "Invalid value for age")
			return
		}
		queries = append(queries, map[string]string{"age " + operator + " ?": age})
	}
	number := c.Query("number")
	if number != "" {
		number, operator, valid := getFilterOperator([]byte(number))
		if !valid {
			ErrorReply(c, http.StatusBadRequest, "Invalid operator for number")
			return
		}
		_, err := ParseNumberFromString(number)
		if err != nil {
			ErrorReply(c, http.StatusBadRequest, "Invalid value for number")
			return
		}
		queries = append(queries, map[string]string{"number " + operator + " ?": number})
	}
	date := c.Query("date")
	if date != "" {
		date, operator, valid := getFilterOperator([]byte(date))
		if !valid {
			ErrorReply(c, http.StatusBadRequest, "Invalid operator for date")
			return
		}
		_, err := ParseDateFromString(date)
		if err != nil {
			ErrorReply(c, http.StatusBadRequest, "Invalid value for date")
			return
		}
		queries = append(queries, map[string]string{"date " + operator + " ?": date})
	}

	c.Set("filter", queries)

	c.Next()
}

func getFilterOperator(v []byte) (string, string, bool) {

	isOperator := func(b byte) bool {
		return b == '>' || b == '<' || b == '='
	}

	isValid := func(s string) bool {
		if s == ">=" || s == "<=" || s == "<>" {
			return true
		}
		return false
	}

	if len(v) == 1 {
		return string(v), "=", true
	}

	if len(v) >= 1 && isOperator(v[0]) {
		if len(v) >= 2 && isOperator(v[1]) {
			s := string(v[0:2])
			return string(v[2:]), s, isValid(s)
		}
		return string(v[1:]), string(v[0]), true
	}
	return string(v), "=", true
}

func defaultAuthenticatedID(c *gin.Context) {

	c.Set("id", c.MustGet("authenticatedID").(uint))

	c.Next()
}

func defaultGetID(c *gin.Context) {
	param := c.Param("id")
	if param == "" {
		ErrorReply(c, http.StatusBadRequest, "Invalid id")
		return
	}

	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		ErrorReply(c, http.StatusBadRequest, "Invalid id")
		return
	}

	c.Set("id", uint(id))

	c.Next()
}

func defaultFindOne(c *gin.Context, db *gorm.DB) {
	var component Model

	r := db.First(&component, c.MustGet("id"))
	if r.RecordNotFound() {
		ErrorReply(c, http.StatusNotFound, "Not found")
		return
	}
	PanicIf(c, r.Error)

	c.Set("one", &component)

	c.Next()
}

func defaultOrder(c *gin.Context) {
	validFields := map[string]struct{}{
		"id":     {},
		"name":   {},
		"age":    {},
		"number": {},
		"date":   {},
	}

	var orderField string = "id"
	var orderDir string = "DESC"

	orderParam := c.Query("Order")
	orderDirParam := c.Query("OrderDir")

	if orderParam != "" {
		if _, ok := validFields[orderParam]; !ok {
			ErrorReply(c, http.StatusBadRequest, "Invalid order field")
			return
		}
		orderField = orderParam
	}

	if orderDirParam != "" {
		if orderDirParam != "ASC" && orderDirParam != "DESC" {
			ErrorReply(c, http.StatusBadRequest, "Invalid order direction")
			return
		}
		orderDir = orderDirParam
	}

	c.Set("order", map[string]string{orderField: orderDir})

	c.Next()
}

func defaultPaginate(c *gin.Context) {
	var offset int = -1
	var limit int = -1

	limitParam := c.Query("Limit")
	offsetParam := c.Query("Offset")

	if limitParam != "" {
		tmp, err := strconv.ParseInt(limitParam, 10, 32)
		limit = int(tmp)
		if err != nil || !genIsLimitValid(limit) {
			ErrorReply(c, http.StatusBadRequest, "Invalid limit")
			return
		}
	}
	c.Set("limit", limit)

	if limit != -1 && offsetParam != "" {
		tmp, err := strconv.ParseInt(offsetParam, 10, 32)
		offset = int(tmp)
		if err != nil || !genIsOffsetValid(offset) {
			ErrorReply(c, http.StatusBadRequest, "Invalid offset")
			return
		}
	}
	c.Set("offset", offset)

	c.Next()
}

func genIsOffsetValid(o int) bool {
	if o >= 0 {
		return true
	}
	return false
}

func genIsLimitValid(l int) bool {
	if l > 0 {
		return true
	}
	return false
}
