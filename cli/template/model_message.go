/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/8/29
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package template

import (
	"sort"
	"strings"
)

type Type int

type ModelMessage struct {
	models map[string]*Model //模型
}

type Model struct {
	Name string //
	//Props map[string]
}

type ServiceMessage struct {
	PackageName string
	modules     map[string]*Module
}

func (this *ServiceMessage) EnsureModule(name string) *Module {
	if this.modules == nil {
		this.modules = make(map[string]*Module)
	}
	module := this.modules[name]
	if module == nil {
		module = &Module{Name: name, UName: strFirstToUpper(name), Handlers: make(map[int]*ProtoHandler), Pushs: make(map[int]string)}
		this.modules[name] = module
	}
	return module
}

/**
 * 字符串首字母转化为大写 ios_bbbbbbbb -> iosBbbbbbbbb
 */
func strFirstToUpper(str string) string {
	f := str[0:1]
	t := str[1:]

	return strings.ToUpper(f) + t
}

type Module struct {
	Name     string
	UName    string
	Handlers map[int]*ProtoHandler
	Pushs    map[int]string
}

func (this *Module) Foreach(callback func(handler *ProtoHandler) bool) {
	var keys []int
	for k := range this.Handlers {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		if !callback(this.Handlers[k]) {
			return
		}
	}
}

func (this *Module) ForeachPush(callback func(push string) bool) {
	var keys []int
	for k := range this.Pushs {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		if !callback(this.Pushs[k]) {
			return
		}
	}
}

type ProtoHandler struct {
	Name      string
	Desc      string
	ORequest  string
	OResponse string
	ORequestType string
	OResponseType string
	//OPush string
}

func (this *ProtoHandler) IsSession() bool {
	return this.ORequest != "" && this.OResponse != ""
	//return this.ORequest != "" && this.OResponse != ""
}

//纯请求
func (this *ProtoHandler) IsRequest() bool {
	return this.ORequest != "" && this.OResponse == ""
	//return this.ORequest != "" && this.OResponse != ""
}

func (this *ProtoHandler) GetName() string {
	if this.ORequest != "" && this.Name == "" {
		this.Name = strings.Replace(this.ORequest, "C2S_", "", 1)
	}
	return this.Name
}
