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

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "user need to regist",
	Run: func(cmd *cobra.Command, args []string) {
		tmpu, _ := cmd.Flags().GetString("username")
		tmppw, _ := cmd.Flags().GetString("password")
		tmpe, _ := cmd.Flags().GetString("email")
		tmpc, _ := cmd.Flags().GetString("cellphone")

		//judge the validity
		if tmpu == "" || tmppw == "" || tmpe == "" || tmpc == "" {
			fmt.Println("Input the username[-u],password[-p],email[-e],cellphone[-c]")
			return;
		}
		notexist, err := service.UserRegister(tmpu,tmppw,tmpe,tmpc)
		if notexist == false {
			fmt.Println("Username has existed, please try again")
			return
		} else {
			if err != nil {
				fmt.Println("Error happened, please check the error log")
				return
			} else {
				fmt.Println("Regist Successfully!")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
	//define personal information
	registerCmd.Flags().StringP("username","u","","username haven't been rigisted")
	registerCmd.Flags().StringP("password","p","","password no shorter than 8 characters")
	registerCmd.Flags().StringP("email","e","","personal email address")
	registerCmd.Flags().StringP("cellphone","c","","personal phone number")
}