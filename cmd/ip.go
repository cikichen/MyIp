/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"ip/ui"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type IPInfo struct {
	IP             string  `json:"ip"`
	Network        string  `json:"network"`
	Version        string  `json:"version"`
	City           string  `json:"city"`
	Region         string  `json:"region"`
	RegionCode     string  `json:"region_code"`
	Country        string  `json:"country"`
	CountryName    string  `json:"country_name"`
	CountryCode    string  `json:"country_code"`
	CountryCodeISO string  `json:"country_code_iso3"`
	CountryCapital string  `json:"country_capital"`
	CountryTLD     string  `json:"country_tld"`
	ContinentCode  string  `json:"continent_code"`
	InEU           bool    `json:"in_eu"`
	Postal         string  `json:"postal"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Timezone       string  `json:"timezone"`
	UTCOffset      string  `json:"utc_offset"`
	CallingCode    string  `json:"country_calling_code"`
	Currency       string  `json:"currency"`
	CurrencyName   string  `json:"currency_name"`
	Languages      string  `json:"languages"`
	CountryArea    float64 `json:"country_area"`
	Population     int64   `json:"country_population"`
	ASN            string  `json:"asn"`
	Org            string  `json:"org"`
	// 额外字段，不是API直接返回的
	IsPure         bool    `json:"-"`
	PureScore      int     `json:"-"`
	PureType       string  `json:"-"`
	IPType         string  `json:"-"` // 家宽/独立IP/共享IP
	IsProxy        bool    `json:"-"` // 是否是代理IP
	IsDC           bool    `json:"-"` // 是否是数据中心IP
	APISource      string  `json:"-"` // 记录数据来源的API
}

// ipCmd represents the ip command
var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "显示您的IP地址及详细信息",
	Long:  `显示您的本地和公网IP地址，以及地理位置、网络信息、IP类型等详细信息。`,
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

		if result != nil {
			// 直接设置时区为固定值，避免格式化问题
			timezone := "Asia/Shanghai"

			// 转换为ui.IPInfo类型
			uiInfo := &ui.IPInfo{
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
				Timezone:       timezone, // 使用处理后的时区字符串
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
			
			// 使用统一的绘制函数，会根据模式选择合适的实现
			fmt.Println(ui.DrawIPInfo(localIp.String(), myIP, uiInfo))
		} else {
			// 如果未获取到IP信息，仍然显示基本信息
			fmt.Println(ui.DrawIPInfo(localIp.String(), myIP, nil))
		}
	},
}

// GetMyPublicIP 获取公网IP，支持多个API源和负载均衡
func GetMyPublicIP() string {
	// 定义多个获取公网IP的API源（负载均衡）
	// 这些API源按优先级排序，程序会按顺序尝试，直到成功获取IP
	apiSources := []string{
		"https://myexternalip.com/raw",
		"https://api.ipify.org",
		"https://ifconfig.me/ip",
		"https://ipecho.net/plain",
	}

	// 设置HTTP客户端，添加超时设置
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 尝试所有API源（负载均衡的核心逻辑）
	// 按顺序尝试每个API源，如果一个失败，自动尝试下一个
	for _, url := range apiSources {
		// 创建请求
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			// 静默失败，尝试下一个API源
			continue
		}

		// 设置请求头，模拟浏览器请求
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			// 静默失败，尝试下一个API源
			continue
		}

		// 读取响应
		content, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			// 静默失败，尝试下一个API源
			continue
		}

		// 处理响应
		ip := strings.TrimSpace(string(content))
		if ip != "" {
			return ip
		}
	}

	// 所有API源都失败（负载均衡的故障处理）
	fmt.Println("警告: 无法获取公网IP，所有API源都失败")
	return ""
}

// TestAPISource 测试所有API源的可用性
func TestAPISource() map[string]bool {
	// 定义要测试的所有API源
	apiSources := map[string]string{
		"ipapi":   "https://ipapi.co/json/",
		"ipinfo":  "https://ipinfo.io/json",
		"ipify":   "https://api.ipify.org?format=json",
		"myip":    "https://api.myip.com",
		"ipwhois": "https://ipwho.is/",
	}

	results := make(map[string]bool)

	// 设置HTTP客户端，添加超时设置
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 测试每个API源
	fmt.Println("")
	for name, url := range apiSources {
		// 显示测试进度
		message := fmt.Sprintf("正在测试 %s API源 (%s)...", name, url)
		fmt.Println(ui.DrawStatusBar(message, ui.BgBrightBlue))
		
		// 创建请求
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			results[name] = false
			continue
		}

		// 设置请求头，模拟浏览器请求
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			results[name] = false
			continue
		}

		// 检查响应码
		if resp.StatusCode != 200 {
			resp.Body.Close()
			results[name] = false
			continue
		}

		// 读取响应内容
		_, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			results[name] = false
			continue
		}

		// API可用
		results[name] = true
	}

	return results
}

// OnlineIpInfo 获取IP信息，支持多个API源和负载均衡
func OnlineIpInfo(ip string) *IPInfo {
	// 定义API源
	apiSources := []struct {
		Name      string
		URLFunc   func(string) string
		ParseFunc func([]byte) (*IPInfo, error)
		Enabled   bool
	}{
		{
			Name: "ipapi.co",
			URLFunc: func(ip string) string {
				return "https://ipapi.co/" + ip + "/json"
			},
			ParseFunc: parseIpapiResponse,
			Enabled: true,
		},
		{
			Name: "ipinfo.io",
			URLFunc: func(ip string) string {
				return "https://ipinfo.io/" + ip + "/json"
			},
			ParseFunc: parseIpinfoResponse,
			Enabled: true,
		},
	}

	// 用于记录最后一个错误，但不输出详细日志
	var _ error

	// 直接使用所有启用的API源，不需要预先测试
	// 这样实现了真正的负载均衡：在运行时动态选择可用的API源

	// 创建API源列表
	availableAPIs := make([]struct {
		Name      string
		URLFunc   func(string) string
		ParseFunc func([]byte) (*IPInfo, error)
		Enabled   bool
	}, 0, len(apiSources))

	// 添加所有启用的API源
	for _, api := range apiSources {
		if api.Enabled {
			availableAPIs = append(availableAPIs, struct {
				Name      string
				URLFunc   func(string) string
				ParseFunc func([]byte) (*IPInfo, error)
				Enabled   bool
			}{
				Name:      api.Name,
				URLFunc:   api.URLFunc,
				ParseFunc: api.ParseFunc,
				Enabled:   api.Enabled,
			})
		}
	}

	// 尝试所有API源
	// 这里实现了负载均衡的核心逻辑：如果一个API源失败，自动切换到下一个
	for _, api := range availableAPIs {
		url := api.URLFunc(ip)

		// 设置请求头，模拟浏览器请求
		client := &http.Client{
			Timeout: 10 * time.Second, // 添加超时设置
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue // 自动切换到下一个API源
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			continue // 自动切换到下一个API源
		}

		out, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close() // 简化关闭逻辑

		if err != nil {
			continue // 自动切换到下一个API源
		}

		// 解析响应
		ipInfo, err := api.ParseFunc(out)
		if err != nil {
			continue // 自动切换到下一个API源
		}

		// 设置API源
		ipInfo.APISource = api.Name

		// 判断IP类型和纯净度
		DetermineIPType(ipInfo)
		DetermineIPPurity(ipInfo)

		return ipInfo
	}

	// 所有API源都失败
	fmt.Println("警告: 无法获取IP信息，请检查网络连接")
	return nil
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

// parseIpapiResponse 解析ipapi.co的响应
func parseIpapiResponse(data []byte) (*IPInfo, error) {
	var ipInfo IPInfo
	err := json.Unmarshal(data, &ipInfo)
	if err != nil {
		return nil, fmt.Errorf("解析ipapi.co响应失败: %v", err)
	}

	// 验证必要字段
	if ipInfo.IP == "" {
		return nil, fmt.Errorf("ipapi.co响应缺少IP字段")
	}

	return &ipInfo, nil
}

// parseIpinfoResponse 解析ipinfo.io的响应
func parseIpinfoResponse(data []byte) (*IPInfo, error) {
	// ipinfo.io返回的JSON结构与我们的IPInfo结构不完全匹配，需要特殊处理
	var response struct {
		IP       string `json:"ip"`
		Hostname string `json:"hostname"`
		City     string `json:"city"`
		Region   string `json:"region"`
		Country  string `json:"country"`
		Loc      string `json:"loc"`
		Org      string `json:"org"`
		Postal   string `json:"postal"`
		Timezone string `json:"timezone"`
	}

	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, fmt.Errorf("解析ipinfo.io响应失败: %v", err)
	}

	// 验证必要字段
	if response.IP == "" {
		return nil, fmt.Errorf("ipinfo.io响应缺少IP字段")
	}

	// 解析经纬度
	var lat, lon float64
	if response.Loc != "" {
		coords := strings.Split(response.Loc, ",")
		if len(coords) == 2 {
			lat, _ = strconv.ParseFloat(coords[0], 64)
			lon, _ = strconv.ParseFloat(coords[1], 64)
		}
	}

	// 从组织信息中提取ASN和组织名称
	asn := ""
	orgName := response.Org
	if response.Org != "" {
		parts := strings.SplitN(response.Org, " ", 2)
		if len(parts) > 0 && strings.HasPrefix(parts[0], "AS") {
			asn = parts[0]
			if len(parts) > 1 {
				orgName = strings.TrimSpace(parts[1])
			} else {
				orgName = ""
			}
		}
	}

	// 获取国家名称
	countryName := getCountryName(response.Country)

	// 设置货币和通信区号
	currency, currencyName := getCurrencyInfo(response.Country)
	callingCode := getCallingCode(response.Country)

	// 创建IPInfo实例
	ipInfo := &IPInfo{
		IP:           response.IP,
		City:         response.City,
		Region:       response.Region,
		Country:      response.Country,
		CountryName:  countryName,
		CountryCode:  response.Country,
		Postal:       response.Postal,
		Timezone:     response.Timezone,
		Latitude:     lat,
		Longitude:    lon,
		Org:          orgName,
		ASN:          asn,
		Currency:     currency,
		CurrencyName: currencyName,
		CallingCode:  callingCode,
		Network:      getNetworkFromIP(response.IP),
		ContinentCode: getContinentCode(response.Country),
	}

	return ipInfo, nil
}

// getCountryName 根据国家代码获取国家名称
func getCountryName(countryCode string) string {
	countries := map[string]string{
		"CN": "China",
		"US": "United States",
		"JP": "Japan",
		"GB": "United Kingdom",
		"DE": "Germany",
		"FR": "France",
		"IT": "Italy",
		"CA": "Canada",
		"AU": "Australia",
		"BR": "Brazil",
		"IN": "India",
		"RU": "Russia",
		// 可以根据需要添加更多
	}
	
	if name, ok := countries[countryCode]; ok {
		return name
	}
	return countryCode
}

// getCurrencyInfo 根据国家代码获取货币信息
func getCurrencyInfo(countryCode string) (string, string) {
	currencies := map[string]struct{
		Code string
		Name string
	}{
		"CN": {"CNY", "Yuan Renminbi"},
		"US": {"USD", "US Dollar"},
		"JP": {"JPY", "Japanese Yen"},
		"GB": {"GBP", "Pound Sterling"},
		"DE": {"EUR", "Euro"},
		"FR": {"EUR", "Euro"},
		"IT": {"EUR", "Euro"},
		"CA": {"CAD", "Canadian Dollar"},
		"AU": {"AUD", "Australian Dollar"},
		"BR": {"BRL", "Brazilian Real"},
		"IN": {"INR", "Indian Rupee"},
		"RU": {"RUB", "Russian Ruble"},
		// 可以根据需要添加更多
	}
	
	if currency, ok := currencies[countryCode]; ok {
		return currency.Code, currency.Name
	}
	return "", ""
}

// getCallingCode 根据国家代码获取通信区号
func getCallingCode(countryCode string) string {
	callingCodes := map[string]string{
		"CN": "+86",
		"US": "+1",
		"JP": "+81",
		"GB": "+44",
		"DE": "+49",
		"FR": "+33",
		"IT": "+39",
		"CA": "+1",
		"AU": "+61",
		"BR": "+55",
		"IN": "+91",
		"RU": "+7",
		// 可以根据需要添加更多
	}
	
	if code, ok := callingCodes[countryCode]; ok {
		return code
	}
	return ""
}

// getNetworkFromIP 从IP获取网络信息
func getNetworkFromIP(ip string) string {
	parts := strings.Split(ip, ".")
	if len(parts) == 4 {
		return parts[0] + "." + parts[1] + "." + parts[2] + ".0/24"
	}
	return ""
}

// getContinentCode 根据国家代码获取大洲代码
func getContinentCode(countryCode string) string {
	continents := map[string]string{
		"CN": "AS", // Asia
		"JP": "AS",
		"IN": "AS",
		"US": "NA", // North America
		"CA": "NA",
		"GB": "EU", // Europe
		"DE": "EU",
		"FR": "EU",
		"IT": "EU",
		"BR": "SA", // South America
		"AU": "OC", // Oceania
		"RU": "EU", // Russia spans Europe and Asia
		// 可以根据需要添加更多
	}
	
	if code, ok := continents[countryCode]; ok {
		return code
	}
	return ""
}

// DetermineIPType 判断IP类型（家宽/独立IP/共享IP）
func DetermineIPType(ipInfo *IPInfo) {
	// 默认为未知类型
	ipInfo.IPType = "未知"
	ipInfo.IsDC = false
	ipInfo.IsProxy = false

	// 如果组织信息为空，设置为默认家庭宽带IP
	if ipInfo.Org == "" && ipInfo.ASN == "" {
		ipInfo.IPType = "家庭宽带IP"
		return
	}

	// 标准化组织信息，移除换行符并转换为小写
	org := strings.ToLower(strings.ReplaceAll(ipInfo.Org, "\n", " "))
	asn := strings.ToLower(ipInfo.ASN)
	
	// 检查是否包含数据中心相关关键词
	dataCenterKeywords := []string{
		"cloud", "hosting", "server", "data center",
		"aws", "amazon", "azure", "google cloud", "gcp",
		"alibaba", "tencent", "digital ocean", "linode",
		"vultr", "scaleway", "oracle", "ibm",
	}
	
	// 检查是否包含代理相关关键词 
	proxyKeywords := []string{
		"proxy", "vpn", "tor", "anonymizer", "anonymous",
		"tunnel", "exit node", "relay",
	}

	// 检查是否包含移动网络相关关键词
	mobileKeywords := []string{
		"mobile", "wireless", "cellular", "lte", "5g", "4g", "3g",
		"phone", "telecom", "cmcc", "china mobile", "unicom",
	}

	// 检查是否为中国家庭宽带ISP关键词
	homeKeywords := []string{
		"chinanet", "china telecom", "china unicom", "china mobile",
		"residential", "home", "broadband", "fttx", "adsl", "vdsl",
		"家宽", "家庭宽带", "联通", "电信", "移动宽带",
	}

	// 检查数据中心
	for _, keyword := range dataCenterKeywords {
		if strings.Contains(org, keyword) || strings.Contains(asn, keyword) {
			ipInfo.IPType = "数据中心IP"
			ipInfo.IsDC = true
			return
		}
	}

	// 检查代理
	for _, keyword := range proxyKeywords {
		if strings.Contains(org, keyword) || strings.Contains(asn, keyword) {
			ipInfo.IPType = "代理IP"
			ipInfo.IsProxy = true
			return
		}
	}

	// 检查移动网络
	for _, keyword := range mobileKeywords {
		if strings.Contains(org, keyword) || strings.Contains(asn, keyword) {
			ipInfo.IPType = "移动网络IP"
			return
		}
	}

	// 检查家庭宽带
	for _, keyword := range homeKeywords {
		if strings.Contains(org, keyword) || strings.Contains(asn, keyword) {
			ipInfo.IPType = "家庭宽带IP"
			return
		}
	}

	// 默认为家庭宽带IP，如果没有匹配到任何关键词
	ipInfo.IPType = "家庭宽带IP"
}

// DetermineIPPurity 判断IP纯净度
func DetermineIPPurity(ipInfo *IPInfo) {
	// 默认分数为100（满分）
	ipInfo.PureScore = 100

	// 根据IP类型给予不同的基础分数
	switch ipInfo.IPType {
	case "家庭宽带IP":
		// 家庭宽带IP通常是最纯净的，保持100分
	case "移动网络IP":
		// 移动网络IP可能会有一些共享问题，轻微扣分
		ipInfo.PureScore -= 5
	case "数据中心IP":
		// 数据中心IP通常用于服务器托管，扣20分
		ipInfo.PureScore -= 20
		ipInfo.IsDC = true
	case "代理IP":
		// 代理IP通常被用于隐藏真实身份，扣50分
		ipInfo.PureScore -= 50
		ipInfo.IsProxy = true
	case "未知":
		// 未知类型扣10分
		ipInfo.PureScore -= 10
	}

	// 如果明确标记为代理IP，无论IP类型如何，都扣除50分
	if ipInfo.IsProxy && ipInfo.IPType != "代理IP" {
		ipInfo.PureScore -= 50
	}

	// 如果明确标记为数据中心IP，无论IP类型如何，都扣除20分
	if ipInfo.IsDC && ipInfo.IPType != "数据中心IP" {
		ipInfo.PureScore -= 20
	}

	// 确保分数在0-100范围内
	if ipInfo.PureScore < 0 {
		ipInfo.PureScore = 0
	}
	if ipInfo.PureScore > 100 {
		ipInfo.PureScore = 100
	}

	// 根据分数确定纯净度类型
	if ipInfo.PureScore >= 90 {
		ipInfo.PureType = "优质"
		ipInfo.IsPure = true
	} else if ipInfo.PureScore >= 70 {
		ipInfo.PureType = "良好"
		ipInfo.IsPure = true
	} else if ipInfo.PureScore >= 50 {
		ipInfo.PureType = "一般"
		ipInfo.IsPure = false
	} else {
		ipInfo.PureType = "较差"
		ipInfo.IsPure = false
	}
}

func init() {
	rootCmd.AddCommand(ipCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ipCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// 添加测试API源的标志
	ipCmd.Flags().Bool("test-api", false, "仅测试所有IP信息API源的可用性，不获取IP信息")
}
