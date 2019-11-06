/*
Copyright ? 2019 NAME HERE <EMAIL ADDRESS>
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"demo/service"
	"github.com/spf13/cobra"
)

var (
	username *string
	password *string
)
// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "For User To Login",
	Run: func(cmd *cobra.Command, args []string) {
		tmpu, _ := cmd.Flags().GetString("username")
		tmppw, _ := cmd.Flags().GetString("password")

		if tmpu == "" || tmppw == "" {
			fmt.Println("Input the username[-u],password[-p]")
			return;
		}
		if _, flag := service.GetCurUser(); flag == true {
			fmt.Println("Please logout firstly")
			return
		}
		if tf := service.UserLogin(tmpu,tmppw); tf == true {
			fmt.Println("Login Successfully!")
		} else {
			fmt.Println("Login Fail: wront username or password")
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().StringP("username","u","","agenda account username")
	loginCmd.Flags().StringP("password","p","","agenda account password")
}