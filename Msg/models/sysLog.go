// sysLog
package models

import (
	"time"
)

type SysLog struct {
	Id        int64
	Ctable    string `xorm:"varchar(40)"`
	Cfield    string `xorm:"varchar(50)"`
	Cline     int
	Odata     string
	Ndata     string
	Ntype     string `xorm:"varchar(20)"`
	Njob      int
	Nnode     string `xorm:"varchar(200)"`
	Nnote     string
	Note      string
	Uid       int
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

func (SysLog) TableName() string {
	return "mc_sys_log"
}
