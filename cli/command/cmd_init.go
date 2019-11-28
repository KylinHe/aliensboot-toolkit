/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/10/29
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>k
 *******************************************************************************/
package command

import (
	"fmt"
	"github.com/KylinHe/aliensboot-cli/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	initCmd.Flags().StringVarP(&defaultModuleName, "source", "s", "defaultmodule", "module source")
	initCmd.Flags().StringVarP(&templateName, "template", "t", "aliensboot-server", "init template")
	initCmd.Flags().StringSliceVarP(&moduleNames, "modules", "m", []string{}, "add modules")
	RootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init [package path], Ex. aliensboot init github.com/KylinHe/aliensboot-server",
	Short: "initial aliensboot project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.CommandPath())
		//if len(args)
		if len(args) == 0 {
			cmd.Help()
			return
		}
		initProject("", args[0])
		//fmt.Println(moduleNames)
	},
}

func getFilterModules(projectPath string, moduleNames []string) []string{
	var filterModules []string
	srcConfigPath := getPath(projectPath, "config", "modules")
	dir, err := ioutil.ReadDir(srcConfigPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		srcModuleName := strings.Split(fi.Name(),".")[0]
		isFilter := true
		for _, moduleName := range moduleNames {
			if srcModuleName == moduleName {
				isFilter = false
			}
		}
		if isFilter {
			filterModules = append(filterModules, srcModuleName)
		}
	}
	return filterModules
}

func initProject(targetHomePath string, packagePath string) {
	//util.AddFileFilter()

	//projectConfig := &model.ProjectConfig{
	//	Name: packagePath,
	//}

	projectPath := getPath(AliensBootHome, templateName)
	sourceProjectConfig := EnsureSourceProjectConfig()

	//writeProjectConfig(targetHomePath, projectConfig)

	srcSrcPath := getPath(projectPath, "src")

	targetSrcPath := getPath(targetHomePath, "src")

	srcCopyPath := getPath(projectPath)

	targetCopyPath := getPath(targetHomePath, getCurrentPath())

	srcProtocolPath := getPath(projectPath, "src", "protocol")

	targetProtocolPath := getPath(targetHomePath, "src", "protocol")

	replaceContent := make(map[string]string)

	replaceContent[sourceProjectConfig.Name] = packagePath

	if moduleNames == nil || len(moduleNames) == 0 {
		util.CopyDir(srcSrcPath, targetSrcPath, replaceContent, defaultModuleName)
		util.CopyDir(srcCopyPath, targetCopyPath, replaceContent, append(sourceProjectConfig.Template.Exclude, defaultModuleName, "src")...)
	} else {
		filterModules := getFilterModules(srcCopyPath, moduleNames)
		// copy src dir to target (filter modules and protocol)
		filterModules1 := append(filterModules, "protocol")
		util.CopyDir(srcProtocolPath, targetProtocolPath, replaceContent)
		util.CopyDir(srcSrcPath, targetSrcPath, replaceContent, filterModules1...)

		// copy copy dir to target (filter modules and src)
		filterModules2 := append(filterModules, "src")
		util.CopyDir(srcCopyPath, targetCopyPath, replaceContent, append(sourceProjectConfig.Template.Exclude,filterModules2...)...)
	}

}

func getCurrentPath() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err.Error())
	}
	return path
}
