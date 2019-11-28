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

func init() {
	RootCmd.AddCommand(moduleCmd)
}

var moduleCmd = &cobra.Command{
	Use:   "module [add gen]",
	Short: "operation for module",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
}
