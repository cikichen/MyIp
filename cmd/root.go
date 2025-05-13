/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"ip/ui"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ip",
	Short: "显示您的IP地址及详细信息",
	Long:  `显示您的本地和公网IP地址，以及地理位置、网络信息、IP类型等详细信息。
您可以使用此工具获取以下信息：
- 本地和公网IP地址
- 地理位置信息（国家、城市、经纬度等）
- 网络提供商(ISP)和网络类型
- IP类型和质量评分
- 其他相关信息（货币、时区等）`,
	Run: func(cmd *cobra.Command, args []string) {
		// 检查是否需要测试API源
		testAPI, _ := cmd.Flags().GetBool("test-api")
		if testAPI {
			fmt.Println(ui.DrawNotice("正在测试所有API源的可用性...", ui.IconInfo, ui.BgBrightBlue))
			results := TestAPISource()

			// 显示测试结果
			fmt.Println("")
			// 构建漂亮的API测试结果表格
			var tableContent strings.Builder
			validCount := 0
			
			tableContent.WriteString(fmt.Sprintf("%s%-20s%s | %s%-10s%s\n", ui.Bold, "API 源", ui.Reset, ui.Bold, "状态", ui.Reset))
			tableContent.WriteString(strings.Repeat("─", 33) + "\n")
			
			// 按名称排序显示结果
			names := make([]string, 0, len(results))
			for name := range results {
				names = append(names, name)
			}
			sort.Strings(names)
			
			for _, name := range names {
				valid := results[name]
				status := ui.BrightRed + "不可用 " + ui.IconCross + ui.Reset
				if valid {
					status = ui.BrightGreen + "可用 " + ui.IconCheck + ui.Reset
					validCount++
				}
				tableContent.WriteString(fmt.Sprintf("%-20s | %s\n", name, status))
			}
			
			tableContent.WriteString(strings.Repeat("─", 33) + "\n")
			// 添加总结行
			tableContent.WriteString(fmt.Sprintf(
				"%s%-20s%s | %s%d/%d%s",
				ui.Bold, "总计", ui.Reset,
				ui.BrightYellow + ui.Bold, validCount, len(results), ui.Reset,
			))
			
			// 创建结果卡片
			resultCard := ui.DrawCard("API源测试结果", ui.IconServer, tableContent.String(), 50, ui.BrightBlue)
			fmt.Println(resultCard)
			
			// 显示测试完成消息
			fmt.Println(ui.DrawNotice("API源测试完成！", ui.IconCheck, ui.BgBrightGreen))
			return
		}

		// 获取本地IP
		localIp, err := externalIP()
		if err != nil {
			fmt.Println(ui.DrawNotice("无法获取本地IP地址: "+err.Error(), ui.IconWarning, ui.BgBrightRed))
			return
		}

		// 显示获取公网IP的状态栏
		fmt.Println(ui.DrawStatusBar("正在获取公网IP地址...", ui.BgBrightBlue))
		
		// 获取公网IP
		myIP := GetMyPublicIP()

		// 显示获取IP信息的状态栏
		fmt.Println(ui.DrawStatusBar("正在获取IP详细信息...", ui.BgBrightBlue))
		
		// 获取IP信息（使用负载均衡机制）
		result := OnlineIpInfo(myIP)

		// 使用卡片式UI显示结果
		var uiInfo *ui.IPInfo
		if result != nil {
			// 转换为ui.IPInfo类型
			uiInfo = &ui.IPInfo{
				IP:             result.IP,
				Network:        result.Network,
				Version:        result.Version,
				City:           result.City,
				Region:         result.Region,
				RegionCode:     result.RegionCode,
				Country:        result.Country,
				CountryName:    result.CountryName,
				CountryCode:    result.CountryCode,
				CountryCodeISO: result.CountryCodeISO,
				CountryCapital: result.CountryCapital,
				CountryTLD:     result.CountryTLD,
				ContinentCode:  result.ContinentCode,
				InEU:           result.InEU,
				Postal:         result.Postal,
				Latitude:       result.Latitude,
				Longitude:      result.Longitude,
				Timezone:       result.Timezone,
				UTCOffset:      result.UTCOffset,
				CallingCode:    result.CallingCode,
				Currency:       result.Currency,
				CurrencyName:   result.CurrencyName,
				Languages:      result.Languages,
				CountryArea:    result.CountryArea,
				Population:     result.Population,
				ASN:            result.ASN,
				Org:            result.Org,
				IsPure:         result.IsPure,
				PureScore:      result.PureScore,
				PureType:       result.PureType,
				IPType:         result.IPType,
				IsProxy:        result.IsProxy,
				IsDC:           result.IsDC,
				APISource:      result.APISource,
			}
		}
		
		// 显示IP信息
		fmt.Println(ui.DrawIPInfo(localIp.String(), myIP, uiInfo))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ip.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "帮助信息示例")

	// 添加测试API源的标志
	rootCmd.Flags().Bool("test-api", false, "仅测试所有IP信息API源的可用性，不获取IP信息")
}
