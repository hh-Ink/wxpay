package wxpay

import (
	"encoding/xml"
	"fmt"
)

type RefundQuery struct {
	XMLName        xml.Name `xml:"xml"`
	Appid          string   `xml:"appid"`
	Mch_id         string   `xml:"mch_id"`
	Nonce_str      string   `xml:"nonce_str"`
	Out_trade_no   string   `xml:"out_trade_no"`
	Transaction_id string   `xml:"transaction_id"`
	Out_refund_no  string   `xml:"out_refund_no"`
	Sign           string   `xml:"sign"`
	Sign_type      string   `xml:"sign_type"`
	Refund_id      string   `xml:"refund_id"`
	Offset         string   `xml:"offset"`
}

/**
 * 构造方法
 */
func NewRefundQuery(appid string, mch_id string) (refundQuery *RefundQuery) {
	refundQuery = &RefundQuery{}
	refundQuery.Appid = appid
	refundQuery.Mch_id = mch_id
	return
}

/**
 * 生成签名
 */
func (refundQuery *RefundQuery) MakeSign(key string) {
	refundQuery.Nonce_str = MakeNoter()
	refundQuery.Sign = WxpayCalcSign(refundQuery, key)
}

/**
 * 发送请求
 */
func (refundQuery *RefundQuery) Post() (map[string]string, error) {
	result := make(map[string]string)
	// 检测必填参数
	if len(refundQuery.Out_trade_no) == 0 &&
		len(refundQuery.Transaction_id) == 0 &&
		len(refundQuery.Out_refund_no) == 0 &&
		len(refundQuery.Refund_id) == 0 {
		return result, fmt.Errorf("%s", "out_trade_no, transaction_id, refund_id, out_refund_no 四个必选一个")
	}
	postData, _err := XmlEncode(refundQuery)
	if _err != nil {
		return result, _err
	}
	url := "https://api.mch.weixin.qq.com/secapi/pay/refundquery"
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
func (refundQuery *RefundQuery) Format() string {
	return StringFormat(refundQuery)
}
