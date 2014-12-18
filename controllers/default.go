package controllers

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"sort"
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
	ArticleCount int
	Articles     []*item
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
	wresp, err := dealwith(wreq)
	if err != nil {
		beego.Error(err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	data, err := wresp.Encode()
	if err != nil {
		beego.Error(err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		return
	}
	c.Ctx.WriteString(string(data))
	return
}

func dealwith(req *Request) (resp *TextResponse, err error) {
	resp = &TextResponse{}
	resp.CreateTime = time.Duration(time.Now().Unix())
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.MsgType = MsgTypeText
	beego.Info(req.MsgType)
	beego.Info(req.Content)
	if req.MsgType == MsgTypeText {
		resp.Content = "Yes! You got it."
		return resp, nil
	} else {
		resp.Content = "Not supported yet."
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

func NewResponse() (resp *msgBaseResp) {
	resp = &msgBaseResp{}
	resp.CreateTime = time.Duration(time.Now().Unix())
	return
}

func (resp TextResponse) Encode() (data []byte, err error) {
	resp.CreateTime = time.Second
	data, err = xml.Marshal(resp)
	return
}
