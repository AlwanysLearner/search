package DB

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

var mysqldb *gorm.DB

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

func InitDatabase() {
	var wg sync.WaitGroup
	conf, _ := os.Open("mysql.json")
	defer conf.Close() //执行完毕后关闭连接
	var config DBConfig
	jsonParser := json.NewDecoder(conf)
	if err := jsonParser.Decode(&config); err != nil {
		panic(err)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.User, config.Password, config.Host, config.Port, config.DBName)
	// 初始化日志配置，设置为Info级别，这样可以打印所有SQL语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // 慢查询阈值
			LogLevel:      logger.Silent, // Log level
			Colorful:      true,          // 彩色打印
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed creating database:%w", err)
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get database: %v", err)
		return
	}
	sqlDB.SetMaxIdleConns(300)
	sqlDB.SetMaxOpenConns(500)
	for i := 0; i < 300; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := db.Raw("SELECT 1").Rows(); err != nil {
				log.Println("error executing query:", err)
			}
		}()
	}
	wg.Wait()
	mysqldb = db
}
func DataBaseSessoin() *gorm.DB {
	return mysqldb.Session(&gorm.Session{PrepareStmt: true})
}
