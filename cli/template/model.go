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
	"github.com/KylinHe/aliensboot-cli/util"
	"sort"
	"strings"
)

type Type int

type ModelData struct {
	PackageName string //包名
	models map[string]*Model //所有模型
}

func (this *ModelData) AddModel(model *Model) {
	this.models[model.Name] = model
}

type Model struct {
	Name string // 模型名
	Tags []string // 模型标签
	fields map[string]*Field //模型的所有字段
}

func (this *Model) AddField(field *Field) {
	this.fields[field.Name] = field
}

type ServiceData struct {
	PackageName string
	modules     map[string]*Module
}

func (this *ServiceData) EnsureModule(name string) *Module {
	if this.modules == nil {
		this.modules = make(map[string]*Module)
	}
	module := this.modules[name]
	if module == nil {
		module = &Module{Name: name, UName: util.FirstToUpper(name), Handlers: make(map[int]*ProtoHandler), Pushs: make(map[int]string)}
		this.modules[name] = module
	}
	return module
}


type Module struct {
	Name     string
	UName    string
	Handlers map[int]*ProtoHandler
	Pushs    map[int]string
}

type Field struct {
	Name string
	Type string
	Desc string
	Repeated bool
	//Field []Field
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
