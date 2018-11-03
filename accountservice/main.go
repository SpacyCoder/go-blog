package main

import (
	"flag"
	"fmt"

	"github.com/spacycoder/go-blog/accountservice/config"
	"github.com/spacycoder/go-blog/accountservice/dbclient"
	"github.com/spacycoder/go-blog/accountservice/service"
	"github.com/spf13/viper"
)

var appName = "accountservice"

func init() {
	profile := flag.String("profile", "test", "Environment profile, something similar to spring profiles")
	configServerURL := flag.String("configServerURL", "http://configserver:8888", "Address to config server")
	configBranch := flag.String("configBranch", "master", "git branch to fetch configuration from")
	flag.Parse()

	viper.Set("profile", *profile)
	viper.Set("configServerURL", *configServerURL)
	viper.Set("configBranch", *configBranch)
}

func main() {

	fmt.Printf("Starting %v\n", appName)

	config.LoadConfigurationFromBranch(viper.GetString("configServerURL"), appName, viper.GetString("profile"), viper.GetString("configBranch"))
	initializeBoltClient()
	service.StartWebServer(viper.GetString("server_port"))

}

func initializeBoltClient() {
	service.DBClient = &dbclient.BoltClient{}
	service.DBClient.OpenBoltDb()
	service.DBClient.Seed()
}
