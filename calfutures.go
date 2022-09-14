package main

import (
	"calfutures/getsinadata"
	"calfutures/calcoper"
	"encoding/json"
	"fmt"
	"strings"
)

func expcode(resp []byte)(resstr string){
	var codes []getsinadata.NameList
	err:=json.Unmarshal(resp,&codes)
	if err!=nil{
		fmt.Println("解析json数据错误")
	}
	for _,code:=range codes{
		resstr=resstr+"nf_"+code.Symbol+","
	}
//	fmt.Println(resstr)
	return resstr
}
func main() {
	resp,err:=getsinadata.GetCode()
	if err!=nil{return}
	fmt.Println(string(resp))
	resstr:=expcode(resp) //获取估值期货代码IC0,IC2205,IC2206,IC2212,
	resdetail,err:=getsinadata.GetCodeValue(resstr)
	strarr:=strings.Split(resdetail,";")
	var details []calcoper.Futuresdetail
	fmt.Println("正在对详细接口信息进行解析...")
	for _,strtmp:=range strarr{
		detail:=strings.Split(strtmp,",")
		fsd:=calcoper.Calcstruct(detail)
		if fsd.Code !="0" && fsd.Code !="IC0"{
			details=append(details, fsd)
		}
	}
	fmt.Println("接口获取的详细数据: \n",details)
	basedetail,errint:=calcoper.FindBase(details)
	if errint==-1{
		fmt.Println("未找到当月期指")
		return
	}
	fmt.Println("计算结果:")
	for _,detail:=range details{
		if !detail.Isbase{
			daynum,revenue:=calcoper.CalcRevenue(basedetail,detail)
			fmt.Printf("期指:%s  最新报价:%.02f 交割日期:%s 间隔天数:%d 每日收益%.02f\n",
				detail.Code,detail.Price,detail.Tradeday.Format("20060102"),daynum,revenue)
		}else{
			fmt.Printf("期指当前月:%s  最新报价:%.02f 交割日期:%s\n",detail.Code,detail.Price,detail.Tradeday.Format("20060102"))
		}
	}
}
