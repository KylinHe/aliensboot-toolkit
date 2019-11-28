/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/10/29
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package command

import (
	"github.com/spf13/cobra"
)


//var path string

func init() {
	//dataCmd.Flags().StringVarP(&path, "path", "p", "/root/cluster/config/default", "-path")
	RootCmd.AddCommand(dataCmd)
}

var dataCmd = &cobra.Command{
	Use:   "data [upload download]",
	Short: "operation for data",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
}
