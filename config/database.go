package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// DB 全局数据库连接对象
var DB *sql.DB

func init() {
	// 在包初始化时就加载 .env（也可在 main.go 中加载）
	err := godotenv.Load()
	if err != nil {
		log.Println("未找到 .env 文件，使用默认环境变量")
	}
}

// ConnectDB 连接数据库
func ConnectDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	// 测试连通性
	if err = db.Ping(); err != nil {
		log.Fatalf("数据库无法Ping通: %v", err)
	}

	log.Println("数据库连接成功!")
	DB = db

	// 这里可执行建表等操作
	initTables()
}

// initTables 可选：初始化表结构
func initTables() {
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL
	);
	`
	_, err := DB.Exec(createUserTable)
	if err != nil {
		log.Printf("创建 users 表出错: %v\n", err)
	}

	createTaskTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		is_done BOOLEAN NOT NULL DEFAULT false,
	    due_date DATETIME NULL,
      reminded BOOLEAN NOT NULL DEFAULT FALSE,
		user_id BIGINT,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`
	_, err = DB.Exec(createTaskTable)
	if err != nil {
		log.Printf("Create tasks tables error: %v\n", err)

		// 新增 notifications 表，用于存储提醒
		createNotifTable := `
        CREATE TABLE IF NOT EXISTS notifications (
          id BIGINT AUTO_INCREMENT PRIMARY KEY,
          user_id BIGINT NOT NULL,
          task_id BIGINT NOT NULL,
          message VARCHAR(255) NOT NULL,
          is_read BOOLEAN NOT NULL DEFAULT FALSE,
          created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
          FOREIGN KEY (user_id) REFERENCES users(id),
          FOREIGN KEY (task_id) REFERENCES tasks(id)
    );
    `
		if _, err := DB.Exec(createNotifTable); err != nil {
			log.Printf("Create notifications tables error: %v\n", err)
		}
	}

}
