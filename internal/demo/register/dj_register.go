package register

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ParamsRegister struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Platform int32  `json:"platform"`
	DeviceId string `json:"DeviceId"`
}

func Register(c *gin.Context) {
	log.NewDebug("RegisterWithName api is statrting...")

	params := ParamsRegister{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}

	var account = params.UserName

	queryParams := im_mysql_model.GetRegisterParams{
		Account: account,
	}
	_, err, rowsAffected := im_mysql_model.GetRegister(&queryParams)

	if err == nil && rowsAffected != 0 {
		log.ErrorByKv("The user name has been registered", queryParams.Account, "err")
		c.JSON(http.StatusOK, gin.H{"errCode": constant.LogicalError, "errMsg": "The user name has been registered"})
		return
	}

	log.InfoByKv("openIM register begin", account)
	resp, err := OpenIMRegister(account)

	log.InfoByKv("openIM register resp", account, resp, err)
	if err != nil {
		log.ErrorByKv("request openIM register error", account, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.HttpError, "errMsg": err.Error()})
		return
	}
	response, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errCode": constant.IoErrot, "errMsg": err.Error()})
		return
	}
	imrep := IMRegisterResp{}
	err = json.Unmarshal(response, &imrep)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errCode": constant.FormattingError, "errMsg": err.Error()})
		return
	}
	if imrep.ErrCode != 0 {
		c.JSON(http.StatusOK, gin.H{"errCode": constant.HttpError, "errMsg": imrep.ErrMsg})
		return
	}

	setQueryParams := im_mysql_model.SetPasswordParams{
		Account:  account,
		Password: params.Password,
	}

	log.InfoByKv("begin store mysql", account, params.Password)
	_, err = im_mysql_model.SetPassword(&setQueryParams)
	if err != nil {
		log.ErrorByKv("set phone number password error", account, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": constant.DatabaseError, "errMsg": err.Error()})
		return
	}

	log.InfoByKv("end setPassword", account)
	c.JSON(http.StatusOK, gin.H{"errCode": constant.NoError, "errMsg": "", "data": imrep.Data})

}
