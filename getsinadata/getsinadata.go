package getsinadata

import (
	"fmt"
	"io/ioutil"
	"net/http"
)
type NameList struct {
	Symbol string `json:"symbol"`
}
//获取期指的代码
func GetCode()(body []byte,err error){
	url:="https://vip.stock.finance.sina.com.cn/quotes_service" +
		"/api/json_v2.php/Market_Center.getNameList?page=1&num=40" +
		"&sort=symbol&asc=1&node=zzgz_qh&_s_r_a=init"
	fmt.Println("正在从接口获取期指代码信息...")
	resp,err:=http.Get(url)
	if err!=nil{
		fmt.Println("无法获取期指代码:",err.Error())
		return body,err
	}
	defer resp.Body.Close()
	body,err=ioutil.ReadAll(resp.Body)
	if err!=nil{
		fmt.Println("无法获取期指代码：",err.Error())
		return body,err
	}
	return body,err
}
//根据代码获取期指的详细信息
func GetCodeValue(resstr string)(resdetail string,err error){
	url:="https://hq.sinajs.cn/rn=udcjh&list="+resstr
	fmt.Println("正在从接口获取期指纤细信息...")
	client:=&http.Client{}
	request,err:=http.NewRequest("GET",url,nil)
	request.Header.Add("Referer","https://vip.stock.finance.sina.com.cn")
	request.Header.Add("Host","hq.sinajs.cn")
	if err!=nil{
		fmt.Println("无法获取期指详细信息:",err.Error())
	}
	resp,_:=client.Do(request)
	body,err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		fmt.Println("无法获取期指详细信息：",err.Error())
	}
	resdetail=string(body)
	return resdetail,err
}
