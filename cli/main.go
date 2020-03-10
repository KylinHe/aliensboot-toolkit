/*******************************************************************************
 * Copyright (c) 2015, 2017 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2018/10/29
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package main

import "github.com/KylinHe/aliensboot-cli/command"

func main() {
	//a := int64(math.Pow(2, 1))
	//b := int64(math.Pow(2, 2))
	//c := int64(math.Pow(2, 50))
	//
	//d := int64(math.Pow(2, 7))
	//result := a + b + c
	//fmt.Println(d&result == d)
	//fmt.Println(b&result == b)
	//fmt.Println(c&result == c)


	//parser := &template.ProtoParser{}
	//parser.Parse("/Users/hejialin/git/aliens/aliensboot/aliensboot-custom-servers/slg_server/src/protocol/game_model.proto")
	//fmt.Println(parser.GetModelData())
	//

	//converter := template.ModelConverter{}
	//converter.Convert(&conf.ModelGenConfig{
	//	ProtoPath:"/Users/hejialin/git/aliens/aliensboot/aliensboot-custom-servers/slg_server/src/protocol/game_model.proto",
	//	TemplatePath:"/Users/hejialin/git/aliens/aliensboot/aliensboot-custom-servers/slg_server/src/protocol/game_model.template",
	//	OutputPath:"/Users/hejialin/git/aliens/aliensboot/aliensboot-custom-servers/slg_server/src/protocol/game_model.output",
	//})

	command.Execute()
}
