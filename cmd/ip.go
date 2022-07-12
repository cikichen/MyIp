/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)

type IPInfo struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	AS          string  `json:"as"`
	QueryIp     string  `json:"query"`
}

// ipCmd represents the ip command
var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "show your ip address",
	Long:  `Show all information about your current ip.`,
	Run: func(cmd *cobra.Command, args []string) {

		localIp, err := externalIP()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("IP:         ", localIp.String())

		//interfaces, err := net.Interfaces()
		//if err != nil {
		//	panic("Error:" + err.Error())
		//}
		//for _, inter := range interfaces {
		//	fmt.Println("%1s", inter.Name)
		//	fmt.Println(inter.Index)
		//	fmt.Println(inter.HardwareAddr)
		//}

		myIP := GetMyPublicIP()
		fmt.Println("PublicIP:   ", myIP)
		result := OnlineIpInfo(myIP)

		if result != nil {
			fmt.Println("Country:    ", result.Country)
			fmt.Println("CountryCode:", result.CountryCode)
			fmt.Println("Region:     ", result.RegionName)
			fmt.Println("City:       ", result.City)
			fmt.Println("lat:        ", result.Lat)
			fmt.Println("lon:        ", result.Lon)
			fmt.Println("TZ:         ", result.Timezone)
			fmt.Println("ISP:        ", result.ISP)
			fmt.Println("AS:         ", result.AS)
		}
	},
}

func GetMyPublicIP() string {
	resp, err := http.Get("https://myexternalip.com/raw")
	if err != nil {
		return ""
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}

func OnlineIpInfo(ip string) *IPInfo {
	url := "http://ip-api.com/json/" + ip + "?lang=zh-CN"
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	var result IPInfo
	if err := json.Unmarshal(out, &result); err != nil {
		return nil
	}
	return &result
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

func init() {
	rootCmd.AddCommand(ipCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ipCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
