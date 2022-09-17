package calcoper
//计算类操作都放入此类
import (
	"fmt"
	"strconv"
	"strings"
	"time"
)
//节假日列表
var holidays =[...]string{ "20220101", "20220102","20220103","20220131",
	"20220201",	"20220202",	"20220203",	"20220204",	"20220205",	"20220206",
	"20220403",	"20220404",	"20220503",	"20220502",	"20220501",	"20220430",
	"20220405",	"20220504",	"20220912",	"20220605",	"20220604",	"20220603",
	"20220911",	"20220910",	"20221005",	"20221004",	"20221003",	"20221002",
	"20221001",	"20221007",	"20221006"}

type Futuresdetail struct {
	Code     string //代码
	Price    float64 //最新报价
	Isbase   bool  //是否当月
	Tradeday time.Time //交割日期
}

//根据获取的接口数据计算转换成Futuresdetail结构数据
func Calcstruct(detail []string)(fd Futuresdetail){
	if len(detail)<20 {
		fmt.Println("calcstruct时字符串数组长度不对，剔除")
		fd.Code ="0"
		return fd
	}
	index:=strings.Index(detail[0],"nf_IC")
	index2:=strings.Index(detail[0],"=")
	if index>0 && index2 >0 {
		if detail[0][index+3:index2]=="IC0"{
			fmt.Println("剔除IC0数据")
			fd.Code ="0"
			return fd
		}
		//fmt.Printf("代码: %s 最新价: %s 今开: %s 最高: %s 最低: %s 昨结算: %s\n",
		//	detail[0][index+3:index2],detail[3],detail[0][index2+2:],detail[1],detail[2],detail[14])
		fd.Code =detail[0][index+3:index2] //代码比如IC2205
		fd.Price,_=strconv.ParseFloat(detail[3],32) //最新价
		fd.Tradeday=getWeek(fd.Code)  //获取交割日期
		fd.Isbase=judgebase(fd.Code)  //是否当月
	}else {
		fmt.Println("calcstruct获取code错误")
		fd.Code ="0"
		return fd
	}
	return fd
}

func getWeek(yearmon string)(pday time.Time){
	k:=0
	//从1号到31号循环判断第3个周五
	for i:=0;i<32;i++{
		tmpstr:=fmt.Sprintf("20%s%02d",yearmon[2:],i)
		pday,_:=time.Parse("20060102",tmpstr)
		if pday.Weekday()==5 {
			k=k+1
			//判断是否第三个周五
			if k==3{
				for judgeholiday(tmpstr){
					//pday=pday+86400
					pday=pday.Add(time.Hour*24)
					tmpstr=pday.Format("20060102")
				}
				//fmt.Println(tmpstr)
				return pday
			}
		}
	}
	//tm,_:=time.Parse("20060102","20"+yearmon[2:]+"01")
	return pday
}
//判断是否为节假日，节假日顺延
func judgeholiday(pday string)(isholiday bool){
	for _,tempstr:=range holidays{
		if strings.Compare(tempstr,pday)==0{
			return true
		}else{
			tm,_:=time.Parse("20060102",pday)
			//判断是否周六周日
			if tm.Weekday()==6 || tm.Weekday()==0{
				return true
			}
		}
	}
	return false
}
//判断是否当月
func judgebase(yearmon string)(isbase bool){
	tm:=time.Now()
	res:=tm.Format("200601")
	if strings.Compare("20"+yearmon[2:],res)==0{
		return true
	}
	return false
}
//计算收益
func CalcRevenue(f0,f1 Futuresdetail)(daynum int ,revenue float64){
//	fmt.Println("f0: ",f0)
//	fmt.Println("f1: ",f1)
	daynum=int(f1.Tradeday.Sub(f0.Tradeday).Hours()/24)
	if daynum>0{
		revenue=(f0.Price-f1.Price)*200/float64(daynum)
	}else{
		revenue=0
	}
	return daynum,revenue
}

func FindBase(details []Futuresdetail)(detail Futuresdetail,err int){
	for _,detail :=range details{
		if detail.Isbase {
			return detail,0
		}
	}
	return detail,-1
}
