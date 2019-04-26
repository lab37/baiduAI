package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var ocrErrorCode = map[float64]string{
	100: "无效参数",
	110: "Token过期失效",
}

var ocrTypeCode = map[float64]string{
	216015: "模块关闭",
	216100: "非法参数",
	216101: "参数数量不够",
	216102: "业务不支持",
	216103: "参数太长",
	216110: "ID不存在",
	216111: "非法用户ID",
	216200: "空的图片",
	216201: "图片格式错误",
	216202: "图片大小错误",
	216300: "DB错误",
	216400: "后端系统错误",
	216401: "内部错误",
	216500: "未知错误",
	216600: "身份证的ID格式错误",
	216601: "身份证的ID和名字不匹配",
	216611: "用户不存在",
	216613: "用户查找不到",
	216614: "图片信息不完整",
	216615: "处理图片信息失败",
	216616: "图片已存在",
	216617: "添加用户失败",
	216618: "群组里没有用户",
	216630: "识别错误",
	216631: "识别银行卡错误",
}

/**
 * 身份证识别
 */
func OcrIdCard(imgPath string, isFront bool) map[string]interface{} {
	imgbytes, err := ioutil.ReadFile(imgPath)
	if err != nil {
		log.Println(err)
	}
	postArgs := url.Values{}
	postArgs.Set("access_token", currentAccessToken.AccessToken)
	postArgs.Set("image", base64.StdEncoding.EncodeToString(imgbytes))
	postArgs.Set("id_card_side", "front")     // front 正面  back 背面
	postArgs.Set("detect_direction", "false") // 是否检测图像朝向[true/false]，默认不检测，即：false。朝向是指输入图像是正常方向、逆时针旋转90/180/270度。

	if !isFront {
		postArgs.Set("id_card_side", "back")
	}

	resp, _ := http.PostForm(IDCARD_API_URI, postArgs)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		panic(err)
	}

	map_result := make(map[string]interface{})
	json.Unmarshal(data, &map_result)

	error_msg, ok := map_result["error_msg"]
	if ok {
		panic(error_msg)
	}

	return map_result
}

/**
 * 身份证识别
 */
func OcrBankCard(imgPath string) map[string]interface{} {
	imgbytes, err := ioutil.ReadFile(imgPath)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	postArgs := url.Values{}
	postArgs.Set("access_token", currentAccessToken.AccessToken)
	postArgs.Set("image", base64.StdEncoding.EncodeToString(imgbytes))

	resp, _ := http.PostForm(BANKCARD_API_URI, postArgs)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		panic(err)
	}

	map_result := make(map[string]interface{})
	json.Unmarshal(data, &map_result)

	error_msg, ok := map_result["error_msg"]
	if ok {
		panic(error_msg)
	}

	return map_result
}

/**
 * 通用ocr识别
 */
func OcrGeneral(imgPath string) map[string]interface{} {
	imgbytes, err := ioutil.ReadFile(imgPath)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	postArgs := url.Values{}
	postArgs.Set("access_token", currentAccessToken.AccessToken)
	postArgs.Set("image", base64.StdEncoding.EncodeToString(imgbytes))
	postArgs.Set("recognize_granularity", "big")  // 是否定位单字符位置，big：不定位单字符位置，默认值；small：定位单字符位置
	postArgs.Set("mask", "")                      // 是否检测图像朝向[true/false]，默认不检测，即：false。朝向是指输入图像是正常方向、逆时针旋转90/180/270度。
	postArgs.Set("language_type", "CHN_ENG")      // CHN_ENG：中英文混合； ENG：英文； POR：葡萄牙语； FRE：法语； GER：德语； ITA：意大利语； SPA：西班牙语； RUS：俄语； JAP：日语
	postArgs.Set("detect_direction", "false")     // 是否检测图像朝向[true/false]，默认不检测，即：false。朝向是指输入图像是正常方向、逆时针旋转90/180/270度。
	postArgs.Set("detect_language", "false")      // 是否检测语言，默认不检测。当前支持（中文、英语、日语、韩语）
	postArgs.Set("classify_dimension", "lottery") // 分类维度（根据OCR结果进行分类）
	postArgs.Set("vertexes_location", "false")    //是否返回文字外接多边形顶点位置

	resp, _ := http.PostForm(GENERALOCR_API_URI, postArgs)

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		panic(err)
	}

	map_result := make(map[string]interface{})
	json.Unmarshal(data, &map_result)

	error_msg, ok := map_result["error_msg"]
	if ok {
		panic(error_msg)
	}

	return map_result
}

func OcrTable(imgPath string) string {
	imgbytes, err := ioutil.ReadFile(imgPath)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	postArgs1 := url.Values{}
	postArgs1.Set("access_token", currentAccessToken.AccessToken)
	postArgs1.Set("image", base64.StdEncoding.EncodeToString(imgbytes))

	resp, _ := http.PostForm(TABLE_API_REQ_URI, postArgs1)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		panic(err)
	}

	map_result := make(map[string]interface{})
	json.Unmarshal(data, &map_result)
	if v, ok := map_result["error_code"].(float64); ok {
		if v == 110 {
			log.Println("access_token过期，准备重新获取accessToken")
			updateToken()
			postArgs1.Set("access_token", currentAccessToken.AccessToken)
			resp, _ = http.PostForm(TABLE_API_REQ_URI, postArgs1)
			defer resp.Body.Close()
			data, err = ioutil.ReadAll(resp.Body)
		} else {
			log.Println(v, map_result["error_code"].(string))
			panic("系统退出")
		}

	}

	requestID := ""
	if v, ok := map_result["result"].([]interface{})[0].(map[string]interface{}); ok {
		requestID = v["request_id"].(string)

	}

	postArgs2 := url.Values{}
	postArgs2.Set("access_token", currentAccessToken.AccessToken)
	postArgs2.Set("request_id", requestID)

	downLoadUrl := "http://aaa.bbb.ccc"
	for range time.Tick(5000 * time.Millisecond) {

		resp, _ := http.PostForm(TABLE_API_RESP_URI, postArgs2)
		defer resp.Body.Close()
		data2, err := ioutil.ReadAll(resp.Body)
		if nil != err {
			panic(err)
		}
		map_result2 := make(map[string]interface{})
		json.Unmarshal(data2, &map_result2)

		if vv, ok := map_result2["result"].(map[string]interface{}); ok {
			if vv["ret_code"].(float64) == 3 {
				downLoadUrl = vv["result_data"].(string)
				break
			}

		}

	}

	return downLoadUrl
}
