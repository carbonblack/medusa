package main

import (
	"os"
	"fmt"
	"flag"
	"bufio"
	"strings"
	"io/ioutil"
	"encoding/json"
	"path/filepath"

	"gopkg.in/resty.v1"
	"github.com/BurntSushi/toml"	
)

type MedusaConfig struct {
	Org string
	ApiKey string
	BaseOrgUrl string
}

// A subset of https://developer.github.com/v3/repos/#list-organization-repositories
type Repo struct {
	Name    string `json:"name"`
	Description string `json:"description"`
	Url string `json:"url"`
	Language string `json:"language"`
	Private bool `json:private"`
	Fork bool `json:fork"`
}

// A subset of https://developer.github.com/v3/users/#get-a-single-user
type User struct {
}

// A subset of https://developer.github.com/v3/teams/#get-team
type Team struct {
}


func main() {
	//The init command
	initCommandPtr := flag.Bool("init", true, "Set/initialize the Medusa configs e.g. org and API key")

	//The repos command
	reposCommand := flag.NewFlagSet("repos", flag.ExitOnError)
	repoTypePtr := reposCommand.String("type", "all", "Repo type all|private|public, defaults to all")
	reposVerbosePtr := reposCommand.Bool("verbose", false, "Verbose mode i.e. full detais")
	reposCsvPtr := reposCommand.Bool("csv", false, "Report results in CSV format")

	//The repo command
	repoCommand := flag.NewFlagSet("repo", flag.ExitOnError)
	repoNamePtr := repoCommand.String("name", "", "The repo's name (required)")
	repoVerbosePtr := repoCommand.Bool("verbose", false, "Verbose mode i.e. full detais")
	repoCsvPtr := repoCommand.Bool("csv", false, "Report results in CSV format")

	//TODO
	//setApiKey := flag.NewFlagSet("api_key", flag.ExitOnError)
	//infoCommand := flag.NewFlagSet("info", flag.ExitOnError)
	//medusa users --filters -2fa -verbose -csv
	//medusa user <user_name> -verbose -teams -repos -csv
	//medusa teams|groups -verbose -users -repos -csv
	//medusa team|group <team_group_name> --filters
	//medusa collaborators --filters
	//medusa collaborator <collborator_name> --filters
	
	flag.Parse()
	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	homeDir := os.Getenv("HOME")
	dot_medusa := filepath.Join(homeDir, ".medusa")
	config := loadConfig(dot_medusa)
	
	switch os.Args[1] {
	case "init":
		fmt.Println(*initCommandPtr)
		if *initCommandPtr == true {
			setConfig(dot_medusa)
		}
	case "repos":
		reposCommand.Parse(os.Args[2:])
		repos(&config, repoTypePtr, reposVerbosePtr, reposCsvPtr)
	case "repo":
		repoCommand.Parse(os.Args[2:])
		if *repoNamePtr == "" {
			repoCommand.PrintDefaults()
			os.Exit(1)
		}
		repo(&config, repoNamePtr, repoVerbosePtr, repoCsvPtr)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}	
}

func setConfig(confFilePath string) (MedusaConfig) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the name of your GitHub organization: ")
	org, _ := reader.ReadString('\n')
	org = strings.TrimSpace(org)
	fmt.Print("Copy/paste your GitHub API key: ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)
	confData := []byte(fmt.Sprintf("Org=\"%s\"\nApiKey=\"%s\"\n", org, apiKey))
	err := ioutil.WriteFile(confFilePath, confData, 0644)
	if err != nil {
		panic(err)
	}
	BaseOrgUrl := fmt.Sprintf("https://api.github.com/orgs/%s", org)
	return MedusaConfig{org, apiKey, BaseOrgUrl}
}

func loadConfig(confFilePath string) (MedusaConfig) {
	var config MedusaConfig
	confExists, _ := confFileExists(confFilePath)
	if confExists {
		if _, err := toml.DecodeFile(confFilePath, &config); err != nil {
			panic(err)
		}
	} else {
		config = setConfig(confFilePath)
	}
	baseOrgUrl := fmt.Sprintf("https://api.github.com/orgs/%s", config.Org)
	config.BaseOrgUrl = baseOrgUrl
	fmt.Println(fmt.Sprintf("config: %s", config))
	return config
}

func confFileExists(confFilePath string) (bool, error){
	exists := true
	var existsError error
	if _, err := os.Stat(confFilePath); err != nil {
		if os.IsNotExist(err) {
			exists = false
		} else {
			existsError = err
		}
	}
	return exists, existsError
}

func repos(config *MedusaConfig, repoType *string, verbose *bool, csv *bool){
	//curl -s -H "Authorization: token TOKEN" 'https://api.github.com/orgs/carbonblack/repos?type=private&per_page=100'
	fmt.Println(config.BaseOrgUrl)
	resp, err := resty.R().
		SetQueryParams(map[string]string{
		"type": *repoType,
		"per_page": "100",
	}).
		SetHeader("Authorization", fmt.Sprintf("token %s", config.ApiKey)).
		Get(fmt.Sprintf("%s/repos", config.BaseOrgUrl))
	if (err != nil) {
		panic(err)
	}
	var repos []Repo
	err = json.Unmarshal(resp.Body(), &repos)
	//jd := json.NewDecoder(resp.Body())
	//err = jd.Decode(&repos)
	if err != nil {
		panic(err)
	}
	for _, r := range repos {
		fmt.Printf("%s\n\n",r.Name)
	}
}	

func repo(config *MedusaConfig, repoName *string, verbose *bool, csv *bool){
	//TODO
}	

func users(){
	fmt.Println("users")
}	

func user(){
	fmt.Println("user")
}	

func teams(){
	fmt.Println("teams")
}	

func team(){
	fmt.Println("team")
}	

func groups(){
	fmt.Println("groups")
}	

func group(){
	fmt.Println("group")
}	

func collaborator(){
	fmt.Println("collaborator")
}	

func collaborators(){
	fmt.Println("collaborators")
}	




