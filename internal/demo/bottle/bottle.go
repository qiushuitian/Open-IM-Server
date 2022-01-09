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
