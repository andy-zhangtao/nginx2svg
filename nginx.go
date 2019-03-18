package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/andy-zhangtao/gogather/tools/znginx"
)

var nginxinfo = []string{"nginx:"}
var nginxSplit = []string{"# configuration file"}

const nginxpath = "NGINX_COMMAND"

func getNginxConfigure() (configrue string, err error) {
	var command []string
	if os.Getenv(nginxpath) == "" {
		command = []string{
			"nginx",
			"-T",
		}
	} else {
		command = strings.Split(os.Getenv(nginxpath), " ")
		command = append(command, "-T")
	}

	cmd := exec.Command(command[0], command[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), errors.New(string(out))
	}

	return string(out), nil
}

func extractNginxMeta(content string) (nginx map[string]map[string]nginxMeta, err error) {
	// nginx 保存解析出的数据, key是server_name, value另外保存 location和dest的映射关系
	nginx = make(map[string]map[string]nginxMeta)

	// 剥离配置文件
	configrue, _ := parseNginxContent(content)

	if len(configrue) == 0 {
		err = errors.New("There no valid nginx data! ")
		return
	}

	// 对每个配置文件进行解析处理
	// 分析出server_name ， location 和 upsteam中的server
	for _, c := range configrue {
		locationmap := make(map[string]nginxMeta)
		server, location, upstream := extractServer(c)

		for _, l := range location {
			dest, isroot, loc := znginx.ExtractLocationDest(l)
			if strings.HasPrefix(dest, "http://") {
				dest = dest[7:]
			}
			if strings.HasPrefix(dest, "https://") {
				dest = dest[8:]
			}

			if isroot {
				// root 模式
				locationmap[loc] = nginxMeta{
					Dest: dest,
				}
			} else {
				for _, u := range upstream {
					if strings.Contains(u, "upstream") && strings.Contains(u, dest) {
						server, err := znginx.ExtractUpstreamValue(u)
						if err != nil {
							locationmap[loc] = nginxMeta{
								Dest: err.Error(),
							}
						} else {
							locationmap[loc] = nginxMeta{
								Dest: strings.Join(server, ";"),
							}
						}
						break
					}
				}
			}
		}

		for _, s := range server {
			nginx[s] = locationmap
		}
	}

	return
}

// extractServer 从Server片段里面抽取server name 和location数据
// 不处理type{}数据段
// server 抽取到的server_name
// location 抽取到的location数据
// upstream 抽取到的upstream数据
func extractServer(content string) (server []string, location []string, upstream []string) {
	if strings.HasPrefix(strings.TrimSpace(content), "type") {
		return
	}

	element := znginx.ExtractLocation(content)
	if len(element) > 0 {
		for key, value := range element {
			server = append(server, key)
			location = append(location, value...)
		}
	}

	upstream = znginx.ExtractUpstream(strings.TrimSpace(content))
	return
}

// parseNginxContent 解析Nginx配置文件并提取每个单独配置文件
// content 配置文件内容
// configrue 解析到的配置文件内容, otherInfo 其它需要关注的数据
func parseNginxContent(content string) (configrue []string, otherInfo []string) {

	dirtyConfigure := strings.Split(content, nginxSplit[0])

	for _, d := range dirtyConfigure {
		if isNginxInfo(strings.TrimSpace(d)) {
			otherInfo = append(otherInfo, d)
			continue
		}

		if strings.Contains(d, ":") {
			configrue = append(configrue, strings.SplitN(d, ":", 2)[1])
		}
	}

	return
}

// parseNginxContent 解析Nginx配置文件并提取每个单独配置文件
// content 配置文件内容
// configrue 解析到的配置文件内容, otherInfo 其它需要关注的数据
// func parseNginxContent(content string) (configrue []string, otherInfo []string) {

// 	inNginxConfigure := false
// 	nginxConfigrue := ""

// 	for _, s := range strings.Split(content, "\n") {

// 		if inNginxConfigure {
// 			nginxConfigrue += s
// 		} else {
// 			configrue = append(configrue, nginxConfigrue)
// 			nginxConfigrue = ""
// 		}

// 		if isNginxInfo(strings.TrimSpace(s)) {
// 			otherInfo = append(otherInfo, s)
// 			continue
// 		}

// 		if isNginxSplit(strings.TrimSpace(s)) {
// 			inNginxConfigure = !inNginxConfigure
// 		}
// 	}

// 	return
// }

// isNginxInfo 是否为Nginx自身数据
func isNginxInfo(content string) bool {
	for _, s := range nginxinfo {
		if strings.HasPrefix(content, s) {
			return true
		}
	}

	return false
}

// isNginxSplit 是否为Nginx分割线
func isNginxSplit(content string) bool {
	for _, s := range nginxSplit {
		if strings.HasPrefix(content, s) {
			return true
		}
	}

	return false
}

func extractServerConfigure(content string) (server string) {
	return
}
