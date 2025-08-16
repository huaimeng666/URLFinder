package crawler

import (
	"compress/gzip"
	"fmt"
	"github.com/pingc0y/URLFinder/cmd"
	"github.com/pingc0y/URLFinder/config"
	"github.com/pingc0y/URLFinder/result"
	"github.com/pingc0y/URLFinder/util"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// 蜘蛛抓取页面内容
func Spider(u string, num int) {
	is := true
	defer func() {
		config.Wg.Done()
		if is {
			<-config.Ch
		}

	}()
	config.Mux.Lock()
	fmt.Printf("\rStart %d Spider...", config.Progress)
	config.Progress++
	config.Mux.Unlock()
	//标记完成

	u, _ = url.QueryUnescape(u)
	if num > 1 && cmd.D != "" && !regexp.MustCompile(cmd.D).MatchString(u) {
		return
	}
	if GetEndUrl(u) {
		return
	}
	if cmd.M == 3 {
		for _, v := range config.Risks {
			if strings.Contains(u, v) {
				return
			}
		}
	}
	AppendEndUrl(u)
	request, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return
	}

	request.Header.Set("Accept-Encoding", "gzip") //使用gzip压缩传输数据让访问更快
	request.Header.Set("User-Agent", util.GetUserAgent())
	request.Header.Set("Accept", "*/*")
	// 增加header选项
	if cmd.C != "" {
		request.Header.Set("Cookie", cmd.C)
	}
	// 添加Referer
	config.Lock.Lock()
	if referer, ok := result.Jstourl[u]; ok {
		request.Header.Set("Referer", referer)
	} else if referer, ok := result.Urltourl[u]; ok {
		request.Header.Set("Referer", referer)
	}
	config.Lock.Unlock()
	// 加载yaml配置(headers)
	if cmd.I {
		util.SetHeadersConfig(&request.Header)
	}

	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	var resultBody string
	//解压
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(response.Body) // gzip解压缩
		if err != nil {
			return
		}
		defer reader.Close()
		con, err := io.ReadAll(reader)
		if err != nil {
			return
		}
		resultBody = string(con)
	} else {
		//提取url用于拼接其他url或js
		dataBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return
		}
		//字节数组 转换成 字符串
		resultBody = string(dataBytes)
	}

	//处理base标签
	var baseHref string
	baseRegex := regexp.MustCompile(`(?i)<base\s+href\s*=\s*["']([^"']+)["']`)
	baseMatch := baseRegex.FindStringSubmatch(resultBody)
	if len(baseMatch) > 1 {
		baseHref = baseMatch[1]
	} else {
		baseVarRegex := regexp.MustCompile(`(?i)(?:base|baseUrl|basePath)\s*[:=]\s*["']([^"']+)["']`)
		baseVarMatch := baseVarRegex.FindStringSubmatch(resultBody)
		if len(baseVarMatch) > 1 {
			baseHref = baseVarMatch[1]
		}
	}

	baseURL, err := url.Parse(u)
	if err != nil {
		return
	}

	if baseHref != "" {
		base, err := url.Parse(baseHref)
		if err == nil {
			baseURL = baseURL.ResolveReference(base)
		}
	}

	path := baseURL.Path
	host := baseURL.Host
	scheme := baseURL.Scheme
	source := scheme + "://" + host + path
	is = false
	<-config.Ch
	//提取js
	jsFind(resultBody, host, scheme, path, u, num)
	//提取url
	urlFind(resultBody, host, scheme, path, u, num)
	//提取信息
	infoFind(resultBody, source)

}

// 打印Validate进度
func PrintProgress() {
	config.Mux.Lock()
	num := len(result.ResultJs) + len(result.ResultUrl)
	fmt.Printf("\rValidate %.0f%%", float64(config.Progress+1)/float64(num+1)*100)
	config.Progress++
	config.Mux.Unlock()
}