package main

import (
	"calfutures/getsinadata"
	"calfutures/calcoper"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"net/http"
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
func calfutures(c *gin.Context) {
	resp,err:=getsinadata.GetCode()
	if err!=nil{return}
	//fmt.Println(string(resp))
	resstr:=expcode(resp) //获取估值期货代码IC0,IC2205,IC2206,IC2212,
	resdetail,err:=getsinadata.GetCodeValue(resstr)
	strarr:=strings.Split(resdetail,";")
	var details []calcoper.Futuresdetail
	fmt.Println("正在对详细接口信息进行解析...")
	for _,strtmp:=range strarr{
		detail:=strings.Split(strtmp,",")
		fsd:=calcoper.Calcstruct(detail)
		if fsd.Code !="0"{
			details=append(details, fsd)
		}
	}
//	fmt.Println("接口获取的详细数据: \n",details)
	basedetail,errint:=calcoper.FindBase(details)
	if errint==-1{
		fmt.Println("未找到当月期指")
		return
	}
	//fmt.Println("计算结果:")
	c.HTML(http.StatusOK,"head.html",gin.H{})
	for _,detail:=range details{
		if !detail.Isbase{
			daynum,revenue:=calcoper.CalcRevenue(basedetail,detail)
			contentstr:=fmt.Sprintf("期指: %s  最新报价： %.02f  间隔天数：%d 每日收益%.02f\n",
				detail.Code,detail.Price,daynum,revenue)
			c.HTML(http.StatusOK,"content.html",gin.H{"content":contentstr})
		}else{
			contentstr:=fmt.Sprintf("期指当前月: %s  最新报价： %.02f\n",detail.Code,detail.Price)
			c.HTML(http.StatusOK,"content.html",gin.H{"content":contentstr})
		}
	}
	c.HTML(http.StatusOK,"tail.html",gin.H{})
}

func main() {
	route:=gin.Default()
	gin.SetMode(gin.ReleaseMode)
	route.LoadHTMLGlob("./html/*")
	route.GET("/calfutures",calfutures)
	route.Run(":18838")
}
