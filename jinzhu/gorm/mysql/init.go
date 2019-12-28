package mysql

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"os"
	"time"
)

var MysqlPool *gorm.DB

var MysqlInitLogSlice []*string

var (
	User      = "root"
	Password  = "f9938689565b2ffc6b0d3f9db9487353"
	Host      = "127.0.0.1"
	Port      = 33067
	Database  = "siya"
	Collation = "utf8mb4"
)

// mysql 连接池 连接 mysql-server，需要设置 mysql 连接最大存活时间[ConnMaxLifetime]，要比 mysql-server thread_pool_idle_timeout 略小即可。这样可以避免 invalid connection 的错误。
// 一般情况下 mysql-server 默认 thread_pool_idle_timeout 为 60s。
// 查看 thread_pool_idle_timeout sql 语句 show variables like '%thread_pool_idle_timeout%';
// 问题参考 https://blog.letsgo.tech/gorm-go-mysql-driver-invalid-connection/
var (
	ConnMaxLifetime = time.Second * 59                             // mysql 最大连接存活时间
	MaxIdleConns    = 50                                           // mysql 最大活跃连接数
	MaxOpenConns    = 50                                           // mysql 最大 open 连接数
	mysqlLogFile    = "/data/logs/go-package-demo/mysql/mysql.log" // mysql 文件路径
)

type Processlist struct {
	ID      int    `gorm:"column:ID"`
	USER    string `gorm:"column:USER"`
	HOST    string `gorm:"column:HOST"`
	DB      string `gorm:"column:DB"`
	COMMAND string `gorm:"column:COMMAND"`
	TIME    int    `gorm:"column:TIME"`
	STATE   string `gorm:"column:STATE"`
	INFO    string `gorm:"column:INFO"`
}

type ThreadVariable struct {
	Value        string `gorm:"column:Value"`
	VariableName string `gorm:"column:Variable_name"`
}

func InitMysql() error {
	var err error
	if MysqlPool, err = getMysqlPool(); err != nil {
		return err
	}

	return nil
}

func getMysqlPool() (*gorm.DB, error) {
	//dsn := "root:root@tcp(127.0.0.1:33067)/siya?timeout=90s&collation=utf8mb4_general_ci"
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", User, Password, Host, Port, Database, Collation)

	log.Println(dns)

	pool, err := gorm.Open("mysql", dns)
	pool.DB().SetConnMaxLifetime(ConnMaxLifetime)
	pool.DB().SetMaxIdleConns(MaxIdleConns)
	pool.DB().SetMaxOpenConns(MaxOpenConns)

	MysqlInitLog(Database, pool)

	if err = logger(mysqlLogFile, MysqlInitLogSlice); err != nil {
		return nil, err
	}

	if err = pool.DB().Ping(); err != nil {
		return pool, errors.New("failed to connect database")
	}

	return pool, nil
}

func MysqlInitLog(database string, pool *gorm.DB) {


	logLinef(database,"SetConnMaxLifetime: %s", ConnMaxLifetime)
	logLinef(database,"SetMaxIdleConns: %d", MaxIdleConns)
	logLinef(database,"SetMaxOpenConns: %d", MaxOpenConns)

	// TODO 服务器环境编译不通过
	// 错误原因如下
	// db.Stats().MaxOpenConnections undefined (type sql.DBStats has no field or method MaxOpenConnections)
	// db.Stats().InUse undefined (type sql.DBStats has no field or method InUse)
	// db.Stats().Idle undefined (type sql.DBStats has no field or method Idle)
	// db.Stats().WaitCount undefined (type sql.DBStats has no field or method WaitCount)
	// db.Stats().WaitDuration undefined (type sql.DBStats has no field or method WaitDuration)
	// db.Stats().MaxIdleClosed undefined (type sql.DBStats has no field or method MaxIdleClosed)
	// db.Stats().MaxLifetimeClosed undefined (type sql.DBStats has no field or method MaxLifetimeClosed)
	// TODO 注释开始
	db := pool.DB()
	logLinef(database, "=========== db.Stats() ===========")
	logLinef(database, "MaxOpenConnections: %d [Maximum number of open connections to the database]", db.Stats().MaxOpenConnections)
	logLinef(database, "Pool Status:OpenConnections: %d [The number of established connections both in use and idle]", db.Stats().OpenConnections)
	logLinef(database, "Pool Status:InUse: %d [The number of connections currently in use]", db.Stats().InUse)
	logLinef(database, "Pool Status:Idle: %d [The number of idle connections]", db.Stats().Idle)
	logLinef(database, "Counters:WaitCount: %d [The total number of connections waited for]", db.Stats().WaitCount)
	logLinef(database, "Counters:WaitDuration: %d [The total time blocked waiting for a new connection]", db.Stats().WaitDuration)
	logLinef(database, "Counters:MaxIdleClosed: %d [The total number of connections closed due to SetMaxIdleConns]", db.Stats().MaxIdleClosed)
	logLinef(database, "Counters:MaxLifetimeClosed: %d [The total number of connections closed due to SetConnMaxLifetime]", db.Stats().MaxLifetimeClosed)
	// TODO 注释结束


	// log mysql-server processlist
	sql := fmt.Sprintf("select * from information_schema.PROCESSLIST where db = '%s';", Database)
	result := make([]Processlist, 0)
	if db := pool.Raw(sql).Scan(&result); db.Error != nil {
		logLinef(database, "err:%s", db.Error)
	}

	if len(result) > 0 {
		logLinef(database, "=========== mysql-server processfulllist ===========")
		for _, v := range result {
			logLinef(database," process id:%d, user:%s, host:%s, db:%s, command:%s, time:%d, state:%s, info:%s", v.ID, v.USER, v.HOST, v.DB, v.COMMAND, v.TIME, v.STATE, v.INFO)
		}
	}

	sql = "show VARIABLES like '%thread_pool_idle_timeout%'"

	result2 := make([]ThreadVariable, 0)
	if db := pool.Raw(sql).Scan(&result2); db.Error != nil {
		logLinef(database,"err:%s", db.Error)
	}

	if len(result2) > 0 {
		logLinef(database, "=========== mysql-server threadTimeOut ===========")
		for _, v := range result2 {
			logLinef(" thread_pool_idle_timeout Value:%s, VariableName: %s", v.Value, v.VariableName)
		}
	}

}

func logLinef(database, format string, args ...interface{}) {

	format = fmt.Sprintf("[%s],%s", database, format)
	msg := format
	if msg == "" && len(args) > 0 {
		msg = fmt.Sprint(args...)
	} else if msg != "" && len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	MysqlInitLogSlice = append(MysqlInitLogSlice, &msg)
}

func logger(filename string, logInfo []*string) error {

	var err error
	logFile, logErr := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Fail to find", *logFile, "cServer start Failed")
		return err
	}

	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	for _, v := range logInfo {
		log.Printf(*v)
	}

	return nil
}
