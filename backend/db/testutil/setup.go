// Package testutil 提供集成测试的数据库初始化（仅 insert 种子数据）。
// 假设数据库和表已由 docker-compose 或外部工具创建。
package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	TestMeta  = "test_course_scheduler_meta"
	TestCalA  = "calendar_999_a"
	TestCalID = 999
)

type TestCfg struct {
	Host string
	Port string
	User string
	Pass string
}

func DefaultCfg() TestCfg {
	return TestCfg{
		Host: getEnv("DB_HOST", "xk-mysql"),
		Port: getEnv("DB_PORT", "3306"),
		User: getEnv("DB_R_USER", "root"),
		Pass: getEnv("DB_R_PASSWORD", ""),
	}
}

func (c TestCfg) DSN(db string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		c.User, c.Pass, c.Host, c.Port, db)
}

func (c TestCfg) RootDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=true",
		c.User, c.Pass, c.Host, c.Port)
}

// Setup 创建测试库 + 种子数据。重复调用幂等。
func Setup(cfg TestCfg) error {
	rootDB, err := connectRetry(cfg.RootDSN())
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer rootDB.Close()

	for _, db := range []string{TestMeta, TestCalA} {
		if _, err := rootDB.Exec(fmt.Sprintf(
			"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", db,
		)); err != nil {
			return fmt.Errorf("create db %s: %w", db, err)
		}
	}

	if err := execDDL(cfg.DSN(TestMeta), metaDDL); err != nil {
		return err
	}
	if err := execDDL(cfg.DSN(TestCalA), courseDDL); err != nil {
		return err
	}
	if err := seedMeta(cfg.DSN(TestMeta)); err != nil {
		return err
	}
	if err := seedCalendar(cfg.DSN(TestCalA)); err != nil {
		return err
	}
	return nil
}

// SetEnv 设置环境变量使 config.Load 指向测试库。
func SetEnv(cfg TestCfg) {
	os.Setenv("DB_HOST", cfg.Host)
	os.Setenv("DB_PORT", cfg.Port)
	os.Setenv("DB_R_USER", cfg.User)
	os.Setenv("DB_R_PASSWORD", cfg.Pass)
	os.Setenv("DB_META", TestMeta)
}

func connectRetry(dsn string) (*sql.DB, error) {
	var lastErr error
	for i := 0; i < 10; i++ {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			lastErr = err
		} else if err := db.Ping(); err != nil {
			db.Close()
			lastErr = err
		} else {
			return db, nil
		}
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("retry exhausted: %w", lastErr)
}

func execDDL(dsn, ddl string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	// MariaDB 支持 CREATE IF NOT EXISTS + OR REPLACE VIEW，直接跑
	for _, s := range []string{ddl} {
		for _, stmt := range splitSQL(s) {
			if _, err := db.Exec(stmt); err != nil {
				return fmt.Errorf("%w\nstmt: %s", err, stmt[:min(len(stmt), 80)])
			}
		}
	}
	return nil
}

func splitSQL(s string) []string {
	var result []string
	var cur string
	for _, line := range splitLines(s) {
		line = trim(line)
		if line == "" || line[0] == '-' {
			continue
		}
		cur += line + " "
		if line[len(line)-1] == ';' {
			result = append(result, trim(cur))
			cur = ""
		}
	}
	return result
}

func seedMeta(dsn string) error {
	db, _ := sql.Open("mysql", dsn)
	if db == nil {
		return fmt.Errorf("cannot connect to meta")
	}
	defer db.Close()
	_, err := db.Exec("INSERT IGNORE INTO calendar_registry (calendarId, calendarIdI18n) VALUES (?, '测试学期')", TestCalID)
	return err
}

func seedCalendar(dsn string) error {
	db, _ := sql.Open("mysql", dsn)
	if db == nil {
		return fmt.Errorf("cannot connect to calendar")
	}
	defer db.Close()
	db.Exec("SET FOREIGN_KEY_CHECKS=0")
	db.Exec("INSERT IGNORE INTO assessment VALUES ('1','考试'),('2','考查')")
	db.Exec("INSERT IGNORE INTO campus VALUES ('1','四平路校区'),('3','嘉定校区')")
	db.Exec("INSERT IGNORE INTO coursenature VALUES (325,'专业课必修'),(958,'科学探索与生命关怀')")
	db.Exec("INSERT IGNORE INTO faculty VALUES ('000034','电子与信息工程学院')")
	db.Exec("INSERT IGNORE INTO language VALUES ('1','中文')")

	db.Exec("INSERT IGNORE INTO coursedetail (id,code,name,courseLabelId,assessmentMode,campus,courseCode,courseName,credit,teachingLanguage,faculty) VALUES " +
		"(1001,'34001201','01班',325,'1','1','340012','编译原理',3.0,'1','000034')," +
		"(1002,'14007601','01班',958,'2','3','140076','公共营养学',1.5,'1','000034')")

	db.Exec("INSERT IGNORE INTO teacher VALUES " +
		"(200000000001,1001,'13060','李华','李华(13060) 星期一3-4节 [1-16] A楼101\n')," +
		"(200000000002,1002,'05222','关佶红','关佶红(05222) 星期一3-4节 [1-17] 南129\n')")

	db.Exec("INSERT IGNORE INTO major (id,code,grade,name) VALUES (1,'10054',2023,'2023(10054 计算机科学与技术)')")
	db.Exec("INSERT IGNORE INTO majorandcourse (majorId,courseId) VALUES (1,1001)")
	db.Exec("SET FOREIGN_KEY_CHECKS=1")
	return nil
}

