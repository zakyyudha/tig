package usecase

import (
	"fmt"
	"github.com/zakyyudha/tig/config"
	"github.com/zakyyudha/tig/pkg/google"
	"github.com/zakyyudha/tig/src/commit/domain"
)

type CommitUsecase interface {
	ToSpreadsheet(params *domain.Params) error
}

type CommitUsecaseImpl struct {
	GoogleSpreadsheet *google.Spreadsheet
	Config            *config.Config
}

func NewCommitUsecaseImpl(config *config.Config) *CommitUsecaseImpl {
	spreadsheet, err := google.NewSpreadsheet(config)
	if err != nil {
		panic(fmt.Sprintf("Something wrong while load a spreadsheet: %v", err))
	}

	return &CommitUsecaseImpl{
		GoogleSpreadsheet: spreadsheet,
		Config:            config,
	}
}
