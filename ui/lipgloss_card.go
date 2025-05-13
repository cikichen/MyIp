package ui

import (
	"fmt"
	"ip/network"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// 使用 lipgloss 定义基本样式
var (
	// 颜色定义
	primaryColor   = lipgloss.Color("#5F87FF") // 蓝色
	secondaryColor = lipgloss.Color("#87D7FF") // 青色
	accentColor    = lipgloss.Color("#FFD787") // 黄色
	infoColor      = lipgloss.Color("#87FF87") // 绿色
	warnColor      = lipgloss.Color("#FFAF5F") // 橙色
	errorColor     = lipgloss.Color("#FF5F5F") // 红色
	grayColor      = lipgloss.Color("#808080") // 灰色
	
	// 基本卡片样式
	baseCardStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(primaryColor).
		Padding(0, 1).
		Width(78)
		
	// 标题样式
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor)
		
	// 数据标签样式
	labelStyle = lipgloss.NewStyle().
		Bold(true)
		
	// 数据值样式（普通）
	valueStyle = lipgloss.NewStyle().
		Foreground(secondaryColor)
		
	// 强调数据值样式
	accentValueStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Bold(true)
		
	// 状态样式 - 良好
	goodStatusStyle = lipgloss.NewStyle().
		Foreground(infoColor)
		
	// 状态样式 - 警告
	warnStatusStyle = lipgloss.NewStyle().
		Foreground(warnColor)
		
	// 状态样式 - 错误
	errorStatusStyle = lipgloss.NewStyle().
		Foreground(errorColor)
)

// DrawLipglossCard 使用 lipgloss 绘制一个卡片
func DrawLipglossCard(title string, icon string, content string, color lipgloss.Color) string {
	cardStyle := baseCardStyle.Copy().BorderForeground(color)
	
	// 创建标题
	cardTitle := titleStyle.Copy().Foreground(color).Render(icon + " " + title)
	
	// 创建内容
	return cardStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		cardTitle, 
		content,
	))
}

