package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ReturnData struct {
	Code int `json:"code"`
	Data map[string]interface{} `json:"data"`
	Msg string `json:"msg"`
	Success bool `json:"success"`
	Time string `json:"time"`

	//"code": 200,
	//"data": {},
	//"msg": "",
	//"success": true,
	//"time": ""
}

func HttpPost(url string, params string,token string) string {
	client := &http.Client{}

	req, err := http.NewRequest("POST", SrvUrl + url, strings.NewReader(params))
	if err != nil {
		// handle error
		ErrorLog(err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")
	if token!="" {
		req.Header.Set("accessToken", token)
	}

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		ErrorLog(err)
		return ""
	}

	//fmt.Println(string(body))
	return string(body)
}

func GetToken() string{
	payload:=`{"appId":"`+appId+`", "appSecret": "`+appSecret+`"}`

	str := HttpPost(accessTokenUrl,payload,"")
	var rtn ReturnData
	err:=json.Unmarshal([]byte(str),&rtn)
	if err!=nil{
		ErrorLog("GetToken error" ,err)
		return ""
	}else{
		if v,ok:=rtn.Data["accessToken"];ok{
			return v.(string)
		}else{
			ErrorLog("GetToken error" ,str)
			return ""
		}
	}
}

func QueryDbData(sql string ,params []interface{}) []map[string]string {
	db := getDB()
	defer db.Close()
	lst,err := QueryDbMap(db, sql, params)
	if err!=nil{
		ErrorLog("查询数据失败",err)
		return nil
	} else{
		return lst
	}
}



func Convert2Yuque(source []map[string]string, keyname string) string{
	rtn := make([]map[string]interface{},0)
	for _,val := range source{
		item := make(map[string]interface{})
		for key,v := range convertDefine[keyname]{
			if vv,ok := val[v];ok{
				//定义为列
				item[key] = vv
			}else{
				//定义为常数
				item[key] = v
			}
		}
		rtn = append(rtn,item)
	}
	bt,_ := json.Marshal(rtn)
	return `{"voList": `+string(bt)+`}`
}




func PushData2YuQue(token string, dataType string,counter string ) string {
	//查询数据
	sqlStr := SqlDefine[dataType]
	params := make([]interface{}, 0)
	params = append(params, sql.Named("cnt", counter))
	lst := QueryDbData(sqlStr,params)
	LLog(dataType,"数据库查询",counter,"共",len(lst),"条数据")

	//转换格式
	payload:=Convert2Yuque(lst,dataType)

	//发送数据
	LLog(dataType,"发送",URLDefine[dataType])
	rtn := HttpPost(URLDefine[dataType],payload,token)
	LLog(dataType,"发送完成",rtn)

	return rtn
}



func main(){
	getConfig()
	fmt.Println("接口地址",SrvUrl)
	fmt.Println("接口内容",apiList)
	fmt.Println("接口数据",rangeOfPost)
	token := GetToken()
	if token!="" {
		//循环类型
		for _, t := range apiList {
			//循环分组数据
			for _, cnt := range rangeOfPost {
				PushData2YuQue(token, t, cnt)
			}
		}
	}
}


