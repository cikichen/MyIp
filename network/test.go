package network

import (
	"context"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 定义全局HTTP超时时间
var globalTimeout = 10 * time.Second

// SetGlobalTimeout 设置全局HTTP请求超时时间
func SetGlobalTimeout(timeout time.Duration) {
	globalTimeout = timeout
}

// SiteTestResult 存储站点测试结果
type SiteTestResult struct {
	Name         string        // 站点名称
	URL          string        // 站点URL
	Accessible   bool          // 是否可访问
	ResponseTime time.Duration // 响应时间
	StatusCode   int           // HTTP状态码
	Error        string        // 错误信息
	DNSTime      time.Duration // DNS解析时间
	ConnectTime  time.Duration // 连接建立时间
	PingTime     time.Duration // Ping延迟时间
	PingLoss     float64       // Ping丢包率 (0-1)
	Generate204  time.Duration // Generate_204延迟测试
}

// CommonSites 定义要测试的常用站点列表
var CommonSites = []struct {
	Name string
	URL  string
}{
	{"Google", "https://www.google.com"},
	{"GitHub", "https://github.com"},
	{"YouTube", "https://www.youtube.com"},
	{"Twitter", "https://twitter.com"},
	{"Facebook", "https://www.facebook.com"},
	{"Baidu", "https://www.baidu.com"},
	{"Alibaba", "https://www.alibaba.com"},
	{"Tencent", "https://www.qq.com"},
	{"Microsoft", "https://www.microsoft.com"},
	{"Apple", "https://www.apple.com"},
}

// TestGenerate204 测试generate_204延迟
func TestGenerate204(host string) (time.Duration, error) {
	// 从URL中提取主机名
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")
	host = strings.Split(host, "/")[0]

	// 构建generate_204 URL
	var generate204URL string
	if strings.Contains(host, "google") {
		generate204URL = "https://www.google.com/generate_204"
	} else if strings.Contains(host, "cloudflare") {
		generate204URL = "https://www.cloudflare.com/cdn-cgi/trace"
	} else {
		// 对于其他站点，使用通用的方法
		if strings.HasSuffix(host, ".com") || strings.HasSuffix(host, ".org") || strings.HasSuffix(host, ".net") {
			generate204URL = "https://" + host + "/generate_204"
		} else {
			// 如果不是常见域名后缀，尝试添加www前缀
			generate204URL = "https://www." + host + "/generate_204"
		}
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: globalTimeout,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	// 发送请求并测量响应时间
	startTime := time.Now()
	resp, err := client.Get(generate204URL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// 计算延迟时间
	latency := time.Since(startTime)
	return latency, nil
}

// PingHost 测试主机的延迟和丢包率
func PingHost(host string) (time.Duration, float64, error) {
	// 从URL中提取主机名
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")
	host = strings.Split(host, "/")[0]

	// 执行ping命令
	cmd := exec.Command("ping", "-c", "5", "-i", "0.2", host)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, 1.0, err // 返回最大丢包率1.0表示100%丢包
	}

	// 解析ping输出
	outputStr := string(output)

	// 提取平均延迟
	avgTimeRegex := regexp.MustCompile(`min/avg/max/(?:stddev|mdev) = [\d.]+/([\d.]+)/[\d.]+/[\d.]+`)
	avgTimeMatch := avgTimeRegex.FindStringSubmatch(outputStr)

	var avgTime float64
	if len(avgTimeMatch) >= 2 {
		avgTime, _ = strconv.ParseFloat(avgTimeMatch[1], 64)
	}

	// 提取丢包率
	packetLossRegex := regexp.MustCompile(`(\d+)% packet loss`)
	packetLossMatch := packetLossRegex.FindStringSubmatch(outputStr)

	var packetLoss float64
	if len(packetLossMatch) >= 2 {
		lossPercent, _ := strconv.ParseFloat(packetLossMatch[1], 64)
		packetLoss = lossPercent / 100.0
	}

	return time.Duration(avgTime * float64(time.Millisecond)), packetLoss, nil
}

// TestSite 测试单个站点的可访问性、响应时间和延迟
func TestSite(site struct {
	Name string
	URL  string
}) SiteTestResult {
	result := SiteTestResult{
		Name: site.Name,
		URL:  site.URL,
	}

	// 创建自定义的Transport，以便测量DNS和连接时间
	var dnsStart, dnsEnd, connectStart, connectEnd time.Time

	dialer := &net.Dialer{
		Timeout:   globalTimeout,
		KeepAlive: globalTimeout / 2,
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dnsStart = time.Now()
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}

			// 解析IP地址
			ips, err := net.LookupIP(host)
			if err != nil {
				return nil, err
			}
			dnsEnd = time.Now()

			// 连接到服务器
			connectStart = time.Now()
			conn, err := dialer.DialContext(ctx, network, net.JoinHostPort(ips[0].String(), port))
			connectEnd = time.Now()
			return conn, err
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   globalTimeout / 2,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   globalTimeout,
	}

	// 发送请求并测量响应时间
	startTime := time.Now()
	resp, err := client.Get(site.URL)

	// 记录DNS和连接时间
	result.DNSTime = dnsEnd.Sub(dnsStart)
	result.ConnectTime = connectEnd.Sub(connectStart)

	if err != nil {
		result.Accessible = false
		result.Error = err.Error()
	} else {
		defer resp.Body.Close()
		// 计算总响应时间
		result.ResponseTime = time.Since(startTime)
		result.StatusCode = resp.StatusCode
		result.Accessible = resp.StatusCode >= 200 && resp.StatusCode < 400
	}

	// 执行Ping测试
	pingTime, pingLoss, _ := PingHost(site.URL)
	result.PingTime = pingTime
	result.PingLoss = pingLoss

	// 执行Generate_204测试
	generate204Time, _ := TestGenerate204(site.URL)
	result.Generate204 = generate204Time

	return result
}

// TestCommonSites 测试所有常用站点
func TestCommonSites() []SiteTestResult {
	results := make([]SiteTestResult, 0, len(CommonSites))

	// 创建通道用于并发测试
	resultChan := make(chan SiteTestResult, len(CommonSites))

	// 并发测试所有站点
	for _, site := range CommonSites {
		go func(s struct {
			Name string
			URL  string
		}) {
			resultChan <- TestSite(s)
		}(site)
	}

	// 收集结果
	for i := 0; i < len(CommonSites); i++ {
		results = append(results, <-resultChan)
	}

	return results
}
