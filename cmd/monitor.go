/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"math"
	. "mystake/lib"
	"mystake/market"
	"os/exec"
	"regexp"
	"time"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor up|down|notify",
	Short: "monitor the deviation from the stocks",
	Long: `For example:
monitor up -i true -f /Users/Yy/Desktop/Table.txt -u 3 -d -5`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		action := args[0]
		for {
			isOpeningTime := IsWeekday(time.Now()) && (IsInTimeRange(time.Now(), "09:40", "11:30") || IsInTimeRange(time.Now(), "13:00", "14:55"))

			if !isOpeningTime && !ignoreCloseTime {
				fmt.Println(IsInTimeRange(time.Now(), "09:40", "11:30"))
				fmt.Println("Take a break")
				time.Sleep(time.Millisecond * 10000)
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
	//if _, ok := stocks[code]; !ok {
	//	stocks[code] = stock_info
	//}
	//stocks[code] = stock_info

	name := Substring(bind.stockInfo.Name, 0, 2)
	//c := ReverseString(name)
	//notifyCode := "0." + c[0:3] + c[3:]
	var msg string
	if bind.stockInfo.ChangePercent > upPercent {
		msg = name + " up"
	}
	bind.notify(msg)
}

func (bind *monitor) changePercentDown() {
	name := Substring(bind.stockInfo.Name, 0, 2)
	var msg string
	if bind.stockInfo.ChangePercent < downPercent {
		msg = name + " down"
	}
	bind.notify(msg)
}

func (bind *monitor) notify(msg string) {
	if msg == "" {
		return
	}
	if _, ok := monitorStocks[bind.stockInfo.Code]; ok {
		//if monitorStocks[bind.stockInfo.Code].spokeTime.After(time.Now().Add(-10 * time.Minute)) {
		//	return
		//}
		if math.Abs(monitorStocks[bind.stockInfo.Code].stockInfo.ChangePercent-bind.stockInfo.ChangePercent) < 2 {
			return
		}
	}
	fmt.Println(msg, bind.stockInfo.ChangePercent)
	log.Infof(msg+" %v", bind.stockInfo.ChangePercent)
	//speak
	cmd := exec.Command("python3", GetCurrentAbPathByCaller()+"/../lib/speak.py", msg)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		log.Errorf("Speak error: %v", err)

	}
	monitorStocks[bind.stockInfo.Code] = monitorStock{
		spokeTime: time.Now(),
		stockInfo: bind.stockInfo,
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
