/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"ip/network"
	"ip/ui"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// nettestCmd 代表网络测试命令
var nettestCmd = &cobra.Command{
	Use:   "nettest",
	Short: "测试常用站点的连通性",
	Long: `测试常用站点的连通性，包括响应时间、Ping延迟和丢包率等指标。
例如:
  ip nettest
  ip nettest --url github.com,google.com
  ip nettest --timeout 15
  ip nettest --detailed`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取参数
		urlsFlag, _ := cmd.Flags().GetString("url")
		timeout, _ := cmd.Flags().GetInt("timeout")
		detailed, _ := cmd.Flags().GetBool("detailed")
		
		// 显示网络测试的状态栏
		fmt.Println(ui.DrawStatusBar("正在测试站点连通性...", ui.BgBrightBlue))
		
		// 处理自定义URL
		var siteResults []network.SiteTestResult
		
		if urlsFlag != "" {
			// 用户提供了自定义URL
			customUrls := strings.Split(urlsFlag, ",")
			sites := make([]struct{
				Name string
				URL  string
			}, 0, len(customUrls))
			
			for _, url := range customUrls {
				url = strings.TrimSpace(url)
				if url == "" {
					continue
				}
				
				// 确保URL格式正确
				if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
					url = "https://" + url
				}
				
				// 从URL中提取名称
				name := url
				name = strings.TrimPrefix(name, "http://")
				name = strings.TrimPrefix(name, "https://")
				name = strings.TrimPrefix(name, "www.")
				name = strings.Split(name, "/")[0]
				name = strings.Split(name, ".")[0]
				name = strings.Title(name)
				
				sites = append(sites, struct{
					Name string
					URL  string
				}{
					Name: name,
					URL:  url,
				})
			}
			
			// 设置全局HTTP超时
			if timeout > 0 {
				network.SetGlobalTimeout(time.Duration(timeout) * time.Second)
			}
			
			// 测试自定义站点
			results := make([]network.SiteTestResult, 0, len(sites))
			resultChan := make(chan network.SiteTestResult, len(sites))
			
			for _, site := range sites {
				go func(s struct{
					Name string
					URL  string
				}) {
					resultChan <- network.TestSite(s)
				}(site)
			}
			
			// 收集结果
			for i := 0; i < len(sites); i++ {
				results = append(results, <-resultChan)
			}
			
			siteResults = results
		} else {
			// 使用默认的常用站点
			// 设置全局HTTP超时
			if timeout > 0 {
				network.SetGlobalTimeout(time.Duration(timeout) * time.Second)
			}
			
			// 测试常用站点连通性
			siteResults = network.TestCommonSites()
		}

		// 使用新的 lipgloss 布局显示网络测试结果
		fmt.Println(ui.RenderNetworkTestWithLipgloss(siteResults))
		
		// 如果需要详细信息，则显示额外的测试细节
		if detailed {
			detailedInfo := getDetailedTestInfo(siteResults)
			fmt.Println(detailedInfo)
		}
	},
}

// getDetailedTestInfo 返回详细的测试信息
func getDetailedTestInfo(results []network.SiteTestResult) string {
	var sb strings.Builder
	
	sb.WriteString(ui.DrawCard("详细测试信息", ui.IconNetwork, "以下是每个站点的详细测试数据：", 60, ui.BrightMagenta))
	sb.WriteString("\n\n")
	
	for _, result := range results {
		// 构建每个站点的详细信息卡片
		var siteInfo strings.Builder
		
		// 基本信息
		siteInfo.WriteString(fmt.Sprintf("%s站点URL:%s %s\n", ui.Bold, ui.Reset, result.URL))
		
		// 状态信息
		statusText := ui.BrightGreen + "成功" + ui.Reset
		if !result.Accessible {
			statusText = ui.BrightRed + "失败 - " + result.Error + ui.Reset
		}
		siteInfo.WriteString(fmt.Sprintf("%s访问状态:%s %s\n", ui.Bold, ui.Reset, statusText))
		
		if result.Accessible {
			siteInfo.WriteString(fmt.Sprintf("%s状态码:%s %d\n", ui.Bold, ui.Reset, result.StatusCode))
		}
		
		// 时间指标
		if result.Accessible {
			siteInfo.WriteString(fmt.Sprintf("%s总响应时间:%s %.2f秒\n", ui.Bold, ui.Reset, result.ResponseTime.Seconds()))
		}
		siteInfo.WriteString(fmt.Sprintf("%sDNS解析时间:%s %.2f秒\n", ui.Bold, ui.Reset, result.DNSTime.Seconds()))
		siteInfo.WriteString(fmt.Sprintf("%s连接建立时间:%s %.2f秒\n", ui.Bold, ui.Reset, result.ConnectTime.Seconds()))
		
		// Ping和Generate204
		pingTimeText := "超时"
		if result.PingTime > 0 {
			pingTimeText = fmt.Sprintf("%.1f毫秒", float64(result.PingTime)/float64(time.Millisecond))
		}
		siteInfo.WriteString(fmt.Sprintf("%sPing延迟:%s %s\n", ui.Bold, ui.Reset, pingTimeText))
		
		siteInfo.WriteString(fmt.Sprintf("%sPing丢包率:%s %.1f%%\n", ui.Bold, ui.Reset, result.PingLoss*100))
		
		gen204Text := "超时"
		if result.Generate204 > 0 {
			gen204Text = fmt.Sprintf("%.1f毫秒", float64(result.Generate204)/float64(time.Millisecond))
		}
		siteInfo.WriteString(fmt.Sprintf("%sGenerate_204延迟:%s %s", ui.Bold, ui.Reset, gen204Text))
		
		// 创建站点卡片
		cardColor := ui.BrightGreen
		if !result.Accessible {
			cardColor = ui.BrightRed
		}
		siteCard := ui.DrawCard(result.Name, ui.IconGlobe, siteInfo.String(), 60, cardColor)
		sb.WriteString(siteCard)
		sb.WriteString("\n\n")
	}
	
	return sb.String()
}

func init() {
	rootCmd.AddCommand(nettestCmd)

	// 添加参数
	nettestCmd.Flags().StringP("url", "u", "", "要测试的站点URL，多个URL用逗号分隔")
	nettestCmd.Flags().IntP("timeout", "t", 0, "设置HTTP请求超时时间(秒)")
	nettestCmd.Flags().BoolP("detailed", "d", false, "显示详细的测试信息")
} 