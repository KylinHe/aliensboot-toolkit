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
	"encoding/json"
	"github.com/KylinHe/aliensboot-cli/util"
	"github.com/KylinHe/aliensboot-core/cluster/center"
	"github.com/KylinHe/aliensboot-core/config"
	"github.com/KylinHe/aliensboot-core/log"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

var endpoint string
var dir bool
var remote string
var local string
var username string
var password string
var md5 string
var suffix string
var check string

func init() {
	dataUploadCmd.Flags().BoolVarP(&dir, "dir", "d", true, "true/false")
	dataUploadCmd.Flags().StringVarP(&check, "check", "c", "true", "true/false")
	dataUploadCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "127.0.0.1:2379", "-endpoint 127.0.0.1:2379")
	dataUploadCmd.Flags().StringVarP(&remote, "remote", "r", "/root/my-cluster/config", "-")
	dataUploadCmd.Flags().StringVarP(&username, "username", "u", "", "-")
	dataUploadCmd.Flags().StringVarP(&password, "password", "p", "", "-")
	dataUploadCmd.Flags().StringVarP(&local, "local", "l", "/usr/local/data/config", "-")
	dataUploadCmd.Flags().StringVarP(&md5, "md5", "m", "", "-md5 md5.sum")
	dataUploadCmd.Flags().StringVarP(&suffix, "suffix", "s", ".json", "-suffix .json")

	dataCmd.AddCommand(dataUploadCmd)
}

var dataUploadCmd = &cobra.Command{
	Use:   "upload Ex. aliens boot upload -r %remotePath% -l %localPath%",
	Short: "upload",
	Run: func(cmd *cobra.Command, args []string) {
		if dir {
			uploadDir(remote, local)
		} else {
			log.Warnf("un not support")
		}
	},
}

func uploadDir(remoteRoot string, localRoot string) {
	registry := center.ETCDServiceCenter{}
	registry.ConnectCluster(config.ClusterConfig{
		Servers:[]string{endpoint},
		Username:username,
		Password:password,
	})

	rd, err := ioutil.ReadDir(localRoot)
	if err != nil {
		log.Fatalf("invalid param localPath %v", localRoot)
	}

	var md5Sum map[string]string
	if md5 != "" {
		md5 = localRoot + string(filepath.Separator) + md5
		md5Sum = util.ReadMD5File(md5)
	}

	count := 0
	success := 0

	haveChange := false
	uploadData := make(map[string][]byte)
	for _, fi := range rd {
		if path.Ext(fi.Name()) != suffix {
			continue
		}
		fileName := strings.TrimSuffix(fi.Name(), suffix)
		filePath := localRoot + string(filepath.Separator) + fi.Name()
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("invalid data file %v", filePath)
			return
		}

		if md5Sum != nil {
			md5Old := md5Sum[fileName]
			md5New := util.GetBytesMD5(data)
			//
			if md5Old == md5New {
				continue
			} else {
				md5Sum[fileName] = md5New
				haveChange = true
				log.Infof("check file change %v - [%v-%v]", fileName, md5New, md5Old)
			}
		}

		// 开启文件验证
		if check == "true" {
			if suffix == ".json" {
				var verifyData interface{}
				err := json.Unmarshal(data, &verifyData)
				if err != nil {
					log.Fatalf("invalid data file %v - %v", filePath, err)
					return
				}
			}
		}
		remotePath := remoteRoot + center.NodeSplit + fileName
		uploadData[remotePath] = data
	}

	for path, data := range uploadData {
		count++
		log.Infof("uploading data %v", path)
		//上传数据
		result := registry.UploadData(path, data)
		if result {
			success ++
		}
	}

	if haveChange {
		util.WriteMD5File(md5, md5Sum)
	}

	log.Infof("upload done, count %v  success %v", count, success)
	//center.UploadData()
}
