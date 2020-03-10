package template

import (
	"fmt"
	"github.com/KylinHe/aliensboot-cli/conf"
	"github.com/KylinHe/aliensboot-cli/util"
	"io/ioutil"
	"strings"
)

const (
	ModelStartTag     = "<model>" //模型开始标签

	ModelEndTag     = "</model>" //模型结束标签

	FieldStartTag     = "<field>" //字段开始标签

	FieldEndTag     = "</field>" //字段结束标签
)

type ModelConverter struct {
	data *ModelData
}

func (s *ModelConverter) Convert(config *conf.ModelGenConfig) {
	parser := &ProtoParser{}
	parser.Parse(config.ProtoPath)
	if parser.modelData == nil {
		fmt.Println("parse model error")
		return
	}
	s.data = parser.modelData
	templatePath := config.TemplatePath
	b, err := ioutil.ReadFile(templatePath)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	templateContent := string(b)
	modelContents := util.ParseTags(templateContent, ModelStartTag, ModelEndTag)
	content := ""

	for i:=0; i<len(modelContents); i++ {
		modelContent := modelContents[i]
		content += s.replacePackage(modelContent.Header)
		modelTemplate := s.replacePackage(modelContent.Content)
		for _, model := range s.data.models {
			content += s.ConvertModel(modelTemplate, model)
		}
		content += s.replacePackage(modelContent.Tail)
	}
	writeFile(config.OutputPath, content, true)
}

func (s *ModelConverter) ConvertModel(templateContent string, model *Model) string {
	content := ""
	fieldContents := util.ParseTags(templateContent, FieldStartTag, FieldEndTag)
	for i:=0; i<len(fieldContents); i++ {
		fieldContent := fieldContents[i]
		content +=  s.replaceModel(fieldContent.Header, model)

		fieldTemplate := s.replaceModel(fieldContent.Content, model)
		for _, field := range model.fields {
			content += s.replaceField(fieldTemplate, field)
		}
		content +=  s.replaceModel(fieldContent.Tail, model)
	}
	return content
}


func (s *ModelConverter) replaceField(content string, field *Field) string {
	//content = s.replaceModel(content, model)
	content = strings.Replace(content, "${fieldName}", field.Name, -1)
	content = strings.Replace(content, "${FIELDNAME}", strings.ToUpper(field.Name), -1)
	content = strings.Replace(content, "${fieldname}", strings.ToLower(field.Name), -1)
	content = strings.Replace(content, "${FieldName}", util.FirstToUpper(field.Name), -1)
	return content
}

func (s *ModelConverter) replaceModel(content string, model *Model) string {
	content = strings.Replace(content, "${modelName}", model.Name, -1)
	content = strings.Replace(content, "${MODELNAME}", strings.ToUpper(model.Name), -1)
	content = strings.Replace(content, "${modelname}", strings.ToLower(model.Name), -1)
	content = strings.Replace(content, "${ModelName}", util.FirstToUpper(model.Name), -1)
	return content
}

func (s *ModelConverter) replacePackage(content string) string {
	content = strings.Replace(content, "${packageName}", s.data.PackageName, -1)
	content = strings.Replace(content, "${PACKAGENAME}", strings.ToUpper(s.data.PackageName), -1)
	content = strings.Replace(content, "${packagename}", strings.ToLower(s.data.PackageName), -1)
	content = strings.Replace(content, "${PackageName}", util.FirstToUpper(s.data.PackageName), -1)
	return content
}
