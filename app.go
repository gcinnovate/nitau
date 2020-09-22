package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
	// Import godotenv for .env variables
	"github.com/gcinnovate/nitau/controllers"
	"github.com/gcinnovate/nitau/helpers"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

// Config is a application configuration structure
type Config struct {
	API struct {
		Username  string `yaml:"username" env:"NITAU_API_USER" env-description:"API user name"`
		Password  string `yaml:"password" env:"NITAU_API_PASSWORD" env-description:"API user password"`
		Email     string `yaml:"email" env:"NITAU_API_EMAIL" env-description:"API user email address"`
		AuthToken string `yaml:"authtoken" env:"NITAU_API_AUTH_TOKEN" env-description:"API JWT authorization token"`
		RootURI   string `yaml:"rooturi" env:"NITAU_API_ROOT_URI" env-description:"API ROOT URI"`
		MOurl     string `yaml:"mourl" env:"NITAU_API_MO_URL" env-description:"MO URL to POST incoming MO messages"`
		SmsURL    string `yaml:"smsurl" env:"NITAU_API_SMS_URL" env-description:"API SMS endpoint"`
	} `yaml:"api"`
	Server struct {
		Port string `yaml:"port" env:"NITAU_SERVER_PORT" env-description:"Server port" env-default:"5000"`
	} `yaml:"server"`
}

// Args command-line parameters
type Args struct {
	ConfigPath string
}

func init() {
	// Log error if .env file does not exist
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}
	// helpers.RefreshToken()
}

func myTask() {
	fmt.Println("This task will run periodically")
}
func executeCronJob() {
	// gocron.Every(1).Minute().Do(helpers.RefreshToken)
	gocron.Every(1).Day().At("23:00").Do(helpers.RefreshToken)
	<-gocron.Start()
}

func main() {
	var cfg Config
	args := ProcessArgs(&cfg)

	// read configuration from the file and environment variables
	if err := cleanenv.ReadConfig(args.ConfigPath, &cfg); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	go executeCronJob()
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		sms := new(controllers.SMSController)
		v1.POST("/sms", sms.Default)

		mo := new(controllers.MOController)
		v1.POST("/mo", mo.Default)

		pl := new(controllers.TestController)
		v1.POST("/test", pl.Default)
	}

	// Handle error response when a route is not defined
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not found"})
	})

	// Init our Server
	router.Run(":" + cfg.Server.Port)
}

// ProcessArgs processes and handles CLI arguments
func ProcessArgs(cfg interface{}) Args {
	var a Args

	f := flag.NewFlagSet("Example server", 1)
	f.StringVar(&a.ConfigPath, "c", "config.yml", "Path to config file")

	fu := f.Usage
	f.Usage = func() {
		fu()
		envHelp, _ := cleanenv.GetDescription(cfg, nil)
		fmt.Fprintln(f.Output())
		fmt.Fprintln(f.Output(), envHelp)
	}

	f.Parse(os.Args[1:])
	return a
}
