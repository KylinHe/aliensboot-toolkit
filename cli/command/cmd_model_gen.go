/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/11/5
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package command

import (
	"fmt"
	"github.com/KylinHe/aliensboot-cli/conf"
	"github.com/KylinHe/aliensboot-cli/template"
	"github.com/KylinHe/aliensboot-cli/util"
	"github.com/go-yaml/yaml"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(modelCmd)
}

var modelCmd = &cobra.Command{
	Use:   "model Ex. model %configPath%",
	Short: "auto generate model code",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		config := readModelConfig(args[0])
		converter := &template.ModelConverter{}
		converter.Convert(config)
	},
}


func readModelConfig(path string) *conf.ModelGenConfig {
	data := util.ReadFile(path)
	if data == nil {
		return nil
	}
	result := &conf.ModelGenConfig{}
	err := yaml.Unmarshal(data, result)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return result
}