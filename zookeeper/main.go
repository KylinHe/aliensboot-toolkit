/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/9/13
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 *******************************************************************************/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	from := *flag.String("from", "", "copy from address")
	to := *flag.String("to", "", "copy to address")
	fromPath := *flag.String("fromPath", "", "fromPath")
	toPath := *flag.String("toPath", "", "node name")
	flag.Parse()

	Transfer(from, to, fromPath, toPath)
	//Transfer("120.77.213.111:2182", "39.108.220.37:2181", "/gok")

	//result, _ := ip.GetLocalPublicIpUseDnspod()
	//log.Debug(result)
}

func uploadWithConfig() {
	address := flag.String("address", "", "address")
	config := flag.String("config", "", "file")
	node := flag.String("node", "", "node name")
	flag.Parse()

	f, err := os.Open(*config)
	if err != nil {
		return
	}
	defer f.Close()

	zkConn, err := Connect(*address)
	if err != nil {
		return
	}
	nodePath := "/" + *node
	Create(zkConn, nodePath)
	Create(zkConn, nodePath+"/config")
	Create(zkConn, nodePath+"/service")

	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}

		var strLine, name, value string
		strLine = string(line)
		strArr := strings.Split(strLine, "    ")
		if len(strArr) >= 2 {
			name = strArr[0]                     //第一个字段是表名
			value = strings.Join(strArr[1:], "") //第二个以上的字段有可能被拆分多个,就需要连接到一起
		} else {
			fmt.Println("##error##", strArr[:1], "表格解析失败，请检查")
		}

		//fmt.Println("name => " + name)
		//fmt.Println("value => " + value)

		path := nodePath + "/config/" + name
		UpdateByPath(zkConn, path, []byte(value))
	}
}

//迁移节点数据数据

func Transfer(oldAddress string, newAddress string, oldPath string, newPath string) {
	oldCon, _, err1 := zk.Connect([]string{oldAddress}, time.Second)
	newCon, _, err2 := zk.Connect([]string{newAddress}, time.Second)
	if err1 != nil {
		fmt.Errorf("old connection error : %v", err1)
		return
	}

	if err2 != nil {
		fmt.Errorf("new connection error : %v", err2)
		return
	}
	copy(oldCon, newCon, oldPath, newPath)
}

func Update(address string, path string, data []byte) {
	con, _, err1 := zk.Connect([]string{address}, time.Second)
	if err1 != nil {
		fmt.Errorf("connection error : %v", err1)
		return
	}

	UpdateByPath(con, path, data)
}

func Connect(address string) (*zk.Conn, error) {
	con, _, err := zk.Connect([]string{address}, time.Second)
	if err != nil {
		return nil, err
	}
	return con, nil
}

func Create(con *zk.Conn, path string) {
	result, _, _ := con.Exists(path)
	if !result {
		con.Create(path, []byte(""), 0, zk.WorldACL(zk.PermAll))
	}
}

func UpdateByPath(con *zk.Conn, path string, data []byte) {
	result, _, _ := con.Exists(path)
	if !result {
		_, err2 := con.Create(path, data, 0, zk.WorldACL(zk.PermAll))
		if err2 != nil {
			fmt.Sprintf("%v", err2)
			return
		}
	}
	_, err := con.Set(path, data, -1)
	if err != nil {
		fmt.Errorf("%v", err)
		return
	}
	fmt.Sprintf("update success : %v", path)
}

func copy(oldConn *zk.Conn, newConn *zk.Conn, oldPath string, newPath string) bool {
	data, _, err1 := oldConn.Get(oldPath)
	if err1 != nil {
		fmt.Errorf("get data error : %v", err1)
		return false
	}
	UpdateByPath(newConn, newPath, data)

	childNames, _, err3 := oldConn.Children(oldPath)
	if err3 != nil {
		fmt.Errorf("get children error : %v", err3)
		return false
	}
	for _, name := range childNames {
		copy(oldConn, newConn, oldPath+"/"+name, newPath+"/"+name)
	}
	return true
}

