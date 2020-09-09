package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type ReturnData struct {
	Code    int                    `json:"code"`
	Data    map[string]interface{} `json:"data"`
	Msg     string                 `json:"msg"`
	Success bool                   `json:"success"`
	Time    string                 `json:"time"`

	//"code": 200,
	//"data": {},
	//"msg": "",
	//"success": true,
	//"time": ""
}

func HttpPost(url string, params string, token string) string {
	client := &http.Client{}

	if isdebug {
		LLog(params)
	}
	req, err := http.NewRequest("POST", SrvUrl+url, strings.NewReader(params))
	if err != nil {
		// handle error
		ErrorLog(err)
		return ""
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Token", token)
	}

	if isdebug {
		LLog(req.Header)
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

func GetToken() string {
	payload := `{"appId":"` + appId + `", "appSecret": "` + appSecret + `"}`

	str := HttpPost(accessTokenUrl, payload, "")
	var rtn ReturnData
	err := json.Unmarshal([]byte(str), &rtn)
	if err != nil {
		ErrorLog("GetToken error", err)
		return ""
	} else {
		if v, ok := rtn.Data["token"]; ok {
			return v.(string)
		} else {
			ErrorLog("GetToken error", str)
			return ""
		}
	}
}

func QueryDbData(sql string, params []interface{}) []map[string]string {
	db := getDB()
	defer db.Close()
	lst, err := QueryDbMap(db, sql, params)
	if err != nil {
		ErrorLog("查询数据失败", err)
		return nil
	} else {
		return lst
	}
}

func StripComments(data []byte) ([]byte, error) {
	data = bytes.Replace(data, []byte("\r"), []byte(""), 0) // Windows
	lines := bytes.Split(data, []byte("\n"))                //split to muli lines
	filtered := make([][]byte, 0)

	for _, line := range lines {
		match, err := regexp.Match(`^\s*#`, line)
		if err != nil {
			return nil, err
		}
		if !match {
			filtered = append(filtered, line)
		}
	}

	return bytes.Join(filtered, []byte("\n")), nil
}

func LoadConfig(path string) map[string]interface{} {

	config_file, err := os.Open(path)
	if err != nil {
		ErrorLog("Failed to open config file：", path, err)
		return nil
	}

	fi, _ := config_file.Stat()

	buffer := make([]byte, fi.Size())
	_, err = config_file.Read(buffer)

	buffer, err = StripComments(buffer) //去掉注释
	if err != nil {
		ErrorLog("Failed to strip comments from json: ", err, string(buffer))
		return nil
	}

	buffer = []byte(os.ExpandEnv(string(buffer))) //特殊
	var f interface{}
	err = json.Unmarshal(buffer, &f) //解析json格式数据

	if err != nil {
		ErrorLog("Failed unmarshalling json: ", err, string(buffer))
		return nil
	}

	m := f.(map[string]interface{})
	return m
}

func Convert2Yuque(source []map[string]string, keyname string) string {
	rtn := make([]map[string]interface{}, 0)

	def := LoadConfig(keyname + ".json")

	for _, val := range source {
		item := make(map[string]interface{})
		for key, v := range def {
			strv := fmt.Sprint(v)
			if vv, ok := val[strv]; ok {
				//定义为列
				item[key] = vv
			} else {
				//定义公共常数
				if cvv, ok := convertDefine[strv]; ok {
					item[key] = cvv
				} else {
					//参数
					item[key] = v
				}
			}
		}
		rtn = append(rtn, item)
	}
	bt, _ := json.Marshal(rtn)
	return `{"voList": ` + string(bt) + `}`
}

func PushData2YuQue(token string, dataType string, counter string) string {
	//查询数据
	sqlStr := SqlDefine[dataType]
	LLog(sqlStr)
	sqlStr = strings.ReplaceAll(sqlStr, "@cnt", counter)
	sqlStr = strings.ReplaceAll(sqlStr, "@ydt", time.Now().Add(-24*time.Hour).Local().Format("20060102"))
	sqlStr = strings.ReplaceAll(sqlStr, "@tdt", time.Now().Local().Format("20060102"))
	//params := make([]interface{}, 0)
	//params = append(params, sql.Named("cnt", counter))
	//lst := []map[string]string{{"a":"1"}}
	lst := QueryDbData(sqlStr, nil)
	LLog(dataType, "数据库查询", counter, "共", len(lst), "条数据")

	//转换格式
	payload := Convert2Yuque(lst, dataType)

	//发送数据
	LLog(dataType, "发送", URLDefine[dataType])
	rtn := HttpPost(URLDefine[dataType], payload, token)
	LLog(dataType, "发送完成", rtn)

	return rtn
}

func main() {
	getConfig()
	fmt.Println("接口地址", SrvUrl)
	fmt.Println("接口内容", apiList)
	fmt.Println("接口数据", rangeOfPost)
	token := GetToken()
	LLog("token", token)
	if token != "" {
		//循环类型
		for _, t := range apiList {
			//循环分组数据
			for _, cnt := range rangeOfPost {
				LLog("开始处理：", t, cnt)
				PushData2YuQue(token, t, cnt)
			}
		}
	}
}
