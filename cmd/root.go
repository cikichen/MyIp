/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ip",
	Short: "show your ip address",
	Long:  `Show all information about your current ip.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
