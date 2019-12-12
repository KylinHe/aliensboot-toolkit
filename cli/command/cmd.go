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
	"fmt"
	"github.com/KylinHe/aliensboot-cli/conf"
	"github.com/KylinHe/aliensboot-cli/util"
	"github.com/go-yaml/yaml"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var AliensBootHome = ""

var templateName string
var defaultModuleName string
var moduleNames []string

var _sourceProjectConfig *conf.ProjectConfig
var _targetProjectConfig *conf.ProjectConfig

var RootCmd = &cobra.Command{
	Use: "aliensboot",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
}

func Execute() {
	//util.AddFileFilter(".git")
	AliensBootHome = os.Getenv("ALIENSBOOT_HOME")
	if AliensBootHome == "" {
		fmt.Println("can not found env ALIENSBOOT_HOME")
		os.Exit(1)
	}
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		os.Exit(1)
	}
}

func EnsureSourceProjectConfig() *conf.ProjectConfig {
	_sourceProjectConfig = readProjectConfig(getPath(AliensBootHome, templateName))
	if _sourceProjectConfig == nil {
		fmt.Println("invalid source project description file 'project.yml'")
		os.Exit(1)
	}
	fmt.Printf("project config %+v \n", _sourceProjectConfig)
	return _sourceProjectConfig
}

func EnsureTargetProjectConfig() *conf.ProjectConfig {
	_targetProjectConfig = readProjectConfig("")
	if _targetProjectConfig == nil {
		fmt.Println("invalid project description file 'project.yml'")
		os.Exit(1)
	}
	fmt.Printf("project config %+v \n", _targetProjectConfig)
	return _targetProjectConfig
}

func getProjectFilePath(projectPath string) string {
	return getPath(projectPath, ".aliensboot", "project.yml")
}

func writeProjectConfig(projectPath string, projectConfig *conf.ProjectConfig) {
	projectFilePath := getProjectFilePath(projectPath)
	data, _ := yaml.Marshal(projectConfig)
	util.WriteFile(projectFilePath, data)
}

func readProjectConfig(projectPath string) *conf.ProjectConfig {
	projectFilePath := getProjectFilePath(projectPath)
	data := util.ReadFile(projectFilePath)
	if data == nil {
		return nil
	}
	result := &conf.ProjectConfig{}
	err := yaml.Unmarshal(data, result)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return result
}

func getPath(basePath string, packages ...string) string {
	result := basePath
	for _, name := range packages {
		if result == "" {
			result = name
		} else {
			result = result + string(filepath.Separator) + name
		}
	}
	return result
}
