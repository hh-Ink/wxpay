package wxpay

import (
	"encoding/xml"
	"fmt"
)

type UnifyOrder struct {
	XMLName          xml.Name `xml:"xml"`
	Appid            string   `xml:"appid"`
	Body             string   `xml:"body"`
	Mch_id           string   `xml:"mch_id"`
	Nonce_str        string   `xml:"nonce_str"`
	Notify_url       string   `xml:"notify_url"`
	Trade_type       string   `xml:"trade_type"`
	Spbill_create_ip string   `xml:"spbill_create_ip"`
	Total_fee        int      `xml:"total_fee"`
	Out_trade_no     string   `xml:"out_trade_no"`
	Sign             string   `xml:"sign"`
	Device_info      string   `xml:"device_info"`
	Sign_type        string   `xml:"sign_type"`
	Detail           string   `xml:"detail"`
	Attach           string   `xml:"attach"`
	Fee_type         string   `xml:"fee_type"`
	Time_start       string   `xml:"time_start"`
	Time_expire      string   `xml:"time_expire"`
	Goods_tag        string   `xml:"goods_tag"`
	Limit_pay        string   `xml:"limit_pay"`
	Scene_info       string   `xml:"scene_info"`
	Openid           string   `xml:"openid"`
	Product_id       string   `xml:"product_id"`
}

/**
 * 构造方法
 */
func NewUor(appid string, mch_id string) (unifyOrder *UnifyOrder) {
	unifyOrder = &UnifyOrder{}
	unifyOrder.Appid = appid
	unifyOrder.Mch_id = mch_id
	return
}

/**
 * 生成签名
 */
func (unifyOrder *UnifyOrder) MakeSign(key string) {
	unifyOrder.Nonce_str = MakeNoter()
	unifyOrder.Sign = WxpayCalcSign(unifyOrder, key)
}

/**
 * 发送请求
 */
func (unifyOrder *UnifyOrder) Post() (map[string]string, error) {
	result := make(map[string]string)
	// 检测必填参数
	if len(unifyOrder.Out_trade_no) == 0 {
		return result, fmt.Errorf("%s", "缺少统一支付接口必填参数out_trade_no！")
	} else if len(unifyOrder.Body) == 0 {
		return result, fmt.Errorf("%s", "缺少统一支付接口必填参数body！")
	} else if unifyOrder.Total_fee == 0 {
		return result, fmt.Errorf("%s", "缺少统一支付接口必填参数total_fee！")
	} else if len(unifyOrder.Trade_type) == 0 {
		return result, fmt.Errorf("%s", "缺少统一支付接口必填参数trade_type！")
	}

	//关联参数
	if unifyOrder.Trade_type == "JSAPI" && len(unifyOrder.Openid) == 0 {
		return result, fmt.Errorf("%s", "统一支付接口中，缺少必填参数openid！trade_type为JSAPI时，openid为必填参数！")
	}
	if unifyOrder.Trade_type == "NATIVE" && len(unifyOrder.Product_id) == 0 {
		return result, fmt.Errorf("%s", "统一支付接口中，缺少必填参数product_id！trade_type为JSAPI时，product_id为必填参数！")
	}
	postData, _err := XmlEncode(unifyOrder)
	if _err != nil {
		return result, _err
	}
	url := "https://api.mch.weixin.qq.com/pay/unifiedorder"
	r, err := Post(postData, url)
	if err != nil {
		return result, err
	}
	result, err = XmlDecode(r)
	return result, err
}

/**
 * 格式化struct
 */
func (unifyOrder *UnifyOrder) Format() string {
	return StringFormat(unifyOrder)
}
