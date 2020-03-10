/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved. 
 * Date:
 *     2019/5/31
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package conf


type ProjectConfig struct {
	Name         string
	CodeTemplatePath  string `yaml:"code.template.path"`
	Template 	 TemplateConfig
}

type TemplateConfig struct {
	Exclude []string
}

type CodeGenConfig struct {
	Package      string `yaml:"package"`       //proto包名
	ProtoPath    string `yaml:"path.proto"`    //proto文件路径
	TemplatePath string `yaml:"path.template"` //模板文件路径
	Modules      []*ModuleConfig
}

type ModelGenConfig struct {
	ProtoPath    string `yaml:"path.proto"`    //proto文件路径
	TemplatePath string `yaml:"path.template"` //模板文件路径
	OutputPath   string `yaml:"path.output"`   //脚本输出路径
}

type ModuleConfig struct {
	Name    string
	Outputs []*Output
}

type Output struct {
	Template  string
	Output    string
	Prefix    string
	Overwrite bool
}