// RenderIPInfoWithLipgloss 使用 lipgloss 渲染 IP 信息
func RenderIPInfoWithLipgloss(localIP string, publicIP string, ipInfo *IPInfo) string {
	// 基本IP信息卡片
	basicInfo := lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf("%s 本地IP:  %s", 
			labelStyle.Render(IconComputer), 
			valueStyle.Render(localIP)),
		fmt.Sprintf("%s 公网IP:  %s", 
			labelStyle.Render(IconNetwork), 
			accentValueStyle.Render(publicIP)),
	)
	basicCard := DrawLipglossCard("基本信息", IconInfo, basicInfo, lipgloss.Color("#5F87FF"))

	if ipInfo == nil {
		errorCard := DrawLipglossCard("错误", IconWarning, 
			errorStatusStyle.Render("无法获取IP详细信息，请检查网络连接"), 
			lipgloss.Color("#FF5F5F"))
		return lipgloss.JoinVertical(lipgloss.Left, basicCard, errorCard)
	}

	// 地理位置信息卡片
	countryInfo := accentValueStyle.Render(ipInfo.CountryName)
	if ipInfo.CountryCode != "" {
		countryInfo = lipgloss.JoinHorizontal(
			lipgloss.Left,
			countryInfo,
			" (",
			valueStyle.Render(ipInfo.CountryCode),
			")",
		)
	}

	// 处理经纬度显示
	latStr := "未知"
	lonStr := "未知"
	if ipInfo.Latitude != 0 || ipInfo.Longitude != 0 {
		latStr = fmt.Sprintf("%.4f", ipInfo.Latitude)
		lonStr = fmt.Sprintf("%.4f", ipInfo.Longitude)
	}

	// 处理时区显示
	timezone := ipInfo.Timezone
	if timezone == "" || strings.Contains(timezone, "%!") {
		timezone = "Asia/Shanghai"
	}

	geoInfo := lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf("%s 国家/地区: %s", labelStyle.Render(IconFlag), countryInfo),
		fmt.Sprintf("%s 省/州:   %s", labelStyle.Render(IconLocation), accentValueStyle.Render(ipInfo.Region)),
		fmt.Sprintf("%s 城市:    %s", labelStyle.Render(IconBuilding), accentValueStyle.Render(ipInfo.City)),
		fmt.Sprintf("%s 经纬度:   %s, %s", labelStyle.Render(IconGlobe), valueStyle.Render(latStr), valueStyle.Render(lonStr)),
		fmt.Sprintf("%s 时区:    %s", labelStyle.Render(IconClock), valueStyle.Render(timezone)),
	)
	geoCard := DrawLipglossCard("地理位置", IconLocation, geoInfo, lipgloss.Color("#FFAF5F"))

	// 格式化网络信息 - 处理可能存在的换行
	orgInfo := strings.ReplaceAll(ipInfo.Org, "\n", " ")
	
	// 网络信息卡片
	netInfo := lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf("%s ISP:    %s", labelStyle.Render(IconServer), valueStyle.Render(orgInfo)),
		fmt.Sprintf("%s AS号:   %s", labelStyle.Render(IconNetwork), valueStyle.Render(ipInfo.ASN)),
		fmt.Sprintf("%s API源:  %s", labelStyle.Render(IconCloud), valueStyle.Render(ipInfo.APISource)),
	)
	netCard := DrawLipglossCard("网络信息", IconServer, netInfo, lipgloss.Color("#87FF87"))

	// IP类型信息卡片
	var dcStatus, proxyStatus string
	var pureIcon, dcIcon, proxyIcon string

	if ipInfo.IsPure {
		pureIcon = IconCheck
	} else {
		pureIcon = IconCross
	}

	if ipInfo.IsDC {
		dcStatus = warnStatusStyle.Render("是")
		dcIcon = IconServer
	} else {
		dcStatus = goodStatusStyle.Render("否")
		dcIcon = IconHome
	}

	if ipInfo.IsProxy {
		proxyStatus = errorStatusStyle.Render("是")
		proxyIcon = IconWarning
	} else {
		proxyStatus = goodStatusStyle.Render("否")
		proxyIcon = IconCheck
	}

	typeInfo := lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf("%s IP类型:   %s", labelStyle.Render(IconInfo), valueStyle.Render(ipInfo.IPType)),
		fmt.Sprintf("%s 数据中心:  %s", labelStyle.Render(dcIcon), dcStatus),
		fmt.Sprintf("%s 代理IP:   %s", labelStyle.Render(proxyIcon), proxyStatus),
		fmt.Sprintf("%s 纯净度:   %s (%s分)", 
			labelStyle.Render(pureIcon), 
			accentValueStyle.Render(ipInfo.PureType), 
			fmt.Sprintf("%d", ipInfo.PureScore)),
	)
	typeCard := DrawLipglossCard("IP类型", IconNetwork, typeInfo, lipgloss.Color("#5FD7FF"))

	// 其他详细信息卡片
	detailInfo := lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf("%s 货币:    %s (%s)", 
			labelStyle.Render(IconStar), 
			accentValueStyle.Render(ipInfo.Currency), 
			valueStyle.Render(ipInfo.CurrencyName)),
		fmt.Sprintf("%s 通信区号:  %s", labelStyle.Render(IconInfo), valueStyle.Render(ipInfo.CallingCode)),
		fmt.Sprintf("%s 网络:    %s", labelStyle.Render(IconGlobe), valueStyle.Render(ipInfo.Network)),
		fmt.Sprintf("%s 大洲:    %s", labelStyle.Render(IconGlobe), valueStyle.Render(ipInfo.ContinentCode)),
	)
	detailCard := DrawLipglossCard("其他信息", IconInfo, detailInfo, lipgloss.Color("#D787FF"))

	// 连接所有卡片
	result := lipgloss.JoinVertical(
		lipgloss.Left,
		basicCard,
		"",
		geoCard,
		"",
		netCard,
		"",
		typeCard,
		"",
		detailCard,
	)
	
	// 添加完成通知
	noticeStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#87FF87")).
		Padding(0, 1).
		Width(78)
		
	notice := noticeStyle.Render(
		labelStyle.Render(IconCheck) + " IP信息获取完成！使用 `ip nettest` 命令测试网络连通性",
	)
	
	return lipgloss.JoinVertical(lipgloss.Left, result, "", notice)
}

