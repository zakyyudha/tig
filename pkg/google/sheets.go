package google

import (
	"context"
	"errors"
	"fmt"
	"github.com/zakyyudha/tig/config"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"net/http"
)

type Spreadsheet struct {
	Service *sheets.Service
	Client  *http.Client
	Config  *config.Config
}

func NewSpreadsheet(config *config.Config) (*Spreadsheet, error) {
	ctx := context.Background()
	client := getClient(config)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to retrieve Sheets client: %v", err))
	}
	return &Spreadsheet{
		Service: srv,
		Client:  client,
		Config:  config,
	}, nil
}

func (g *Spreadsheet) ReadData(readRange string) ([][]interface{}, error) {
	resp, err := g.Service.Spreadsheets.Values.Get(g.Config.SheetConfig.SpreadsheetId, readRange).Do()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to retrieve data from sheet: %v", err))
	}

	if len(resp.Values) == 0 {
		return nil, errors.New(fmt.Sprintf("No data found."))
	}

	return resp.Values, nil
}

func (g *Spreadsheet) UpdateData(updateRange string, values [][]interface{}) error {
	ctx := context.Background()
	valueInputOptions := "USER_ENTERED"
	_, err := g.Service.Spreadsheets.Values.Update(
		g.Config.SheetConfig.SpreadsheetId,
		updateRange,
		&sheets.ValueRange{
			Range:  updateRange,
			Values: values,
		},
	).
		ValueInputOption(valueInputOptions).
		Context(ctx).
		Do()
	if err != nil {
		return errors.New(fmt.Sprintf("Error while updating data to spreadsheet: %v", err))
	}
	return nil
}

func getClient(config *config.Config) *http.Client {
	return config.OAuthConfig.Client(context.Background(), config.OAuthToken)
}
