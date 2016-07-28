package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Unknwon/macaron"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/robfig/cron"
	model "shangqu-shop/shopback"
	"strconv"
	"strings"
	//"shangqu-shop/tool/log"
	"os"
	"time"
)

const (
	//SITE_CODE =3

	AUTH      = "WB00001"
	SITE_CODE = 110
)

var (
	PRICE_CHANNEL      = ""
	PRICE              = ""
	STOCK              = ""
	PRICE_CHANNEL_PROD = "http://10.20.10.100/bms-snk-price/api/v1/pd/pd_priceway"
	PRICE_PROD         = "http://10.20.10.100/bms-snk-price/api/v1/pd/pd_pricecycle"
	STOCK_PROD         = "http://10.20.10.100/bms/api/v1/pd/pd_stock"
	PRICE_CHANNEL_UAT  = "http://10.45.21.248/msk-snk-price/api/v1/pd/pd_priceway"
	PRICE_UAT          = "http://10.45.21.248/msk-snk-price/api/v1/pd/pd_pricecycle"
)

func init() {

	env := os.Getenv("RUN_ENV")

	// 域名变为读取配置文件随配置文件不同而改变
	if env == "prod" {
		PRICE_CHANNEL = PRICE_CHANNEL_PROD
		PRICE = PRICE_PROD
		STOCK = STOCK_PROD
	} else if env == "local" {
		PRICE_CHANNEL = PRICE_CHANNEL_PROD
		PRICE = PRICE_PROD
		STOCK = STOCK_PROD
	} else if env == "fuji" {
		PRICE_CHANNEL = PRICE_CHANNEL_UAT
		PRICE = PRICE_UAT
		STOCK = STOCK_PROD
	} else {
		PRICE_CHANNEL = PRICE_CHANNEL_UAT
		PRICE = PRICE_UAT
		STOCK = STOCK_PROD
	}
	c := cron.New()
	c.AddFunc("00 00 00  1,6,11,16,21,26 * *", func() { SyncProductRankPricess(SITE_CODE) })
	c.AddFunc("00 02 00  1,6,11,16,21,26 * *", func() { GetFlagInfosssss() })

	c.AddFunc("00 00 00  1,6,11,16,21,26 * *", func() { SyncProductInfossssssss(SITE_CODE) })
	c.AddFunc("00 02 00  1,6,11,16,21,26 * *", func() { CheckSyncStockInfosssss() })

	c.AddFunc("@every 10m", func() { SyncFujiStockInfosssss(SITE_CODE) })
	c.AddFunc("@every 11m", func() { CheckSyncStockInfosssss() })
	//c.AddFunc("0 */10 * * * *", func (){SyncProductRankPricess(SITE_CODE)})
	//c.AddFunc("0 */15 * * * *",func(){GetFlagInfosssss()})
	//
	//c.AddFunc("0 */10 * * * *", func (){SyncProductInfossssssss(SITE_CODE)})
	//c.AddFunc("0 */15 * * * *",func (){CheckSyncStockInfosssss()})
	//
	//c.AddFunc("0 */10 * * * *", func(){SyncFujiStockInfosssss(SITE_CODE)})
	//c.AddFunc("0 */15 * * * *",func (){CheckSyncStockInfosssss()})
	c.Start()
}

var flag int64

var stock map[string]int64

type SendRequest struct {
	SiteCode int64  `json:"siteCode"`
	Auth     string `json:"auth"`
}

type PriceRankResponse struct {
	Status string     `json:"status"`
	Msg    string     `json:"msg"`
	Result DataResult `json:"result"`
}

type DataResult struct {
	TotalCount int64           `json:"totalCount"`
	TotalPage  int64           `json:"totalPage"`
	PageNo     int64           `json:"pageNo"`
	List       []SingleProduct `json:"searchList"`
}
type SingleProduct struct {
	ProductId    string    `json:"productId"`
	GradeCode    string    `json:"gradeCode"`
	LogiAreaCode string    `json:"logiAreaCode"`
	List         []WayList `json:"wayList"`
}
type WayList struct {
	OrderLevel string `json:"orderLevel"`
	BoxMin     string `json:"boxMin"`
	Boxmax     string `json:"boxMax"`
}

type QueryRequest struct {
	SiteCode int64 `json:"siteCode"`
}

