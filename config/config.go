package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// NitauConf is the global conf
var NitauConf Config

func init() {
	args := ProcessArgs(&NitauConf)

	// read configuration from the file and environment variables
	if err := cleanenv.ReadConfig(args.ConfigPath, &NitauConf); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

// Config is the top level cofiguration object
type Config struct {
	API struct {
		Username     string `yaml:"username" env:"NITAU_API_USER" env-description:"API user name"`
		Password     string `yaml:"password" env:"NITAU_API_PASSWORD" env-description:"API user password"`
		Email        string `yaml:"email" env:"NITAU_API_EMAIL" env-description:"API user email address"`
		AuthToken    string `yaml:"authtoken" env:"NITAU_API_AUTH_TOKEN" env-description:"API JWT authorization token"`
		RootURI      string `yaml:"rooturi" env:"NITAU_API_ROOT_URI" env-description:"API ROOT URI"`
		MOurl        string `yaml:"mourl" env:"NITAU_API_MO_URL" env-description:"MO URL to POST incoming MO messages"`
		SmsURL       string `yaml:"smsurl" env:"NITAU_API_SMS_URL" env-description:"API SMS endpoint"`
		ShortCodeUID string `yaml:"shortcode_uid" env:"NITAU_API_SHORTCODE_UID" env-description:"NITAU short code UID"`
	} `yaml:"api"`
	Server struct {
		Port string `yaml:"port" env:"NITAU_SERVER_PORT" env-description:"Server port" env-default:"5000"`
	} `yaml:"server"`
	Database struct {
		URI string `yaml:"uri" env:"NITAU_DB" env-default:"postgres://postgres:postgres@localhost/nitau?sslmode=disable"`
	}
}

// Args command-line parameters
type Args struct {
	ConfigPath string
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
