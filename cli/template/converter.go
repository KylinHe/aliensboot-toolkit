/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/3/30
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package template

import (
	"fmt"
	"github.com/KylinHe/aliensboot-cli/conf"
	"io/ioutil"
	"os"
	"strings"
)

const (
	MsgSplitStr     = "<message>"

	RequestSplitStr = "<request>" //单请求

)



type ServiceConverter struct {
	data *ServiceData
}

func (s *ServiceConverter) Convert(config *conf.CodeGenConfig) {
	parser := &ProtoParser{}
	parser.Parse(config.ProtoPath)
	s.data = parser.serviceData
	//fmt.Printf("proto data %v", message.modules["passport"].Handlers[6])
	for _, moduleConfig := range config.Modules {
		module := s.data.modules[moduleConfig.Name]
		if module == nil {
			fmt.Printf("module %v is not found in proto file %v \n", moduleConfig.Name, config.ProtoPath)
			continue
		}
		s.convertModule(moduleConfig, module)
	}
}

func (s *ServiceConverter) ConvertModel(config *conf.CodeGenConfig) {
	parser := &ProtoParser{}
	parser.Parse(config.ProtoPath)
	s.data = parser.serviceData
	//fmt.Printf("proto data %v", message.modules["passport"].Handlers[6])
	for _, moduleConfig := range config.Modules {
		module := s.data.modules[moduleConfig.Name]
		if module == nil {
			fmt.Printf("module %v is not found in proto file %v \n", moduleConfig.Name, config.ProtoPath)
			continue
		}
		s.convertModule(moduleConfig, module)
	}
}

func (s *ServiceConverter) convertModule(moduleConfig *conf.ModuleConfig, module *Module) {
	for _, outputConfig := range moduleConfig.Outputs {
		templatePath := outputConfig.Template
		////配置模板根目录需要加上根目录
		//if templateRoot != "" {
		//	templatePath = templateRoot + "/" + templatePath
		//}

		b, err := ioutil.ReadFile(templatePath)
		if err != nil {
			fmt.Printf(err.Error())
			return
		}

		content := string(b)

		if outputConfig.Prefix == "" {
			//写一个文件
			content = s.convertService(content, module, MsgSplitStr)
			content = s.convertService(content, module, RequestSplitStr)

			writeFile(outputConfig.Output, content, outputConfig.Overwrite)
		} else {

			s.convertHandle(MsgSplitStr, RequestSplitStr, content, module, outputConfig)
			s.convertHandle(RequestSplitStr, MsgSplitStr, content, module, outputConfig)
		}

	}

}

func (s *ServiceConverter) convertService(templateContent string, module *Module, split string) string {
	results := strings.Split(templateContent, split)
	header := ""
	content := ""
	tailf := ""

	if len(results) == 3 {
		header = s.replaceMessage(results[0], module)
		tailf = s.replaceMessage(results[2], module)
		module.Foreach(func(handler *ProtoHandler) bool {
			handleStr := s.replaceMessage(results[1], module)
			if split == MsgSplitStr && !handler.IsSession() {
				return true
			}
			if split == RequestSplitStr && !handler.IsRequest() {
				return true
			}
			handleStr = s.replaceHandle(handleStr, handler)
			content += handleStr
			return true
		})

	} else {
		header = s.replaceMessage(templateContent, module)
	}
	return header + content + tailf
}

func (s *ServiceConverter) convertHandle(rp1 string, rp2 string, content string, module *Module, outputConfig *conf.Output) {
	handlers := s.convertHandler(content, module, rp1)

	if handlers != nil {
		for handler, content := range handlers {
			//templateContent string, module *Module, outputConfig *model.Output, split string, handlerName string
			newContent := s.convertHandler1(content, module, rp2, handler)
			filePath := outputConfig.Output + "/" + strings.Replace(outputConfig.Prefix, "${}", strings.ToLower(handler), -1)
			writeFile(filePath, newContent, outputConfig.Overwrite)
		}
	}
}

func (s *ServiceConverter) convertHandler(templateContent string, module *Module, split string) map[string]string {
	results := strings.Split(templateContent, split)

	handlers := make(map[string]string)

	header := ""
	//content := ""
	tailf := ""

	if len(results) == 3 {
		header = s.replaceMessage(results[0], module)
		tailf = s.replaceMessage(results[2], module)

		module.Foreach(func(handler *ProtoHandler) bool {
			handleStr := s.replaceMessage(results[1], module)
			if split == MsgSplitStr && !handler.IsSession() {
				return true
			}
			if split == RequestSplitStr && !handler.IsRequest() {
				return true
			}
			handleStr = s.replaceHandle(handleStr, handler)
			handlers[handler.GetName()] = header + handleStr + tailf
			return true
		})
	}
	return handlers
}

func (s *ServiceConverter) convertHandler1(templateContent string, module *Module, split string, handlerName string) string {
	results := strings.Split(templateContent, split)

	header := ""
	content := ""
	tailf := ""

	if len(results) == 3 {
		header = s.replaceMessage(results[0], module)
		tailf = s.replaceMessage(results[2], module)

		content := ""
		module.Foreach(func(handler *ProtoHandler) bool {
			if handlerName != handler.GetName() {
				return true
			}
			handleStr := s.replaceMessage(results[1], module)

			if split == MsgSplitStr && !handler.IsSession() {
				return true
			}

			if split == RequestSplitStr && !handler.IsRequest() {
				return true
			}
			handleStr = s.replaceHandle(handleStr, handler)
			content = header + handleStr + tailf
			return false
		})
		if content != "" {
			return content
		}
	} else {
		header = s.replaceMessage(templateContent, module)
	}
	return header + content + tailf
}

func writeFile(filePath string, content string, overwrite bool) {
	if !overwrite {
		//文件存在不允许覆盖
		_, err := os.Stat(filePath)
		if err == nil {
			fmt.Printf("file " + filePath + " alread exist, skip it! \n")
			return
		}
	}
	f, err := os.Create(filePath) //创建文件
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	defer f.Close()
	_, err1 := f.Write([]byte(content)) //写入文件(字节数组)
	if err1 != nil {
		fmt.Printf(err1.Error())
		return
	}
	fmt.Printf("gen code file " + filePath + " success! \n")
}

func (s *ServiceConverter) replaceMessage(content string, module *Module) string {
	content = strings.Replace(content, "${package}", s.data.PackageName, -1)
	content = strings.Replace(content, "${module}", module.Name, -1)
	content = strings.Replace(content, "${Module}", module.UName, -1)
	content = strings.Replace(content, "${request}", module.Name, -1)
	content = strings.Replace(content, "${response}", module.Name, -1)
	return content
}

func (s *ServiceConverter) replaceHandle(content string, handler *ProtoHandler) string {
	content = strings.Replace(content, "${o_desc}", handler.Desc, -1)
	content = strings.Replace(content, "${o_request}", handler.ORequest, -1)
	content = strings.Replace(content, "${o_response}", handler.OResponse, -1)
	content = strings.Replace(content, "${o_request_type}", handler.ORequestType, -1)
	content = strings.Replace(content, "${o_response_type}", handler.OResponseType, -1)
	content = strings.Replace(content, "${o_name}", handler.GetName(), -1)
	return content
}
