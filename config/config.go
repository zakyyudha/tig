package config

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"os"
	"path"
	"strings"
)

var configDir string

const (
	CredentialsFileName = "credentials.json"
	TokenFileName       = "token.json"
	SheetConfigFileName = "sheet_config.json"
	Scope               = "https://www.googleapis.com/auth/spreadsheets"
	AppConfigPath       = ".config/tig"
	JiraURL             = "https://telkomdds.atlassian.net/browse/%s"
	HelpUsage           = `
Usage: tig <command> [args...]

Tig is a CLI app that allows you to push your commit message to a Google Spreadsheet.

The following commands are available:

commit:
  Push commit message to the Google Spreadsheet.
  Arguments:
  -m, --message: Commit message - for column Activity
  -tJC, --tigJiraCode: Jira Code - for column ReferTo
  -tN, --tigNote: Note - for column Note

Note: Tig supports all commit flags and args based on official git. Tig only gets the commit message and converts it to activity on the spreadsheet. For more information, please refer to the documentation at https://github.com/zakyyudha/tig.`
)

type Config struct {
	OAuthToken  *oauth2.Token
	OAuthConfig *oauth2.Config
	SheetConfig *sheetConfig
}

type sheetConfig struct {
	SpreadsheetId string `json:"spreadsheetId"`
	SheetName     string `json:"sheetName"`
	SquadName     string `json:"squadName"`
	UserName      string `json:"userName"`
}

func New() *Config {
	d, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Unable to find $HOME dir: %v", err))
	}

	configDir = path.Join(d, AppConfigPath)
	credPath := path.Join(configDir, CredentialsFileName)
	b, err := os.ReadFile(credPath)
	if err != nil {
		panic(fmt.Sprintf("Unable to find credentials.json on %v: %v", credPath, err))
	}

	config, err := google.ConfigFromJSON(b, Scope)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse client secret file to config: %v", err))
	}

	token, err := getTokenFromFile(TokenFileName)
	if err != nil {
		token = getTokenFromWeb(config)
		saveToken(TokenFileName, token)
	}

	configPath := path.Join(configDir, SheetConfigFileName)
	sheetCfg, err := readSheetConfig()
	if err != nil {
		panic(fmt.Sprintf("Unable to parse config.json on %v please delete it first: %v ", configPath, err))
	}

	return &Config{
		OAuthToken:  token,
		OAuthConfig: config,
		SheetConfig: sheetCfg,
	}
}

func getTokenFromFile(fileName string) (*oauth2.Token, error) {
	tokenPath := path.Join(configDir, fileName)
	f, err := os.Open(tokenPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Use json.Decoder to read the token file in chunks
	decoder := json.NewDecoder(f)
	token := &oauth2.Token{}
	for {
		if err := decoder.Decode(token); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}

	return token, nil
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
	fmt.Printf("Input authorization code: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		panic(fmt.Sprintf("Unable to read authorization code: %v", err))
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		panic(fmt.Sprintf("Unable to retrieve token from web: %v", err))
	}
	return tok
}

func saveToken(fileName string, token *oauth2.Token) {
	tokenPath := path.Join(configDir, fileName)
	fmt.Printf("Saving credential file to: %s\n", tokenPath)
	f, err := os.OpenFile(tokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		panic(fmt.Sprintf("Unable to cache oauth token: %v", err))
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func readSheetConfig() (*sheetConfig, error) {
	configPath := path.Join(configDir, SheetConfigFileName)
	sheetCfg := &sheetConfig{}
	c, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, create a new one
			return createNewSheetConfig(sheetCfg)
		}
		return nil, fmt.Errorf("unable to open config file %s: %v", configPath, err)
	}
	defer c.Close()

	err = json.NewDecoder(c).Decode(sheetCfg)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config file %s: %v", configPath, err)
	}

	return sheetCfg, nil
}

func createNewSheetConfig(sheetCfg *sheetConfig) (*sheetConfig, error) {
	configPath := path.Join(configDir, SheetConfigFileName)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("If your input contains a space, please add double quotes (e.g., \"Some Input\").")
	sheetCfg.SpreadsheetId = getInput("Input Spreadsheet ID: ", *reader)
	sheetCfg.SheetName = getInput("Input Sheet Name: ", *reader)
	sheetCfg.SquadName = getInput("Input Squad Name: ", *reader)
	sheetCfg.UserName = getInput("Input Your Name: ", *reader)

	f, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("unable to create config file %s: %v", configPath, err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(sheetCfg)
	if err != nil {
		return nil, fmt.Errorf("unable to encode config file %s: %v", configPath, err)
	}

	return sheetCfg, nil
}

func getInput(message string, reader bufio.Reader) string {
	for {
		fmt.Print(message)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("An error occurred: %s\n", err)
			continue
		}
		input = strings.TrimSpace(input)
		input = strings.Replace(input, "\n", "", -1)
		input = strings.Trim(input, "\"")
		if input != "" {
			return input
		}
	}
}