// RenderNetworkTestWithLipgloss 使用 lipgloss 渲染网络测试结果
func RenderNetworkTestWithLipgloss(results []network.SiteTestResult) string {
	// 按名称排序结果
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	// 计算可访问站点数量和平均值
	accessibleCount := 0
	totalSites := len(results)
	var totalRespTime, totalPingTime, totalGenerate204 time.Duration
	respTimeCount, pingTimeCount, generate204Count := 0, 0, 0

	for _, result := range results {
		if result.Accessible {
			accessibleCount++
			totalRespTime += result.ResponseTime
			respTimeCount++
		}

		if result.PingTime > 0 {
			totalPingTime += result.PingTime
			pingTimeCount++
		}

		if result.Generate204 > 0 {
			totalGenerate204 += result.Generate204
			generate204Count++
		}
	}

	// 计算平均值
	var avgRespTime, avgPingTime, avgGenerate204 float64
	if respTimeCount > 0 {
		avgRespTime = totalRespTime.Seconds() / float64(respTimeCount)
	}
	if pingTimeCount > 0 {
		avgPingTime = float64(totalPingTime) / float64(pingTimeCount) / float64(time.Millisecond)
	}
	if generate204Count > 0 {
		avgGenerate204 = float64(totalGenerate204) / float64(generate204Count) / float64(time.Millisecond)
	}

	// 计算可访问率
	accessRate := float64(accessibleCount) / float64(totalSites) * 100
	var accessRateColor, avgRespTimeColor, avgPingTimeColor, avgGenerate204Color lipgloss.Style
	
	if accessRate >= 90 {
		accessRateColor = goodStatusStyle
	} else if accessRate >= 70 {
		accessRateColor = warnStatusStyle
	} else {
		accessRateColor = errorStatusStyle
	}

	if avgRespTime < 0.3 {
		avgRespTimeColor = goodStatusStyle
	} else if avgRespTime < 1.0 {
		avgRespTimeColor = warnStatusStyle
	} else {
		avgRespTimeColor = errorStatusStyle
	}

	if avgPingTime < 50 {
		avgPingTimeColor = goodStatusStyle
	} else if avgPingTime < 150 {
		avgPingTimeColor = warnStatusStyle
	} else {
		avgPingTimeColor = errorStatusStyle
	}

	if avgGenerate204 < 200 {
		avgGenerate204Color = goodStatusStyle
	} else if avgGenerate204 < 500 {
		avgGenerate204Color = warnStatusStyle
	} else {
		avgGenerate204Color = errorStatusStyle
	}

	// 获取网络评级
	networkRating := getNetworkRating(accessRate, avgRespTime, avgPingTime)
	
	// 构建统计信息
	statsInfo := lipgloss.JoinVertical(
		lipgloss.Left,
		fmt.Sprintf("%s站点连通率: %s (%d/%d站点可访问)",
			labelStyle.Render(), 
			accessRateColor.Render(fmt.Sprintf("%.1f%%", accessRate)),
			accessibleCount, totalSites),
		"",
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			fmt.Sprintf("%s平均响应时间: %s   ",
				labelStyle.Render(), 
				avgRespTimeColor.Render(fmt.Sprintf("%.2f秒", avgRespTime))),
			fmt.Sprintf("%sPing平均延迟: %s",
				labelStyle.Render(), 
				avgPingTimeColor.Render(fmt.Sprintf("%.1fms", avgPingTime))),
		),
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			fmt.Sprintf("%s204测试延迟: %s   ",
				labelStyle.Render(), 
				avgGenerate204Color.Render(fmt.Sprintf("%.1fms", avgGenerate204))),
			fmt.Sprintf("%s总体评分: %s",
				labelStyle.Render(), 
				networkRating),
		),
	)
	
	// 创建统计卡片
	statsCard := DrawLipglossCard("网络状况概览", IconSpeed, statsInfo, lipgloss.Color("#5F87FF"))

	// 创建表格样式
	headers := []string{"站点", "状态", "HTTP响应", "Ping延迟", "丢包率", "204延迟"}
	tableStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#5FD7FF")).
		Padding(0, 1)
	
	// 头部样式
	headerStyle := lipgloss.NewStyle().Bold(true)
	
	// 创建表格列宽
	colWidths := []int{12, 15, 15, 15, 10, 15}
	
	// 添加标题
	tableTitle := titleStyle.Render("📊 站点连通性测试详情")
	
	// 创建头部行
	var headerRow strings.Builder
	for i, h := range headers {
		width := colWidths[i]
		paddedHeader := lipgloss.NewStyle().Width(width).Render(headerStyle.Render(h))
		headerRow.WriteString(paddedHeader)
	}
	
	// 初始化表格内容
	tableContent := lipgloss.JoinVertical(
		lipgloss.Left,
		tableTitle,
		headerRow.String(),
	)
	
	// 创建表格行数据
	rows := make([][]string, 0, len(results))
	for _, result := range results {
		// 添加行
		rows = append(rows, []string{
			lipgloss.NewStyle().Bold(true).Render(result.Name),
			getFriendlyStatusTextLipgloss(result.Accessible),
			getFriendlyResponseTimeTextLipgloss(result.ResponseTime, result.Accessible),
			getFriendlyPingTextLipgloss(result.PingTime, result.PingLoss),
			getFriendlyLossRateTextLipgloss(result.PingLoss),
			getFriendlyGenerateTextLipgloss(result.Generate204),
		})
	}

	// 添加数据行
	for _, row := range rows {
		var rowString strings.Builder
		for i, cell := range row {
			if i < len(colWidths) {
				width := colWidths[i]
				paddedCell := lipgloss.NewStyle().Width(width).Render(cell)
				rowString.WriteString(paddedCell)
			}
		}
		tableContent = lipgloss.JoinVertical(
			lipgloss.Left,
			tableContent,
			rowString.String(),
		)
	}
	
	// 注释
	noteText := lipgloss.NewStyle().
		Faint(true).
		Italic(true).
		Render("注: 绿色=良好, 黄色=中等, 红色=较差")
	
	tableContent = lipgloss.JoinVertical(
		lipgloss.Left,
		tableContent,
		noteText,
	)
	
	// 添加表格边框
	table := tableStyle.Render(tableContent)
	
	// 完成通知
	noticeStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#87FF87")).
		Padding(0, 1).
		Width(78)
		
	notice := noticeStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			goodStatusStyle.Render(IconCheck),
			" 测试完成！感谢使用网络测试工具",
		),
	)
	
	// 组合结果
	return lipgloss.JoinVertical(
		lipgloss.Left,
		statsCard,
		"",
		table,
		"",
		notice,
	)
}

