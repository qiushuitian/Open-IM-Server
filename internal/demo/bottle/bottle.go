package bottle

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// var bottle_pool = make(map[queryKey]BottleInfo)
var bottle_pool sync.Map

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
