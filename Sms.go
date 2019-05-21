package main

import (
	"Sms/SmsFactory"
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmApiResolver"
	"github.com/alfredyang1986/BmServiceDef/BmConfig"
	"github.com/alfredyang1986/BmServiceDef/BmPodsDefine"
	"github.com/julienschmidt/httprouter"
	"github.com/manyminds/api2go"
	"net/http"
	"os"
)

func main() {
	version := "v0"
	prodEnv := "SMS_HOME"
	fmt.Println("Sms Server begins, version =", version)

	fac := SmsFactory.SmsTable{}
	var pod = BmPodsDefine.Pod{Name: "new SMS", Factory: fac}
	envHome := os.Getenv(prodEnv)
	pod.RegisterSerFromYAML(envHome + "/resource/service-def.yaml")

	var bmRouter BmConfig.BmRouterConfig
	bmRouter.GenerateConfig(prodEnv)

	addr := bmRouter.Host + ":" + bmRouter.Port
	fmt.Println("Listening on ", addr)
	api := api2go.NewAPIWithResolver(version, &BmApiResolver.RequestURL{Addr: addr})
	pod.RegisterAllResource(api)

	pod.RegisterAllFunctions(version, api)
	pod.RegisterAllMiddleware(api)

	handler := api.Handler().(*httprouter.Router)
	pod.RegisterPanicHandler(handler)
	http.ListenAndServe(":"+bmRouter.Port, handler)

	fmt.Println("SMS Server ends, version =", version)
}
