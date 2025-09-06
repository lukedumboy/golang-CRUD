package main

//gin + mysql + gorm
import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "luke:0714@tcp(127.0.0.1:3306)/crud-list?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database:%v", err)
	} else {
		log.Println("Connected to database")
	}

	sqlDB, err := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了可以重新使用连接的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	type List struct {
		Name    string
		State   string
		Phone   string
		Email   string
		Address string
	}
	db.AutoMigrate(List{})

	r := gin.Default()
	_ = r.SetTrustedProxies([]string{"127.0.0.1"})

	PORT := "3000"
	_ = r.Run(":" + PORT)

}
