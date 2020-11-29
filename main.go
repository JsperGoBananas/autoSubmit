package main

import (
	"autosubmit/notifier"
	"autosubmit/utils"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

var jar http.CookieJar

var email string
var lxdh string

var username = flag.String("username", "", "学号")
var password = flag.String("password", "", "portal密码")
var reason = flag.String("reason", "探亲", "出入校事由")
var track = flag.String("track", "北大西门-畅春园-北大西门", "出校行动轨迹")
var ftKey string

func errorNotifier(title string) {
	if err := recover(); err != nil {
		notifier.Ft(ftKey, title, *username)
		panic(title)
	}
}

func initCookies() {
	jar, _ = cookiejar.New(nil)
}

func initFlags() {
	ftKey = os.Getenv("FT_SCKEY")
	defer errorNotifier("备案失败：参数配置错误")
	flag.Parse()
	if *username == "" {
		*username = os.Getenv("USERNAME")
	}
	if *password == "" {
		*password = os.Getenv("PASSWORD")
	}
	if *username == "" || *password == "" {
		panic("Must specify username and password")
	}
}

func oauthLogin() string {
	defer errorNotifier("备案失败：iaaa登录失败")
	payloadStr := fmt.Sprintf("appid=portal2017&userName=%s&password=%s", *username, *password) + "&randCode=&smsCode=&otpCode=&redirUrl=https%3A%2F%2Fportal.pku.edu.cn%2Fportal2017%2FssoLogin.do"
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	payload := strings.NewReader(payloadStr)
	req, err := http.NewRequest("POST", "https://iaaa.pku.edu.cn/iaaa/oauthlogin.do", payload)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.80 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://iaaa.pku.edu.cn")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")

	client := http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	token := utils.OauthLoginResp{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		panic(err)
	}
	if !token.Success {
		panic("iaaa密码错误")
	}
	return token.Token
}

func ssoLogin(token string) {
	defer errorNotifier("备案失败：iaaa登录失败")
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	req, err := http.NewRequest("GET", "https://portal.pku.edu.cn/portal2017/ssoLogin.do?_rand=0.6223201749662104&token=" + token, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.80 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Referer", "https://iaaa.pku.edu.cn/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")
	req.Header.Set("Cookie", "fromURL=/pub/life")

	client := http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}

func testPortal() {
	defer errorNotifier("备案失败：iaaa登录失败")
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	req, err := http.NewRequest("GET", "https://portal.pku.edu.cn/portal2017/isUserLogged.do", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.80 Safari/537.36")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://portal.pku.edu.cn/portal2017/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")

	client := http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	result := utils.PortalCheckResp{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
	if result.Success {
		fmt.Println("portal登录成功")
	}
}


func testSimso() {
	defer errorNotifier("备案失败：备案小程序登录失败")

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	req, err := http.NewRequest("GET", "https://simso.pku.edu.cn/ssapi/getLoginInfo", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.80 Safari/537.36")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://portal.pku.edu.cn/portal2017/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")

	client := http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	result := utils.SimsoCheckResp{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
	if result.Success {
		fmt.Println("simso登录成功")
	}
}

func appSysRedir() string {
	defer errorNotifier("备案失败：备案小程序登录失败")

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	req, err := http.NewRequest("GET", "https://portal.pku.edu.cn/portal2017/util/appSysRedir.do?appId=stuCampusExEn", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.80 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")

	client := http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	values, err := url.ParseQuery(resp.Request.URL.RawQuery)
	if err != nil {
		panic(err)
	}
	return values["token"][0]
}

func simsoLogin(token string) string {
	defer errorNotifier("备案失败：备案小程序登录失败")

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	req, err := http.NewRequest("GET", "https://simso.pku.edu.cn/ssapi/simsoLogin?token=" + token, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.80 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")

	client := http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	result := utils.SimsoLoginResp{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
	return result.Sid
}


func saveOut(sid string) string {
	defer errorNotifier("备案失败：出校备案保存失败")
	timeLocal, _ := time.LoadLocation("Asia/Chongqing")
	time.Local = timeLocal
	t := time.Now().Local()
	today := t.Format("20060102")

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	data := utils.SaveSqxxReq{
		Crxsy: *reason,
		Cxcs: 1,
		Cxmdd: "北京",
		Cxrq: today,
		Cxxdgj: *track,
		Cxxm: "",
		Dfx14Qrbz: "y",
		Email: email,
		Jzdbjdjrq: "",
		Jzdbjjd: "",
		Jzdbjqx: "",
		Jzdbjyzzj14: "",
		Jzdjwdjrq: "",
		Jzdjwdjsdm: "",
		Jzdjwgjdq: "156",
		Jzdjwssdm: "",
		Lxdh: lxdh,
		Rxcs: 1,
		Rxjzd: "北京",
		Rxrq: "",
		Rxxm: "",
		Sfkcx: true,
		Sfqdcxrq: "",
		Sfqdhxrq: "",
		Sfyxtycj: "",
		Shbz: "",
		Sqbh: "",
		Sqlb: "出校",
		Szxq: "燕园",
		Tjbz: "",
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://simso.pku.edu.cn/ssapi/stuaffair/epiAccess/saveSqxx?sid=" + sid, body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.80 Safari/537.36")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Origin", "https://simso.pku.edu.cn")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://simso.pku.edu.cn/pages/epidemicAccess.html")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")

	client := http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	result := utils.SaveSqxxResp{}
	err = json.Unmarshal(body2, &result)
	if err != nil {
		panic(err)
	}
	return result.Row
}

func saveIn(sid string) string {
	defer errorNotifier("备案失败：入校备案保存失败")
	timeLocal, _ := time.LoadLocation("Asia/Chongqing")
	time.Local = timeLocal
	t := time.Now().Local()
	today := t.Format("20060102")

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	data := utils.SaveSqxxReq{
		Crxsy: *reason,
		Cxcs: 1,
		Cxmdd: "北京",
		Cxrq: "",
		Cxxdgj: "",
		Cxxm: "",
		Dfx14Qrbz: "y",
		Email: email,
		Jzdbjdjrq: "",
		Jzdbjjd: "燕园街道",
		Jzdbjqx: "08",
		Jzdbjyzzj14: "y",
		Jzdjwdjrq: "",
		Jzdjwdjsdm: "",
		Jzdjwgjdq: "156",
		Jzdjwssdm: "",
		Lxdh: lxdh,
		Rxcs: 1,
		Rxjzd: "北京",
		Rxrq: today,
		Rxxm: "",
		Sfkcx: true,
		Sfqdcxrq: "",
		Sfqdhxrq: "",
		Sfyxtycj: "",
		Shbz: "",
		Sqbh: "",
		Sqlb: "入校",
		Szxq: "燕园",
		Tjbz: "",
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://simso.pku.edu.cn/ssapi/stuaffair/epiAccess/saveSqxx?sid=" + sid, body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.80 Safari/537.36")
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Origin", "https://simso.pku.edu.cn")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://simso.pku.edu.cn/pages/epidemicAccess.html")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")

	client := http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	result := utils.SaveSqxxResp{}
	err = json.Unmarshal(body2, &result)
	if err != nil {
		panic(err)
	}
	return result.Row
}


func submitSqxx(sid string, row string) bool {
	defer errorNotifier("备案失败：提交失败")

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	req, err := http.NewRequest("GET", "https://simso.pku.edu.cn/ssapi/stuaffair/epiAccess/submitSqxx?sid="+ sid + "&sqbh=" + row, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.80 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://simso.pku.edu.cn/pages/epidemicAccess.html")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")

	client := http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	result := utils.SaveSqxxResp{}
	err = json.Unmarshal(body2, &result)
	if err != nil {
		panic(err)
	}
	if result.Success {
		return true
	}
	return false
}

func getSqzt(sid string) (string, string) {
	defer errorNotifier("备案失败：获取上次填报信息失败")

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	req, err := http.NewRequest("GET", "https://simso.pku.edu.cn/ssapi/stuaffair/epiAccess/getSqzt?sid=" + sid, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.80 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://simso.pku.edu.cn/pages/epidemicAccess.html")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,zh-TW;q=0.7")

	client := http.Client{Jar: jar}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	result := utils.SqztResp{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
	return result.Row.LastSqxx.Email, result.Row.LastSqxx.Lxdh
}


func main() {
	initCookies()
	initFlags()

	portalToken := oauthLogin()
	ssoLogin(portalToken)
	testPortal()
	simsoToken := appSysRedir()
	sid := simsoLogin(simsoToken)
	testSimso()

	email, lxdh = getSqzt(sid)

	row := saveOut(sid)
	if submitSqxx(sid, row) {
		notifier.Ft(ftKey, "出校备案成功", *username)
	}
	row = saveIn(sid)
	if submitSqxx(sid, row) {
		notifier.Ft(ftKey, "入校备案成功", *username)
	}
}
