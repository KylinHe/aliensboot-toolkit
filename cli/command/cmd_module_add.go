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
	"github.com/KylinHe/aliensboot-cli/util"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	moduleCmd.AddCommand(moduleAddCmd)
}

var moduleAddCmd = &cobra.Command{
	Use:   "add Ex. add %module%",
	Short: "add initial module code in current path",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		config := EnsureTargetProjectConfig()
		addModule("", config.Name, args[0])
	},
}

func addModule(targetHomePath string, packagePath string, moduleName string) {
	projectPath := getPath(AliensBootHome, templateName)

	sourceProjectConfig := EnsureSourceProjectConfig()

	sourceSrcPath := getPath(projectPath, "src")
	targetSrcPath := getPath(targetHomePath, "src")


	srcModulePath := getPath(sourceSrcPath, "module", defaultModuleName)

	targetModulePath := getPath(targetSrcPath, "module", moduleName)

	srcInfo, err := os.Stat(targetModulePath)
	if err == nil && srcInfo.IsDir() {
		fmt.Println(fmt.Errorf("module path already exists : %v", targetModulePath))
		return
	}

	srcConfigPath := getPath(projectPath, "config", "modules", defaultModuleName+".yml.bak")

	targetConfigPath := getPath(targetHomePath, "config", "modules", moduleName+".yml")

	srcPublicPath := getPath(sourceSrcPath, "public", defaultModuleName+".go")

	targetPublicPath := getPath(targetSrcPath, "public", moduleName+".go")

	replaceContent := make(map[string]string)
	replaceContent[defaultModuleName] = moduleName
	replaceContent[sourceProjectConfig.Name] = packagePath

	util.CopyDir(srcModulePath, targetModulePath, replaceContent)
	util.CopyFile(srcConfigPath, targetConfigPath, replaceContent)
	util.CopyFile(srcPublicPath, targetPublicPath, replaceContent)

}
