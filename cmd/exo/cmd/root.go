package cmd

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var gConfigFolder string
var gConfigFilePath string

//current Account informations
var gAccountName string
var gCurrentAccount *account

var gAllAccount *config

//egoscale client
var cs *egoscale.Client

//Aliases
var gListAlias = []string{"ls"}
var gRemoveAlias = []string{"rm"}
var gDeleteAlias = []string{"del"}
var gShowAlias = []string{"get"}
var gCreateAlias = []string{"add"}
var gUploadAlias = []string{"up"}

type account struct {
	Name        string
	Account     string
	Endpoint    string
	Key         string
	Secret      string
	DefaultZone string
}

type config struct {
	DefaultAccount string
	Accounts       []account
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "exo",
	Short: "A simple CLI to use CloudStack using egoscale lib",
	//Long:  `A simple CLI to use CloudStack using egoscale lib`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&gConfigFilePath, "config", "", "Specify an alternate config file [env EXOSCALE_CONFIG]")
	RootCmd.PersistentFlags().StringVarP(&gAccountName, "account", "a", "", "Account to use in config file [env EXOSCALE_ACCOUNT]")

	cobra.OnInitialize(initConfig, buildClient)

}

var ignoreClientBuild = false

func buildClient() {
	if ignoreClientBuild {
		return
	}

	if cs != nil {
		return
	}

	cs = egoscale.NewClient(gCurrentAccount.Endpoint, gCurrentAccount.Key, gCurrentAccount.Secret)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	envs := map[string]string{
		"EXOSCALE_CONFIG":  "config",
		"EXOSCALE_ACCOUNT": "account",
	}

	for env, flag := range envs {
		flag := RootCmd.Flags().Lookup(flag)
		if value := os.Getenv(env); value != "" {
			if err := flag.Value.Set(value); err != nil {
				log.Fatal(err)
			}
		}
	}

	envEndpoint := os.Getenv("EXOSCALE_ENDPOINT")
	envKey := os.Getenv("EXOSCALE_KEY")
	envSecret := os.Getenv("EXOSCALE_SECRET")

	if envEndpoint != "" && envKey != "" && envSecret != "" {
		cs = egoscale.NewClient(envEndpoint, envKey, envSecret)
		return
	}

	config := &config{}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	gConfigFolder = path.Join(usr.HomeDir, ".exoscale")

	if gConfigFilePath != "" {
		// Use config file from the flag.
		viper.SetConfigFile(gConfigFilePath)
	} else {
		// Search config in home directory with name ".cobra_test" (without extension).
		viper.SetConfigName("exoscale")
		viper.AddConfigPath(path.Join(usr.HomeDir, ".exoscale"))
		viper.AddConfigPath(usr.HomeDir)
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil && getCmdPosition("config") == 1 {
		ignoreClientBuild = true
		return
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := viper.Unmarshal(config); err != nil {
		log.Fatal(fmt.Errorf("couldn't read config: %s", err))
	}

	if config.DefaultAccount == "" && gAccountName == "" {
		log.Fatalf("default account not defined")
	}

	if gAccountName == "" {
		gAccountName = config.DefaultAccount
	}

	gAllAccount = config
	gAllAccount.DefaultAccount = gAccountName

	for i, acc := range config.Accounts {
		if acc.Name == gAccountName {
			gCurrentAccount = &config.Accounts[i]
			return
		}
	}
	log.Fatalf("Could't find any account with name: %q", gAccountName)
}

// return a command position by fetching os.args and ignoring flags
//
//example: "$ exo -r preprod vm create" vm position is 1 and create is 2
//
func getCmdPosition(cmd string) int {

	count := 1

	isFlagParam := false

	for _, arg := range os.Args[1:] {

		if strings.HasPrefix(arg, "-") {

			flag := RootCmd.Flags().Lookup(strings.Trim(arg, "-"))
			if flag == nil {
				flag = RootCmd.Flags().ShorthandLookup(strings.Trim(arg, "-"))
			}

			if flag != nil && (flag.Value.Type() != "bool") {
				isFlagParam = true
			}
			continue
		}

		if isFlagParam {
			isFlagParam = false
			continue
		}

		if arg == cmd {
			break
		}
		count++
	}

	return count
}
