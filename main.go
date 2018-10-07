package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"log"
	"errors"
	"strings"
	"github.com/axgle/mahonia"
	"io/ioutil"
)

// 词法分析

type Config struct {
	ApiKey    string // client id
	SecretKey string // client secret
}

type AccessTokenResponse struct {
	ExpiresIn        int64  `json:"expires_in"`        // 过期时间
	AccessToken      string `json:"access_token"`      // 访问码
	Error            string `json:"error"`             // 错误码
	ErrorDescription string `json:"error_description"` // 错误信息
}

func getToken(cfg *Config) (string, error) {
	var b bytes.Buffer
	b.WriteString("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=")
	b.WriteString(cfg.ApiKey)
	b.WriteString("&client_secret=")
	b.WriteString(cfg.SecretKey)
	res, err := http.Post(b.String(), "application/json; charset=UTF-8", nil)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	ret := new(AccessTokenResponse)
	if err = json.NewDecoder(res.Body).Decode(ret); err != nil {
		return "", err
	}
    
	if ret.Error != "" {
		return "", errors.New(ret.Error + ": " + ret.ErrorDescription)
	}

	return ret.AccessToken, nil
}


func analyzeWords(tk string, ex string ) (*wordResult, error) {
	enc:=mahonia.NewEncoder("gbk")
	str:=enc.ConvertString(ex)
	var a bytes.Buffer
	a.WriteString("{")
	a.WriteString(`"text":"`)
	a.WriteString(str)
	a.WriteString(`"}`)
	 
    queryStr :=a.String()
	var b bytes.Buffer
	b.WriteString("https://aip.baidubce.com/rpc/2.0/nlp/v1/lexer?access_token=")
	b.WriteString(tk)

	resp, err := http.Post(b.String(), "application/json", strings.NewReader(queryStr))
	if err != nil {
		log.Println("Get server response err:",err)
	}
	defer resp.Body.Close()
	respRst := new(wordResult)
	body, err := ioutil.ReadAll(resp.Body)
	ostr := mahonia.NewDecoder("gbk").ConvertString(string(body))
	log.Println(ostr)
	if strings.Contains(ostr, `"error_code":110`) || strings.Contains(ostr, `"error_code":111`) {
	errc := errors.New("bad token")
	return respRst, errc	
	}
	
	if err = json.Unmarshal([]byte(ostr),&respRst); err != nil {
		log.Println("Parse response to json err:",err)
	}
	return respRst, nil
	
}


type address struct {
 AddType string                  `json:"type"`
 ByteOffset int64                `json:"byte_offset"`
 ByteLength int64                `json:"byte_length"`
}
type wordItem struct {
 Item string                     `json:"item"`
 Ne string                       `json:"ne"`
 Pos string                      `json:"pos"`
 ByteOffset int64                `json:"byte_offset"`
 ByteLength int64                `json:"byte_length"`
 Uri string                      `json:"uri"`
 Formal string                   `json:"formal"`
 BasicWord []string              `json:"base_words"`
 LocalDetails []address          `json:"loc_details"`
}
type wordResult struct {
 Text string                      `json:"text"`
 Items []wordItem                 `json:"items"`
}

func main() {
    cfg := new(Config)
	cfg.ApiKey = "bwG08GfxuwBIhIjTGtC3e4na"
	cfg.SecretKey = "lhPE5TjKKT0VXVDGMvF3WGuRao4SgMuh"
	
	istrc := "山东泰安泰山区泰玻大街中段峪龙佳苑55号楼3单元1603，周京成收，15853873773"
	
	
	tokenb, err := ioutil.ReadFile("token.cfg")
    if err != nil {
	     
        log.Println("Get token from file err:",err)
    }
    token := string(tokenb)
	
	respRst2, err := analyzeWords(token, istrc)
	if err != nil {
	log.Println(err)
	accessToken, err := getToken(cfg)
	if err != nil {
	log.Println(err)
	}
	token = accessToken
	err = ioutil.WriteFile("token.cfg", []byte(token), 0644)
	respRst2, err = analyzeWords(token, istrc)
	}
	log.Println(respRst2.Text)
	
}
