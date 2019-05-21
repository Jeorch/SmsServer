package SmsHandler

import (
	"Sms/util"
	"encoding/json"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmSms"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"
)

type SmsSendHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
	rd         *BmRedis.BmRedis
	sm		   *BmSms.BmSms
}

type Sms struct {
	Phone	string	`json:"phone"`
	Code	string	`json:"code"`
}

type SmsResponse struct {
	Message		string 	`json:"Message"`
	RequestID	string	`json:"RequestId"`
	BizID		string	`json:"BizId"`
	Code		string	`json:"Code"`
}

func (h SmsSendHandler) NewSmsHandler(args ...interface{}) SmsSendHandler {
	var m *BmMongodb.BmMongodb
	var r *BmRedis.BmRedis
	var sms *BmSms.BmSms
	var hm string
	var md string
	var ag []string
	for i, arg := range args {
		if i == 0 {
			sts := arg.([]BmDaemons.BmDaemon)
			for _, dm := range sts {
				tp := reflect.ValueOf(dm).Interface()
				tm := reflect.ValueOf(tp).Elem().Type()
				if tm.Name() == "BmMongodb" {
					m = dm.(*BmMongodb.BmMongodb)
				} else if tm.Name() == "BmRedis" {
					r = dm.(*BmRedis.BmRedis)
				} else if tm.Name() == "BmSms" {
					sms = dm.(*BmSms.BmSms)
				}
			}
		} else if i == 1 {
			md = arg.(string)
		} else if i == 2 {
			hm = arg.(string)
		} else if i == 3 {
			lst := arg.([]string)
			for _, str := range lst {
				ag = append(ag, str)
			}
		} else {
		}
	}

	return SmsSendHandler{Method: md, HttpMethod: hm, Args: ag, db: m, rd: r, sm: sms }
}

func (h SmsSendHandler) SendSms(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	w.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)

	response := map[string]interface{}{}

	if err != nil {
		log.Printf("解析Body出错：%v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}

	sms := Sms{}
	err = json.Unmarshal(body, &sms)

	if err != nil {
		log.Printf("解析Json出错：%v", err)
		http.Error(w, "can't convert Sms struct", http.StatusBadRequest)
		return 1
	}

	code := util.SixRandomNumberByPhone()
	err, result := h.sm.SendMsg(sms.Phone, code)

	if err != nil {
		log.Printf("Error SendMsg %v", err)
		response["status"] = "error"
		response["msg"] = "短信发送失败！"
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	}

	aliResponse := SmsResponse{}
	err = json.Unmarshal(result.GetHttpContentBytes(), &aliResponse)

	if aliResponse.Code != "OK" || aliResponse.Message != "OK" {
		response["status"] = "error"
		response["msg"] = "短信发送失败！"
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	}
	response["status"] = "success"
	response["msg"] = "短信发送成功！"
	err = h.rd.PushPhoneCode(sms.Phone, code, time.Minute * 5)
	if err != nil {
		panic(err.Error())
	}
	enc := json.NewEncoder(w)
	enc.Encode(response)
	return 0
}

func (h SmsSendHandler) VerifySmsCode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	w.Header().Add("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	response := map[string]interface{}{}

	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}

	sms := Sms{}

	err = json.Unmarshal(body, &sms)

	if err != nil {
		log.Printf("解析Json出错：%v", err)
		http.Error(w, "can't convert Sms struct", http.StatusBadRequest)
		return 1
	}

	if len(sms.Code) == 0 {
		response["status"] = "error"
		response["msg"] = "验证码为空"
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	}

	code, err := h.rd.GetPhoneCode(sms.Phone)

	if err != nil && err.Error() == "phoneCode expired" {
		response["status"] = "error"
		response["msg"] = "验证码已过期"
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	}

	if sms.Code != code {
		response["status"] = "error"
		response["msg"] = "输入的验证码不一致"
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	}
	response["status"] = "success"
	response["msg"] = "验证成功"

	h.removeCodeByPhone(sms.Phone)
	enc := json.NewEncoder(w)
	enc.Encode(response)

	return 0
}

func (h SmsSendHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h SmsSendHandler) GetHandlerMethod() string {
	return h.Method
}

func (h SmsSendHandler) removeCodeByPhone(phone string) {
	client := h.rd.GetRedisClient()
	defer client.Close()

	client.Del(phone)
}
