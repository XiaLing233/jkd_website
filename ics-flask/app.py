from flask import Flask, request, jsonify
import mysql.connector
from flask_cors import CORS
import configparser
import logging
from logging.handlers import TimedRotatingFileHandler
import os
from datetime import datetime

app = Flask(__name__)
CORS(app)

# 读取配置文件
CONFIG = configparser.ConfigParser()
CONFIG.read('config.ini', encoding='utf-8')

# 设置数据库连接
DB_HOST = CONFIG['database']['host']
DB_USER_READ_ONLY = CONFIG['database']['user-read-only']
DB_PASSWORD_READ_ONLY = CONFIG['database']['password-read-only']
DB_DATABASE = CONFIG['database']['database']
DB_PORT = int(CONFIG['database']['port'])
DB_CHARSET = CONFIG['database']['charset']

def get_client_ip():
    if 'X-Forwarded-For' in request.headers:
        return request.headers.getlist('X-Forwarded-For')[0].split(',')[0]
    return request.remote_addr

# MySQL 数据库配置
db_config = {
    'host': DB_HOST,
    'user': DB_USER_READ_ONLY,
    'password': DB_PASSWORD_READ_ONLY,
    'database': DB_DATABASE,
    'port': DB_PORT,
    'charset': DB_CHARSET
}

# 返回体格式
response_data = {
    "message": "",
    "status": 200,
    "data": []
}

# 配置日志
def setup_logger():
    os.makedirs('./logs/', exist_ok=True)
    logger = logging.getLogger('main')
    logger.setLevel(logging.INFO)
    
    # 使用 TimedRotatingFileHandler 按天轮转日志
    # 每天午夜(midnight)创建新文件，保留所有备份(backupCount=0表示保留所有)
    handler = TimedRotatingFileHandler(
        filename='./logs/app.log',
        when='midnight',
        interval=1,
        backupCount=0,
        encoding='utf-8'
    )
    # 设置日志文件名后缀格式
    handler.suffix = "%Y-%m-%d.log"
    handler.setLevel(logging.INFO)
    
    formatter = logging.Formatter('%(asctime)s - %(levelname)s\n%(message)s', datefmt='%Y-%m-%d %H:%M:%S')
    handler.setFormatter(formatter)
    logger.addHandler(handler)
    
    return logger

LOGGER = setup_logger()

# API: 获取所有可检索的字段
@app.route('/api/get_searchable_fields', methods=['GET'])
def get_searchable_fields():
    # 定义字段及其类型
    # type: "select" 表示选项式（有对应的 get_field_options 查询）
    # type: "input" 表示填写式（用户自由输入）
    fields = [
        {"field": "课程序号", "type": "input"},
        {"field": "课程名称", "type": "input"},
        {"field": "授课教师", "type": "input"},
        {"field": "教师工号", "type": "input"},
        {"field": "课程性质", "type": "select"},
        {"field": "校区", "type": "select"},
        {"field": "开课学院", "type": "select"},
        {"field": "排课信息", "type": "input"},
        {"field": "听课专业", "type": "select"}
    ]
    
    response_data['data'] = fields
    response_data['message'] = "获取字段成功！"
    response_data['status'] = 'OK'
    
    return jsonify(response_data), 200

# API: 获取字段的可选项（下拉菜单）
@app.route('/api/get_field_options', methods=['POST'])
def get_field_options():
    data = request.json
    field_name = data.get('field_name')
    select_term = data.get('select_term')
    
    conn = mysql.connector.connect(**db_config)
    cursor = conn.cursor()
    
    # 构建学期的 WHERE 条件
    term_conditions = []
    term_params = []
    if select_term and len(select_term) > 0:
        term_conditions = [f"cal.calendarIdI18n = %s" for _ in select_term]
        term_params = [f"{term}" for term in select_term]
    
    term_where = f"({' OR '.join(term_conditions)})" if term_conditions else ""
    
    # 根据字段名查询不同的表
    if field_name == '校区':
        query = """
            SELECT DISTINCT c.campusI18n 
            FROM campus c
            INNER JOIN coursedetail cd ON c.campus = cd.campus
            WHERE c.campusI18n IS NOT NULL
        """
        params = []
    elif field_name == '开课学院':
        query = """
            SELECT DISTINCT f.facultyI18n 
            FROM faculty f
            INNER JOIN coursedetail cd ON f.faculty = cd.faculty
            WHERE f.facultyI18n IS NOT NULL
        """
        params = []
    elif field_name == '课程性质':
        if not term_where:
            response_data['message'] = "空"
            response_data['data'] = []
            response_data['status'] = 'OK'
            return jsonify(response_data), 200
        
        query = f"""
            SELECT DISTINCT cn.courseLabelName 
            FROM coursenature cn
            INNER JOIN coursedetail cd ON cn.courseLabelId = cd.courseLabelId
            INNER JOIN calendar cal ON cd.calendarId = cal.calendarId
            WHERE {term_where} AND cn.courseLabelName IS NOT NULL
        """
        params = term_params
    elif field_name == '听课专业':
        if not term_where:
            response_data['message'] = "空"
            response_data['data'] = []
            response_data['status'] = 'OK'
            return jsonify(response_data), 200
            
        query = f"""
            SELECT DISTINCT m.name, m.grade
            FROM major m
            INNER JOIN majorandcourse mc ON m.id = mc.majorId
            INNER JOIN coursedetail cd ON mc.courseId = cd.id
            INNER JOIN calendar cal ON cd.calendarId = cal.calendarId
            WHERE {term_where}
            ORDER BY m.grade DESC
        """
        params = term_params
    else:
        response_data['message'] = "空"
        response_data['data'] = []
        response_data['status'] = 'OK'
        return jsonify(response_data), 200
    
    print("[DEBUG] Executing query for field options:", query)
    print("[DEBUG] With parameters:", params)

    cursor.execute(query, params)
    options = [row[0] for row in cursor.fetchall()]
    
    cursor.close()
    conn.close()
    
    # 排除空的 options
    options = [option for option in options if option]
    
    # 限制返回数量
    if len(options) > 800:
        options = options[:800]
        response_data['message'] = f"{field_name}过多，为了保证浏览体验，只返回了前 800 个"
        response_data['status'] = 'WARNING'
        response_data['data'] = options
    else:
        response_data['data'] = options
        response_data['message'] = "获取字段成功！"
        response_data['status'] = 'OK'
    
    return jsonify(response_data), 200

