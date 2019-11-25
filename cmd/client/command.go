package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
)

type ClientParam struct {
	medallion  *string
	pickupdate *string
	useCache   *string
}

func main() {
	var host, port, selectedCommand string
	cmdClient := kingpin.New("client", "Cab data service command line client")
	cmdParams := configureCmdline(cmdClient)
	viper.SetConfigName("command")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Warnf("Config file not found...")
		log.Warnf("Using Defaults")
		host = "localhost"
		port = "8080"
	} else {
		host = viper.GetString("server.host")
		port = viper.GetString("server.port")
	}
	if selectedCommand, err = cmdClient.Parse(os.Args[1:]); err != nil {
		log.Error(err)
		return
	}
	switch selectedCommand {
	case "get":
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		url := fmt.Sprintf("http://%s:%s/cabtrip/%s", host, port, *cmdParams.medallion)
		if *cmdParams.pickupdate != "" {
			url += "/date/" + *cmdParams.pickupdate
		}
		if *cmdParams.useCache == "true" {
			url += "?cache=true"
		}
		log.Info(url)
		response, err := client.Get(url)
		if err != nil {
			log.Error("error while accessing the service")
			return
		}
		err = response.Write(log.StandardLogger().Writer())
		if err != nil {
			log.Error("error while accessing the service")
			return
		}
	case "refresh":
		client := &http.Client{
			Timeout: time.Second * 200,
		}
		url := fmt.Sprintf("http://%s:%s/cache/refresh_cache", host, port)
		log.Info(url)
		response, err := client.Get(url)
		if err != nil {
			log.Error("error while accessing the service")
			return
		}
		err = response.Write(log.StandardLogger().Writer())
		if err != nil {
			log.Error("error while accessing the service")
			return
		}
	default:
		log.Error("unexpected command")
		return
	}
}

func configureCmdline(cmdClient *kingpin.Application) *ClientParam {
	cabCmd := cmdClient.Command("get", "Get cab data")
	returnValues := &ClientParam{}
	returnValues.useCache = cabCmd.Flag("cache", "use cache (default: false)").Default("false").String()
	returnValues.medallion = cabCmd.Flag("medallion", "enter the cab medallion").Required().String()
	returnValues.pickupdate = cabCmd.Flag("date", "enter the cab pickupdate").String()
	cmdClient.Command("refresh", "Refresh cache")
	return returnValues
}
