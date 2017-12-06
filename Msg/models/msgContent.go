// msg
package models

import (
	"time"
)

type MsgContent struct {
	Id        int64
	Mid       int64 `xorm:"notnull index"`
	Title     string
	Content   string `xorm:"varchar(1024)"`
	Uri       string `xorm:"varchar(200)"`
	Status    int
	Creater   int
	CreatedAt time.Time `xorm:"created"`
	Updater   int
	UpdatedAt time.Time `xorm:"updated"`
}
