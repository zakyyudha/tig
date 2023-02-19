package usecase

import (
	"fmt"
	"github.com/zakyyudha/tig/src/commit/domain"
	"time"
)

func (uc *CommitUsecaseImpl) ToSpreadsheet(params *domain.Params) error {
	//ReadData to get lastRows
	readRange := fmt.Sprintf("%v!D:G", uc.Config.SheetConfig.SheetName)
	data, err := uc.GoogleSpreadsheet.ReadData(readRange)
	if err != nil {
		return err
	}

	// Update the data
	lastRows := len(data) + 1
	err = uc.GoogleSpreadsheet.UpdateData(
		fmt.Sprintf("%v!B%v:G", uc.Config.SheetConfig.SheetName, lastRows), [][]interface{}{
			{
				uc.Config.SheetConfig.UserName,
				uc.Config.SheetConfig.SquadName,
				params.Activity,
				params.ReferTo,
				time.Now().Format("02-Jan-2006"),
				params.Note,
			},
		})
	if err != nil {
		return err
	}

	return nil
}
