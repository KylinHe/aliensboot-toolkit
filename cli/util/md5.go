//create file md5sum.md5
//check file's md5sum from md5sum.md5

package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/go-yaml/yaml"
)


// 获取文件的md5码
func GetMD5(fileName string) string {
	data := ReadFile(fileName)
	return GetBytesMD5(data)
}

func GetBytesMD5(data []byte) string {
	return hex.EncodeToString(byte2string(md5.Sum(data)))
}

func WriteMD5File(md5File string, data map[string]string) {
	wData, err := yaml.Marshal(data)
	if err != nil {
		log.Errorf("write md5 file err: %v", err)
	} else {
		WriteFile(md5File, wData)
	}
}

//获取MD5检查文件内容
func ReadMD5File(md5File string) map[string]string {
	data := ReadFile(md5File)
	result := make(map[string]string)
	if data == nil {
		return result
	}
	err := yaml.Unmarshal(data, result)
	if err != nil {
		log.Errorf("read md5 file err: %v", err)
	}
	return result
}

func byte2string(in [16]byte) []byte {
	return in[:16]
}