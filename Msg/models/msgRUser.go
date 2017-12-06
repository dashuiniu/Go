// msg
package models

import (
	"time"
)

type MsgRUser struct {
	Id        int64
	Mid       int64 `xorm:"notnull index"`
	Uid       int
	IsRead    int       `xorm:"null default 0"`
	Rtime     time.Time `xorm:"null default null"`
	Status    int
	Creater   int
	CreatedAt time.Time `xorm:"created"`
	Updater   int
	UpdatedAt time.Time `xorm:"updated"`
}
