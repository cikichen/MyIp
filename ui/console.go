package ui

import (
	"fmt"
	"ip/network"
	"sort"
	"strings"
	"time"
)

// 定义ANSI颜色代码
const (
	Reset      = "\033[0m"
	Bold       = "\033[1m"
	Dim        = "\033[2m"
	Italic     = "\033[3m"
	Underline  = "\033[4m"
	Blink      = "\033[5m"
	Reverse    = "\033[7m"
	Hidden     = "\033[8m"
	Black      = "\033[30m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Magenta    = "\033[35m"
	Cyan       = "\033[36m"
	White      = "\033[37m"
	BgBlack    = "\033[40m"
	BgRed      = "\033[41m"
	BgGreen    = "\033[42m"
	BgYellow   = "\033[43m"
	BgBlue     = "\033[44m"
	BgMagenta  = "\033[45m"
	BgCyan     = "\033[46m"
	BgWhite    = "\033[47m"

	// 明亮颜色
	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"

	// 明亮背景色
	BgBrightBlack   = "\033[100m"
	BgBrightRed     = "\033[101m"
	BgBrightGreen   = "\033[102m"
	BgBrightYellow  = "\033[103m"
	BgBrightBlue    = "\033[104m"
	BgBrightMagenta = "\033[105m"
	BgBrightCyan    = "\033[106m"
	BgBrightWhite   = "\033[107m"
)

// 定义常用图标
const (
	IconCheck       = "✓"
	IconCross       = "✗"
	IconWarning     = "⚠"
	IconInfo        = "ℹ"
	IconStar        = "★"
	IconHeart       = "♥"
	IconArrowRight  = "→"
	IconArrowLeft   = "←"
	IconArrowUp     = "↑"
	IconArrowDown   = "↓"
	IconGlobe       = "🌐"
	IconComputer    = "💻"
	IconNetwork     = "📡"
	IconLocation    = "📍"
	IconClock       = "🕒"
	IconSpeed       = "⚡"
	IconLock        = "🔒"
	IconUnlock      = "🔓"
	IconServer      = "🖥️"
	IconCloud       = "☁️"
	IconHome        = "🏠"
	IconBuilding    = "🏢"
	IconFlag        = "🚩"
	IconPing        = "📶"
	IconLoading     = "⏳"
)

// IPInfo 结构体，与cmd/ip.go中的IPInfo保持一致
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

// DrawIPInfo 绘制IP信息
func DrawIPInfo(localIP string, publicIP string, ipInfo *IPInfo) string {
	return RenderIPInfoWithLipgloss(localIP, publicIP, ipInfo)
}

// DrawBox 绘制一个带标题的框
func DrawBox(title string, content string, width int, color string) string {
	lines := strings.Split(content, "\n")

	// 计算框的宽度
	boxWidth := width
	for _, line := range lines {
		if len(line) > boxWidth {
			boxWidth = len(line)
		}
	}
	boxWidth += 4 // 左右各加2个空格

	// 绘制顶部边框和标题
	result := color + "╭"
	titleStart := (boxWidth - len(title) - 2) / 2
	for i := 0; i < boxWidth-2; i++ {
		if i == titleStart && title != "" {
			result += "┤ " + Bold + title + Reset + color + " ├"
			i += len(title) + 3
		} else {
			result += "─"
		}
	}
	result += "╮\n"

	// 绘制内容
	for _, line := range lines {
		result += "│ " + line
		for i := len(line); i < boxWidth-4; i++ {
			result += " "
		}
		result += " │\n"
	}

	// 绘制底部边框
	result += "╰"
	for i := 0; i < boxWidth-2; i++ {
		result += "─"
	}
	result += "╯" + Reset

	return result
}

// DrawFancyBox 绘制一个更美观的带图标和渐变色的框
func DrawFancyBox(title string, icon string, content string, width int, primaryColor string, secondaryColor string) string {
	lines := strings.Split(content, "\n")

	// 计算框的宽度
	boxWidth := width
	for _, line := range lines {
		if len(line) > boxWidth {
			boxWidth = len(line)
		}
	}
	boxWidth += 4 // 左右各加2个空格

	// 绘制顶部边框和标题
	result := primaryColor + "╭"
	titleWithIcon := icon + " " + title + " " + icon
	titleStart := (boxWidth - len(titleWithIcon)) / 2
	for i := 0; i < boxWidth-2; i++ {
		if i == titleStart && title != "" {
			result += "┤" + Bold + BgBrightBlack + " " + titleWithIcon + " " + Reset + primaryColor + "├"
			i += len(titleWithIcon) + 1
		} else {
			result += "━"
		}
	}
	result += "╮\n"

	// 绘制内容
	for i, line := range lines {
		// 交替使用主色和次色，创建渐变效果
		if i % 2 == 0 {
			result += primaryColor
		} else {
			result += secondaryColor
		}

		result += "│ " + line
		for j := len(line); j < boxWidth-4; j++ {
			result += " "
		}
		result += " │\n"
	}

	// 绘制底部边框
	result += primaryColor + "╰"
	for i := 0; i < boxWidth-2; i++ {
		result += "━"
	}
	result += "╯" + Reset

	return result
}

// DrawProgressBar 绘制一个进度条
func DrawProgressBar(value float64, max float64, width int, color string) string {
	percent := value / max
	if percent > 1.0 {
		percent = 1.0
	}

	// 计算填充的字符数
	fillWidth := int(percent * float64(width))

	// 构建进度条
	result := "["
	for i := 0; i < width; i++ {
		if i < fillWidth {
			result += color + "█" + Reset
		} else {
			result += "░"
		}
	}
	result += "] " + fmt.Sprintf("%.1f%%", percent*100)

	return result
}

// DrawSpeedometer 绘制一个速度计
func DrawSpeedometer(value float64, max float64, good float64, medium float64) string {
	// 根据值选择颜色
	var color string
	if value <= good {
		color = Green
	} else if value <= medium {
		color = Yellow
	} else {
		color = Red
	}

	// 计算百分比
	percent := value / max
	if percent > 1.0 {
		percent = 1.0
	}

	// 构建速度计
	result := "["
	segments := 10
	for i := 0; i < segments; i++ {
		threshold := float64(i) / float64(segments)
		if percent >= threshold {
			result += color + "■" + Reset
		} else {
			result += "□"
		}
	}
	result += "] " + color + fmt.Sprintf("%.1f", value) + Reset

	return result
}

// DrawTable 绘制一个表格
func DrawTable(headers []string, rows [][]string, color string) string {
	// 计算每列的最大宽度
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// 绘制表头
	result := color + "┌"
	for i, width := range colWidths {
		for j := 0; j < width+2; j++ {
			result += "─"
		}
		if i < len(colWidths)-1 {
			result += "┬"
		}
	}
	result += "┐\n"

	// 绘制表头内容
	result += "│"
	for i, header := range headers {
		result += " " + Bold + header + Reset + color
		for j := len(header); j < colWidths[i]+1; j++ {
			result += " "
		}
		result += "│"
	}
	result += "\n"

	// 绘制表头和内容的分隔线
	result += "├"
	for i, width := range colWidths {
		for j := 0; j < width+2; j++ {
			result += "─"
		}
		if i < len(colWidths)-1 {
			result += "┼"
		}
	}
	result += "┤\n"

	// 绘制表格内容
	for _, row := range rows {
		result += "│"
		for i, cell := range row {
			if i < len(colWidths) {
				result += " " + cell
				for j := len(cell); j < colWidths[i]+1; j++ {
					result += " "
				}
				result += "│"
			}
		}
		result += "\n"
	}

	// 绘制底部边框
	result += "└"
	for i, width := range colWidths {
		for j := 0; j < width+2; j++ {
			result += "─"
		}
		if i < len(colWidths)-1 {
			result += "┴"
		}
	}
	result += "┘" + Reset

	return result
}

// PrintIPInfo 美观地打印IP信息
func PrintIPInfo(localIP string, publicIP string, ipInfo *IPInfo) {
	// 绘制标题
	titleBar := BgBrightBlue + BrightWhite + Bold + "  " + IconGlobe + " IP信息查询结果 " + IconGlobe + "  " + Reset
	fmt.Println("\n" + titleBar)

	// 绘制分隔线
	separator := BrightCyan
	for i := 0; i < 60; i++ {
		separator += "━"
	}
	separator += Reset
	fmt.Println(separator)

	// 基本IP信息
	basicInfo := fmt.Sprintf("%s%s 本地IP:%s %s\n%s%s 公网IP:%s %s",
		Bold, IconComputer, Reset, BrightCyan + localIP + Reset,
		Bold, IconNetwork, Reset, BrightYellow + Bold + publicIP + Reset)
	fmt.Println(DrawFancyBox("基本信息", IconInfo, basicInfo, 50, BrightBlue, Blue))

	if ipInfo == nil {
		fmt.Println(BgRed + White + Bold + "\n  " + IconWarning + " 无法获取IP详细信息，请检查网络连接  " + Reset)
		return
	}

	// 地理位置信息
	countryInfo := BrightYellow + Bold + ipInfo.CountryName + Reset
	if ipInfo.CountryCode != "" {
		countryInfo += " (" + BrightWhite + ipInfo.CountryCode + Reset + ")"
	}

	geoInfo := fmt.Sprintf("%s%s 国家/地区:%s %s\n%s%s 省/州:%s %s\n%s%s 城市:%s %s\n%s%s 经纬度:%s %.4f, %.4f\n%s%s 时区:%s %s",
		Bold, IconFlag, Reset, countryInfo,
		Bold, IconLocation, Reset, BrightYellow + ipInfo.Region + Reset,
		Bold, IconBuilding, Reset, BrightYellow + ipInfo.City + Reset,
		Bold, IconGlobe, Reset, BrightCyan + fmt.Sprintf("%.4f", ipInfo.Latitude) + Reset + ", " + BrightCyan + fmt.Sprintf("%.4f", ipInfo.Longitude) + Reset,
		Bold, IconClock, Reset, BrightMagenta + ipInfo.Timezone + Reset)
	fmt.Println("\n" + DrawFancyBox("地理位置", IconLocation, geoInfo, 50, BrightYellow, Yellow))

	// 网络信息
	netInfo := fmt.Sprintf("%s%s ISP:%s %s\n%s%s AS号:%s %s\n%s%s API来源:%s %s",
		Bold, IconServer, Reset, BrightGreen + ipInfo.Org + Reset,
		Bold, IconNetwork, Reset, BrightCyan + ipInfo.ASN + Reset,
		Bold, IconCloud, Reset, BrightBlue + ipInfo.APISource + Reset)
	fmt.Println("\n" + DrawFancyBox("网络信息", IconServer, netInfo, 50, BrightGreen, Green))

	// IP类型信息
	var dcStatus, proxyStatus string
	var pureIcon, dcIcon, proxyIcon string

	if ipInfo.IsPure {
		pureIcon = IconCheck
	} else {
		pureIcon = IconCross
	}

	if ipInfo.IsDC {
		dcStatus = Yellow + Bold + "是" + Reset
		dcIcon = IconServer
	} else {
		dcStatus = Green + Bold + "否" + Reset
		dcIcon = IconHome
	}

	if ipInfo.IsProxy {
		proxyStatus = Red + Bold + "是" + Reset
		proxyIcon = IconUnlock
	} else {
		proxyStatus = Green + Bold + "否" + Reset
		proxyIcon = IconLock
	}

	// 根据纯净度分数选择颜色和图标
	var scoreColor string
	var scoreIcon string
	if ipInfo.PureScore >= 90 {
		scoreColor = BrightGreen
		scoreIcon = IconStar + IconStar + IconStar
	} else if ipInfo.PureScore >= 70 {
		scoreColor = BrightYellow
		scoreIcon = IconStar + IconStar
	} else {
		scoreColor = BrightRed
		scoreIcon = IconStar
	}

	// 绘制纯净度进度条
	pureBar := DrawProgressBar(float64(ipInfo.PureScore), 100.0, 20, scoreColor)

	// 根据IP类型选择图标
	var ipTypeIcon string
	if ipInfo.IPType == "数据中心IP" {
		ipTypeIcon = IconServer
	} else if ipInfo.IPType == "代理IP" {
		ipTypeIcon = IconUnlock
	} else {
		ipTypeIcon = IconHome
	}

	typeInfo := fmt.Sprintf("%s%s IP类型:%s %s\n%s%s 数据中心IP:%s %s\n%s%s 代理IP:%s %s\n%s%s 纯净度:%s %s %s\n%s 纯净IP评分: %s",
		Bold, ipTypeIcon, Reset, BrightMagenta + Bold + ipInfo.IPType + Reset,
		Bold, dcIcon, Reset, dcStatus,
		Bold, proxyIcon, Reset, proxyStatus,
		Bold, scoreIcon, Reset, ipInfo.PureType, scoreColor + fmt.Sprintf("(%d分)", ipInfo.PureScore) + Reset,
		pureIcon, pureBar)
	fmt.Println("\n" + DrawFancyBox("IP类型分析", IconInfo, typeInfo, 50, BrightMagenta, Magenta))

	// 底部
	fmt.Println("\n" + BgBrightBlack + BrightWhite + "  " + IconClock + " 数据更新时间: " + time.Now().Format("2006-01-02 15:04:05") + "  " + Reset)
	fmt.Println("")
}

// PrintNetworkTests 美观地打印网络测试结果
func PrintNetworkTests(results []network.SiteTestResult) {
	// 绘制标题
	titleBar := BgBrightGreen + Black + Bold + "  " + IconNetwork + " 常用站点连通性测试 " + IconNetwork + "  " + Reset
	fmt.Println("\n" + titleBar)

	// 绘制分隔线
	separator := BrightGreen
	for i := 0; i < 80; i++ {
		separator += "━"
	}
	separator += Reset
	fmt.Println(separator)

	// 按名称排序结果
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	// 创建表格头
	headers := []string{
		BgBrightBlack + White + " 站点 " + Reset,
		BgBrightBlack + White + " 状态 " + Reset,
		BgBrightBlack + White + " 响应时间 " + Reset,
		BgBrightBlack + White + " Ping延迟 " + Reset,
		BgBrightBlack + White + " 丢包率 " + Reset,
		BgBrightBlack + White + " DNS解析 " + Reset,
		BgBrightBlack + White + " 连接时间 " + Reset,
	}

	// 创建表格行
	rows := make([][]string, 0, len(results))

	// 计算可访问站点数量
	accessibleCount := 0
	totalSites := len(results)

	// 计算平均响应时间和Ping延迟
	var totalRespTime, totalPingTime time.Duration
	respTimeCount, pingTimeCount := 0, 0

	for _, result := range results {
		// 状态和图标
		var status string
		var statusIcon string
		if result.Accessible {
			accessibleCount++
			status = BgGreen + Black + Bold + " 可访问 " + Reset
			statusIcon = Green + IconCheck + Reset
		} else {
			status = BgRed + White + Bold + " 不可访问 " + Reset
			statusIcon = Red + IconCross + Reset
		}

		// 响应时间
		var respTime string
		if result.Accessible {
			// 根据响应时间选择颜色
			var timeColor string
			if result.ResponseTime < 300*time.Millisecond {
				timeColor = BrightGreen
			} else if result.ResponseTime < 1000*time.Millisecond {
				timeColor = BrightYellow
			} else {
				timeColor = BrightRed
			}
			respTime = fmt.Sprintf("%s%.2f秒%s", timeColor, result.ResponseTime.Seconds(), Reset)

			// 累计总响应时间
			totalRespTime += result.ResponseTime
			respTimeCount++
		} else {
			respTime = BrightRed + "超时" + Reset
		}

		// Ping延迟 - 使用速度计显示
		var pingTime string
		if result.PingTime > 0 {
			// 根据Ping延迟选择阈值
			var goodThreshold float64 = 50
			var mediumThreshold float64 = 150

			// 使用速度计显示Ping延迟
			pingValue := float64(result.PingTime) / float64(time.Millisecond)
			pingMeter := DrawSpeedometer(pingValue, 300, goodThreshold, mediumThreshold)
			pingTime = fmt.Sprintf("%s (%.1fms)", pingMeter, pingValue)

			// 累计总Ping延迟
			totalPingTime += result.PingTime
			pingTimeCount++
		} else {
			pingTime = BrightRed + "超时" + Reset
		}

		// 丢包率 - 使用进度条显示
		var lossRate string
		if result.PingLoss == 0 {
			lossRate = BrightGreen + "0%" + Reset
		} else {
			lossPercent := result.PingLoss * 100
			lossBar := DrawProgressBar(lossPercent, 100, 10, BrightRed)
			lossRate = fmt.Sprintf("%s (%.0f%%)", lossBar, lossPercent)
		}

		// DNS解析时间
		var dnsTime string
		if result.Accessible && result.DNSTime > 0 {
			dnsValue := result.DNSTime.Seconds()
			var dnsColor string
			if dnsValue < 0.1 {
				dnsColor = BrightGreen
			} else if dnsValue < 0.3 {
				dnsColor = BrightYellow
			} else {
				dnsColor = BrightRed
			}
			dnsTime = fmt.Sprintf("%s%.2f秒%s", dnsColor, dnsValue, Reset)
		} else {
			dnsTime = Dim + "-" + Reset
		}

		// 连接时间
		var connTime string
		if result.Accessible && result.ConnectTime > 0 {
			connValue := result.ConnectTime.Seconds()
			var connColor string
			if connValue < 0.1 {
				connColor = BrightGreen
			} else if connValue < 0.3 {
				connColor = BrightYellow
			} else {
				connColor = BrightRed
			}
			connTime = fmt.Sprintf("%s%.2f秒%s", connColor, connValue, Reset)
		} else {
			connTime = Dim + "-" + Reset
		}

		// 添加行
		rows = append(rows, []string{
			BrightWhite + Bold + result.Name + Reset,
			statusIcon + " " + status,
			respTime,
			pingTime,
			lossRate,
			dnsTime,
			connTime,
		})
	}

	// 显示表格
	fmt.Println(DrawTable(headers, rows, BrightCyan))

	// 显示统计信息
	fmt.Println()

	// 计算平均值
	var avgRespTime, avgPingTime float64
	if respTimeCount > 0 {
		avgRespTime = totalRespTime.Seconds() / float64(respTimeCount)
	}
	if pingTimeCount > 0 {
		avgPingTime = float64(totalPingTime) / float64(pingTimeCount) / float64(time.Millisecond)
	}

	// 可访问率
	accessRate := float64(accessibleCount) / float64(totalSites) * 100
	var accessRateColor string
	if accessRate >= 90 {
		accessRateColor = BrightGreen
	} else if accessRate >= 70 {
		accessRateColor = BrightYellow
	} else {
		accessRateColor = BrightRed
	}

	// 平均响应时间颜色
	var avgRespTimeColor string
	if avgRespTime < 0.3 {
		avgRespTimeColor = BrightGreen
	} else if avgRespTime < 1.0 {
		avgRespTimeColor = BrightYellow
	} else {
		avgRespTimeColor = BrightRed
	}

	// 平均Ping延迟颜色
	var avgPingTimeColor string
	if avgPingTime < 50 {
		avgPingTimeColor = BrightGreen
	} else if avgPingTime < 150 {
		avgPingTimeColor = BrightYellow
	} else {
		avgPingTimeColor = BrightRed
	}

	// 构建统计信息
	statsInfo := fmt.Sprintf(
		"%s站点可访问率:%s %s%.1f%%%s (%d/%d)\n%s平均响应时间:%s %s%.2f秒%s\n%s平均Ping延迟:%s %s%.1fms%s",
		Bold, Reset, accessRateColor, accessRate, Reset, accessibleCount, totalSites,
		Bold, Reset, avgRespTimeColor, avgRespTime, Reset,
		Bold, Reset, avgPingTimeColor, avgPingTime, Reset,
	)

	// 显示统计信息框
	fmt.Println(DrawFancyBox("网络统计", IconSpeed, statsInfo, 50, BrightBlue, Blue))

	// 底部
	fmt.Println("\n" + BgBrightBlack + BrightWhite + "  " + IconClock + " 测试完成时间: " + time.Now().Format("2006-01-02 15:04:05") + "  " + Reset)
	fmt.Println("")
}

// DrawStatusBar 绘制一个状态栏（确保右侧边框对齐）
func DrawStatusBar(message string, color string) string {
	// 固定宽度
	totalWidth := 80
	
	// 准备图标和消息
	iconAndMessage := IconInfo + " " + message
	cleanText := stripANSI(iconAndMessage)
	
	// 计算填充
	padding := totalWidth - len(cleanText) - 2  // 减去前后两个空格
	if padding < 0 {
		padding = 0
	}
	
	// 构建状态栏
	return " " + Bold + White + iconAndMessage + Reset + strings.Repeat(" ", padding) + " "
}

// DrawNotice 绘制一个通知框（确保右侧边框对齐）
func DrawNotice(message string, icon string, color string) string {
	// 固定宽度
	totalWidth := 80
	
	// 创建结果
	var result []string
	
	// 添加顶部边框
	result = append(result, color + "+" + strings.Repeat("-", totalWidth-2) + "+" + Reset)
	
	// 准备消息内容
	iconAndMessage := icon + " " + message
	cleanText := stripANSI(iconAndMessage)
	
	// 截断过长内容
	displayText := iconAndMessage
	if len(cleanText) > totalWidth-4 {
		maxVisible := totalWidth - 7  // 为省略号保留空间
		if maxVisible > 0 {
			visibleText := cleanText[:maxVisible] + "..."
			displayText = strings.Replace(iconAndMessage, cleanText, visibleText, 1)
			cleanText = stripANSI(displayText)
		}
	}
	
	// 计算填充
	padding := totalWidth - 3 - len(cleanText)
	if padding < 0 {
		padding = 0
	}
	
	// 添加消息行
	messageLine := color + "|" + Reset + " " + Bold + White + displayText + Reset + strings.Repeat(" ", padding) + color + "|" + Reset
	result = append(result, messageLine)
	
	// 添加底部边框
	result = append(result, color + "+" + strings.Repeat("-", totalWidth-2) + "+" + Reset)
	
	// 组合所有行
	return strings.Join(result, "\n")
}
