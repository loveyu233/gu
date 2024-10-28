package excel

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

// GetCell 从1开始,0返回空字符串 例如：0:"" 1:"A" 2:"B" 3:"C" ...
func GetCell(column int) string {
	result := ""
	for column > 0 {
		column--
		result = string('A'+cast.ToInt32(column%26)) + result
		column /= 26
	}
	return result
}

func NewExcelFile(sheetName []string, defaultSheetIndex ...int) (*excelize.File, error) {
	file := excelize.NewFile()
	for _, sheet := range sheetName {
		_, err := file.NewSheet(sheet)
		if err != nil {
			return nil, fmt.Errorf("创建工作表 %s 失败：%s", sheet, err.Error())
		}
	}

	if err := file.DeleteSheet("Sheet1"); err != nil {
		return nil, fmt.Errorf("删除Sheet1工作表失败：%s", err.Error())
	}

	index := 0
	if len(defaultSheetIndex) > 0 {
		if defaultSheetIndex[0] < 0 {
			return nil, fmt.Errorf("默认工作表索引不能小于0，当前默认工作表索引：%d", defaultSheetIndex[0])
		} else if defaultSheetIndex[0] > len(sheetName) {
			return nil, fmt.Errorf("默认工作表索引不能大于总工作表数量，当前默认工作表索引：%d", defaultSheetIndex[0])
		}
		index = defaultSheetIndex[0]
	} else {
		index = 0
	}
	file.SetActiveSheet(index)
	return file, nil
}

func FillExcelFile(file *excelize.File, sheetName string, data [][]any) error {
	for i := range data {
		for j := range data[i] {
			err := file.SetCellValue(sheetName, GetCell(j+1)+cast.ToString(i+1), data[i][j])
			if err != nil {
				return err
			}
		}
	}
	return nil
}
