package app

import (
	"encoding/json"
	"io/ioutil"
)

// -------------小程序配置 start-------------

type AppletConfigJson struct{
	AppId	string
	Secret  string
}

var AppletConfig = make(map[string]AppletConfigJson)

// -------------小程序配置 end-------------
// -------------mysql配置 start-------------

type MysqlConfigJson struct {
	Database string
	Port     int64
	Charset  string
	Protocol string
	Master   MysqlConn
	Slaves   []MysqlConn
}

var MysqlConfig *MysqlConfigJson

// -------------mysql配置 end-------------
// -------------websocket配置 start-------------

type WebsocketConfigJson struct {
	Addr string
	Ip string
	Port string
	Path string
}

var WebsocketConfig *WebsocketConfigJson

// -------------websocket配置 end-------------
// -------------mongodb配置 start-------------

type MongoConfigJson struct {
	Ip 	       string
	Port       string
	Username   string
	Password   string
}

var MongoConfig *MongoConfigJson

// -------------mongodb配置 end-------------



func LoadConfig() {
	var basePath = Path + "/config/" + ENV + "/"
	files, _ := ioutil.ReadDir(basePath)
	for _, f := range files {
		configJson, _ := ioutil.ReadFile(basePath + f.Name())
		switch f.Name() {
		case "applet.json":
			var listApplet []AppletConfigJson
			json.Unmarshal(configJson, &listApplet)
			for _, v := range listApplet {
				AppletConfig[v.AppId] = v
			}
		case "db.json":
			json.Unmarshal(configJson, &MysqlConfig)
		case "websocket.json":
			json.Unmarshal(configJson, &WebsocketConfig)
		case "mongo.json":
			json.Unmarshal(configJson, &MongoConfig)
		}
	}
}
