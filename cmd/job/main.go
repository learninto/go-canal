package job

import (
	"context"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/learninto/go-canal/pkg/oslib"
	"github.com/learninto/goutil/conf"
	"github.com/siddontang/go-log/log"
	"math/rand"
	"time"
)

type MyEventHandler struct {
	canal.DummyEventHandler
}

// OnRow 监听数据记录
func (h *MyEventHandler) OnRow(ev *canal.RowsEvent) error {
	/*//record := fmt.Sprintf("%s %v %v %v %s\n",e.Action,e.Rows,e.Header,e.Table,e.String())

	//库名，表名，行为，数据记录
	record := fmt.Sprintf("%v %v %s %v\n", ev.Table.Schema, ev.Table.Name, ev.Action, ev.Rows)
	fmt.Println(record)

	//此处是参考 https://github.com/gitstliu/MysqlToAll 里面的获取字段和值的方法
	for columnIndex, currColumn := range ev.Table.Columns {
		//字段名，字段的索引顺序，字段对应的值
		row := fmt.Sprintf("%v %v %v\n", currColumn.Name, columnIndex, ev.Rows[len(ev.Rows)-1][columnIndex])
		fmt.Println("row info:", row)
	}*/

	m := make(map[string]interface{})
	for columnIndex, currColumn := range ev.Table.Columns {
		m[currColumn.Name] = ev.Rows[len(ev.Rows)-1][columnIndex]
	}

	return nil
}

// OnTableChanged 创建、更改、重命名或删除表时触发，通常会需要清除与表相关的数据，如缓存。It will be called before OnDDL.
func (h *MyEventHandler) OnTableChanged(e *replication.EventHeader, schema string, table string) error {
	/*//库，表
	record := fmt.Sprintf("%s %s \n", schema, table)
	fmt.Println(record)*/
	return nil
}

// OnPosSynced 监听binlog日志的变化文件与记录的位置
func (h *MyEventHandler) OnPosSynced(e *replication.EventHeader, pos mysql.Position, set mysql.GTIDSet, force bool) error {
	/*//源码：当force为true，立即同步位置
	record := fmt.Sprintf("%v %v \n", pos.Name, pos.Pos)
	fmt.Println("OnPosSynced", record)*/
	return nil
}

// OnRotate 当产生新的binlog日志后触发(在达到内存的使用限制后（默认为 1GB），会开启另一个文件，每个新文件的名称后都会有一个增量。)
func (h *MyEventHandler) OnRotate(e *replication.EventHeader, r *replication.RotateEvent) error {
	/*//record := fmt.Sprintf("On Rotate: %v \n",&mysql.Position{Name: string(r.NextLogName), Pos: uint32(r.Position)})
	//binlog的记录位置，新binlog的文件名
	record := fmt.Sprintf("On Rotate %v %v \n", r.Position, r.NextLogName)
	fmt.Println(record)*/
	return nil
}

// OnDDL create alter drop truncate(删除当前表再新建一个一模一样的表结构)
func (h *MyEventHandler) OnDDL(e *replication.EventHeader, nextPos mysql.Position, queryEvent *replication.QueryEvent) error {
	/*//binlog日志的变化文件与记录的位置
	record := fmt.Sprintf("%v %v\n", nextPos.Name, nextPos.Pos)
	query_event := fmt.Sprintf("%v\n %v\n %v\n %v\n %v\n",
		queryEvent.ExecutionTime,         //猜是执行时间，但测试显示0
		string(queryEvent.Schema),        //库名
		string(queryEvent.Query),         //变更的sql语句
		string(queryEvent.StatusVars[:]), //测试显示乱码
		queryEvent.SlaveProxyID)          //从库代理ID？
	fmt.Println("OnDDL:", record, query_event)*/
	return nil
}

func (h *MyEventHandler) String() string {
	return "MyEventHandler"
}

func run(ctx context.Context) (err error) {
	// 获取conf路径
	confPath, err := oslib.GetConfPath(ctx)
	if err != nil {
		log.Error(ctx, "os.Getwd Error：", err)
		return
	}

	// 读取toml文件格式
	cfg, err := canal.NewConfigWithFile(confPath + "/sniper.toml")
	if err != nil {
		log.Error(ctx, "canal.NewConfigWithFile error", err)
		return
	}
	cfg.ServerID = uint32(rand.New(rand.NewSource(time.Now().Unix())).Intn(1000)) + 1001

	// 获取监听
	c, err := canal.NewCanal(cfg)
	if err != nil {
		log.Error(ctx, "canal.NewCanal error", err)
		return
	}

	// Register a handler to handle RowsEvent
	c.SetEventHandler(&MyEventHandler{})

	binName := conf.Get("mysql.position.name")
	binPos := conf.GetInt("mysql.position.pos")

	log.Info(ctx, "Go run")

	if binName == "" || binPos <= 0 { // 从头开始监听
		err = c.Run()
	} else { // 根据位置监听
		startPos := mysql.Position{Name: binName, Pos: uint32(binPos)} // mysql-bin.000004, 1027
		err = c.RunFrom(startPos)
	}

	if err != nil {
		log.Error(ctx, "Go run Error：", err)
	}

	return
}
