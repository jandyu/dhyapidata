package main

import (
	"ngrok/src/github.com/Unknwon/goconfig"
	"strings"
)

var SrvUrl = "http://202.101.190.110:9009"
var appId = "16078425"
var appSecret = "c286ec62e31a7804026c8b1433ceec0779a7a31e"
var enterpriseCode="734514772"

var rangeOfPost=[]string{"88010101","88010102","88010103","88010104","88010105","88010106","88010107","88010108","88010109","88010110","88010111","88040101","88040101","88040101","88040101"}
var apiList =[]string{"supplier","product","category","brand","purchase","sale"}

var accessTokenUrl="/nfwlApi/auth/getAccessToken"

//接口配置
var(
	convertDefine = map[string]map[string]string{
		"supplier":map[string]string{
			"enterpriseCode":enterpriseCode,
			"code":"supplier_id",
		},
	}
	URLDefine = map[string]string{
		"supplier":"/api/supplier/save",
		"product":"/api/product/save",
		"category":"/api/category/save",
		"brand":"/api/brand/save",
		"purchase":"/api/pdtPurchase/save",
		"sale":"/api/pdtSale/save",
	}
	SqlDefine= map[string]string{
		"supplier":"/api/supplier/save",
		"product":"/api/product/save",
		"category":"/api/category/save",
		"brand":"/api/brand/save",
		"purchase":"/api/pdtPurchase/save",
		"sale":"/api/pdtSale/save",
	}
)

func getConfig(){
	//读取配置文件：
	cfg, _ := goconfig.LoadConfigFile("api.conf")

	SrvUrl, _ = cfg.GetValue("server", "SrvUrl")
	URLDefine["supplier"],_= cfg.GetValue("server","supplier")
	URLDefine["product"],_= cfg.GetValue("server","product")
	URLDefine["category"],_= cfg.GetValue("server","category")
	URLDefine["brand"],_= cfg.GetValue("server","brand")
	URLDefine["purchase"],_= cfg.GetValue("server","purchase")
	URLDefine["sale"],_= cfg.GetValue("server","sale")
	accessTokenUrl,_= cfg.GetValue("server","accessToken")

	enterpriseCode,_= cfg.GetValue("server","enterpriseCode")
	appId,_= cfg.GetValue("server","appId")
	appSecret,_= cfg.GetValue("server","appSecret")

	r,_:= cfg.GetValue("data","range")
	rangeOfPost = strings.Split(r,",")

	api,_:= cfg.GetValue("data","api")
	apiList = strings.Split(api,",")
}