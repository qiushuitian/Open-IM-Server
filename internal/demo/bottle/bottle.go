package bottle

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yanyiwu/gojieba"
)

// var bottle_pool = make(map[queryKey]BottleInfo)
var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../../..")

	bottle_pool     sync.Map
	sensiveWordsMap map[string]string
	fenChi          *gojieba.Jieba
)

// var bottle_keys

type ParamsPickBottle struct {
	UserName string `json:"userName"`
}

type BottleInfo struct {
	bId        int64
	UserName   string `json:"userName"`
	BottleType string `json:"bottleType"`
	Text       string `json:"text"`
}

type ParamsCheckUpgrade struct {
	Version  string `json:"version"`
	Platform int    `json:"platform"`
	AppName  string `json:"appName"`
}

type UpgradeInfo struct {
	NeedForceUpdate        bool   `json:"needForceUpdate"`
	DownloadURL            string `json:"downloadURL"`
	BuildVersion           string `json:"buildVersion"`
	BuildUpdateDescription string `json:"buildUpdateDescription"`
}

type TermsInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Url     string `json:"url"`
	Content string `json:"content"`
}

type ParamsReport struct {
	UserName string `json:"userName"`
	Type     string `json:"type"`
	Content  string `json:"content"`
}

type ParamsFeedback struct {
	UserName string `json:"userName"`
	Type     string `json:"type"`
	Content  string `json:"content"`
}

type ParamsOnlineConfig struct {
	Platform int    `json:"platform"`
	Version  string `json:"version"`
}

func init() {
	sensiveWordsMap = readSensiveWords()
	fenChi = gojieba.NewJieba()

	for _, v := range sensiveWordsMap {
		fenChi.AddWord(v)
	}
}

func readSensiveWords() map[string]string {
	file, err := os.Open(Root + "/config/sensitive_words.txt")
	if err != nil {
		log.NewError("Cannot open text file: %s, err: [%v]", "sensitive_words.txt", err)
		return nil
	}
	defer file.Close()

	var swmap = make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		swmap[line] = line
	}
	return swmap
}

func hasMatchSensiveWords(str string) bool {
	var words = fenChi.CutForSearch(str, true)
	fmt.Println("分词结果:", strings.Join(words, "/"))
	for _, v := range words {
		if vIn, ok := sensiveWordsMap[v]; ok {
			fmt.Println("敏感词匹配词:", vIn)
			return true
		} else {
			continue
		}
	}
	return false
}

func CheckUpgrade(c *gin.Context) {
	log.NewDebug("CheckUpgrade api is statrting...")
	params := ParamsCheckUpgrade{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}

	log.InfoByKv("api CheckUpgrade get params", params.Version, params.Platform, params.AppName)

	upgradeInfo := &UpgradeInfo{
		NeedForceUpdate:        false,
		BuildVersion:           "1.1.0",
		DownloadURL:            "http://baidu.com",
		BuildUpdateDescription: "欢迎来到漂流瓶新版本"}
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": upgradeInfo})
}

func PickBottle(c *gin.Context) {
	log.NewDebug("PickBottle api is statrting...")
	params := ParamsPickBottle{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}

	log.InfoByKv("api PickBottle get params", params.UserName)

	// 先不做这么复杂，只是从map中取第一个非自己的瓶子
	bottleInfo := &BottleInfo{UserName: "_sys_default_bottle", Text: "欢迎来到漂流瓶"}
	bottle_pool.Range(func(key, value interface{}) bool {
		if value.(*BottleInfo).UserName != params.UserName {
			bottleInfo = value.(*BottleInfo)
			bottle_pool.Delete(key)
			return false
		}
		return true
	})

	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": bottleInfo})
}

func ThrowBottle(c *gin.Context) {
	log.NewDebug("ThrowBottle api is statrting...")
	params := BottleInfo{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}

	// 文本审核
	var hassensitive = hasMatchSensiveWords(params.Text)
	if hassensitive {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.ContentIllegal, "errMsg": "包含敏感内容。请遵守法律法规，文明聊天"})
		return

	}

	params.bId = time.Now().UnixMicro()
	bottle_pool.Store(params.bId, &params)

	data := make(map[string]interface{})
	data["msg"] = "扔出成功"

	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": data})

}

func Terms(c *gin.Context) {
	log.NewDebug("Terms api is statrting...")

	termsResp := make(map[string]*TermsInfo)
	termsResp["user_terms"] = &TermsInfo{
		Name:    "用户协议",
		Version: "1.0.0",
		Content: "",
		Url:     "https://www.baidu.com",
	}
	termsResp["privacy_terms"] = &TermsInfo{
		Name:    "隐私协议",
		Version: "1.0.0",
		Content: "",
		Url:     "https://www.sohu.com/",
	}

	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": termsResp})

}

// 举报
func Report(c *gin.Context) {
	log.NewDebug("Report api is statrting...")

	params := ParamsReport{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}

	// 记录下举报内容
	log.InfoByKv("bottle_Report:", params.Type, params.UserName, params.Content)

	data := make(map[string]interface{})
	data["msg"] = "举报成功"

	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": data})

}

// 意见反馈
// 内容
func Feedback(c *gin.Context) {
	log.NewDebug("Feedback api is statrting...")

	params := ParamsFeedback{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}

	// 记录下反馈内容
	log.InfoByKv("bottle_Feedback:", params.Type, params.UserName, params.Content)

	data := make(map[string]interface{})
	data["msg"] = "反馈成功"

	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": data})

}

// 在线参数
func OnlineConfig(c *gin.Context) {
	log.NewDebug("OnlineConfig api is statrting...")

	params := ParamsOnlineConfig{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}

	data := make(map[string]interface{})
	if params.Version == "1.2.0" && params.Platform == 1 { //1:iOS 2:Android
		data["ad_on"] = "off"
	} else {
		data["ad_on"] = "on"
	}

	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": data})

}