# API: 获取所有学期
@app.route('/api/get_terms', methods=['GET'])
def get_terms():
    conn = mysql.connector.connect(**db_config)
    cursor = conn.cursor()
    
    cursor.execute("SELECT DISTINCT calendarId, calendarIdI18n FROM calendar ORDER BY calendarId ASC")
    options = [row[1] for row in cursor.fetchall() if row[1]]
    
    cursor.close()
    conn.close()
    
    response_data['data'] = options
    response_data['message'] = "获取学期成功！"
    response_data['status'] = 'OK'
    
    return jsonify(response_data), 200

# API: 搜索功能
@app.route('/api/search', methods=['POST'])
def search():
    print(f"IP: {get_client_ip()}")
    LOGGER.info(f"IP 地址：{get_client_ip()} 的用户进行了检索。检索的条件为：{request.json}\n")
    
    data = request.json
    
    # 新的请求格式：{ groups: [...], terms: [...] }
    groups = data.get('groups', [])
    terms = data.get('terms', [])
    
    # 字段白名单
    allowed_fields = ["学期", "课程序号", "课程名称", "授课教师", "教师工号", "课程性质", "校区", "开课学院", "排课信息", "听课专业"]
    
    # 字段到数据库列的映射
    field_mapping = {
        "课程序号": "cd.code",
        "课程名称": "cd.courseName",
        "授课教师": "t.teacherName",
        "教师工号": "t.teacherCode",
        "课程性质": "cn.courseLabelName",
        "校区": "cp.campusI18n",
        "开课学院": "f.facultyI18n",
        "排课信息": "t.arrangeInfoText",
        "听课专业": "m.name",
        "学期": "cal.calendarIdI18n",
    }
    
    query_params = []
    group_clauses = []
    
    # 处理每个条件组（组间用 AND 连接）
    for group in groups:
        conditions = group.get('conditions', [])
        condition_clauses = []
        
        for condition in conditions:
            field = condition.get('field', '')
            match_type = condition.get('matchType', 'contains')
            value = condition.get('value', '')
            
            # 验证字段
            if field not in allowed_fields:
                response_data['message'] = f"非法字段：{field}"
                response_data['data'] = []
                response_data['status'] = 'ERROR'
                return jsonify(response_data), 400
            
            if not value:
                continue
            
            db_field = field_mapping.get(field, field)
            
            # 根据匹配类型构建条件
            if match_type == 'not_contains':
                condition_clauses.append(f"{db_field} NOT LIKE %s")
            else:  # contains
                condition_clauses.append(f"{db_field} LIKE %s")
            
            query_params.append(f"%{value}%")
        
        # 组内条件用 OR 连接
        if condition_clauses:
            group_clauses.append(f"({' OR '.join(condition_clauses)})")
    
    # 处理学期条件
    if not terms:
        response_data['message'] = "请选择至少一个学期！"
        response_data['data'] = []
        response_data['status'] = 'ERROR'
        return jsonify(response_data), 400
    
    # 如果没有其他条件且学期超过2个，拒绝请求
    if not group_clauses and len(terms) > 2:
        response_data['message'] = "至少选择1个检索条件！不允许在不选择条件的情况下查看超过两个学期的课程！"
        response_data['data'] = []
        response_data['status'] = 'ERROR'
        return jsonify(response_data), 400
    
    # 构建学期条件
    term_clauses = [f"cal.calendarIdI18n LIKE %s" for _ in terms]
    for term in terms:
        query_params.append(f"%{term}%")
    
    term_where = f"({' OR '.join(term_clauses)})"
    
    # 组合所有条件：组间用 AND，最后 AND 学期条件
    if group_clauses:
        where_clause = f"{' AND '.join(group_clauses)} AND {term_where}"
    else:
        where_clause = term_where
    
    # 构建复杂的 JOIN 查询来获取所有需要的信息
    # 使用子查询来筛选符合条件的课程，然后在外层查询中获取所有专业
    query = f"""
        SELECT DISTINCT
            cal.calendarId,
            cal.calendarIdI18n AS term,
            cd.code AS courseCode,
            cd.courseName AS courseName,
            cp.campusI18n AS campus,
            f.facultyI18n AS faculty,
            (SELECT GROUP_CONCAT(DISTINCT m2.name ORDER BY m2.name SEPARATOR ', ')
             FROM majorandcourse mc2
             LEFT JOIN major m2 ON mc2.majorId = m2.id
             WHERE mc2.courseId = cd.id) AS majors,
            cd.period AS totalHours,
            cn.courseLabelName AS courseType,
            a.assessmentModeI18n AS assessmentMode,
            cd.number AS capacity,
            cd.elcNumber AS enrolled,
            GROUP_CONCAT(DISTINCT CONCAT(t.teacherName, ' (', t.teacherCode, ')') ORDER BY t.teacherName SEPARATOR ', ') AS teachers,
            GROUP_CONCAT(DISTINCT t.arrangeInfoText SEPARATOR '\n') AS schedule
        FROM coursedetail cd
        LEFT JOIN calendar cal ON cd.calendarId = cal.calendarId
        LEFT JOIN coursenature cn ON cd.courseLabelId = cn.courseLabelId
        LEFT JOIN assessment a ON cd.assessmentMode = a.assessmentMode
        LEFT JOIN campus cp ON cd.campus = cp.campus
        LEFT JOIN faculty f ON cd.faculty = f.faculty
        LEFT JOIN teacher t ON cd.id = t.teachingClassId
        LEFT JOIN majorandcourse mc ON cd.id = mc.courseId
        LEFT JOIN major m ON mc.majorId = m.id
        WHERE {where_clause}
        GROUP BY cd.id, cal.calendarId, cal.calendarIdI18n, cd.code, cd.courseName, cp.campusI18n, 
                 f.facultyI18n, cd.period, cn.courseLabelName, a.assessmentModeI18n, 
                 cd.number, cd.elcNumber
        ORDER BY cal.calendarId ASC, cd.code ASC
    """
    
    print(query)
    print(query_params)
    
    conn = mysql.connector.connect(**db_config)
    cursor = conn.cursor()
    try:
        cursor.execute(query, query_params)
        results = cursor.fetchall()
        
        # 将结果转换为字典列表
        results = [dict(zip(cursor.column_names, row)) for row in results]
        
        # 处理 schedule 字段，去除重复的排课信息
        for result in results:
            if result.get('schedule'):
                # 按行分割
                schedule_lines = result['schedule'].split('\n')
                # 去重并保持顺序（使用 dict.fromkeys 去重同时保持顺序）
                unique_lines = list(dict.fromkeys(line for line in schedule_lines if line.strip()))
                # 重新组合
                result['schedule'] = '\n'.join(unique_lines)
        
    except Exception as e:
        print(f"SQL Error: {e}")
        response_data['message'] = f"检索出错啦！<br>错误信息：{str(e)}"
        response_data['data'] = []
        response_data['status'] = 'ERROR'
        return jsonify(response_data), 500
    finally:
        cursor.close()
        conn.close()
    
    response_data['data'] = results
    response_data['message'] = "检索成功！"
    response_data['status'] = 'OK'
    
    return jsonify(response_data), 200

# API: 获取更新时间
@app.route('/api/get_last_update', methods=['GET'])
def get_last_update():
    """获取数据库最后更新时间"""
    conn = mysql.connector.connect(**db_config)
    cursor = conn.cursor()
    
    try:
        cursor.execute("SELECT fetchTime, msg FROM fetchlog ORDER BY fetchTime DESC LIMIT 1")
        result = cursor.fetchone()
        
        if result:
            fetch_time, msg = result
            response_data['data'] = {
                'fetchTime': fetch_time.strftime("%Y-%m-%d") if fetch_time else "未知",
                'message': msg if msg else "数据已更新"
            }
            response_data['message'] = "获取更新时间成功！"
            response_data['status'] = 'OK'
        else:
            response_data['data'] = {
                'fetchTime': "未知",
                'message': "暂无更新记录"
            }
            response_data['message'] = "暂无更新记录"
            response_data['status'] = 'OK'
    except Exception as e:
        response_data['message'] = f"获取更新时间失败：{str(e)}"
        response_data['data'] = {}
        response_data['status'] = 'ERROR'
    finally:
        cursor.close()
        conn.close()
    
    return jsonify(response_data), 200

if __name__ == '__main__':
    app.run()
