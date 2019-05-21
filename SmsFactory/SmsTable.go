package SmsFactory

import (
	"Sms/SmsHandler"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmSms"
)

type SmsTable struct{}

var SMS_MODEL_FACTORY = map[string]interface{}{
}

var SMS_STORAGE_FACTORY = map[string]interface{}{
}

var SMS_RESOURCE_FACTORY = map[string]interface{}{
}

var SMS_FUNCTION_FACTORY = map[string]interface{}{
	"SmsCommonPanicHandle":	SmsHandler.CommonPanicHandle{},
	"SmsSendHandler":		SmsHandler.SmsSendHandler{},
	"SmsVerifyHandler":		SmsHandler.SmsSendHandler{},
}
var SMS_MIDDLEWARE_FACTORY = map[string]interface{}{
}

var SMS_DAEMON_FACTORY = map[string]interface{}{
	"BmMongodbDaemon": BmMongodb.BmMongodb{},
	"BmRedisDaemon":   BmRedis.BmRedis{},
	"BmSmsDaemon":     BmSms.BmSms{},
}

func (t SmsTable) GetModelByName(name string) interface{} {
	return SMS_MODEL_FACTORY[name]
}

func (t SmsTable) GetResourceByName(name string) interface{} {
	return SMS_RESOURCE_FACTORY[name]
}

func (t SmsTable) GetStorageByName(name string) interface{} {
	return SMS_STORAGE_FACTORY[name]
}

func (t SmsTable) GetDaemonByName(name string) interface{} {
	return SMS_DAEMON_FACTORY[name]
}

func (t SmsTable) GetFunctionByName(name string) interface{} {
	return SMS_FUNCTION_FACTORY[name]
}

func (t SmsTable) GetMiddlewareByName(name string) interface{} {
	return SMS_MIDDLEWARE_FACTORY[name]
}