func SyncProductRankPricess(siteCode int64) {
	fmt.Println("=======================================时间点为=======================================", time.Unix(time.Now().Unix(), 0).String()[0:10])
	log.Println("=======================================开始同步架盘通道================================================")
	//fmt.Println()
	merchantId := "b9cb11e5bea000163e00026578e74468"
	//queryRequest,_ :=ctx.Req.Body().String()
	//siteCodeRequest :=new(QueryRequest)
	//json.Unmarshal([]byte(queryRequest),siteCodeRequest)
	//请求体
	req := new(SendRequest)
	req.SiteCode = siteCode
	req.Auth = AUTH

	//前端返回
	response := new(model.GeneralResponse)

	reqStr, _ := json.Marshal(req)

	body := bytes.NewBuffer(reqStr)

	//请求架盘通道
	resPriceWayData, respPriceWayErr := http.Post(PRICE_CHANNEL, "application/json", body)
	if respPriceWayErr != nil {
		fmt.Println(respPriceWayErr)
	}

	defer resPriceWayData.Body.Close()

	resBody, _ := ioutil.ReadAll(resPriceWayData.Body)

	resp := new(PriceRankResponse)

	json.Unmarshal(resBody, resp)

	//查询架盘
	fmt.Println("请求查询架盘的状态码=========>>>>>>>>>", resp.Status)
	fmt.Println("请求回来的数据总条数=========>>>>>>>>>", resp.Result.TotalCount)
	if resp.Status != "S" {
		response.Code = 10001
		response.Msg = resp.Msg
		ret_str, _ := json.Marshal(response)
		fmt.Println("=======>>>>>>>>", string(ret_str))
	}

	for v, k := range resp.Result.List {
		//str :=strconv.FormatInt(v,64)
		fmt.Println("正在更新商品序列号为=========>>>>>>>>>", k.ProductId, "是第", v, "个商品")
		fmt.Println("时间点为========>>>>>", time.Unix(time.Now().Unix(), 0).String()[0:10])
		go SyncInfo(merchantId, k)
	}

	response.Code = 10000
	response.Msg = "ok"
	ret_str, _ := json.Marshal(response)
	fmt.Println(string(ret_str))
	//return ""
}

