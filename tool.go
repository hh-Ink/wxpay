package wxpay

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"
)

/**
 * 微信支付计算签名的函数
 */
func WxpayCalcSign(s interface{}, key string) (sign string) {
	// 将结构转化为map
	m := StructToMap(s)

	//STEP 1, 对key进行升序排序.
	sorted_keys := make([]string, 0)
	for k, _ := range m {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)

	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sorted_keys {
		value := fmt.Sprintf("%v", m[k])
		// 跳过 xml 顶标签设置项
		if k == `XMLName` {
			continue
		}
		if value != "" {
			// 将参数名转为小写拼接
			signStrings = signStrings + strings.ToLower(k) + "=" + value + "&"
		}
	}

	//STEP3, 在键值对的最后加上key=API_KEY
	if key != "" {
		signStrings = signStrings + "key=" + key
	}

	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign
}

/**
 * 将 struct 转化为 map
 */
func StructToMap(s interface{}) (m map[string]interface{}) {
	m = make(map[string]interface{})
	j, _ := json.Marshal(s)
	json.Unmarshal(j, &m)
	return m
}

/**
 * 发送post请求
 */
func Post(bytes_req []uint8, url string) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(bytes_req))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")
	c := http.Client{}
	resp, _err := c.Do(req)
	if _err != nil {
		return "", _err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

/**
 * 申请退款 （需要配置证书）
 */
func RefundPost(bytes_req []uint8, url string, rootCa string, rootKey string, caPath string) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(bytes_req))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	var _tlsConfig *tls.Config
	cert, err := tls.LoadX509KeyPair(rootCa, rootKey)
	if err != nil {
		return "", err
	}
	caData, err := ioutil.ReadFile(caPath)
	if err != nil {
		return "", err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)
	_tlsConfig = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            pool,
		InsecureSkipVerify: true,
	}
	tr := &http.Transport{TLSClientConfig: _tlsConfig}
	c := http.Client{Transport: tr}
	resp, _err := c.Do(req)
	if _err != nil {
		return "", _err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

/**
 * xml 编码
 */
func XmlEncode(s interface{}) ([]uint8, error) {
	bytes_req, err := xml.Marshal(s)
	if err != nil {
		return []byte(``), err
	}
	str_req := string(bytes_req)
	bytes_req = []byte(str_req)
	return bytes_req, nil
}

/**
 * xml 解码
 */
func XmlDecode(xmlText string) (map[string]string, error) {
	var (
		d     *xml.Decoder
		start *xml.StartElement
	)
	m := make(map[string]string)
	d = xml.NewDecoder(bytes.NewBuffer([]byte(xmlText)))
	for {
		tok, err := d.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			start = &t
		case xml.CharData:
			if t = bytes.TrimSpace(t); len(t) > 0 {
				m[start.Name.Local] = string(t)
			}
		}
	}
	return m, nil
}

/**
 * 结构类型格式化
 */
func StringFormat(s interface{}) string {
	m := StructToMap(s)
	var buf bytes.Buffer
	for k, v := range m {
		buf.WriteString(string(k))
		buf.WriteString(":")
		if str, ok := v.(string); ok {
			buf.WriteString(str)
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

/**
 * 生成随机串
 */
func MakeNoter() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 8; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	ctx := md5.New()
	ctx.Write([]byte(string(result)))
	return hex.EncodeToString(ctx.Sum(nil))
}

/**
 * 解密微信退款回调
 *（1）对加密串A做base64解码，得到加密串B
 *（2）对商户key做md5，得到32位小写key* ( key设置路径：微信商户平台(pay.weixin.qq.com)-->账户设置-->API安全-->密钥设置 )
 *（3）用key*对加密串B做AES-256-ECB解密（PKCS7Padding）
 */
func AesDecrypt(cryted string, key string) (string, error) {
	// 对商户key做md5，得到32位小写key
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(key))
	cipherStr := md5Ctx.Sum(nil)
	key = hex.EncodeToString(cipherStr)

	// 对加密串A做base64解码，得到加密串B
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	// 用key*对加密串B做AES-256-ECB解密（PKCS7Padding）
	SetAesKey(key)
	plaintext, err := AesECBDecrypt(crytedByte)
	if err != nil {
		return ``, err
	}
	return string(plaintext), nil
}
