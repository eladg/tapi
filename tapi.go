package main

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"bufio"
	"strings"
	"reflect"
	"path/filepath"
	"encoding/json"

	cli "github.com/codegangsta/cli"
	transport "google.golang.org/api/googleapi/transport"
	// translate "google.golang.org/api/translate/v2" -> DISABLE due to: github.com/google/google-api-go-client/issues/79
	translate "github.com/eladg/google-api-go-client/translate/v2"
)

/***************************************************************************
 * globals, codegangsta/cli & config
 *************************************************************************** */
var debug_flag bool
var setup_flag bool

func Commands() []cli.Command {
	return []cli.Command {
		cli.Command {},
	}
}

func Flags() []cli.Flag {
	return []cli.Flag {
		cli.BoolFlag {
			Name: "debug",
			Usage: "Enable debug prints",
			Destination: &debug_flag,
		},
		cli.BoolFlag {
			Name: "setup",
			Usage: "start setup dialog",
			Destination: &setup_flag,
		},
	}
}

type Config struct {
	ApiToken string `json:"api_token"`
	Target   string `json:"target"`
	Origin   string `json:"origin"`
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
	return string(filepath.Join(os.Getenv("HOME"), ".tapirc"))
}

func (c *Config) IsEmpty() bool {
	return c.ApiToken == ""
}

func main() {
	app := cli.NewApp()
	app.Name = "tapi"
	app.Version = "0.1.1"
	app.Usage = "Google Translate API from the terminal"
	app.Author = "Elad Gariany (github.com/eladg)"
	app.Commands = Commands()
	app.Flags = Flags()	
	app.Action = translateAction
	app.HideVersion = true
	app.Run(os.Args)
}

/***************************************************************************
 * setup & user input
 *************************************************************************** */

func setupAction(c *cli.Context) {
	conf := ReadConfig()
	if !conf.IsEmpty() {
		fmt.Printf("Found the following config:\n  API Key: %s\n  Default Target: %s\n  Default Origin: %s\n\n", conf.ApiToken, conf.Target, conf.Origin)
		fmt.Print("Would you like to overwrite it? (y/n) ")
		answer := GetUserInput() ; fmt.Println()
		if strings.ToLower(answer) != "y" {
			return
		}
	}

	fmt.Print("In order to use tapi, you will to setup Google API Key over here: https://goo.gl/6aj3Ha.\nIf you don't have a google developers account setup, follow instructions here: https://cloud.google.com/translate/v2/getting_started\n\n")
	conf = GetConfigFromUserInput() ; fmt.Println()

	err := conf.WriteConfig()
	if err != nil {
		log.Fatal("Error writing config: ",err)
	}

	fmt.Println("config was saved to:", conf.ConfigPath())
}

func GetUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	buff, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %v",err)
	}
	return strings.Trim(buff,"\n")
}

func GetConfigFromUserInput() *Config {
	fmt.Print("Google API Key: "); key := GetUserInput()
	fmt.Print("Default Target: "); target := GetUserInput()
	fmt.Print("Default Origin: "); origin := GetUserInput()
	return &Config{ ApiToken: key, Target: target, Origin: origin, }
}

/***************************************************************************
 * translate
 *************************************************************************** */

func analyzeTranslations(translations []*translate.TranslationsResource) []string {
	res := make([]string, 0)
	for _, t := range translations {
		res = append(res, t.TranslatedText)
	}
	return res
}

func translateRequest(text []string, target, token string) (*translate.TranslationsListMain, error) {
	s, err := translate.New(&http.Client{Transport: &transport.APIKey{Key: token},})
	if err != nil { log.Fatalf("Unable to create translate service: %v", err)	}
	req := s.Translations.List(text,target)
	return req.Do()
}

func checkTranslateRequest(err error) {
	if err != nil {
		log.Fatalf("Google Translate replied with error.\n%v\n\nFor help setting up approprate keys, please follow 'tapi --setup' or https://github.com/eladg/tapi\n\n",err)
	}
}

func printResults(s []string) {
	fmt.Println(strings.Join(s, " "))
}

func translateAction(c *cli.Context) {
	if setup_flag {	setupAction(c); return }
	
	conf := ReadConfig()
	if conf.IsEmpty() {
		fmt.Println("User configuration and Google Translate API tokens were not set.\nRun 'tapi --setup' to get help in setting up the proper keys.\n")
		os.Exit(3)
	}

	// manage args and needed inputs
	text := []string(c.Args())
	target := conf.Target

	// perform a request & analyze its response
	res, err := translateRequest(text, target, conf.ApiToken) ; checkTranslateRequest(err)
	translations := analyzeTranslations(res.Data.Translations)

	// perform a reverse call in case translation equals user's import
	if reflect.DeepEqual(translations,text) {
		target := conf.Origin
		res, err := translateRequest(text, target, conf.ApiToken) ; checkTranslateRequest(err)	
		translations = analyzeTranslations(res.Data.Translations)
	}

	// print results to STDOUT
	printResults(translations)
}
