package market

import (
	"encoding/json"
	"github.com/pkg/errors"
	. "mystake/lib"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type StockInfo struct {
	Code          string  `json:"f57"`
	Name          string  `json:"f58"`
	LastPrice     float64 `json:"f60"`
	Price         float64 `json:"f43"`
	HighestPrice  float64 `json:"f44"`
	LowestPrice   float64 `json:"f45"`
	ChangePercent float64 `json:"f170"`
	Buy1Lots      int64   `json:"f20"`
	SellLots      int64   `json:"f40"`
	Volume        int64   `json:"f47"`
}

func GetStockInfo(code string) (*StockInfo, error) {
	if code == "" || len(code) < 6 || !unicode.IsDigit(rune(code[0])) {
		return nil, errors.Errorf("Invalid stock code: %v", code)
	}
	if strings.HasPrefix(code, "60") ||
		strings.HasPrefix(code, "900") ||
		strings.HasPrefix(code, "11") ||
		strings.HasPrefix(code, "688") {
		code = "1." + code
	} else {
		code = "0." + code
	}
	var stockInfo *StockInfo
	url := "https://push2.eastmoney.com/api/qt/stock/get"
	params := map[string]string{
		"fltt":   "2",
		"invt":   "2",
		"klt":    "1",
		"secid":  code,
		"fields": "f57,f58,f60,f43,f44,f45,f47,f170,f19,f20,f39,f40,f530",
		"ut":     "b2884a393a59ad64002292a3e90d46a5",
		//"cb": "jQuery183003743205523325188_1589197499471",
		"_": strconv.FormatInt(time.Now().Unix(), 10),
	}
	resp, err := RequestUrl(url, params, 3, "GET")
	if err != nil {
		return nil, err
	}
	//stock_info := json.Unmarshal([]byte(resp.String()), &resp_data)
	respData := struct {
		Data StockInfo
	}{}
	json.Unmarshal([]byte(resp.String()), &respData)
	stockInfo = &respData.Data
	return stockInfo, nil
}
