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

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "For User To Logout",
	Run: func(cmd *cobra.Command, args []string) {
		if err := service.UserLogout(); err != true {
			fmt.Println("Error happened, please check the error log")
		} else {
			fmt.Println("Logout Successfully!")
		}
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}