// MsgInfo
// 阅读消息
// 使用方法：/msgInfo/:uid/:mid
package controllers

import (
	. "Msg/libs"
	. "Msg/models"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func MsgInfo(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	runtime.Gosched()
	uid := c.Param("uid")
	mid := c.Param("mid")
	uidNum, _ := strconv.Atoi(uid)
	midNum, _ := strconv.ParseInt(mid, 10, 64)

	if uid != "" && mid != "" {
		session := Engine.NewSession()
		defer session.Close()
		err := session.Begin()

		msgs := make([]Msgs, 0)
		err = session.Table(new(Msgs).TableName()).
			Join("LEFT", "mc_msg_content", "mc_msg_content.mid = mc_msg.id").
			Join("LEFT", "mc_msg_r_user", "mc_msg_r_user.mid = mc_msg.id").
			Where("mc_msg.status = 1 and mc_msg.id = ?", midNum).
			Find(&msgs)
		if err != nil {
			session.Rollback()
			returnMsg(c, "error", err.Error())
		}

		//回写查看状态和时间
		msgRUser := MsgRUser{
			IsRead:    1,
			Rtime:     time.Now(),
			Updater:   uidNum,
			UpdatedAt: time.Now(),
		}
		_, err = session.Where("mid = ? and uid = ? and is_read = 0", midNum, uidNum).Update(&msgRUser)
		if err != nil {
			session.Rollback()
			returnMsg(c, "error", err.Error())
		}

		//添加日志
		//		log := make([]SysLog, 0)
		//		err = session.Table(new(SysLog).TableName()).
		//			Where("ntype = ? and njob = ? and uid = ?", "mcMsg", midNum, uidNum).
		//			Find(&log)
		//		if err != nil {
		//			session.Rollback()
		//			returnMsg(c, "error", err.Error())
		//		}
		//		if len(log) == 0 {
		midNumF, _ := strconv.Atoi(mid)
		sysLog := SysLog{
			Ntype:     "mcMsg",
			Njob:      midNumF,
			Note:      "审阅流程(" + mid + ")时间：" + time.Now().Format("2006-01-02 15:04:05"),
			Uid:       uidNum,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, err = session.Insert(&sysLog)
		if err != nil {
			session.Rollback()
			returnMsg(c, "error", err.Error())
		}
		//		}

		err = session.Commit()
		if err != nil {
			returnMsg(c, "error", err.Error())
		} else {
			c.JSON(200, msgs)
		}

	} else {
		returnMsg(c, "error", "参数个数不足")
	}
}
