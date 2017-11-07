package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {

	config := initConfig()
	db := initDatabase(config)
	//awsSession := initAWS()
	//TODO add recovery handler
	if config.Env != "develop" {
		gin.SetMode(gin.ReleaseMode)
	}

	persistenceHandler := PersistenceHandler{db}
	usecaseHandler := UsecaseHandler{&persistenceHandler, config}
	endpointHandler := EndpointHandler{&usecaseHandler}

	router := gin.New()

	router.POST("/signup", endpointHandler.Signup())
	router.POST("/login", endpointHandler.Login())

	auth := router.Group("/", Authenticate(config))
	{
		//auth.GET("/", AuthenticatedID(), FindOne(db), endpointHandler.GetOne())
		//auth.PUT("/", AuthenticatedID(), FindOne(db), endpointHandler.PutOne())
		//auth.DELETE("/", AuthenticatedID(), FindOne(db), endpointHandler.DeleteOne())

		admin := router.Group("/") //TODO add auth middleware
		{
			admin.POST("/", endpointHandler.Post())

			auth.GET("/:id", GetID(), FindOne(db), endpointHandler.GetOne())
			auth.PUT("/:id", GetID(), FindOne(db), endpointHandler.PutOne())
			auth.DELETE("/:id", GetID(), FindOne(db), endpointHandler.DeleteOne())

			admin.GET("/", Filter(), Order(), Paginate(), endpointHandler.Get())
			admin.PUT("/", Filter(), endpointHandler.Put())
			admin.DELETE("/", Filter(), endpointHandler.Delete())
		}
	}

	router.Run(":" + config.Port)
}

func initConfig() *Config {
	config := Config{}
	err := configor.Load(&config, "config.yml")
	if err != nil {
		panic(err)
	}
	return &config
}

func initDatabase(config *Config) *gorm.DB {
	name := config.DB.Name
	host := config.DB.Host
	port := config.DB.Port
	user := config.DB.Username
	pass := config.DB.Password

	db, err := gorm.Open("mysql", buildMySQLConnectionString(host, port, name, user, pass))
	if err != nil {
		panic(err)
	}
	return db
}

func buildMySQLConnectionString(host, port, name, user, pass string) string {
	str := user + ":" + pass + "@" + "tcp(" + host + ":" + port + ")" + "/" + name + "?"
	str += "charset=utf8&parseTime=True&loc=Local"
	return str
}

func initAWS() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}
