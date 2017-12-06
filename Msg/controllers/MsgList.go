// MsgList
// 使用方法：/msgList/:uid/:mtype/:type/:page
package controllers

import (
	. "Msg/libs"
	. "Msg/models"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func MsgList(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	runtime.Gosched()
	uid := c.Param("uid")
	msgMType := c.Param("mtype")
	//mTypeNum, _ := strconv.Atoi(msgMType)
	msgType := c.Param("type")
	//typeNum, _ := strconv.Atoi(msgType)

	if uid != "" && msgMType != "" {
		msgs := make([]MsgForJobs, 0)
		var err error
		pageENum := 50
		page := c.Param("page")
		pageNum, _ := strconv.Atoi(page)
		if pageNum == 0 {
			pageNum = 1
		}

		/**
		switch mTypeNum {
		case 1:
			//公告
			err = Engine.Table("mc_msg").
				Join("LEFT", "mc_msg_content", "mc_msg_content.mid = mc_msg.id").
				Join("LEFT", "mc_msg_r_user", "mc_msg_r_user.mid = mc_msg.id").
				Where("mc_msg.status = 1 and mc_msg.mtype = ? and mc_msg_r_user.uid <> ?", mTypeNum, uid).
				OrderBy("mc_msg_r_user.is_read asc, mc_msg.created_at desc").
				Limit(pageENum, (pageNum-1)*pageENum).
				Find(&msgs)
		case 2:
			//审批
			err = Engine.Table("mc_msg").
				Join("LEFT", "mc_msg_content", "mc_msg_content.mid = mc_msg.id").
				Join("LEFT", "mc_msg_r_user", "mc_msg_r_user.mid = mc_msg.id").
				Join("LEFT", "mc_oa_workjob", "mc_oa_workjob.job_id = mc_msg.job_id").
				Where("mc_msg.status = 1 and mc_msg.mtype = ? and mc_msg.type = ? and mc_msg_r_user.uid = ?", mTypeNum, typeNum, uid).
				OrderBy("mc_msg_r_user.is_read asc, mc_msg.created_at desc").
				Limit(pageENum, (pageNum-1)*pageENum).
				Find(&msgs)
		}
		**/

		mTypeArr := strings.Split(msgMType, ",")
		typeArr := strings.Split(msgType, ",")
		if pageNum == -1 && uid == "-1" {
			err = Engine.Table("mc_msg").
				Join("LEFT", "mc_msg_content", "mc_msg_content.mid = mc_msg.id").
				Join("LEFT", "mc_msg_r_user", "mc_msg_r_user.mid = mc_msg.id").
				Join("LEFT", "mc_oa_workjob", "mc_oa_workjob.job_id = mc_msg.job_id").
				Where("mc_msg.status = 1 ").
				In("mc_msg.mtype", mTypeArr).
				In("mc_msg.type", typeArr).
				OrderBy("mc_msg_r_user.is_read asc, mc_msg.created_at desc").
				Find(&msgs)
		} else {
			err = Engine.Table("mc_msg").
				Join("LEFT", "mc_msg_content", "mc_msg_content.mid = mc_msg.id").
				Join("LEFT", "mc_msg_r_user", "mc_msg_r_user.mid = mc_msg.id").
				Join("LEFT", "mc_oa_workjob", "mc_oa_workjob.job_id = mc_msg.job_id").
				Where("mc_msg.status = 1 and mc_msg_r_user.uid = ?", uid).
				In("mc_msg.mtype", mTypeArr).
				In("mc_msg.type", typeArr).
				OrderBy("mc_msg_r_user.is_read asc, mc_msg.created_at desc").
				Limit(pageENum, (pageNum-1)*pageENum).
				Find(&msgs)
		}

		if err != nil {
			returnMsg(c, "error", err.Error())
		}

		//根据标题和链接过滤重复消息
		filterMap := make([]string, 0)
		newMsgs := make([]MsgForJobs, 0)
		for _, d := range msgs {
			uidStr := strconv.Itoa(d.MsgRUser.Uid)
			fKey := uidStr + d.MsgContent.Title + d.MsgContent.Uri
			if filterMsg(fKey, filterMap) != true {
				filterMap = append(filterMap, fKey)
				newMsgs = append(newMsgs, d)
			}
		}

		//相同消息同时有审批和通知时过滤通知消息，避免重复
		for _, d := range newMsgs {
			if strings.HasSuffix(d.MsgContent.Uri, "/edit") == true && d.MsgRUser.IsRead == 0 {
				var tmpUri = strings.Replace(d.MsgContent.Uri, "/edit", "", strings.LastIndex(d.MsgContent.Uri, "/edit"))
				for k, data := range newMsgs {
					if data.MsgContent.Uri == tmpUri && data.MsgRUser.IsRead == 0 {
						index := k
						newMsgs = append(newMsgs[:index], newMsgs[index+1:]...)
					}
				}
			}
		}
		c.JSON(200, newMsgs)
	}
}

//验证数组是否包含字符串
func filterMsg(fKey string, filterMap []string) bool {
	for _, k := range filterMap {
		if fKey == k {
			return true
		}
	}
	return false
}
