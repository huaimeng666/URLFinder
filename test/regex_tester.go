package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	// 1. 在这里放入你的测试文本
	text := `
      // 示例1: 标准阿里云/腾讯云/华为云密钥
      accessKey: D25RTPTV9V0U0BQ64Y3O
      secretKey: YAsYqDAu92WBONouBkiNFOsJ82Zx185K6s3xsHqb

      // 示例2: Base64编码风格的长密钥 (长度 > 40)
      ai-model:
        volcano:
          accessKey: "AKLTNTI1ZGYzODdjMzZmNDE2OGI5YzgzNjlhZjdmZWJiNzk"
          secretKey: "TW1JMk5EYzNaRFprTmpZNE5HVTVOamc1T0RZNVlXWmtZV1UwTmpJNFpUSQ=="
      
      // 示例3: Meitu密钥
      meitu:
        accessKey: "6ca542de9492407f96bcc27f82791d81"
        secretKey: "50944fef10234696a8d385f3f9702360"

      // 示例4: 其他各类信息
      email: user123@example.com
      phone: 13800138000
      id_card: 110101199001011234
      ip_address: 192.168.1.1
      jwt_token: "eyJh...<snip>...SflKxw"
      uuid: "123e4567-e89b-12d3-a456-426614174000"
      jdbc_url: "jdbc:mysql://localhost:3306/mydb?user=root&password=password123"
      grafana_key: "glc_aBcDeFgHiJkLmNoPqRsTuVwXyZ1234567890aBcD="
			webhook_url: "https://oapi.dingtalk.com/robot/send?access_token=a1b2c3d4e5f6..."
			swagger_page: "swagger-ui.html"
			other_pass: "password: mysecretpassword"
    `

	// 2. 在这里定义所有可能用到的正则表达式
	// 格式与 config.go 完全一致，方便复制粘贴
	var (
		Phone     = []string{`[^\w]((?:(?:\+|00)86)?1(?:(?:3[\d])|(?:4[5-79])|(?:5[0-35-9])|(?:6[5-7])|(?:7[0-8])|(?:8[\d])|(?:9[189]))\d{8})[^\w]`}
		Email     = []string{`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`}
		IDcard    = []string{`[^0-9]((\d{8}(0\d|10|11|12)([0-2]\d|30|31)\d{3}$)|(\d{6}(18|19|20)\d{2}(0[1-9]|10|11|12)([0-2]\d|30|31)\d{3}(\d|X|x)))[^0-9]`}
		Jwt       = []string{`'"` + "`" + `(ey[A-Za-z0-9_-]{10,}\.[A-Za-z0-9._-]{10,}|ey[A-Za-z0-9_\/+-]{10,}\.[A-Za-z0-9._\/+-]{10,})` + "`" + `'"`}
		UUIDToken = []string{`'"` + "`" + `?([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})` + "`" + `?"`}
		// 注意: AKSK 正则已更新，长度限制放宽到 {20,80}
		AKSK      = []string{`(?i)["']{0,1}(secret|access|security|"ak"|"sk"|'ak'|'sk'|credential|api_key|client_secret|private|access_key|secret_access)[._ ]{0,3}(?i)(Id|Key|token){0,1}["']{0,1}:\s*['"]{0,1}([A-Za-z0-9+/]{20,80}={0,2})['"]{0,1}`}
		TheKey    = []string{`((?i)enc.Utf8.parse|(?i)x-secret-id|(?i)ACCESS KEY SECRET|(?i)headerToSign)`}
		Other     = []string{`(access.{0,1}key|access.{0,1}Key|access.{0,1}Id|access.{0,1}id|.{0,5}密码|.{0,5}账号|默认.{0,5}|加密|解密|password:.{0,10}|username:.{0,10})`}
		Webhook   = []string{`\bhttps://qyapi.weixin.qq.com/cgi-bin/webhook/send\?key=[a-zA-Z0-9\-]{25,50}\b`, `\bhttps://oapi.dingtalk.com/robot/send\?access_token=[a-z0-9]{50,80}\b`, `\bhttps://open.feishu.cn/open-apis/bot/v2/hook/[a-z0-9\-]{25,50}\b`, `\bhttps://hooks.slack.com/services/[a-zA-Z0-9\-_]{6,12}/[a-zA-Z0-9\-_]{6,12}/[a-zA-Z0-9\-_]{15,24}\b`}
		Grafana   = []string{`\bglc_[A-Za-z0-9\-_+/]{32,200}={0,2}\b`, `\bglsa_[A-Za-z0-9]{32}_[A-Fa-f0-9]{8}\b`}
		Ip        = []string{`\b((?:(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9]))\b`}
		Swaggerui = []string{`((swagger-ui.html)|("swagger":)|(Swagger UI)|(swaggerUi)|(swaggerVersion))`}
		Jdbc      = []string{`(jdbc:[a-z:]+://[a-z0-9\.\-_:;=/@?,&]+)`}
	)

	// 将独立的正则变量组装成 map
	Infofind := map[string][]string{
		"Phone":     Phone,
		"Email":     Email,
		"IDcard":    IDcard,
		"Jwt":       Jwt,
		"UUIDToken": UUIDToken,
		"AKSK":      AKSK,
		"TheKey":    TheKey,
		"Other":     Other,
		"Webhook":   Webhook,
		"Grafana":   Grafana,
		"Ip":        Ip,
		"Swaggerui": Swaggerui,
		"Jdbc":      Jdbc,
	}

	// 3. 根据命令行参数决定要测试哪些正则
	var targetCategories []string
	rawArgs := os.Args[1:]
	for _, arg := range rawArgs {
		parts := strings.Split(arg, ",")
		for _, part := range parts {
			if trimmedPart := strings.TrimSpace(part); trimmedPart != "" {
				targetCategories = append(targetCategories, trimmedPart)
			}
		}
	}

	infoToTest := make(map[string][]string)

	if len(targetCategories) > 0 {
		fmt.Println("--- 将只测试以下指定分类:", targetCategories, "---")
		for _, key := range targetCategories {
			if regexps, ok := Infofind[key]; ok {
				infoToTest[key] = regexps
			} else {
				fmt.Printf("警告: 分类 '%s' 未在代码中定义，将被跳过。\n", key)
			}
		}
	} else {
		fmt.Println("--- 未指定分类，将测试所有已定义的正则 ---")
		infoToTest = Infofind
	}

	// 4. 执行匹配和打印结果的逻辑
	fmt.Println("--- 开始正则匹配测试 ---")

	for key, regexps := range infoToTest {
		fmt.Printf("\n--- 测试分类: %s ---\n", key)
		categoryFound := false
		for i, regexpstr := range regexps {
			re, err := regexp.Compile(regexpstr)
			if err != nil {
				fmt.Printf("  正则 #%d 编译失败: %v\n", i+1, err)
				continue
			}

			allMatches := re.FindAllStringSubmatch(text, -1)

			if len(allMatches) > 0 {
				categoryFound = true
				fmt.Printf("  使用正则 #%d: `%s`\n", i+1, regexpstr)
				for _, match := range allMatches {
					var matchedValue string
					if len(match) > 1 {
						matchedValue = match[len(match)-1]
					} else {
						matchedValue = match[0]
					}
					fmt.Printf("    > 匹配到: %s\n", matchedValue)
				}
			}
		}
		if !categoryFound {
			fmt.Println("  未匹配到任何内容。")
		}
	}
	fmt.Println("\n--- 正则匹配测试结束 ---")
}
