package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"github.com/5dao/gdav/server"
)

func main() {

	log.Println("ok")
	rootCmd := bigFlags()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func bigFlags() *cobra.Command {
	var cfgFile string
	var salt, uid, pwd string

	var rootCmd = &cobra.Command{
		Use:   "gdav",
		Short: "A webdav server",
	}

	var svrCmd = &cobra.Command{
		Use:   "server",
		Short: "run webdav server",
		Run: func(cmd *cobra.Command, args []string) {
			serverMain(cfgFile)
		},
	}
	svrCmd.PersistentFlags().StringVar(&cfgFile, "c", "config.toml", "config file (default is ./config.toml)")
	rootCmd.AddCommand(svrCmd)

	var pwdCmd = &cobra.Command{
		Use:   "pwd",
		Short: "make encrypted password by uid、pwd and salt",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(server.UserPwd(salt, uid, pwd))
		},
	}
	pwdCmd.Flags().StringVarP(&salt, "salt", "s", "", "password salt,please same with config.toml")
	pwdCmd.Flags().StringVarP(&uid, "uid", "u", "", "user name")
	pwdCmd.Flags().StringVarP(&pwd, "pwd", "p", "", "password")

	rootCmd.AddCommand(pwdCmd)

	return rootCmd
}

func serverMain(configPath string) {
	cfg := &server.Config{}
	_, err := toml.DecodeFile(configPath, cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	svr, err := server.NewServer(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//log.Println(svr)
	go svr.Start()

	ch := make(chan int)
	<-ch
}
