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
	"github.com/KylinHe/aliensboot-cli/proto"
	"github.com/KylinHe/aliensboot-core/common/util"
	"os"
	"regexp"
	"strings"
)

const (
	RequestTag  = "request"
	ResponseTag = "response"
	PushTag     = "push"
)

type ProtoParser struct {
	serviceData *ServiceData
	modelData *ModelData
}


func (p *ProtoParser) GetModelData() *ModelData {
	return p.modelData
}
//var message = &ServiceData{}

func (p *ProtoParser) Parse(protoPath string) {
	p.serviceData = &ServiceData{modules: make(map[string]*Module)}
	p.modelData = &ModelData{models: make(map[string]*Model)}
	//"/Users/hejialin/git/server/kylin/src/aliens/protocol/scene/protocol.proto"
	reader, _ := os.Open(protoPath)
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, _ := parser.Parse()

	for _, element := range definition.Elements {
		switch element.(type) {
		case *proto.Package:
			p.serviceData.PackageName = element.(*proto.Package).Name
			p.modelData.PackageName = element.(*proto.Package).Name
			break
		case *proto.Message:
			tag := getComment(element.(*proto.Message).Doc())
			if tag == "" {
				break
			}
			p.handleMessage(element.(*proto.Message), tag)
		}
	}
}

func getComment(comment *proto.Comment) string {
	//if document == nil {
	//	return ""
	//}
	//comment := document.Doc()
	if comment == nil {
		return ""
	}
	tag := strings.TrimSpace(comment.Message())
	return tag
}

func (p *ProtoParser)handleMessage(member *proto.Message, tag string) {
	tags := ParseModelTag(tag)
	if tags != nil && len(tags) > 0 {
		p.handleModel(member, tags)
	} else {
		for _, visitee := range member.Elements {
			moduleField, ok := visitee.(*proto.Oneof)
			if !ok {
				continue
			}
			module := p.serviceData.EnsureModule(moduleField.Name)
			for _, moduleHandleField := range moduleField.Elements {
				field, ok := moduleHandleField.(*proto.OneOfField)
				if !ok {
					continue
				}
				handler := module.Handlers[field.Sequence]
				if handler == nil {
					handler = &ProtoHandler{}
					module.Handlers[field.Sequence] = handler
				}
				if field.Doc() != nil {
					handler.Desc = field.Doc().Message()
				}
				if tag == RequestTag {
					handler.ORequest = util.FirstToUpper(field.Name)
					handler.ORequestType = util.FirstToUpper(field.Type)
					//message.RequestName = element.(*proto.Message).Name
				} else if tag == ResponseTag {
					handler.OResponse = util.FirstToUpper(field.Name)
					handler.OResponseType = util.FirstToUpper(field.Type)
					//message.ResponseName = element.(*proto.Message).Name
				}
			}
		}
	}
}

func ParseModelTag(tag string) []string {
	reg := regexp.MustCompile(`^model\[([\S\s]+)\]$`)
	results := reg.FindStringSubmatch(tag)
	if len(results) != 2 {
		return nil
	}
	results = strings.Split(results[1], ",")
	for idx, result := range results {
		results[idx] = strings.TrimSpace(result)
	}
	return results
}

func (p *ProtoParser) handleModel(member *proto.Message, tags []string)  {
	model := &Model{
		Name:   member.Name,
		Tags:   tags,
		fields: make(map[string]*Field),
	}
	for _, visitee := range member.Elements {
		field, ok := visitee.(*proto.NormalField)
		if !ok {
			continue
		}
		field.Doc()
		comment := getComment(field.InlineComment)
		//fmt.Println(field.Options)
		model.AddField(&Field{Name:field.Name, Type:field.Type, Desc:comment, Repeated:field.Repeated})
	}
	p.modelData.AddModel(model)
}

//func handlePushMessage(member *proto.Message) []Field {
//	result := make([]Field, 0)
//	for _, visitee := range member.Elements {
//		field, ok := visitee.(*proto.NormalField)
//		if !ok {
//			continue
//		}
//		comment := getComment(field)
//		result = append(result, Field{Name:field.Name, Type:field.Type, Desc:comment, Repeated:field.Repeated})
//	}
//	return result
//}
