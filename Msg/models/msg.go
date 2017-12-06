// msg
package models

import (
	"time"
)

type Msg struct {
	Id        int64
	Mtype     int `xorm:"notnull"`
	Type      int `xorm:"notnull"`
	Jtype     int `xorm:"null"`
	Recipient string
	JobId     int `xorm:"null"`
	Status    int
	Creater   int
	CreatedAt time.Time `xorm:"created"`
	Updater   int
	UpdatedAt time.Time `xorm:"updated"`
}

type Msgs struct {
	Msg        Msg        `xorm:"extends"`
	MsgContent MsgContent `xorm:"extends"`
	MsgRUser   MsgRUser   `xorm:"extends"`
}

type MsgForJobs struct {
	Msg        Msg        `xorm:"extends"`
	MsgContent MsgContent `xorm:"extends"`
	MsgRUser   MsgRUser   `xorm:"extends"`
	OaWorkjob  OaWorkjob  `xorm:"extends"`
}

func (Msgs) TableName() string {
	return "mc_msg"
}
