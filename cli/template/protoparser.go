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
	"strings"
)

const (
	RequestTag  = "request"
	ResponseTag = "response"
	PushTag     = "push"
)

var message = &ServiceMessage{}

//Handlers:make( map[int]*ProtoHandler)

func ParseProto(protoPath string) *ServiceMessage {
	message = &ServiceMessage{modules: make(map[string]*Module)}
	//"/Users/hejialin/git/server/kylin/src/aliens/protocol/scene/protocol.proto"
	reader, _ := os.Open(protoPath)
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, _ := parser.Parse()

	for _, element := range definition.Elements {
		switch element.(type) {
		case *proto.Package:
			message.PackageName = element.(*proto.Package).Name
			//log.Println()
			break
		case *proto.Message:
			comment := element.(*proto.Message).Doc()
			if comment == nil {
				break
			}

			tag := strings.TrimSpace(comment.Message())
			if tag == "" {
				break
			}

			handleMessage(element.(*proto.Message), tag)

			//proto.Walk(, proto.WithOneof(&messageWalk{tag:tag}.handleMessage))

			//if tag == REQUEST_TAG {
			//	//message.RequestName = element.(*proto.Message).Name
			//} else if tag == RESPONSE_TAG  {
			//	//message.ResponseName = element.(*proto.Message).Name
			//} else if tag == PUSH_TAG{
			//	//message.PushName = element.(*proto.Message).Name
			//}
		}
	}

	return message

}

func handleMessage(member *proto.Message, tag string) {
	for _, visitee := range member.Elements {
		moduleField, ok := visitee.(*proto.Oneof)
		if !ok {
			continue
		}

		module := message.EnsureModule(moduleField.Name)
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
			} else if tag == PushTag {
				module.Pushs[field.Sequence] = util.FirstToUpper(field.Name)
				//message.PushName = element.(*proto.Message).Name
			}
			//if handler.ORequest == "" {
			//	handler.ORequest = util.FirstToUpper(field.Name)
			//} else {
			//	handler.OResponse = util.FirstToUpper(field.Name)
			//}
		}
	}

	//m.
	//proto.Walk(m, proto.WithOneof(handleHandle))

}
