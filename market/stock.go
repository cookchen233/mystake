package market

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
		return nil, errors.Errorf("Invalid stock code: %s", code)
	}
	if strings.HasPrefix(code, "60") ||
		strings.HasPrefix(code, "900") ||
		strings.HasPrefix(code, "11") ||
		strings.HasPrefix(code, "688") {
		code = "1." + code
	} else {
		code = "0." + code
	}
	var stock_info *StockInfo
	url := "https://push2.eastmoney.com/api/qt/stock/get"
	client := resty.New()
	//client.SetProxy("http://127.0.0.1:1080")
	//resp2, err := client.R().Get("https://ip.900cha.com/")
	//fmt.Println(resp2.String())
	var TryRequestTimes int64
TryRequest:
	TryRequestTimes++
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"fltt":   "2",
			"invt":   "2",
			"klt":    "1",
			"secid":  code,
			"fields": "f57,f58,f60,f43,f44,f45,f47,f170,f19,f20,f39,f40,f530",
			"ut":     "b2884a393a59ad64002292a3e90d46a5",
			//"cb": "jQuery183003743205523325188_1589197499471",
			"_": strconv.FormatInt(time.Now().Unix(), 10),
		}).
		SetHeader("Accept", "application/json").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36").
		//SetResult(resp_data).
		//ForceContentType("application/json").
		Get(url)

	if err != nil {
		if strings.Contains(err.Error(), "timeout") && TryRequestTimes < 3 {
			log.Warnf("Timeout and try again")
			goto TryRequest
		}
		log.Errorf("request error, url: %s, msg: %s, body:", url, err, resp.RawBody())
		return nil, err
	}
	//stock_info := json.Unmarshal([]byte(resp.String()), &resp_data)
	resp_data := struct {
		Data StockInfo
	}{}
	json.Unmarshal([]byte(resp.String()), &resp_data)
	stock_info = &resp_data.Data
	return stock_info, nil
}
