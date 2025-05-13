package ui

import (
	"strings"
)

// stripANSI 移除ANSI颜色代码，用于计算实际文本长度
func stripANSI(str string) string {
	// 简单匹配ANSI颜色代码，可能不完全精确，但足够满足这个程序的需求
	result := str
	for {
		i := strings.Index(result, "\x1b[")
		if i == -1 {
			break
		}
		
		j := strings.Index(result[i:], "m")
		if j == -1 {
			break
		}
		
		result = result[:i] + result[i+j+1:]
	}
	
	return result
}

// getNetworkRating 根据各项指标评估网络状况
func getNetworkRating(accessible float64, avgRespTime float64, avgPingTime float64) string {
	if accessible >= 90 && avgRespTime < 1.0 && avgPingTime < 100 {
		return "优秀"
	} else if accessible >= 70 && avgRespTime < 2.0 && avgPingTime < 200 {
		return "良好"
	} else if accessible >= 50 && avgRespTime < 3.0 && avgPingTime < 300 {
		return "一般"
	} else {
		return "较差"
	}
}
