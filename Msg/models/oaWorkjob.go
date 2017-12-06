// msg
package models

import (
	"time"
)

type OaWorkjob struct {
	JobId     int64 `xorm:"notnull pk"`
	FlowId    int   `xorm:"notnull index"`
	FlowName  string
	AppName   string
	CopyTo    string
	Status    int
	Urgency   string
	FlowSysId int
	Creater   int
	CreatedAt time.Time `xorm:"created"`
	Updater   int
	UpdatedAt time.Time `xorm:"updated"`
}
