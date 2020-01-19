package main

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go-package-demo/pkg/jinzhu/gorm/model"
	"go-package-demo/pkg/jinzhu/gorm/mysql"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	if err := mysql.InitMysql(); err != nil {
		panic(err)
	}

	product := new(model.Product)

	//product.Migrade()

	//product.Create()
	//product.First()
	product.Update()
	//product.Delete()

	graceShutDown()

}

// 释放连接池，优雅退出
func graceShutDown() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	<-quit

	log.Println("Shutdown Server ...")

	/*
		在进程退出时释放mysql连接池  在入口处，应用程序退出时 DB().Close()
	*/
	mysql.MysqlPool.Close()
}
