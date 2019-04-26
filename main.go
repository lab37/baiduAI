package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	BAIDU_AI_APPID = 16114169
	BAIDU_AI_KEY   = "GSE6LEdGFvhSFf7vSVC5TIGD"
	BAIDU_AI_CRET  = "sIzN7VRya9hT82NpqGwp67Lbcf7sfaYj"
	API_TIMEOUT    = 120
)

const (
	TOKEN_API_URI = "https://openapi.baidu.com/oauth/2.0/token"

	// 文本相关API
	SEG_API_URI       = "https://aip.baidubce.com/rpc/2.0/nlp/v1/wordseg"
	WORDPOS_API_URI   = "https://aip.baidubce.com/rpc/2.0/nlp/v1/wordpos"
	WORDEMBED_API_URI = "https://aip.baidubce.com/rpc/2.0/nlp/v1/wordembedding"
	DNNL_API_URI      = "https://aip.baidubce.com/rpc/2.0/nlp/v1/dnnlm_cn"
	SIMNET_API_URI    = "https://aip.baidubce.com/rpc/2.0/nlp/v1/simnet"
	COMTAG_API_URI    = "https://aip.baidubce.com/rpc/2.0/nlp/v1/comment_tag"

	// 语音相关API
	TXT2VOICE_API_URI = "http://tsn.baidu.com/text2audio" // 语音合成
	VOICE2TXT_API_URI = "http://vop.baidu.com/server_api" // 语音识别

	// OCR相关API
	IDCARD_API_URI     = "https://aip.baidubce.com/rest/2.0/ocr/v1/idcard"
	BANKCARD_API_URI   = "https://aip.baidubce.com/rest/2.0/ocr/v1/bankcard"
	GENERALOCR_API_URI = "https://aip.baidubce.com/rest/2.0/ocr/v1/general"
	TABLE_API_REQ_URI  = "https://aip.baidubce.com/rest/2.0/solution/v1/form_ocr/request" //表格识别异步端口
	TABLE_API_RESP_URI = "https://aip.baidubce.com/rest/2.0/solution/v1/form_ocr/get_request_result"
	// 人脸检测相关API
	FACEDETECT_API_URI = "https://aip.baidubce.com/rest/2.0/face/v1/detect"
	FACEMATCH_API_URI  = "https://aip.baidubce.com/rest/2.0/faceverify/v1/match"

	// 黄反识别
	ANTIPORN_API_URI = "https://aip.baidubce.com/rest/2.0/antiporn/v1/detect"
)

type accessTokenResponse struct {
	ExpiresIn        int64  `json:"expires_in"`        // 过期时间
	AccessToken      string `json:"access_token"`      // 访问码
	Error            string `json:"error"`             // 错误码
	ErrorDescription string `json:"error_description"` // 错误信息
}

var currentAccessToken = new(accessTokenResponse)

func getToken() *accessTokenResponse {
	postArgs := url.Values{}
	postArgs.Set("client_id",BAIDU_AI_KEY)
	postArgs.Set("client_secret",BAIDU_AI_CRET)
	postArgs.Set("grant_type","client_credentials")

	resp, _ := http.PostForm(TOKEN_API_URI, postArgs)
	defer resp.Body.Close()
	newAccessToken := new(accessTokenResponse)
	if err := json.NewDecoder(resp.Body).Decode(newAccessToken); err != nil {
		panic(err)
	}

	if newAccessToken.Error != "" {
		panic(errors.New(newAccessToken.Error + ": " + newAccessToken.ErrorDescription))
	}

	return newAccessToken
}

func saveToken(token *accessTokenResponse){
	data,err := json.Marshal(token)
	if err !=nil {
		panic("序列化token_response失败")
	}
	if ioutil.WriteFile("token.json",data,0664) == nil {
		log.Println("token写入文件成功！")
	}
}

func main() {

	configFileData, err := ioutil.ReadFile("token.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(configFileData, currentAccessToken)
	if err != nil {
		return
	}
	log.Println("current_token:", currentAccessToken.AccessToken)

	aa := OcrTable("1.jpg")
	fmt.Println(aa)

}

// {"refresh_token":"25.b4347b49afe043a1f523e870df05bb9d.315360000.1871566213.282335-16114169","expires_in":2592000,"session_key":"9mzdDxUrfQTp3AR+9JQvDcZWZ8yyf5zWvluGlx\/2\/x9Mv\/ggDxUMadbCRETNNIevKIVAaBnCEojq5i289TBBy+lojv0iNg==","access_token":"24.4127187fdf9e41c2f7e27541bf20f076.2592000.1558798213.282335-16114169","scope":"public vis-ocr_ocr brain_ocr_scope brain_ocr_general brain_ocr_general_basic brain_ocr_general_enhanced vis-ocr_business_license brain_ocr_webimage brain_all_scope brain_ocr_idcard brain_ocr_driving_license brain_ocr_vehicle_license vis-ocr_plate_number brain_solution brain_ocr_plate_number brain_ocr_accurate brain_ocr_accurate_basic brain_ocr_receipt brain_ocr_business_license brain_solution_iocr brain_ocr_handwriting brain_ocr_passport brain_ocr_vat_invoice brain_numbers brain_ocr_train_ticket brain_ocr_taxi_receipt vis-ocr_\u8f66\u8f86vin\u7801\u8bc6\u522b vis-ocr_\u5b9a\u989d\u53d1\u7968\u8bc6\u522b brain_ocr_vin brain_ocr_quota_invoice wise_adapt lebo_resource_base lightservice_public hetu_basic lightcms_map_poi kaidian_kaidian ApsMisTest_Test\u6743\u9650 vis-classify_flower lpq_\u5f00\u653e cop_helloScope ApsMis_fangdi_permission smartapp_snsapi_base iop_autocar oauth_tp_app smartapp_smart_game_openapi oauth_sessionkey smartapp_swanid_verify smartapp_opensource_openapi smartapp_opensource_recapi","session_secret":"5b662f207c31662ef88387387bbaec6d"}
