package wxpay

import (
	"encoding/xml"
	"fmt"
)

type Notify struct {
	XMLName     xml.Name `xml:"xml"`
	Return_code string   `xml:"return_code"`
	Return_msg  string   `xml:"return_msg"`
}

/**
 * 构造方法
 */
func NewNotify() (notify *Notify) {
	notify = &Notify{}
	return
}

/**
 * 验证微信的回调结果
 */
func (notify *Notify) CheckSign(notifyXML string, key string) (map[string]string, error) {
	m, ero := XmlDecode(notifyXML)
	if ero != nil {
		notify.SetCode("FAIL")
		notify.SetMsg(fmt.Sprintf("%s", ero))
		return m, ero
	}
	if _, ok := m["req_info"]; ok { // 退款回调
		reqXML, err := AesDecrypt(m[`req_info`], key)
		delete(m, "req_info")
		if err != nil {
			notify.SetCode("FAIL")
			notify.SetMsg(fmt.Sprintf("%s", err))
			return make(map[string]string), err
		}
		reqMAP, _err := XmlDecode(reqXML)
		if _err != nil {
			return make(map[string]string), _err
		}
		for k, v := range reqMAP {
			m[k] = v
		}
	} else { // 支付回调
		sign := m["sign"]
		delete(m, "sign")
		nSign := WxpayCalcSign(m, key)
		if sign != nSign {
			notify.SetCode("FAIL")
			notify.SetMsg("签名验证失败")
			return m, fmt.Errorf("%s", "签名验证失败")
		}
	}
	notify.SetCode("SUCCESS")
	notify.SetMsg("OK")
	return m, nil
}

/**
 * 设置 return_code 值
 */
func (notify *Notify) SetCode(code string) {
	notify.Return_code = code
}

/**
 * 设置 return_msg 值
 */
func (notify *Notify) SetMsg(msg string) {
	notify.Return_msg = msg
}

/**
 * 输出回调处理结果
 */
func (notify *Notify) Print() {
	uint8Xml, _err := XmlEncode(notify)
	if _err != nil {
		fmt.Println(_err)
	}
	fmt.Println(string(uint8Xml))
}
