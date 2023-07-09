/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	. "mystake/lib"
	"mystake/market"
	"os"
	"regexp"
	"time"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor up",
	Short: "monitor the deviation from the stocks",
	Long: `For example:
monitor up -i true -f /Users/Yy/Desktop/Table.txt -u 3 -d -5`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for {
			is_openning_time := IsInTimeRange(time.Now(), "9:40", "11:30") || IsInTimeRange(time.Now(), "13:30", "14:55")

			if (!IsWeekday(time.Now()) || !is_openning_time) && !ignore_close_time {
				fmt.Println("Take a break")
				time.Sleep(time.Millisecond * 5000)
				continue
			}
			monitor := monitor{
				speech: htgotts.Speech{
					Folder:   "audio",
					Language: voices.Chinese,
					Handler:  &handlers.MPlayer{},
				},
			}
			codes, err := getCodes()
			if err != nil {
				log.Errorf("get codes failed: %v", err)
			}
			for _, code := range codes {
				stock_info, err := market.GetStockInfo(code)
				if err != nil {
					log.Errorf("get stock info failed: %v", err)
					continue
				}
				monitor.stock_info = stock_info
				switch args[0] {
				case "up":
					monitor.deviationChangePercent()
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
	},
}

var ignore_close_time bool
var filename string
var up_percent float64
var down_percent float64

func init() {
	rootCmd.AddCommand(monitorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// monitorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	monitorCmd.Flags().BoolVarP(&ignore_close_time, "ignore_close_time", "i", false, "If set to true it runs anytime or only in opening time")
	monitorCmd.Flags().StringVarP(&filename, "filename", "f", "./code_list.txt", "Stock code list file")
	monitorCmd.Flags().Float64VarP(&up_percent, "up_percent", "u", 3.0, "Up ChangePercent")
	monitorCmd.Flags().Float64VarP(&down_percent, "down_percent", "d", -3.0, "Down ChangePercent")
}

type monitor struct {
	stock_info *market.StockInfo
	speech     htgotts.Speech
}
type spokeInfo struct {
	SpokeTime    time.Time
	SpokeMessage string
}

var spokes = make(map[string]spokeInfo)

func (bind *monitor) deviationChangePercent() {
	//if _, ok := stocks[code]; !ok {
	//	stocks[code] = stock_info
	//}
	//stocks[code] = stock_info

	name := Substring(bind.stock_info.Name, 0, 2)
	//c := ReverseString(name)
	//notifyCode := "0." + c[0:3] + c[3:]
	var msg string
	if bind.stock_info.ChangePercent > up_percent {
		msg = name + " up"
	} else if bind.stock_info.ChangePercent < down_percent {
		msg = name + " down"
	}
	bind.speak(msg)
}

func (bind *monitor) speak(msg string) {
	if msg == "" {
		return
	}
	if _, ok := spokes[bind.stock_info.Code]; ok {
		if spokes[bind.stock_info.Code].SpokeTime.After(time.Now().Add(-5 * time.Minute)) {
			return
		}
	}
	fmt.Println(msg, bind.stock_info.ChangePercent)
	log.Infof(msg+" %v", bind.stock_info.ChangePercent)
	err := bind.speech.Speak(msg)
	if err != nil {
		fmt.Println(err)
		log.Errorf("Speak error: %v", err)

	}
	spokes[bind.stock_info.Code] = spokeInfo{
		SpokeTime:    time.Now(),
		SpokeMessage: msg,
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
	// 读取原始文件内容
	inputData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// 转换为UTF-8编码
	outputData, _, err := transform.Bytes(simplifiedchinese.GB18030.NewDecoder(), inputData)
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(filename, outputData, 0644)
	if err != nil {
		return nil, err
	}
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
