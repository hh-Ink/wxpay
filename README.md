wxpay - 微信支付 go语言SDK
===

    该 SDK 对微信支付的涉及到的工具做了封装并对常用接口进行了实现。如果想要实现不存在的接口，模仿例子实现即可
    
现有接口调用
---

1. 统一下单

```golang
uor := wxpay.NewUor(`yourAppId`, `yourMchId`)
uor.Trade_type = `APP`
uor.Out_trade_no = `2018082606352323654`
uor.Body = `微信支付测试`
uor.Total_fee = 1
uor.MakeSign(`yourKey`)
u, err := uor.Post()
if err != nil {
    fmt.Println(err)
}
fmt.Println(wxpay.StructFormat(uor))
```
2. 订单查询
```golang
q := wxpay.NewOrderQuery(`yourAppId`, `yourMchId`)
q.Transaction_id = `4200000185201808090657628709`
q.MakeSign(`yourKey`)
u, err := q.Post()
if err != nil {
    fmt.Println(err)
}
fmt.Println(wxpay.StructFormat(u))
```

3. 退款
```golang
r := wxpay.NewRefund(`yourAppId`, `yourMchId`)
r.Transaction_id = `4200000175201808106695291370`
r.Out_refund_no = `2018082606352323654`
r.Total_fee = 1000
r.Refund_fee = 1000
r.MakeSign(`yourKey`)
res, err := r.Post(`apiclient_cert.pem`, `apiclient_key.pem`, `rootca.pem`)
if err != nil {
    fmt.Println(err)
}
fmt.Println(wxpay.StructFormat(res))
```
4. 支付回调及退款回调验证并提取信息
```golang
xmlText := `<xml><appid><![CDATA[yourAppId]]></appid><attach><![CDATA[attach]]></attach><bank_type><![CDATA[CIB_DEBIT]]></bank_type><cash_fee><![CDATA[1]]></cash_fee><fee_type><![CDATA[CNY]]></fee_type><is_subscribe><![CDATA[N]]></is_subscribe><mch_id><![CDATA[yourMchId]]></mch_id><nonce_str><![CDATA[8ldyxnpwnzsua53wx029kuvohf16rsrm]]></nonce_str><openid><![CDATA[otMtJwh51qmLnkHRIk9v0fL8RoN8]]></openid><out_trade_no><![CDATA[2018082606352323654]]></out_trade_no><result_code><![CDATA[SUCCESS]]></result_code><return_code><![CDATA[SUCCESS]]></return_code><sign><![CDATA[DBE0A2EA1B9DD0CA828BAAACB9A8DDD9]]></sign><time_end><![CDATA[20180810152350]]></time_end><total_fee>1000</total_fee><trade_type><![CDATA[APP]]></trade_type><transaction_id><![CDATA[4200000175201808106695291370]]></transaction_id></xml>`
n := wxpay.NewNotify()
m, _ := n.CheckSign(xmlText, `yourKey`)
fmt.Println(m)
n.Print()
```
