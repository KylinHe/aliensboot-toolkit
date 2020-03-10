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
	"github.com/KylinHe/aliensboot-cli/conf"
	"github.com/KylinHe/aliensboot-cli/template"
	"github.com/KylinHe/aliensboot-cli/util"
	"github.com/spf13/cobra"
	"strings"
)

func init() {
	moduleCmd.AddCommand(codeCmd)
}

var codeCmd = &cobra.Command{
	Use:   "gen Ex. gen %module%",
	Short: "auto generate module code",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}

		config := EnsureTargetProjectConfig()
		GenCode(args[0], config.Name, config.CodeTemplatePath, "")
	},
}

func GenCode(module string, packageName string, templatePath string, rootPath string) {
	protocolPath := getPath(rootPath, "src", "protocol", "protocol.proto")
	templatePath = getPath(templatePath, "templates", "protocol")

	config := &conf.CodeGenConfig{
		Package:      packageName,
		ProtoPath:    protocolPath,
		TemplatePath: templatePath,
		Modules:      []*conf.ModuleConfig{},
		//TemplatePath:templatePath,
	}

	moduleName := strings.ToLower(module)

	buildModuleConfig(rootPath, config, moduleName)

	converter := &template.ServiceConverter{}
	converter.Convert(config)
	//fmt.Sprintf("config data %+v", config)
	//template.Convert(config)
}

func buildModuleConfig(rootPath string, config *conf.CodeGenConfig, moduleName string) {
	moduleConfig := &conf.ModuleConfig{
		Name:    moduleName,
		Outputs: []*conf.Output{},
	}

	moduleConfig.Outputs = append(moduleConfig.Outputs, &conf.Output{
		Template:  getModuleTemplatePath(config.TemplatePath, moduleName, "service.template"),
		Output:    getPath(rootPath, "src", "module", moduleName, "service", "service.go"),
		Overwrite: true,
	})

	moduleConfig.Outputs = append(moduleConfig.Outputs, &conf.Output{
		Template:  getModuleTemplatePath(config.TemplatePath, moduleName, "handle.template"),
		Output:    getPath(rootPath, "src", "module", moduleName, "service"),
		Prefix:    "handle_${}.go",
		Overwrite: false,
	})

	moduleConfig.Outputs = append(moduleConfig.Outputs, &conf.Output{
		Template:  getModuleTemplatePath(config.TemplatePath, moduleName, "rpc.template"),
		Output:    getPath(rootPath, "src", "dispatch", "rpc", moduleName+".go"),
		Overwrite: true,
	})

	config.Modules = append(config.Modules, moduleConfig)
}

func getModuleTemplatePath(templateRoot string, module string, name string) string {
	moduleTemplatePath := getPath(templateRoot, module, name)
	exist, _ := util.PathExists(moduleTemplatePath)
	if exist {
		return moduleTemplatePath
	} else {
		return getPath(templateRoot, name)
	}
}
