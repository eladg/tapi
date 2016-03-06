package main

import (
	"os"
	"fmt"
	"log"
	"net/http"

	"io/ioutil"
	"bufio"
	"strings"
	"path/filepath"
	"encoding/json"

	cli "github.com/codegangsta/cli"
	transport "google.golang.org/api/googleapi/transport"
	translate "google.golang.org/api/translate/v2"
)

type Config struct {
	ApiToken string `json:"api_token"`
	Target   string `json:"target"`
	Origin   string `json:"origin"`
}

func commands() []cli.Command {
	return []cli.Command {
		cli.Command {
			Name:        "setup",
			Usage:       "setup usage",
			Description: "setup long description",
			Action:      setup_action,
		},
	}
}

func flags() []cli.Flag {
	return []cli.Flag {
		cli.StringFlag {
			Name: "debug",
			Value: "debug",
			Usage: "prints debug messages",
		},
	}
}

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	buff, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %v",err)
	}
	return strings.Trim(buff,"\n")
}

func getConfigFromUserInput() *Config {
	fmt.Print("Google API Key: "); key := getUserInput()
	fmt.Print("Default Target: "); target := getUserInput()
	fmt.Print("Default Origin: "); origin := getUserInput()
	return &Config{ ApiToken: key, Target: target, Origin: origin, }
}

func ReadConfig() *Config {
	c := new(Config)
	
	buf, err := ioutil.ReadFile(c.ConfigPath())
	if err != nil { return c }

	err = json.Unmarshal(buf, c)
	if err != nil {
		log.Println("corrupted config file was found. Ignoring.")
	}

	return c	
}

func (c *Config) WriteConfig() error {
	path   := c.ConfigPath()
  buf, _ := json.Marshal(c)
  return ioutil.WriteFile(path, buf, 0644)
}

func (c *Config) ConfigPath() string {
	return string(filepath.Join(os.Getenv("HOME"), ".translaterc"))
}

func (c *Config) IsEmpty() bool {
	return c.ApiToken == ""
}

func setup_action(c *cli.Context) {
	conf := ReadConfig()
	if !conf.IsEmpty() {
		fmt.Printf("Found the following config:\n  API Key: %s\n  Default Target: %s\n  Default Origin: %s\n\n", conf.ApiToken, conf.Target, conf.Origin)
		fmt.Print("Would you like to overwrite it? (y/n) ")
		answer := getUserInput() ; fmt.Println()
		if strings.ToLower(answer) != "y" {
			return
		}
	}

	fmt.Print("In order to use translate, you will to setup Google API Key over here: https://goo.gl/6aj3Ha.\nIf you don't have a google developers account setup, follow instuctions here: https://developers.google.com/places/web-service/get-api-key\n\n")
	conf = getConfigFromUserInput() ; fmt.Println()

	err := conf.WriteConfig()
	if err != nil {
		log.Fatal("Error writing config: ",err)
	}

	fmt.Println("config was saved to:", conf.ConfigPath())
}

func translate_action(c *cli.Context) {
	conf := ReadConfig()
	if conf.IsEmpty() {
		fmt.Println("User configuration and Google Translate API tokens were not set.\nRun 'translate setup --help' to get help in setting up the proper keys.\n")
		os.Exit(3)
	}

	text := c.Args()
	target := conf.Target

	client := &http.Client{
		Transport: &transport.APIKey{Key: conf.ApiToken},
	}

	s, err := translate.New(client)
	if err != nil {
		log.Fatalf("Unable to create translate service: %v", err)
	}

	req := s.Translations.List(text,target)
	res, err := req.Do()

	if err != nil {
		fmt.Printf("Google Translate replied with error.\n%v\n\nFor help setting up approprate keys, please follow 'translate setup --help' or https://github.com/eladg/translate\n\n",err)
	}

	fmt.Printf("%s\n",res.Data.Translations[0].TranslatedText)
}

func main() {
	app := cli.NewApp()
	app.Name = "Translate"
	app.Version = "0.1.0"
	app.Usage = "Cli tool for Google Translate API"
	app.Author = "github.com/eladg"
	app.Commands = commands()
	app.Flags = flags()	
	app.Action = translate_action
	app.Run(os.Args)
}