// ─── helpers ───

func splitLines(s string) []string {
	var r []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			r = append(r, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		r = append(r, s[start:])
	}
	return r
}

func trim(s string) string {
	lo, hi := 0, len(s)
	for lo < hi && (s[lo] == ' ' || s[lo] == '\t' || s[lo] == '\r') {
		lo++
	}
	for hi > lo && (s[hi-1] == ' ' || s[hi-1] == '\t' || s[hi-1] == '\r') {
		hi--
	}
	return s[lo:hi]
}

func getEnv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}

const metaDDL = `
CREATE TABLE IF NOT EXISTS calendar_registry (
  calendarId int NOT NULL, calendarIdI18n varchar(200) NOT NULL,
  active_suffix char(1) NOT NULL DEFAULT 'a',
  created_at datetime DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (calendarId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE OR REPLACE VIEW active_calendars AS
SELECT calendarId, calendarIdI18n,
  CONCAT('calendar_', calendarId, '_', active_suffix) AS db_name,
  active_suffix, updated_at FROM calendar_registry;

CREATE TABLE IF NOT EXISTS fetchlog (
  id int NOT NULL AUTO_INCREMENT, calendarId int NOT NULL,
  startTime datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  endTime datetime(3) DEFAULT NULL,
  status enum('running','completed','failed') NOT NULL DEFAULT 'running',
  totalCourses int DEFAULT 0, totalPages int DEFAULT 0,
  msg varchar(500) DEFAULT NULL, errorMessage text, fullLog mediumtext,
  PRIMARY KEY (id), KEY idx_fetchlog_time (startTime),
  CONSTRAINT fk_fetchlog_calendar FOREIGN KEY (calendarId) REFERENCES calendar_registry (calendarId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
`

const courseDDL = `
CREATE TABLE IF NOT EXISTS assessment (
  assessmentMode varchar(200) NOT NULL, assessmentModeI18n varchar(200) DEFAULT NULL,
  PRIMARY KEY (assessmentMode)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS campus (
  campus varchar(200) NOT NULL, campusI18n varchar(200) DEFAULT NULL,
  PRIMARY KEY (campus)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS coursenature (
  courseLabelId int NOT NULL, courseLabelName varchar(200) DEFAULT NULL,
  PRIMARY KEY (courseLabelId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS faculty (
  faculty varchar(200) NOT NULL, facultyI18n varchar(200) DEFAULT NULL,
  PRIMARY KEY (faculty)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS language (
  teachingLanguage varchar(200) NOT NULL, teachingLanguageI18n varchar(200) DEFAULT NULL,
  PRIMARY KEY (teachingLanguage)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS major (
  id int NOT NULL AUTO_INCREMENT, code varchar(200) DEFAULT NULL,
  grade int DEFAULT NULL, name varchar(200) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS coursedetail (
  id bigint NOT NULL, code varchar(200) DEFAULT NULL, name TEXT DEFAULT NULL,
  courseLabelId int DEFAULT NULL, assessmentMode varchar(200) DEFAULT NULL,
  period int DEFAULT NULL, weekHour int DEFAULT NULL, campus varchar(200) DEFAULT NULL,
  number int DEFAULT NULL, elcNumber int DEFAULT NULL,
  startWeek int DEFAULT NULL, endWeek int DEFAULT NULL,
  courseCode varchar(200) DEFAULT NULL, courseName varchar(200) DEFAULT NULL,
  credit double DEFAULT NULL, teachingLanguage varchar(200) DEFAULT NULL,
  faculty varchar(200) DEFAULT NULL, newCode varchar(200) DEFAULT NULL,
  newCourseCode varchar(200) DEFAULT NULL,
  PRIMARY KEY (id), KEY nature_idx (courseLabelId), KEY assess_idx (assessmentMode),
  KEY campusKey_idx (campus), KEY facultyKey_idx (faculty), KEY langKey_idx (teachingLanguage),
  CONSTRAINT assessKey FOREIGN KEY (assessmentMode) REFERENCES assessment (assessmentMode) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT campusKey FOREIGN KEY (campus) REFERENCES campus (campus) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT facultyKey FOREIGN KEY (faculty) REFERENCES faculty (faculty) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT langKey FOREIGN KEY (teachingLanguage) REFERENCES language (teachingLanguage) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT natureKey FOREIGN KEY (courseLabelId) REFERENCES coursenature (courseLabelId) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS majorandcourse (
  id int NOT NULL AUTO_INCREMENT, majorId int NOT NULL, courseId bigint NOT NULL,
  PRIMARY KEY (id), KEY courseKey_idx (courseId), KEY majorKeyForMajor_idx (majorId),
  CONSTRAINT courseKeyForMajor FOREIGN KEY (courseId) REFERENCES coursedetail (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT majorKeyForMajor FOREIGN KEY (majorId) REFERENCES major (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS teacher (
  id bigint NOT NULL, teachingClassId bigint DEFAULT NULL,
  teacherCode varchar(200) DEFAULT NULL, teacherName varchar(200) DEFAULT NULL,
  arrangeInfoText mediumtext,
  PRIMARY KEY (id), KEY classKey_idx (teachingClassId),
  CONSTRAINT courseKey FOREIGN KEY (teachingClassId) REFERENCES coursedetail (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
`
