/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"math"
	. "mystake/lib"
	"mystake/market"
	"os"
	"os/exec"
	"regexp"
	"time"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor up|down",
	Short: "monitor the deviation from the stocks",
	Long: `For example:
monitor up -i true -f /Users/Yy/Desktop/Table.txt -u 3 -d -5`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		action := args[0]
		for {
			isOpeningTime := IsWeekday(time.Now()) && (IsInTimeRange(time.Now(), "09:25", "11:30") || IsInTimeRange(time.Now(), "13:00", "15:00"))

			if !isOpeningTime && !ignoreCloseTime {
				fmt.Println("Take a break")
				time.Sleep(time.Millisecond * 100000)
				continue
			}
			monitor := monitor{}
			codes, err := getCodes()
			if err != nil {
				log.Errorf("get codes failed: %v", err)
			}
			for _, code := range codes {
				stockInfo, err := market.GetStockInfo(code)
				if err != nil {
					log.Errorf("get stock info failed: %v", err)
					continue
				}
				monitor.stockInfo = stockInfo
				switch action {
				case "up":
					monitor.changePercentUp()
				case "down":
					monitor.changePercentDown()
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
	},
}

var ignoreCloseTime bool
var filename string
var upPercent float64
var downPercent float64

func init() {
	rootCmd.AddCommand(monitorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monitorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	monitorCmd.Flags().BoolVarP(&ignoreCloseTime, "ignore_close_time", "i", false, "If set to true it runs anytime or only in opening time")
	monitorCmd.Flags().StringVarP(&filename, "filename", "f", "./code_list.txt", "Stock code list file")
	monitorCmd.Flags().Float64VarP(&upPercent, "up_percent", "u", 3.0, "Up ChangePercent")
	monitorCmd.Flags().Float64VarP(&downPercent, "down_percent", "d", -3.0, "Down ChangePercent")
}

type monitor struct {
	stockInfo *market.StockInfo
}
type monitorStock struct {
	spokeTime time.Time
	stockInfo *market.StockInfo
}

var monitorStocks = make(map[string]monitorStock)

func (bind *monitor) changePercentUp() {
	var msg string
	if bind.stockInfo.ChangePercent > upPercent {
		msg = "up"
	}
	bind.notify(msg)
}

func (bind *monitor) changePercentDown() {
	var msg string
	if bind.stockInfo.ChangePercent < downPercent {
		msg = "low"
	}
	bind.notify(msg)
}

func (bind *monitor) notify(msg string) {
	if msg == "" {
		return
	}
	if _, ok := monitorStocks[bind.stockInfo.Code]; ok {
		change := math.Abs(monitorStocks[bind.stockInfo.Code].stockInfo.ChangePercent - bind.stockInfo.ChangePercent)
		if change < 2 && time.Now().Before(monitorStocks[bind.stockInfo.Code].spokeTime.Add(time.Hour*12)) {
			return
		}
	}
	name := Substring(bind.stockInfo.Name, 0, 2)
	rcode := ReverseString(bind.stockInfo.Code)
	notifyCode := "0." + rcode[0:3] + rcode[3:]
	fmt.Printf("%v %v %.1f\n", name, msg, bind.stockInfo.ChangePercent)
	log.Infof("%v %v %.1f", name, msg, bind.stockInfo.ChangePercent)
	audio := GetCurrentAudio()
	if audio == "MacBook Air扬声器" {
		//dingding
		bind.ding(notifyCode + ", " + msg)
	} else {
		//speak
		bind.speak(name + " " + msg)
	}
	monitorStocks[bind.stockInfo.Code] = monitorStock{
		spokeTime: time.Now(),
		stockInfo: bind.stockInfo,
	}
}

func (bind *monitor) speak(msg string) {
	cmd := exec.Command("python3", "/Users/Chen/go/bin/speak.py", msg)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		log.Errorf("Speak error: %v", err)

	}
}
func (bind *monitor) ding(msg string) {
	url := "https://oapi.dingtalk.com/robot/send?access_token=" + os.Getenv("DIND_BOT_TOKEN")
	params := fmt.Sprintf(`{
		"msgtype": "text",
		"text":    {"content": "完成率: %v"},
	}`, msg)
	resp, _ := RequestUrl(url, params, 3, "POST")
	var result struct {
		errcode int64
		errmsg  string
	}
	json.Unmarshal([]byte(resp.String()), &result)
	if result.errcode != 0 {
		log.Errorf("ding error, " + result.errmsg)
	}
}

func getCodes() ([]string, error) {
	var codes []string
	//filename := "/Users/Chen/Downloads/2023-07-08.xlsx"
	//f, err := excelize.OpenFile(filename)
	//if err != nil {
	//	return nil, err
	//}
	//defer func() {
	//	if err := f.Close(); err != nil {
	//		fmt.Println(err)
	//	}
	//}()
	//rows, err := f.GetRows("Sheet1")
	CovertToUTF8(filename) //nolint:errcheck
	lines, err := ReadLines(filename)
	if err != nil {
		return nil, err
	}
	for _, line := range lines {
		re := regexp.MustCompile("[0-9]+")
		// 在字符串中查找所有匹配的数字
		matches := re.FindAllString(line, -1)
		if len(matches) > 0 && len(matches[0]) >= 6 {
			code := matches[0][0:6]
			codes = append(codes, code)
		}
		//if len(line) >= 8 && regexp.MustCompile(`^[0-9]+$`).MatchString(line[2:8]) {
		//	codes = append(codes, line[2:8])
		//}
	}

	return codes, nil
}
