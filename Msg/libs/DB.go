// DB
package libs

import (
	"fmt"

	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
)

var Engine *xorm.Engine

func init() {
	//datasource := "user=postgres password=postgres@123 host=192.168.32.85 port=5432 dbname=maccura sslmode=disable"
	datasource := "user=postgres password=111111 host=localhost port=5432 dbname=maccura sslmode=disable"
	var err error
	Engine, err = xorm.NewEngine("postgres", datasource)
	if err != nil {
		fmt.Println("新建引擎错误：", err)
		return
	}
	//	if err := Engine.Ping(); err != nil {
	//		fmt.Println(err)
	//		return
	//	}

	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "mc_")
	Engine.SetTableMapper(tbMapper)
	//Engine.ShowSQL(true)
}
