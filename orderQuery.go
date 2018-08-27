package wxpay

import (
	"encoding/xml"
	"fmt"
)

type OrderQuery struct {
	XMLName        xml.Name `xml:"xml"`
	Appid          string   `xml:"appid"`
	Mch_id         string   `xml:"mch_id"`
	Nonce_str      string   `xml:"nonce_str"`
	Out_trade_no   string   `xml:"out_trade_no"`
	Transaction_id string   `xml:"transaction_id"`
	Sign           string   `xml:"sign"`
	Sign_type      string   `xml:"sign_type"`
}

/**
 * 构造方法
 */
func NewOrderQuery(appid string, mch_id string) (orderQuery *OrderQuery) {
	orderQuery = &OrderQuery{}
	orderQuery.Appid = appid
	orderQuery.Mch_id = mch_id
	return
}

/**
 * 生成签名
 */
func (orderQuery *OrderQuery) MakeSign(key string) {
	orderQuery.Nonce_str = MakeNoter()
	orderQuery.Sign = WxpayCalcSign(orderQuery, key)
}

/**
 * 发送请求
 */
func (orderQuery *OrderQuery) Post() (map[string]string, error) {
	result := make(map[string]string)
	// 检测必填参数
	if len(orderQuery.Out_trade_no) == 0 && len(orderQuery.Transaction_id) == 0 {
		return result, fmt.Errorf("%s", "out_trade_no 和 transaction_id 不能同时为空")
	}
	postData, _err := XmlEncode(orderQuery)
	if _err != nil {
		return result, _err
	}
	url := "https://api.mch.weixin.qq.com/pay/orderquery"
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
func (orderQuery *OrderQuery) Format() string {
	return StringFormat(orderQuery)
}
