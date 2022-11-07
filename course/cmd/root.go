package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/youngshawn/go-project-demo/course/config"
	"github.com/youngshawn/go-project-demo/course/models"
	"github.com/youngshawn/go-project-demo/course/routes"

	_ "github.com/spf13/viper/remote"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "course",
	Short: "A Course Management system",
	Long:  `A Course Management system`,
	Run: func(cmd *cobra.Command, args []string) {

		// viper
		/*out1, err := json.MarshalIndent(viper.AllKeys(), "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(out1))*/

		for {
			out2, err := json.MarshalIndent(&config.Config, "", "    ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out2))
			time.Sleep(5 * time.Second)
		}

		os.Exit(2)

		// get configurations
		address := config.Config.Listen

		// init database and cache
		config.DatabaseConnectAndSetup()
		config.CacheConnectAndSetup()
		models.ModelInit()

		// setup gin server
		router := gin.Default()
		router.SetTrustedProxies([]string{"127.0.0.1"})
		routes.InstallRoutes(router)

		// start gin server
		log.Fatal(router.Run(address))
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.course.yaml)")

	// set pflags from config.config struct
	config.ExposeConfigAsPFlags(rootCmd)

	// bind pflags to viper
	viper.BindPFlags(rootCmd.PersistentFlags())

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".course.yaml".
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".course")
	}

	// load env into viper
	viper.SetEnvPrefix("COURSE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	viper.AutomaticEnv() // read in environment variables that match

	// load local config-file into viper
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	log.Println("Using config file:", viper.ConfigFileUsed())

	// load remote config-file into viper
	viper.AddRemoteProvider("etcd3", "http://127.0.0.1:2379", "/config/prod/cloud/region/github.com/youngshawn/go-proect-demo/course/course.yaml")
	viper.SetConfigType("yaml")
	if err := viper.ReadRemoteConfig(); err != nil {
		log.Fatal(err)
	}

	// Unmarshal
	if err := viper.Unmarshal(&config.Config); err != nil {
		log.Fatal(err)
	}

	// watch local config-file
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Local config file changed:", e.Name)
		if err := viper.ReadInConfig(); err != nil {
			log.Println("Read local config file failed, error:", err)
			return
		}
		log.Println("Read local config file succeed.")

		if err := viper.Unmarshal(&config.Config); err != nil {
			log.Println("Unmarshal config failed, error:", err)
			return
		}
		log.Println("Unmarshal config succeed.")

		if err := config.PublishDynamicConfigs(); err != nil {
			log.Println("Publish config failed, error:", err)
			return
		}
		log.Println("Publish config succeed.")
	})
	viper.WatchConfig()

	// watch remote config-file
	// open a goroutine to watch remote changes forever
	go func() {
		for {
			time.Sleep(time.Second * 10) // delay after each request

			// currently, only tested with etcd support
			if err := viper.WatchRemoteConfig(); err != nil {
				log.Printf("Unable to read remote config: %v", err)
				continue
			}

			// unmarshal new config into our runtime config struct. you can also use channel
			// to implement a signal to notify the system of the changes
			if err := viper.Unmarshal(&config.Config); err != nil {
				log.Println("Unmarshal config failed, error:", err)
				continue
			}
		}
	}()
}