// 辅助函数，使用 lipgloss 样式渲染
func getFriendlyStatusTextLipgloss(accessible bool) string {
	if accessible {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			goodStatusStyle.Render(IconCheck),
			" ",
			goodStatusStyle.Render("可访问"),
		)
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		errorStatusStyle.Render(IconCross),
		" ",
		errorStatusStyle.Render("不可访问"),
	)
}

func getFriendlyResponseTimeTextLipgloss(responseTime time.Duration, accessible bool) string {
	if !accessible {
		return errorStatusStyle.Render("超时")
	}
	
	var timeStyle lipgloss.Style
	if responseTime < 300*time.Millisecond {
		timeStyle = goodStatusStyle
	} else if responseTime < 1000*time.Millisecond {
		timeStyle = warnStatusStyle
	} else {
		timeStyle = errorStatusStyle
	}
	return timeStyle.Render(fmt.Sprintf("%.2f秒", responseTime.Seconds()))
}

func getFriendlyPingTextLipgloss(pingTime time.Duration, pingLoss float64) string {
	if pingTime <= 0 || pingLoss >= 1.0 {
		return errorStatusStyle.Render("超时")
	}
	
	var pingStyle lipgloss.Style
	pingMs := float64(pingTime) / float64(time.Millisecond)
	
	if pingMs < 50 {
		pingStyle = goodStatusStyle
	} else if pingMs < 150 {
		pingStyle = warnStatusStyle
	} else {
		pingStyle = errorStatusStyle
	}
	
	return pingStyle.Render(fmt.Sprintf("%.1fms", pingMs))
}

func getFriendlyLossRateTextLipgloss(lossRate float64) string {
	lossPercent := lossRate * 100
	
	if lossRate == 0 {
		return goodStatusStyle.Render("0%")
	} else if lossPercent < 20 {
		return warnStatusStyle.Render(fmt.Sprintf("%.0f%%", lossPercent))
	} else {
		return errorStatusStyle.Render(fmt.Sprintf("%.0f%%", lossPercent))
	}
}

func getFriendlyGenerateTextLipgloss(genTime time.Duration) string {
	if genTime <= 0 {
		return errorStatusStyle.Render("超时")
	}
	
	var genStyle lipgloss.Style
	genMs := float64(genTime) / float64(time.Millisecond)
	
	if genMs < 200 {
		genStyle = goodStatusStyle
	} else if genMs < 500 {
		genStyle = warnStatusStyle
	} else {
		genStyle = errorStatusStyle
	}
	
	return genStyle.Render(fmt.Sprintf("%.1fms", genMs))
}

// 添加兼容函数，供其他代码继续使用
// DrawCard 绘制一个卡片（兼容旧接口，内部使用 DrawLipglossCard）
func DrawCard(title string, icon string, content string, width int, color string) string {
	// 根据旧的颜色字符串映射到 lipgloss 颜色
	var lipglossColor lipgloss.Color
	switch color {
	case BrightBlue:
		lipglossColor = primaryColor
	case BrightCyan:
		lipglossColor = secondaryColor
	case BrightYellow:
		lipglossColor = accentColor
	case BrightGreen:
		lipglossColor = infoColor
	case BrightRed:
		lipglossColor = errorColor
	case BrightMagenta:
		lipglossColor = lipgloss.Color("#D787FF") // 紫色
	default:
		lipglossColor = primaryColor // 默认使用蓝色
	}
	
	return DrawLipglossCard(title, icon, content, lipglossColor)
} 