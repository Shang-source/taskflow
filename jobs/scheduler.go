package jobs

import (
	"context"
	"log"
	"taskflow/config"
	"taskflow/models"
	"time"
)

// StartScheduler 在应用启动时调用
func StartScheduler(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case now := <-ticker.C:
				checkAndNotify(now)
			}
		}
	}()
}

func checkAndNotify(now time.Time) {
	rows, err := config.DB.Query(`
        SELECT id, title, user_id 
        FROM tasks 
        WHERE due_date <= ? AND is_done = FALSE AND reminded = FALSE
    `, now)
	if err != nil {
		log.Printf("检查到期任务出错: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.UserID); err != nil {
			continue
		}
		msg := "任务「" + t.Title + "」已到期，请及时处理"
		// 插入 notifications 表
		_, err := config.DB.Exec(
			"INSERT INTO notifications (user_id, task_id, message) VALUES (?,?,?)",
			t.UserID, t.ID, msg,
		)
		if err == nil {
			// 标记为已提醒，避免重复
			config.DB.Exec("UPDATE tasks SET reminded = TRUE WHERE id = ?", t.ID)
		}
	}
}
