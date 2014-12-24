package controllers

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

const (
	TOKEN = "weixingis"

	MsgTypeDefault          = ".*"
	MsgTypeText             = "text"
	MsgTypeImage            = "image"
	MsgTypeVoice            = "voice"
	MsgTypeVideo            = "video"
	MsgTypeLocation         = "location"
	MsgTypeLink             = "link"
	MsgTypeNews             = "news"
	MsgTypeEvent            = "event"
	MsgTypeEventSubscribe   = "subscribe"
	MsgTeypEventUnsubscribe = "unsubscribe"
)

type msgBaseReq struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
}

type msgBaseResp struct {
	XMLName      xml.Name `xml:xml`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	FuncFlag     int // 位0x0001被标志时，星标刚收到的消息

}

type Request struct {
	XMLName xml.Name `xml:xml`
	msgBaseReq
	Content    string
	Location_X float32
	Location_Y float32
	Scale      int
	Label      string
	PicUrl     string
	MsgId      int64
}

//回复文本消息
type TextResponse struct {
	XMLName xml.Name `xml:"xml"`
	msgBaseResp
	Content string
}

//回复图片消息
type ImageResponse struct {
	XMLName xml.Name `xml:"xml"`
	msgBaseResp
	MediaId int64 //通过上传多媒体文件，得到的id。
}

//回复语音消息
type VoiceResponse struct {
	XMLName xml.Name `xml:"xml"`
	msgBaseResp
	MediaId int64 //通过上传多媒体文件，得到的id。
}

//回复视频消息
type VideoResponse struct {
	XMLName xml.Name `xml:"xml"`
	msgBaseResp
	MediaId     int64
	Title       string
	Description string
}

//回复音乐消息
type MusicResponse struct {
	XMLName xml.Name `xml:"xml"`
	msgBaseResp
	Title        string
	Description  string
	MusicURL     string
	HQMusicUrl   string
	ThumbMediaId int64 //缩略图的媒体id，通过上传多媒体文件，得到的id
}

//回复图文消息
type NewsResponse struct {
	XMLName xml.Name `xml:"xml"`
	msgBaseResp
	ArticleCount int     `xml:",omitempty"`
	Articles     []*item `xml:"Articles>item,omitempty"`
}

type item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string
	Description string
	PicUrl      string
	Url         string
}

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	signature := c.Input().Get("signature")
	beego.Info(signature)
	timestamp := c.Input().Get("timestamp")
	beego.Info(timestamp)
	nonce := c.Input().Get("nonce")
	beego.Info(nonce)
	echostr := c.Input().Get("echostr")
	beego.Info(echostr)
	beego.Info(Signature(timestamp, nonce))
	if Signature(timestamp, nonce) == signature {
		c.Ctx.WriteString(echostr)
	} else {
		c.Ctx.WriteString("")
	}
}

