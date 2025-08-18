package crawler

import (
	"github.com/huaimeng666/URLFinder/config"
	"github.com/huaimeng666/URLFinder/mode"
	"github.com/huaimeng666/URLFinder/result"
	"github.com/huaimeng666/URLFinder/util"
	"regexp"
)

var (
	// 修改：预编译js路径正则
	jsPathRegex = regexp.MustCompile("(.+/)[^/]+.js")
	// 修改：预编译主机正则
	hostRegex = regexp.MustCompile("(https{0,1}://([a-z0-9\\-]+\\.)*([a-z0-9\\-]+\\.[a-z0-9\\-]+)(:[0-9]+)?/)")
)

func JsFuzz() {

	paths := []string{}
	for i := range result.ResultJs {
		// 修改：使用预编译的正则对象
		re := jsPathRegex.FindAllStringSubmatch(result.ResultJs[i].Url, -1)
		if len(re) != 0 {
			paths = append(paths, re[0][1])
		}
		// 修改：使用预编译的正则对象
		re2 := hostRegex.FindAllStringSubmatch(result.ResultJs[i].Url, -1)
		if len(re2) != 0 {
			paths = append(paths, re2[0][1])
		}
	}
	paths = util.UniqueArr(paths)
	for i := range paths {
		for i2 := range config.JsFuzzPath {
			result.ResultJs = append(result.ResultJs, mode.Link{
				Url:    paths[i] + config.JsFuzzPath[i2],
				Source: "Fuzz",
			})
		}
	}
}