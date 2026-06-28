package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"

	"jkd-website/backend/config"
	"jkd-website/backend/models"
)

type Router struct {
	cfg    *config.Config
	metaDB *sql.DB
	calDBs map[string]*sql.DB
	mu     sync.RWMutex
}

func NewRouter(cfg *config.Config) (*Router, error) {
	metaDB, err := sql.Open("mysql", cfg.MetaDSN())
	if err != nil {
		return nil, fmt.Errorf("连接元数据库失败: %w", err)
	}
	metaDB.SetMaxOpenConns(4)
	metaDB.SetMaxIdleConns(2)

	r := &Router{
		cfg:    cfg,
		metaDB: metaDB,
		calDBs: make(map[string]*sql.DB),
	}

	log.Printf("[DB] 已连接元数据库 %s", cfg.DBMeta)
	return r, nil
}

func (r *Router) Meta() *sql.DB {
	return r.metaDB
}

// Calendars 实时查询 active_calendars 视图（不做内存缓存，每次请求拉最新）。
func (r *Router) Calendars() ([]models.CalendarInfo, error) {
	rows, err := r.metaDB.Query(
		`SELECT calendarId, calendarIdI18n, db_name
		 FROM active_calendars
		 ORDER BY calendarId DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.CalendarInfo
	for rows.Next() {
		var c models.CalendarInfo
		if err := rows.Scan(&c.CalendarID, &c.CalendarName, &c.DbName); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

// CalendarByID 按 ID 查找单个学期。
func (r *Router) CalendarByID(id int) (models.CalendarInfo, error) {
	var c models.CalendarInfo
	err := r.metaDB.QueryRow(
		`SELECT calendarId, calendarIdI18n, db_name
		 FROM active_calendars WHERE calendarId = ?`, id,
	).Scan(&c.CalendarID, &c.CalendarName, &c.DbName)
	return c, err
}

func (r *Router) GetConnection(dbName string) (*sql.DB, error) {
	r.mu.RLock()
	db, ok := r.calDBs[dbName]
	r.mu.RUnlock()
	if ok {
		return db, nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if db, ok = r.calDBs[dbName]; ok {
		return db, nil
	}

	db, err := sql.Open("mysql", r.cfg.CalendarDSN(dbName))
	if err != nil {
		return nil, fmt.Errorf("连接学期库 %s 失败: %w", dbName, err)
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(2)

	r.calDBs[dbName] = db
	return db, nil
}

func (r *Router) Close() {
	for _, db := range r.calDBs {
		db.Close()
	}
	r.metaDB.Close()
}
