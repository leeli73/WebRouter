package main

import(
	"log"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/vulcand/oxy/forward"
  	"github.com/vulcand/oxy/testutils"
)

type Rule struct{
	Domain string `json:Domain`
	Address string `json:Address`
}
type Setting struct{
	MainPort string `json:MainPort`
	Rules []Rule `json:Rules`
}

var setting Setting
var fwd *forward.Forwarder

func redirectHandle(w http.ResponseWriter, r *http.Request){
	log.Println(r.Host)
	address := getAddress(r.Host)
	if address != "" {
		r.URL = testutils.ParseURI("http://"+address)
		log.Println("Info: "+r.Host+" => "+address)
		fwd.ServeHTTP(w, r)
	} else {
		w.WriteHeader(500)
	}
}

func getAddress(host string) string {
	for _,v := range setting.Rules {
		if strings.Compare(host,v.Domain) == 0 {
			return v.Address
		}
	}
	return ""
}

func Init(){
	var Data,err = ioutil.ReadFile("setting.json")
	if err != nil{
		log.Fatal("Read Config File Error！")
		return
	}
	err = json.Unmarshal(Data,&setting)
	if err != nil{
		log.Fatal("Read Config JSON Error！Please Check!")
		return
	}
	log.Println("Main Port: "+setting.MainPort)
	for i:=0;i<len(setting.Rules);i++{
		log.Println("Import Rule: "+setting.Rules[i].Domain+" <----> "+setting.Rules[i].Address)
	}
}

func main(){
	Init()
	fwd, _ = forward.New()
	redirect := http.HandlerFunc(redirectHandle)
	s := &http.Server{
		Addr:           ":"+setting.MainPort,
		Handler:        redirect,
	}
	log.Println("Info: Listening port "+s.Addr)
	s.ListenAndServe()
}