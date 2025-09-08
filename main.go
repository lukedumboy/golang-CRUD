package main

//gin + mysql + gorm
import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	// 设置用户名 密码 并选中数据库，需要该用户对该数据库有完全访问权限
	dsn := "luke:0714@tcp(127.0.0.1:3306)/crud-list?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalf("failed to connect database:%v", err)
	} else {
		log.Println("Connected to database")
	}

	sqlDB, sqlDBErr := db.DB()
	if sqlDBErr != nil {
		log.Fatalf("failed to get sql.DB from gorm: %v", err)
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了可以重新使用连接的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	type List struct {
		gorm.Model
		Name    string `gorm:"type:varchar(20);not null" json:"name" binding:"required"`
		State   string `gorm:"type:varchar(20);not null" json:"state" binding:"required"`
		Phone   string `gorm:"type:varchar(20);not null" json:"phone" binding:"required"`
		Email   string `gorm:"type:varchar(40);not null" json:"email" binding:"required"`
		Address string `gorm:"type:varchar(200);not null" json:"address" binding:"required"`
	}
	_ = db.AutoMigrate(List{})

	router := gin.Default()
	_ = router.SetTrustedProxies([]string{"127.0.0.1"})
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//业务码约定
	//正确:ADD_SUCCESS 错误:ADD_FAILED

	//c create
	router.POST("user/add", func(c *gin.Context) {
		var json List
		err := c.ShouldBindJSON(&json)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "添加失败:" + err.Error(),
				"data":    gin.H{},
				"code":    "ADD_FAILED",
			})
		} else {
			//数据库操作
			db.Create(&json)
			c.JSON(http.StatusOK, gin.H{
				"message": "添加成功",
				"data":    json,
				"code":    "ADD_SUCCESS",
			})
		}
	})
	//r read
	//u update
	//d delete

	PORT := "3000"
	_ = router.Run(":" + PORT)
}