func (c *MainController) Post() {
	body, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		beego.Error(err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	beego.Info(string(body))
	var wreq *Request
	if wreq, err = DecodeRequest(body); err != nil {
		beego.Error(err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	beego.Info(wreq.Content)
	//wresp, err := dealwith(wreq)
	//wresp, err := responseNewsMsg(wreq)
	//if err != nil {
	//	beego.Error(err)
	//	c.Ctx.ResponseWriter.WriteHeader(500)
	//	return
	//}
	//data, err := wresp.Encode()
	//if err != nil {
	//	beego.Error(err)
	//	c.Ctx.ResponseWriter.WriteHeader(500)
	//	return
	//}
	//beego.Info(string(data))
	str, err := dealwith(wreq)
	if err != nil {
		beego.Error(err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	beego.Info(str)
	c.Ctx.WriteString(str)
	return
}

func dealwith(req *Request) (str string, err error) {
	content := strings.Trim(strings.ToLower(req.Content), " ")
	if req.MsgType == MsgTypeText {
		switch content {
		case "arcgis", "arcgisproduct":
			str, _ = responseProduct(req, "README")
		case "desktop":
			str, _ = responseProduct(req, "desktop")
		case "server":
			str, _ = responseProduct(req, "server")
		case "engine":
			str, _ = responseProduct(req, "engine")
		default:
			responseChat(req, content)
		}
	} else if req.MsgType == MsgTypeImage {

	} else if req.MsgType == MsgTypeVoice {

	} else if req.MsgType == MsgTypeVideo {

	} else if req.MsgType == MsgTypeLocation {

	} else if req.MsgType == MsgTypeLink {

	}
	return
}

//回复产品信息
func responseProduct(req *Request, product string) (str string, err error) {
	resp := NewNewsResponse()
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	var resurl string
	var a item
	resurl = "https://raw.github.com/xzdbd/gisproduct/master/arcgisproduct/" + strings.Trim(product, " ") + ".md"
	a.Url = "https://github.com/xzdbd/gisproduct/blob/master/arcgisproduct/" + strings.Trim(product, " ") + ".md"
	rsp, err := http.Get(resurl)
	if err != nil {
		beego.Info("error:" + err.Error())
		resp := NewTextResponse()
		resp.ToUserName = req.FromUserName
		resp.FromUserName = req.ToUserName
		resp.Content = "不存在该产品"
	} else if rsp.StatusCode == 404 {
		beego.Info("error:" + err.Error())
		resp := NewTextResponse()
		resp.ToUserName = req.FromUserName
		resp.FromUserName = req.ToUserName
		resp.Content = "找不到你要查询的产品"
	} else {
		resp.ArticleCount = 1
		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			beego.Error(err)
		}
		a.Description = getProductIntro(string(body))
		a.Title = req.Content
		a.PicUrl = "https://github.com/xzdbd/gisproduct/raw/master/images/desktop1.png?raw=true"
		resp.Articles = append(resp.Articles, &a)
		resp.FuncFlag = 1
	}
	data, err := resp.Encode()
	if err != nil {
		return
	}
	str = string(data)
	//defer rsp.Body.Close()
	return
}

//回复聊天信息
func responseChat(req *Request, content string) (str string, err error) {
	return
}

func responseNewsMsg(req *Request) (resp *NewsResponse, err error) {
	resp = NewNewsResponse()
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	if req.MsgType == MsgTypeText {
		str := strings.ToLower(req.Content)
		if strings.Trim(str, " ") == "desktop" {
			var resurl string
			var a item
			resurl = "https://raw.github.com/xzdbd/gisproduct/master/arcgisproduct/" + strings.Trim(str, " ") + ".md"
			a.Url = "https://github.com/xzdbd/gisproduct/blob/master/arcgisproduct/" + strings.Trim(str, " ") + ".md"
			beego.Info(resurl)
			beego.Info(a.Url)
			rsp, err := http.Get(resurl)
			if err != nil {
				beego.Info("error")
				return nil, err
			}
			defer rsp.Body.Close()
			if rsp.StatusCode == 404 {
				beego.Info("could not found")
				return resp, nil
			}
			resp.ArticleCount = 1
			body, err := ioutil.ReadAll(rsp.Body)
			beego.Info(string(body))
			a.Description = string(body)
			a.Title = req.Content
			a.PicUrl = "https://github.com/xzdbd/gisproduct/raw/master/images/desktop1.png?raw=true"
			resp.Articles = append(resp.Articles, &a)
			resp.FuncFlag = 1
		} else if strings.Trim(str, " ") == "arcgis" {
			var resurl string
			var a item
			resurl = "https://raw.github.com/xzdbd/gisproduct/master/arcgisproduct/README.md"
			a.Url = "https://github.com/xzdbd/gisproduct/blob/master/arcgisproduct/README.md"
			beego.Info(resurl)
			beego.Info(a.Url)
			rsp, err := http.Get(resurl)
			if err != nil {
				beego.Info("error")
				return nil, err
			}
			defer rsp.Body.Close()
			if rsp.StatusCode == 404 {
				beego.Info("could not found")
				return resp, nil
			}
			resp.ArticleCount = 1
			body, err := ioutil.ReadAll(rsp.Body)
			beego.Info(string(body))
			a.Description = string(body)
			a.Title = req.Content
			a.PicUrl = "https://github.com/xzdbd/gisproduct/raw/master/images/desktop1.png?raw=true"
			resp.Articles = append(resp.Articles, &a)
			resp.FuncFlag = 1
		}
	} else {
		beego.Info("not supported")
	}
	return resp, nil
}

func Signature(timestamp, nonce string) string {
	strs := sort.StringSlice{TOKEN, timestamp, nonce}
	sort.Strings(strs)
	str := ""
	for _, s := range strs {
		str += s
	}
	h := sha1.New()
	h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func DecodeRequest(data []byte) (req *Request, err error) {
	req = &Request{}
	if err = xml.Unmarshal(data, req); err != nil {
		return
	}
	req.CreateTime *= time.Second
	return
}

func NewBaseResponse() (resp *msgBaseResp) {
	resp = &msgBaseResp{}
	resp.CreateTime = time.Duration(time.Now().Unix())
	return
}

func NewTextResponse() (resp *TextResponse) {
	resp = &TextResponse{}
	resp.CreateTime = time.Duration(time.Now().Unix())
	return
}

func NewNewsResponse() (resp *NewsResponse) {
	resp = &NewsResponse{}
	resp.CreateTime = time.Duration(time.Now().Unix())
	resp.MsgType = MsgTypeNews
	return
}

func (resp TextResponse) Encode() (data []byte, err error) {
	resp.CreateTime = time.Second
	data, err = xml.Marshal(resp)
	return
}

func (resp NewsResponse) Encode() (data []byte, err error) {
	resp.CreateTime = time.Second
	data, err = xml.Marshal(resp)
	return
}

func SubString(s string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(s)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}

//获取产品信息简介，用于图文信息的decription，截取到第二个#号
func getProductIntro(s string) (subStr string) {
	s = strings.Replace(s, "#", " ", 1)
	l := strings.Index(s, "#")
	subStr = SubString(s, 0, l)
	return

}

type Response interface {
	Encode()
}
