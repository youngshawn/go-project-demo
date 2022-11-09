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

		out2, err := json.MarshalIndent(&config.Config, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(out2))

		//os.Exit(2)

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
	remote_enable := viper.GetBool("remoteconfig.enable")
	remote_provider := viper.GetString("remoteconfig.provider")
	remote_endpoint := viper.GetString("remoteconfig.endpoint")
	remote_path := viper.GetString("remoteconfig.path")
	remote_format := viper.GetString("remoteconfig.format")
	if remote_enable {
		viper.AddRemoteProvider(remote_provider, remote_endpoint, remote_path)
		viper.SetConfigType(remote_format)
		if err := viper.ReadRemoteConfig(); err != nil {
			log.Fatal(err)
		}
	}

	// Unmarshal viper map into config
	if err := viper.Unmarshal(&config.Config); err != nil {
		log.Fatal(err)
	}

	// watch local config-file
	viper.OnConfigChange(func(e fsnotify.Event) {
		config.ViperLocker.Lock()
		defer config.ViperLocker.Unlock()

		log.Println("Local config file changed:", e.Name)
		if err := viper.ReadInConfig(); err != nil {
			log.Println("Read local config file failed, error:", err)
			return
		}
		log.Println("Read local config file succeed.")

		config.ConfigLocker.Lock()
		defer config.ConfigLocker.Unlock()
		if err := viper.Unmarshal(&config.Config); err != nil {
			log.Println("Unmarshal config failed, error:", err)
			return
		}
		log.Println("Unmarshal config succeed.")
	})
	viper.WatchConfig()

	// watch remote config-file
	// open a goroutine to watch remote changes forever
	if remote_enable {
		go func() {
			for {
				time.Sleep(time.Second * 10)

				func() {
					config.ViperLocker.Lock()
					defer config.ViperLocker.Unlock()

					if err := viper.WatchRemoteConfig(); err != nil {
						log.Printf("Read remote config file failed, error: %v", err)
						return
					}

					config.ConfigLocker.Lock()
					defer config.ConfigLocker.Unlock()
					if err := viper.Unmarshal(&config.Config); err != nil {
						log.Println("Unmarshal config failed, error:", err)
						return
					}
				}()
			}
		}()
	}
}
