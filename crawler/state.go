package crawler

import (
	"github.com/huaimeng666/URLFinder/cmd"
	"github.com/huaimeng666/URLFinder/config"
	"github.com/huaimeng666/URLFinder/mode"
	"github.com/huaimeng666/URLFinder/result"
	"github.com/huaimeng666/URLFinder/util"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// 检测js访问状态码
func JsState(u string, i int, sou string) {

	defer func() {
		config.Wg.Done()
		<-config.Jsch
		PrintProgress()
	}()
	if cmd.S == "" {
		result.ResultJs[i].Url = u
		return
	}
	if cmd.M == 3 {
		for _, v := range config.Risks {
			if strings.Contains(u, v) {
				result.ResultJs[i] = mode.Link{Url: u, Status: "疑似危险路由"}
				return
			}
		}
	}

	var redirect string
	ur, err2 := url.Parse(u)
	if err2 != nil {
		return
	}
	request, err := http.NewRequest("GET", ur.String(), nil)
	if err != nil {
		result.ResultJs[i].Url = ""
		return
	}
	if cmd.C != "" {
		request.Header.Set("Cookie", cmd.C)
	}
	request.Header.Set("User-Agent", util.GetUserAgent())
	request.Header.Set("Accept", "*/*")
	if cmd.I {
		util.SetHeadersConfig(&request.Header)
	}

	response, err := client.Do(request)
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout") && cmd.S == "" {
			result.ResultJs[i] = mode.Link{Url: u, Status: "timeout", Size: "0"}

		} else {
			result.ResultJs[i].Url = ""
		}
		return
	}
	defer response.Body.Close()

	code := response.StatusCode
	if strings.Contains(cmd.S, strconv.Itoa(code)) || cmd.S == "all" && (sou != "Fuzz" && code == 200) {
		var length int
		dataBytes, err := io.ReadAll(response.Body)
		if err != nil {
			length = 0
		} else {
			length = len(dataBytes)
		}

		config.Lock.Lock()
		if result.Redirect[ur.String()] {
			code = 302
			redirect = response.Request.URL.String()
		}
		config.Lock.Unlock()

		seenKey := strconv.Itoa(code) + strconv.Itoa(length) + ""
		result.SeenMutex.Lock()
		if _, ok := result.Seen[seenKey]; ok {
			result.ResultJs[i].Url = ""
			result.SeenMutex.Unlock()
			return
		}
		result.Seen[seenKey] = struct{}{}
		result.SeenMutex.Unlock()

		result.ResultJs[i] = mode.Link{Url: u, Status: strconv.Itoa(code), Size: strconv.Itoa(length), Redirect: redirect}
	} else {
		result.ResultJs[i].Url = ""
	}
}

// 检测url访问状态码
func UrlState(u string, i int) {
	defer func() {
		config.Wg.Done()
		<-config.Urlch
		PrintProgress()
	}()
	if cmd.S == "" {
		result.ResultUrl[i].Url = u
		return
	}
	if cmd.M == 3 {
		for _, v := range config.Risks {
			if strings.Contains(u, v) {
				result.ResultUrl[i] = mode.Link{Url: u, Status: "0", Size: "0", Title: "疑似危险路由,已跳过验证"}
				return
			}
		}
	}

	var redirect string
	ur, err2 := url.Parse(u)
	if err2 != nil {
		return
	}
	request, err := http.NewRequest("GET", ur.String(), nil)
	if err != nil {
		result.ResultUrl[i].Url = ""
		return
	}

	if cmd.C != "" {
		request.Header.Set("Cookie", cmd.C)
	}
	request.Header.Set("User-Agent", util.GetUserAgent())
	request.Header.Set("Accept", "*/*")

	if cmd.I {
		util.SetHeadersConfig(&request.Header)
	}

	response, err := client.Do(request)
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout") && cmd.S == "all" {
			result.ResultUrl[i] = mode.Link{Url: u, Status: "timeout", Size: "0"}
		} else {
			result.ResultUrl[i].Url = ""
		}
		return
	}
	defer response.Body.Close()

	code := response.StatusCode
	if strings.Contains(cmd.S, strconv.Itoa(code)) || cmd.S == "all" {
		var length int
		dataBytes, err := io.ReadAll(response.Body)
		if err != nil {
			length = 0
		} else {
			length = len(dataBytes)
		}

		body := string(dataBytes)
		re := regexp.MustCompile("<[tT]itle>(.*?)</[tT]itle>")
		title := re.FindAllStringSubmatch(body, -1)
		config.Lock.Lock()
		if result.Redirect[ur.String()] {
			code = 302
			redirect = response.Request.URL.String()
		}
		config.Lock.Unlock()

		if len(title) != 0 {
			seenKey := strconv.Itoa(code) + strconv.Itoa(length) + title[0][1]
			result.SeenMutex.Lock()
			if _, ok := result.Seen[seenKey]; ok {
				result.ResultUrl[i].Url = ""
				result.SeenMutex.Unlock()
				return
			}
			result.Seen[seenKey] = struct{}{}
			result.SeenMutex.Unlock()

			result.ResultUrl[i] = mode.Link{Url: u, Status: strconv.Itoa(code), Size: strconv.Itoa(length), Title: title[0][1], Redirect: redirect}
		} else {
			seenKey := strconv.Itoa(code) + strconv.Itoa(length) + ""
			result.SeenMutex.Lock()
			if _, ok := result.Seen[seenKey]; ok {
				result.ResultUrl[i].Url = ""
				result.SeenMutex.Unlock()
				return
			}
			result.Seen[seenKey] = struct{}{}
			result.SeenMutex.Unlock()

			result.ResultUrl[i] = mode.Link{Url: u, Status: strconv.Itoa(code), Size: strconv.Itoa(length), Redirect: redirect}
		}
	} else {
		result.ResultUrl[i].Url = ""
	}
}