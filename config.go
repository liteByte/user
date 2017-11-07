package main

type Config struct {

	// Env vars

	Env       string `default:"develop" env:"ENV"`
	Port      string `default:"5001" env:"PORT"`
	JwtSecret string `required:"true" env:"JWT_SECRET"`

	DB struct {
		Host     string `required:"true" env:"DB_HOST"`
		Port     string `required:"true" env:"DB_PORT"`
		Name     string `required:"true" env:"DB_NAME"`
		Username string `required:"true" env:"DB_USERNAME"`
		Password string `required:"true" env:"DB_PASSWORD"`
	}

	AWS struct {
		AccessKeyID     string `required:"true" env:"AWS_ACCESS_KEY_ID"`
		SecretAccessKey string `required:"true" env:"AWS_SECRET_ACCESS_KEY"`
		Bucket          string `required:"true" env:"AWS_BUCKET"`
		Region          string `required:"true" env:"AWS_REGION"`
	}

	//SendgridApiKey string `required:"true" env:"SENDGRID_API_KEY"`

	// Not env vars

	AppName string

	MaxConnectionsAllowed int `default:"10"`

	Cors struct {
		AllowedMethods string
		AllowedHeaders string
	}

	Login struct {
		TokenDurationInDays uint
	}

	Password struct {
		Min uint
		Max uint
	}
}
