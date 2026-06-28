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
	cfg       *config.Config
	metaDB    *sql.DB
	calDBs    map[string]*sql.DB        // db_name → connection
	calendars []models.CalendarInfo
	calByID   map[int]models.CalendarInfo // calendarId → info
	mu        sync.RWMutex
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

	if err := r.refreshCalendars(); err != nil {
		return nil, fmt.Errorf("加载学期列表失败: %w", err)
	}

	return r, nil
}

func (r *Router) Meta() *sql.DB {
	return r.metaDB
}

func (r *Router) Calendars() []models.CalendarInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.calendars
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

	// 双重检查，避免重复创建连接池
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

func (r *Router) refreshCalendars() error {
	rows, err := r.metaDB.Query(
		`SELECT calendarId, calendarIdI18n, db_name
		 FROM active_calendars
		 ORDER BY calendarId DESC`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	var list []models.CalendarInfo
	byID := make(map[int]models.CalendarInfo)
	for rows.Next() {
		var c models.CalendarInfo
		if err := rows.Scan(&c.CalendarID, &c.CalendarName, &c.DbName); err != nil {
			return err
		}
		list = append(list, c)
		byID[c.CalendarID] = c
	}
	if err := rows.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	r.calendars = list
	r.calByID = byID
	r.mu.Unlock()

	log.Printf("[DB] 已加载 %d 个学期", len(list))
	return nil
}

func (r *Router) CalendarByID(id int) (models.CalendarInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.calByID[id]
	return c, ok
}
