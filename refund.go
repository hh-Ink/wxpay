package wxpay

import (
	"encoding/xml"
	"fmt"
)

type Refund struct {
	XMLName         xml.Name `xml:"xml"`
	Appid           string   `xml:"appid"`
	Mch_id          string   `xml:"mch_id"`
	Nonce_str       string   `xml:"nonce_str"`
	Notify_url      string   `xml:"notify_url"`
	Total_fee       int      `xml:"total_fee"`
	Out_trade_no    string   `xml:"out_trade_no"`
	Transaction_id  string   `xml:"transaction_id"`
	Out_refund_no   string   `xml:"out_refund_no"`
	Sign            string   `xml:"sign"`
	Sign_type       string   `xml:"sign_type"`
	Refund_fee      int      `xml:"refund_fee"`
	Refund_fee_type string   `xml:"refund_fee_type"`
	Refund_desc     string   `xml:"refund_desc"`
	Refund_account  string   `xml:"refund_account"`
}

/**
 * 构造方法
 */
func NewRefund(appid string, mch_id string) (refund *Refund) {
	refund = &Refund{}
	refund.Appid = appid
	refund.Mch_id = mch_id
	return
}

/**
 * 生成签名
 */
func (refund *Refund) MakeSign(key string) {
	refund.Nonce_str = MakeNoter()
	refund.Sign = WxpayCalcSign(refund, key)
}

/**
 * 发送请求
 */
func (refund *Refund) Post(rootCa string, rootKey string, wechatCAPath string) (map[string]string, error) {
	result := make(map[string]string)
	// 检测必填参数
	if len(refund.Out_trade_no) == 0 && len(refund.Transaction_id) == 0 {
		return result, fmt.Errorf("%s", "out_trade_no 和 transaction_id 不能同时为空")
	}
	if len(refund.Out_refund_no) == 0 {
		return result, fmt.Errorf("%s", "out_refund_no 不能为空")
	}
	if refund.Total_fee == 0 {
		return result, fmt.Errorf("%s", "total_fee 参数必选")
	}
	if refund.Refund_fee == 0 {
		return result, fmt.Errorf("%s", "refund_fee 参数必选")
	}
	postData, _err := XmlEncode(refund)
	if _err != nil {
		return result, _err
	}
	url := "https://api.mch.weixin.qq.com/secapi/pay/refund"
	r, err := RefundPost(postData, url, rootCa, rootKey, wechatCAPath)
	if err != nil {
		return result, err
	}
	result, err = XmlDecode(r)
	return result, err
}

/**
 * 格式化struct
 */
func (refund *Refund) Format() string {
	return StringFormat(refund)
}