func SyncProductRankPrice(ctx *macaron.Context) string {

	merchantId := "b9cb11e5bea000163e00026578e74468"
	queryRequest, _ := ctx.Req.Body().String()
	siteCodeRequest := new(QueryRequest)
	json.Unmarshal([]byte(queryRequest), siteCodeRequest)
	//请求体
	req := new(SendRequest)
	req.SiteCode = siteCodeRequest.SiteCode
	req.Auth = AUTH

	//前端返回
	response := new(model.GeneralResponse)

	reqStr, _ := json.Marshal(req)

	body := bytes.NewBuffer(reqStr)

	//请求架盘通道
	resPriceWayData, respPriceWayErr := http.Post(PRICE_CHANNEL, "application/json", body)
	if respPriceWayErr != nil {
		fmt.Println(respPriceWayErr)
	}

	defer resPriceWayData.Body.Close()

	resBody, _ := ioutil.ReadAll(resPriceWayData.Body)


	resp := new(PriceRankResponse)

	fmt.Println(string(resBody))

	json.Unmarshal(resBody, resp)

	//查询架盘
	fmt.Println("请求查询架盘的状态码=========>>>>>>>>>", resp.Status)
	fmt.Println("请求回来的数据总条数=========>>>>>>>>>", resp.Result.TotalCount)
	if resp.Status != "S" {
		response.Code = 10001
		response.Msg = resp.Msg
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	for v, k := range resp.Result.List {
		//str :=strconv.FormatInt(v,64)
		fmt.Println("正在更新商品序列号为=========>>>>>>>>>", k.ProductId, "是第", v, "个商品")
		go SyncInfo(merchantId, k)
	}

	response.Code = 10000
	response.Msg = "ok"
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
	//return ""
}

func GetFlagInfos(ctx *macaron.Context) string {
	response := new(model.GeneralResponse)
	if flag == 0 {
		response.Code = 10000
		response.Msg = "ok"
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	} else {
		response.Code = 10001
		response.Msg = "wait"
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
}
func GetFlagInfosssss() {
	//response:=new(model.GeneralResponse)
	if flag == 0 {
		log.Println("=======================================同步结束================================================")
	} else {
		log.Println("=======================================同步中,请稍后================================================")
	}
}
func SyncInfo(merchantId string, single SingleProduct) {
	flag = 1
	uniqueInfo, uniqueErr := model.GetProductUniqueInfoBySerialNumber(merchantId, single.ProductId)
	if uniqueErr != nil {
		fmt.Println(uniqueErr.Error())
	}
	fmt.Println(uniqueInfo.Uuid)
	fmt.Println(uniqueInfo.SerialNumber)
	fmt.Println(uniqueInfo.Stock)
	if uniqueErr != nil && !strings.Contains(uniqueErr.Error(), "record not found") {
		flag = 2
	}
	productUniqueId := uniqueInfo.Uuid
	fmt.Println("===========>>>>>>>>")
	fmt.Println(productUniqueId)
	fmt.Println("===========>>>>>>>>")
	priceList, listErr := model.GetPriceListByProductId(merchantId, productUniqueId)
	if listErr != nil && !strings.Contains(listErr.Error(), "record not found") {
		flag = 3
		//return string(ret_str)
	}
	kiloList, kiloErr := model.GetFujiKiloPriceListById(merchantId, productUniqueId)
	if kiloErr != nil && !strings.Contains(listErr.Error(), "record not found") {
		flag = 3
		//return string(ret_str)
	}

	for _, v := range single.List {
		for _, j := range priceList {
			order := strconv.FormatInt(9-j.Order, 10)

			if strings.TrimSpace(order) == strings.TrimSpace(v.OrderLevel) {
				//找到匹配等级
				startInt, _ := strconv.ParseInt(v.BoxMin, 10, 64)
				var endInt int64
				if len(strings.TrimSpace(v.Boxmax)) == 0 {
					endInt = -1
				} else {
					end, _ := strconv.ParseInt(v.Boxmax, 10, 64)

					endInt = end
				}
				//if order =="2"{
				//	j.IsGeneral= "1"
				//	//这里要同步东西过去
				//}else{
				//	j.IsGeneral = "0"
				//}
				j.Start = startInt
				j.End = endInt
				j.SerialNumber = single.ProductId

				errUpdate := model.EditBoxPriceInfo(j)
				if errUpdate != nil && !strings.Contains(errUpdate.Error(), "record not found") {
					flag = 4
				}
			}
		}
		for _, j := range kiloList {
			order := strconv.FormatInt(9-j.Order, 10)

			if strings.TrimSpace(order) == strings.TrimSpace(v.OrderLevel) {
				//找到匹配等级
				startInt, _ := strconv.ParseInt(v.BoxMin, 10, 64)
				var endInt int64
				if len(strings.TrimSpace(v.Boxmax)) == 0 {
					endInt = -1
				} else {
					end, _ := strconv.ParseInt(v.Boxmax, 10, 64)

					endInt = end
				}
				j.Start = startInt
				j.End = endInt
				j.SerialNumber = single.ProductId
				errUpdate := model.EditKiloPriceInfos(j)
				if errUpdate != nil && !strings.Contains(errUpdate.Error(), "record not found") {
					flag = 4
				}
			}
		}
	}
	//count = count +1
	defer func() {
		flag = 0
	}()

}

type PriceRequest struct {
}

type PriceQueryResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []ProductInfo `json:"result"`
}

type ProductInfo struct {
	ProductId    string      `json:"productId"`
	GradeCode    string      `json:"gradeCode"`
	LogiAreaCode string      `json:"logiAreaCode"`
	PricePeriod  string      `json:"pricePeriod"`
	Pricelist    []PriceInfo `json:"priceList"`
}

type PriceInfo struct {
	OrderLevel string  `json:"orderLevel"`
	PriceOfKg  float64 `json:"priceOfKg"`
	PriceOfBox float64 `json:"priceOfBox"`
}

func SyncProductInfo(ctx *macaron.Context) string {
	fmt.Println("======>>>>>>查询架盘", time.Unix(time.Now().Unix(), 0).String()[0:19])
	merchantId := "b9cb11e5bea000163e00026578e74468"
	//queryRequest,_ :=ctx.Req.Body().String()
	//siteCodeRequest :=new(QueryRequest)
	//json.Unmarshal([]byte(queryRequest),siteCodeRequest)
	//请求体
	req := new(SendRequest)
	req.SiteCode = 3
	req.Auth = AUTH

	//前端返回
	response := new(model.GeneralResponse)

	reqStr, _ := json.Marshal(req)

	body := bytes.NewBuffer(reqStr)
	//请求架盘通道
	resPriceWayData, respPriceWayErr := http.Post(PRICE, "application/json", body)
	if respPriceWayErr != nil {
		fmt.Println("errors are .........")
		fmt.Println(respPriceWayErr)
	}

	defer resPriceWayData.Body.Close()
	resBody, _ := ioutil.ReadAll(resPriceWayData.Body)

	priceStruct := new(PriceQueryResponse)

	a := json.Unmarshal(resBody, priceStruct)
	if a != nil {
		fmt.Println(a.Error())
	}
	//fmt.Println("=======>>>>>>>>")
	fmt.Println("请求数据返回状态码======>>>>>>>", priceStruct.Status)
	fmt.Println("请求数据的条数为======>>>>>>", len(priceStruct.Result))
	//fmt.Println("=======>>>>>>>>")
	if strings.TrimSpace(priceStruct.Status) != "S" {
		response.Code = 10001
		response.Msg = priceStruct.Message
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	for v, k := range priceStruct.Result {
		//str :=strconv.FormatInt(v,64)
		fmt.Println("正在查询商品架盘序列号为=========>>>>>>>>>", k.ProductId, "是第", v, "个商品")
		//fmt.Println("======xxxxxx")
		go SyncPrice(merchantId, k)
	}
	response.Code = 10000
	response.Msg = "ok"
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

func SyncProductInfossssssss(siteCode int64) {
	fmt.Println("=======================================时间点为=======================================", time.Unix(time.Now().Unix(), 0).String()[0:10])
	log.Println("=======================================开始同步查询架盘================================================")
	merchantId := "b9cb11e5bea000163e00026578e74468"

	//请求体
	req := new(SendRequest)
	req.SiteCode = siteCode
	req.Auth = AUTH

	//前端返回
	response := new(model.GeneralResponse)

	reqStr, _ := json.Marshal(req)

	body := bytes.NewBuffer(reqStr)
	//请求架盘通道
	resPriceWayData, respPriceWayErr := http.Post(PRICE, "application/json", body)
	if respPriceWayErr != nil {
		fmt.Println(respPriceWayErr)
	}

	defer resPriceWayData.Body.Close()
	resBody, _ := ioutil.ReadAll(resPriceWayData.Body)

	priceStruct := new(PriceQueryResponse)

	a := json.Unmarshal(resBody, priceStruct)
	if a != nil {
		fmt.Println(a.Error())
	}
	//fmt.Println("=======>>>>>>>>")
	fmt.Println("请求数据返回状态码======>>>>>>>", priceStruct.Status)
	fmt.Println("请求数据的条数为======>>>>>>", len(priceStruct.Result))
	//fmt.Println("=======>>>>>>>>")
	if strings.TrimSpace(priceStruct.Status) != "S" {
		response.Code = 10001
		response.Msg = priceStruct.Message
		ret_str, _ := json.Marshal(response)
		fmt.Println(string(ret_str))
	}

	for v, k := range priceStruct.Result {
		//str :=strconv.FormatInt(v,64)
		fmt.Println("正在查询商品架盘序列号为=========>>>>>>>>>", k.ProductId, "是第", v, "个商品")
		fmt.Println("正在查询架盘,时间点为========>>>>>>", time.Unix(time.Now().Unix(), 0).String()[0:10])
		//fmt.Println("======xxxxxx")
		go SyncPrice(merchantId, k)
	}
	response.Code = 10000
	response.Msg = "ok"
	ret_str, _ := json.Marshal(response)
	fmt.Println(string(ret_str))
}
func SyncPrice(merchantId string, p ProductInfo) {
	flag = 1
	//log.Info("aaaaaaaaaaaaa")
	uniqueInfo, uniqueErr := model.GetProductUniqueInfoBySerialNumber(merchantId, p.ProductId)
	if uniqueErr != nil && !strings.Contains(uniqueErr.Error(), "record not found") {
		flag = 5
	}
	productUniqueId := uniqueInfo.Uuid
	fmt.Println("===========>>>>>>>>")
	fmt.Println(productUniqueId)
	fmt.Println("===========>>>>>>>>")
	priceList, listErr := model.GetPriceListByProductId(merchantId, productUniqueId)
	if listErr != nil && !strings.Contains(listErr.Error(), "record not found") {
		flag = 6
	}
	kiloList, kiloErr := model.GetPriceKiloListById(merchantId, productUniqueId)
	if kiloErr != nil && !strings.Contains(kiloErr.Error(), "record not found") {
		flag = 7
	}
	for _, v := range priceList {
		for _, a := range p.Pricelist {
			order := strconv.FormatInt(9-v.Order, 10)

			if strings.TrimSpace(order) == strings.TrimSpace(a.OrderLevel) {
				//str :=strconv.FormatInt(a.PriceOfBox,10)
				str := strconv.FormatFloat(a.PriceOfBox, 'f', 2, 64)
				fmt.Println("=====>>>>>公斤",a.PriceOfBox)
				fmt.Println("=====>>>>>箱子",a.PriceOfBox)

				v.UniquePrice = str
				fmt.Println("转换后的",str)
				errEditBox := model.EditBoxPriceInfo(v)

				if errEditBox != nil && !strings.Contains(errEditBox.Error(), "record not found") {
					flag = 8
				}
				if a.OrderLevel == "7" {
					//这里要更新单品
					fmt.Println("-=================>>>>>>>>> priceListOrder is ====>>", v.Order)
					fmt.Println("=============", a.OrderLevel, "=====>>>>", a.PriceOfBox)
					info, infoErr := model.GetProductUniqueInfos(merchantId, v.SerialNumber)
					if infoErr != nil && !strings.Contains(infoErr.Error(), "record not found") {
						log.Print(infoErr.Error())
					}
					info.BoxPrice = str
					updateSingleErr := model.EditProductUniqueInfos(info)
					if updateSingleErr != nil && !strings.Contains(updateSingleErr.Error(), "record not found") {
						log.Print(updateSingleErr.Error())
					}
				}
			}
		}
	}
	//0 ---8
	for _, v := range kiloList {
		for _, a := range p.Pricelist {
			order := strconv.FormatInt(9-v.Order, 10)
			if strings.TrimSpace(order) == strings.TrimSpace(a.OrderLevel) {
				str := strconv.FormatFloat(a.PriceOfKg, 'f', 2, 64)
				v.UniquePrice = str
				errEditBox := model.EditKiloPriceInfo(v)
				if errEditBox != nil && !strings.Contains(errEditBox.Error(), "record not found") {
					flag = 9
				}
				if a.OrderLevel == "7" {
					//这里要更新单品
					fmt.Println("price order is .....=====>>>>>>", order)
					fmt.Println("-=================>>>>>>>>> price Kilo ListOrder is ====>>", v.Order)
					fmt.Println("=============", a.OrderLevel, "=====>>>>", a.PriceOfKg)
					info, infoErr := model.GetProductUniqueInfos(merchantId, v.SerialNumber)
					if infoErr != nil && !strings.Contains(infoErr.Error(), "record not found") {
						log.Print(infoErr.Error())
					}
					info.KiloPrice = str
					updateSingleErr := model.EditProductUniqueInfos(info)
					if updateSingleErr != nil && !strings.Contains(updateSingleErr.Error(), "record not found") {
						log.Print(updateSingleErr.Error())
					}
				}
			}
		}
	}
	defer func() {
		flag = 0
	}()

}

type FujiStockRequest struct {
	SiteCode   int64  `json:"siteCode"`
	Auth       string `json:"auth"`
	LoginId    string `json:"loginID"`
	Parameters Params `json:"param"`
}
type Params struct {
	SellerCode   string `json:"sellerCode"`
	PlatformType string `json:"platformType"`
	DistrictCode int64  `json:"districtCode"`
	SellerType   int64  `json:"sellerType"`
	//PdCode string `json:"pdCode"`
}
type FujiProductStockResponse struct {
	Status     string       `json:"status"`
	Message    string       `json:"message"`
	ResultInfo ResultStruct `json:"result"`
}

type ResultStruct struct {
	SellerCode   string     `json:"sellerCode"`
	DistrictCode int64      `json:"districtCode"`
	ProductList  []Products `json:"products"`
}
type Products struct {
	PdCode   string `json:"pdCode"`
	StockCnt int64  `json:"stockCnt"`
}

var isOk int64

func SyncFujiStockInfo(ctx *macaron.Context) string {

	log.Println("=======================================开始同步库存================================================")
	isOk = 1
	merchantId := "b9cb11e5bea000163e00026578e74468"
	req := new(FujiStockRequest)

	bodys, _ := ctx.Req.Body().String()
	fmt.Println(bodys)

	json.Unmarshal([]byte(bodys), req)
	////fmt.Println(string(bodys))
	//reqByte ,_:=json.Marshal(req)
	rs := new(FujiStockRequest)
	rs.Auth = req.Auth
	rs.LoginId = req.LoginId
	rs.SiteCode = req.SiteCode
	p := new(Params)
	//p.DistrictCode = req.Parameters.DistrictCode
	//p.PlatformType = req.Parameters.PlatformType
	//p.SellerCode = req.Parameters.SellerCode
	//p.SellerType = req.Parameters.SellerType
	p.DistrictCode = 41
	p.PlatformType = "1"
	p.SellerCode = "0000000"
	p.SellerType = 1
	rs.Parameters = *p
	psb, _ := json.Marshal(rs)
	body := bytes.NewBuffer(psb)

	response := new(model.GeneralResponse)
	fmt.Println(string(bodys))
	fmt.Println(STOCK)
	fujiStockInfo, fujiStockErr := http.Post(STOCK, "application/json", body)
	if fujiStockErr != nil {
		fmt.Println(fujiStockErr.Error())
	}
	if fujiStockErr != nil {
		response.Code = 10001
		response.Msg = fujiStockErr.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}
	defer fujiStockInfo.Body.Close()
	resBody, _ := ioutil.ReadAll(fujiStockInfo.Body)

	res := new(FujiProductStockResponse)
	errUm := json.Unmarshal(resBody, res)
	fmt.Println(string(resBody))
	fmt.Println("同步商品库存,获取状态码为======>>>>>>", res.Status)
	if errUm != nil {
		response.Code = 10001
		response.Msg = errUm.Error()
		ret_str, _ := json.Marshal(response)
		return string(ret_str)
	}

	for v, k := range res.ResultInfo.ProductList {
		//str :=strconv.FormatInt(v,64)
		//fmt.Println("拉取第"+str+"条数据,商品序列号为======>>>>>>",k.PdCode)
		fmt.Println("正在拉取商品序列号为=========>>>>>>>>>", k.PdCode, "是第", v, "个商品")
		go UpdateProductStockInfo(merchantId, k)
	}

	defer func() {
		isOk = -1
	}()
	response.Code = 10000
	response.Msg = "ok"
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
}

func SyncFujiStockInfosssss(siteCode int64) {

	log.Println("=======================================开始同步库存================================================")
	isOk = 1
	merchantId := "b9cb11e5bea000163e00026578e74468"
	req := new(FujiStockRequest)

	//bodys,_:=ctx.Req.Body().String()
	//fmt.Println(bodys)

	//json.Unmarshal([]byte(bodys),req)
	////fmt.Println(string(bodys))
	//reqByte ,_:=json.Marshal(req)
	rs := new(FujiStockRequest)
	rs.Auth = req.Auth
	rs.LoginId = "wb01"
	rs.SiteCode = SITE_CODE
	p := new(Params)
	p.DistrictCode = req.Parameters.DistrictCode
	p.PlatformType = req.Parameters.PlatformType
	p.SellerCode = req.Parameters.SellerCode
	p.SellerType = req.Parameters.SellerType
	rs.Parameters = *p
	psb, _ := json.Marshal(rs)
	body := bytes.NewBuffer(psb)

	response := new(model.GeneralResponse)
	//fmt.Println(string(bodys))
	fujiStockInfo, fujiStockErr := http.Post(STOCK, "application/json", body)
	if fujiStockErr != nil {
		response.Code = 10001
		response.Msg = fujiStockErr.Error()
		ret_str, _ := json.Marshal(response)
		fmt.Println(string(ret_str))
	}
	defer fujiStockInfo.Body.Close()
	resBody, _ := ioutil.ReadAll(fujiStockInfo.Body)

	res := new(FujiProductStockResponse)
	errUm := json.Unmarshal(resBody, res)
	fmt.Println(string(resBody))
	fmt.Println("同步商品库存,获取状态码为======>>>>>>", res.Status)
	if errUm != nil {
		response.Code = 10001
		response.Msg = errUm.Error()
		ret_str, _ := json.Marshal(response)
		fmt.Println(ret_str)
	}

	for v, k := range res.ResultInfo.ProductList {
		//str :=strconv.FormatInt(v,64)
		//fmt.Println("拉取第"+str+"条数据,商品序列号为======>>>>>>",k.PdCode)
		fmt.Println("正在拉取商品序列号为=========>>>>>>>>>", k.PdCode, "是第", v, "个商品")
		go UpdateProductStockInfo(merchantId, k)
	}

	defer func() {
		isOk = -1
	}()
	response.Code = 10000
	response.Msg = "ok"
	ret_str, _ := json.Marshal(response)
	fmt.Println(ret_str)
}

func CheckSyncStockInfo(ctx *macaron.Context) string {

	merchantId := "b9cb11e5bea000163e00026578e74468"
	response := new(model.GeneralResponse)
	if isOk == -1 {
		idList, errId := model.GetProductIdListByMerchantIdI(merchantId)
		if errId != nil && !strings.Contains(errId.Error(), "record not found") {
			response.Code = 10001
			response.Msg = errId.Error()
			ret_str, _ := json.Marshal(response)
			return string(ret_str)
		}
		for _, k := range idList {
			go UpdateTotalStock(merchantId, k)
		}
		response.Code = 10001
		response.Msg = "处理中...."
		ret_str, _ := json.Marshal(response)
		log.Println("=======================================同步中,请稍后================================================")
		return string(ret_str)

	} else if isOk == -2 {
		response.Code = 10000
		response.Msg = "ok"
		ret_str, _ := json.Marshal(response)
		log.Println("=======================================同步ok===============================================")
		return string(ret_str)
	}
	response.Code = 10000
	response.Msg = "ok"
	ret_str, _ := json.Marshal(response)
	return string(ret_str)
	//fmt.Println(string(ret_str))
	//log.Println("=======================================同步OKssss================================================")
}

func CheckSyncStockInfosssss() {

	merchantId := "b9cb11e5bea000163e00026578e74468"
	response := new(model.GeneralResponse)
	if isOk == -1 {
		idList, errId := model.GetProductIdListByMerchantIdI(merchantId)
		if errId != nil && !strings.Contains(errId.Error(), "record not found") {
			response.Code = 10001
			response.Msg = errId.Error()
			ret_str, _ := json.Marshal(response)
			fmt.Println(ret_str)
		}
		for _, k := range idList {
			go UpdateTotalStock(merchantId, k)
		}
		response.Code = 10001
		response.Msg = "处理中...."
		log.Println("========================================处理中========================================")
		ret_str, _ := json.Marshal(response)
		fmt.Println(ret_str)

	} else if isOk == -2 {
		response.Code = 10000
		response.Msg = "ok"
		ret_str, _ := json.Marshal(response)
		log.Println("========================================处理完成========================================")
		fmt.Println(string(ret_str))
	}
	response.Code = 10000
	response.Msg = "ok"
	ret_str, _ := json.Marshal(response)
	//return string(ret_str)
	log.Println("========================================处理完成========================================")
	fmt.Println(ret_str)
}

func UpdateProductStockInfo(merchantId string, product Products) {
	flag = 1
	log.Println("更新库存============>>>>>>>", product.PdCode)
	uniqueInfo, uniqueErr := model.GetProductUniqueInfoBySerialNumber(merchantId, product.PdCode)
	if uniqueErr != nil && !strings.Contains(uniqueErr.Error(), "record not found") {
		flag = 5
	}
	uniqueInfo.Stock = product.StockCnt
	err := model.EditProductUniqueProductInfo(uniqueInfo)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		flag = 10
	}
	defer func() {
		flag = 0
	}()
}

func UpdateTotalStock(merchantId, productId string) {
	singleList, singleListErr := model.GetProductDataDetailByProductId(productId)

	if singleListErr != nil && !strings.Contains(singleListErr.Error(), "record not found") {
		fmt.Println(singleListErr.Error())
	}
	var stock int64
	for _, v := range singleList {
		fmt.Println("stockInfo =====>>>", stock)
		fmt.Println("商品ID为=====》》》》》》", v.Uuid)
		fmt.Println("单品库存品库存为==========>>>>>>>>>", v.Stock)
		stock = stock + v.Stock

		fmt.Println("商品库存为==========>>>>>>>>>>>", stock)
	}
	fmt.Println("正在更新MS_PRODUCT_BASE_INFO 的库存,商品ID为 ======>>>>>>>", productId, "库存数量为======>>>>>>>", stock)
	editStockErr := model.EditProductStocks(productId, merchantId, stock)
	if editStockErr != nil && !strings.Contains(editStockErr.Error(), "record not found") {
		fmt.Println(editStockErr)
	}
	defer func() {
		isOk = -2
	}()
}

//type FujiProductStockResponse struct{
//	Status string `json:"status"`
//	Message string `json:"message"`
//	ResultInfo ResultStruct `json:"result"`
//}
type StockResponse struct {
	Code  int64  `json:"code"`
	Msg   string `json:"msg"`
	Stock int64  `json:"stock"`
}
type UpdateUniqueRequest struct {
	SiteCode int64       `json:"siteCode"`
	Auth     string      `json:"auth"`
	LoginId  string      `json:"loginId"`
	Pa       UniqueParam `json:"param"`
}
type UniqueParam struct {
	SellerCode   string `json:"sellerCode"`
	PlatformType string `json:"platformType"`
	DistrictCode int64  `json:"districtCode"`
	SellerType   int64  `json:"sellerType"`
	PdCode       string `json:"pdCode"`
}

//func UpdateUniqueProductStockInfo(ctx *macaron.Context)string{
//
//	merchantId:="b9cb11e5bea000163e00026578e74468"
//
//	body,_:=ctx.Req.Body().String()
//	req :=new(UpdateUniqueRequest)
//
//	reqs :=new(UniqueParam)
//	json.Unmarshal([]byte(body),req)
//
//	pam :=new(UniqueParam)
//
//	pam.DistrictCode = reqs.DistrictCode
//	pam.PdCode = reqs.PdCode
//	pam.PlatformType = reqs.PlatformType
//	pam.SellerCode = reqs.SellerCode
//	pam.SellerType = reqs.SellerType
//	upRequest :=new(UpdateUniqueRequest)
//	upRequest.SiteCode = SITE_CODE
//	upRequest.Auth = AUTH
//	upRequest.LoginId = "wb01"
//	upRequest.Pa = *pam
//	psb ,_:=json.Marshal(upRequest)
//	bodys :=bytes.NewBuffer(psb)
//
//	responses :=new(StockResponse)
//	//response:=new(model.GeneralResponse)
//	//fmt.Println(string(bodys))
//	fujiStockInfo ,fujiStockErr :=http.Post(STOCK,"application/json",bodys)
//	if fujiStockErr!=nil{
//		responses.Code = 10001
//		responses.Msg = fujiStockErr.Error()
//		ret_str,_:=json.Marshal(responses)
//		return string(ret_str)
//	}
//	defer fujiStockInfo.Body.Close()
//	resBody ,_:=ioutil.ReadAll(fujiStockInfo.Body)
//	res :=new(FujiProductStockResponse)
//	errUm :=json.Unmarshal(resBody,res)
//	fmt.Println(string(resBody))
//	fmt.Println("同步商品库存,获取状态码为======>>>>>>",res.Status)
//	if errUm!=nil{
//		responses.Code = 10001
//		responses.Msg = errUm.Error()
//		ret_str,_:=json.Marshal(responses)
//		return string(ret_str)
//	}
//	//var flag string
//	var stock int64
//	for v, k :=range res.ResultInfo.ProductList{
//		//str :=strconv.FormatInt(v,64)
//		//fmt.Println("拉取第"+str+"条数据,商品序列号为======>>>>>>",k.PdCode)
//		fmt.Println("正在拉取商品序列号为=========>>>>>>>>>",k.PdCode,"是第",v,"个商品")
//		stock = k.StockCnt
//		str :=UpdateUniqueStockInfo(merchantId,k)
//		if strings.Contains(str,"10001"){
//			return str
//		}
//	}
//
//	responses.Code= 10000
//	responses.Msg = "ok"
//	responses.Stock =stock
//	ret_str,_:=json.Marshal(responses)
//	return string(ret_str)
//
//}
//
//
//func UpdateUniqueStockInfo(merchantId string,product Products)string{
//	response :=new(StockResponse)
//	uniqueInfo ,uniqueErr :=model.GetProductUniqueInfoBySerialNumber(merchantId,product.PdCode)
//	if uniqueErr!=nil&&!strings.Contains(uniqueErr.Error(),"record not found"){
//		response.Code = 10001
//		response.Msg = uniqueErr.Error()
//		ret_str,_:=json.Marshal(response)
//		return string(ret_str)
//	}
//	uniqueInfo.Stock = product.StockCnt
//	err :=model.EditProductUniqueProductInfo(uniqueInfo)
//	if err!=nil&&!strings.Contains(err.Error(),"record not found"){
//		response.Code = 10001
//		response.Msg = uniqueErr.Error()
//		ret_str,_:=json.Marshal(response)
//		return string(ret_str)
//	}
//	response.Code = 10000
//	response.Msg = "ok"
//	ret_str,_:=json.Marshal(response)
//	return string(ret_str)
//}