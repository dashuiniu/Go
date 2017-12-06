// AddMsg
//PHP使用方法
//$url = "http://localhost:8088/addMsg";
//$post = array (
//   'type' => '1',
//   'recipient' => '1,2,3',
//	 'title' => 'title1111',
//	 'content' => 'contentcontentcontent11'
//);
//$res = post($url, $post);
package controllers

import (
	. "Msg/libs"
	. "Msg/models"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AddMsg(c *gin.Context) {
	runtime.Gosched()
	typeMsg := c.PostForm("type")
	mTypeMsg := c.PostForm("mtype")
	jTypeMsg := c.PostForm("jtype")
	recipient := c.PostForm("recipient")
	title := c.PostForm("title")
	content := c.PostForm("content")
	uri := c.PostForm("uri")
	uid := c.PostForm("uid")
	jid := c.PostForm("job_id")
	mTypeNum, _ := strconv.Atoi(mTypeMsg)
	jTypeNum, _ := strconv.Atoi(jTypeMsg)
	typeNum, _ := strconv.Atoi(typeMsg)
	uidNum, _ := strconv.Atoi(uid)
	jidNum, _ := strconv.Atoi(jid)
	users := strings.Split(recipient, ",")

	if typeMsg != "" && mTypeMsg != "" && title != "" && content != "" {

		session := Engine.NewSession()
		defer session.Close()
		err := session.Begin()

		//判断是否已经有该消息
		msgs := make([]Msgs, 0)
		err = session.Table(new(Msgs).TableName()).
			Join("LEFT", "mc_msg_content", "mc_msg_content.mid = mc_msg.id").
			Join("LEFT", "mc_msg_r_user", "mc_msg_r_user.mid = mc_msg.id").
			Where("mc_msg.status = 1 and mc_msg.mtype = ? and mc_msg.job_id = ? and mc_msg_content.title = ? and mc_msg_content.content = ? and mc_msg_content.uri= ? ",
			mTypeNum, jidNum, title, content, uri).
			Find(&msgs)
		if err != nil {
			session.Rollback()
			returnMsg(c, "error", err.Error())
		}

		if len(msgs) > 0 {
			//有主消息时更新发送到每个人
			has := true
			midNum := msgs[0].Msg.Id
			for _, u := range users {
				user, _ := strconv.Atoi(u)
				if user == 0 {
					continue
				}
				msgRUsers := new(MsgRUser)
				has, err = session.Where("mid = ? and uid = ?", midNum, user).Get(msgRUsers)
				if has == true {
					//					msgRUser := MsgRUser{
					//						IsRead:    0,
					//						Status:    1,
					//						Updater:   user,
					//						UpdatedAt: time.Now(),
					//					}
					//					_, err = session.Where("mid = ? and uid = ?", midNum, user).Update(msgRUser)
					_, err = session.Exec("UPDATE mc_msg_r_user SET is_read = ?, status = ?, updater = ?, updated_at = ? WHERE mid = ? and uid = ?",
						0, 1, user, time.Now(), midNum, user)
					if err != nil {
						session.Rollback()
						returnMsg(c, "error", err.Error())
					}
				} else {
					msgRUser := MsgRUser{
						Mid:       midNum,
						Uid:       user,
						IsRead:    0,
						Status:    1,
						Creater:   uidNum,
						CreatedAt: time.Now(),
						Updater:   uidNum,
						UpdatedAt: time.Now(),
					}
					_, err = session.Insert(msgRUser)
					if err != nil {
						session.Rollback()
						returnMsg(c, "error", err.Error())
					}
				}
			}
		} else {
			//主表
			msg := Msg{
				Mtype: mTypeNum,
				Jtype: jTypeNum,
				Type:  typeNum,
				//			Recipient: recipient,
				JobId:     jidNum,
				Status:    1,
				Creater:   uidNum,
				CreatedAt: time.Now(),
				Updater:   uidNum,
				UpdatedAt: time.Now(),
			}
			_, err = session.Insert(&msg)
			if err != nil {
				session.Rollback()
				returnMsg(c, "error", err.Error())
			}
			//内容
			msgContent := MsgContent{
				Mid:       msg.Id,
				Title:     title,
				Content:   content,
				Uri:       uri,
				Status:    1,
				Creater:   uidNum,
				CreatedAt: time.Now(),
				Updater:   uidNum,
				UpdatedAt: time.Now(),
			}
			_, err = session.Insert(msgContent)
			if err != nil {
				session.Rollback()
				returnMsg(c, "error", err.Error())
			}

			//消息对发送对象
			for _, u := range users {
				user, _ := strconv.Atoi(u)
				if user == 0 {
					continue
				}
				msgRUser := MsgRUser{
					Mid:       msg.Id,
					Uid:       user,
					IsRead:    0,
					Status:    1,
					Creater:   uidNum,
					CreatedAt: time.Now(),
					Updater:   uidNum,
					UpdatedAt: time.Now(),
				}
				_, err = session.Insert(msgRUser)
				if err != nil {
					session.Rollback()
					returnMsg(c, "error", err.Error())
				}
			}
		}

		err = session.Commit()
		if err != nil {
			returnMsg(c, "error", err.Error())
		}
		returnMsg(c, "success", "true")
	} else {
		returnMsg(c, "error", "参数个数不足")
	}

}

// 返回消息
func returnMsg(c *gin.Context, title string, contents string) {
	c.JSON(200, gin.H{
		title: contents,
	})
}
