from flask import Flask, request, jsonify
import mysql.connector
from flask_cors import CORS  # 导入 CORS

app = Flask(__name__)

CORS(app) # 允许跨域请求

# 配置 MySQL 数据库连接
# MySQL 数据库配置
db_config = {
    'host': 'localhost',      # 数据库主机地址，通常是 'localhost' 或远程服务器 IP
    'user': 'root',  # 连接数据库的用户名
    'password': 'PNxe3LuNCjx9LT*',  # 用户对应的密码
    'database': 'curriculum',    # 数据库名
    'port': 3306,             # MySQL 服务的端口号，默认为 3306
    'charset': 'utf8mb4',     # 字符集设置，保证支持中文等多语言字符
}

# 返回体格式
response_data = {
    "message": "",
    "status": 200,
    "data": []
}

# 去重函数
def get_unique_majors(result):
    unique_majors = set() # 用集合去重

    for record in result:
        majors = record.split(",") # 假设每条记录是以逗号分隔的
        for major in majors:
            unique_majors.add(major.strip()) # 去除空格后添加到集合中

    # 按照字段的前四位年份降序排序
    unique_majors = sorted(list(unique_majors), key=lambda x: x[:4], reverse=True)

    return list(unique_majors) # 不返回状态码，因为不是 HTTP 请求的返回

# 以下是 Flask 的路由函数

# 用于获取所有可检索的字段
@app.route('/api/get_searchable_fields', methods=['GET'])
def get_searchable_fields():
    # 示例字段，可以从数据库获取
    fields = ["课程序号", "课程名称", "授课教师", "教师工号", "课程性质", "校区", "开课学院", "排课信息", "听课专业"]

    response_data['data'] = fields
    response_data['message'] = "获取字段成功！"
    response_data['status'] = 'OK'

    return jsonify(response_data), 200

# 有一些字段，提供给用户下拉菜单选择
@app.route('/api/get_field_options', methods=['POST'])
def get_field_options():
    data = request.json  # 获取请求的 JSON 数据
    field_name = data.get('field_name')  # 获取字段名
    select_term = data.get('select_term')  # 获取学期

    # 连接数据库，查询字段的可选值
    conn = mysql.connector.connect(**db_config)  # 连接数据库
    cursor = conn.cursor()  # 创建游标

    # 构建学期的 LIKE 语句
    term_like = " OR ".join(["学期 LIKE %s" for _ in select_term])
    term_like = f"({term_like})"
    term_params = [f"%{term}%" for term in select_term]

    print(term_like)

    # 根据字段名查询不同的表
    if field_name == '校区':  # 如果字段名是“校区”
        query = "SELECT DISTINCT 校区 FROM course_all"
        params = []
    elif field_name == '开课学院':
        query = "SELECT DISTINCT 开课学院 FROM course_all"
        params = []
    elif field_name == '课程性质':
        if (term_like == "()"): # 鬼知道还有人不选学期？没错，就是我！
            response_data['message'] = "空"
            response_data['data'] = []
            response_data['status'] = 'OK'
            return jsonify(response_data), 200  # 返回空数组，如果没有匹配的字段
        query = f"SELECT DISTINCT 课程性质 FROM course_all WHERE {term_like}"
        params = term_params
    elif field_name == '听课专业':
        if (term_like == "()"):
            response_data['message'] = "空"
            response_data['data'] = []
            response_data['status'] = 'OK'
            return jsonify(response_data), 200  # 返回空数组，如果没有匹配的字段
        query = f"SELECT DISTINCT 听课专业 FROM course_all WHERE {term_like} ORDER BY SUBSTRING(听课专业, 1, 4) DESC"
        params = term_params
    else:
        response_data['message'] = "空"
        response_data['data'] = []
        response_data['status'] = 'OK'
        return jsonify(response_data), 200  # 返回空数组，如果没有匹配的字段

    print(query)  # 打印查询语句

    cursor.execute(query, params)  # 执行查询
    options = [row[0] for row in cursor.fetchall()]  # 获取查询结果

    if field_name == '听课专业':
        options = get_unique_majors(options)

    cursor.close()  # 关闭游标
    conn.close()  # 关闭数据库连接

    # 排除空的 options
    options = [option for option in options if option]

    # 如果个数多于 800 个，只返回前 800 个，并且给出 warning
    print(len(options))
    if len(options) > 800:
        options = options[:800]
        response_data['message'] = "%s过多，为了保证浏览体验，只返回了前 800 个，可能漏掉某些字段。建议缩小检索的学期范围，或者先选定一个学期，找到目标后进一步检索。" % field_name
        response_data['status'] = 'WARNING'
        response_data['data'] = options
    else:
        response_data['data'] = options
        response_data['message'] = "获取字段成功！"
        response_data['status'] = 'OK'

    return jsonify(response_data), 200

@app.route('/api/get_terms', methods=['GET'])
def get_terms():
    conn = mysql.connector.connect(**db_config) # 连接数据库
    cursor = conn.cursor()

    cursor.execute("SELECT DISTINCT 学期 FROM course_all ORDER BY 学期 ASC") # 查询 course_all 表的学期字段
    # 获取不为空的学期
    options = [row[0] for row in cursor.fetchall() if row[0] != ''] # 获取查询结果

    cursor.close() # 关闭游标
    conn.close() # 关闭数据库连接

    response_data['data'] = options
    response_data['message'] = "获取学期成功！"
    response_data['status'] = 'OK'

    return jsonify(response_data), 200


# 搜索功能
@app.route('/api/search', methods=['POST'])
def search():
    conditions = request.json # 获取请求的 JSON 数据
    # 构建查询条件
    query_conditions = []
    query_params = []
    current_field = None
    field_conditions = []
    first_connector = "Empty"  # 第一个连接词
    
    # 黑白名单
    allowed_fields = ["学期", "课程序号", "课程名称", "授课教师", "教师工号", "课程性质", "校区", "开课学院", "排课信息", "听课专业"]
    allowed_connectors = ["AND", "OR", "NOT"] # 允许的连接词

    for condition in conditions: # 遍历所有查询条件
        try:
            field = condition['selectedItem'] # 获取字段名
        except: # 字段名这一栏压根就不存在
            response_data['message'] = "字段名不能为空！"
            response_data['data'] = []
            response_data['status'] = 'ERROR'
            return jsonify(response_data), 400
        
        # 如果 field 为空，返回错误信息
        if not field:
            response_data['message'] = "字段名不能为空！"
            response_data['data'] = []
            response_data['status'] = 'ERROR'
            return jsonify(response_data), 400 # 返回错误信息

        try:
            value = condition['searchWord'] # 获取搜索关键词
        except:# 关键词这一栏压根就不存在
            response_data['message'] = "搜索关键词不能为空！"
            response_data['data'] = []
            response_data['status'] = 'ERROR'
            return jsonify(response_data), 400 # 返回错误信息
        
        connector = condition.get('connector', '') # 获取连接词

        if field not in allowed_fields: # 如果字段名不在允许的字段列表中
            response_data['message'] = "非法字段！"
            response_data['data'] = []
            response_data['status'] = 'ERROR'
            return jsonify(response_data), 400 # 返回错误信息
        
        if not value: # 如果搜索关键词为空
            response_data['message'] = "搜索关键词不能为空！"
            response_data['data'] = []
            response_data['status'] = 'ERROR'
            return jsonify(response_data), 400 # 返回错误信息
        
        if connector and connector not in allowed_connectors: # 如果连接词存在，但不在允许的连接词列表中
            response_data['message'] = "非法连接词！"
            response_data['data'] = []
            response_data['status'] = 'ERROR'
            return jsonify(response_data), 400 # 返回错误信息
        
        if not connector and field_conditions: # 如果连接词不存在，并且临时条件列表不为空
            response_data['message'] = "连接词不能为空！"
            response_data['data'] = []
            response_data['status'] = 'ERROR'
            return jsonify(response_data), 400 # 返回错误信息

        if current_field is None:
            current_field = field

        # 如果当前字段与前一个字段不同
        if field != current_field:
            # 将当前字段的条件括起来，并添加到查询条件列表中
            if field_conditions:
                if current_field == "学期":
                    query_conditions.append(f"AND ({' '.join(field_conditions)})")
                else:
                    if first_connector == "NOT":
                        if not query_conditions:
                            query_conditions.append(f"({' '.join(field_conditions)})")
                        else:
                            query_conditions.append(f"AND ({' '.join(field_conditions)})")
                    else:
                        query_conditions.append(f"{first_connector} ({' '.join(field_conditions)})")
            current_field = field
            field_conditions = []
            first_connector = connector  # 更新新的字段的第一个连接词

        # 构建单个条件
        if connector == "NOT":
            field_condition = f"AND {field} NOT LIKE %s"
            if not field_conditions:
                field_condition = f"{field} NOT LIKE %s"
        else:
            field_condition = f"{connector} {field} LIKE %s"

        if not field_conditions:
            if connector == "NOT":
                field_condition = f"{field} NOT LIKE %s"
            else:
                field_condition = f"{field} LIKE %s"

        field_conditions.append(field_condition)
        query_params.append(f"%{value}%")

        if first_connector == "Empty":
            first_connector = connector

    if field_conditions:
        if current_field == "学期":
            if not query_conditions:
                if len(field_conditions) > 2:
                    response_data['message'] = "至少选择1个检索条件！不允许在不选择条件的情况下查看超过两个学期的课程喔！"
                    response_data['data'] = []
                    response_data['status'] = 'ERROR'
                    return jsonify(response_data), 400
                else:
                    query_conditions.append(f"({' '.join(field_conditions)})")
            else:
                query_conditions.append(f"AND ({' '.join(field_conditions)})")
        else:
            if first_connector == "NOT":
                if not query_conditions:
                    query_conditions.append(f"{' '.join(field_conditions)}")
                else:
                    query_conditions.append(f"AND ({' '.join(field_conditions)})")
            else:
                query_conditions.append(f"{first_connector} ({' '.join(field_conditions)})")

    # For Debugging Purpose
    print(query_conditions)
    print(query_params)


    where_clause = ' '.join(query_conditions)
    query = f"SELECT * FROM course_all WHERE {where_clause} ORDER BY `学期` ASC, `课程序号` ASC"

    # For Debugging Purpose
    print(query)

    conn = mysql.connector.connect(**db_config)
    cursor = conn.cursor()
    try:
        cursor.execute(query, query_params)
        results = cursor.fetchall()
    except Exception as e:
        response_data['message'] = "检索出错啦！<br>生成的 SQL 语句为：<br>" + query
        response_data['data'] = []
        response_data['status'] = 500
        return jsonify(response_data), 500

    results = [dict(zip(cursor.column_names, row)) for row in results]

    cursor.close()
    conn.close()

    response_data['data'] = results
    response_data['message'] = "检索成功！"
    response_data['status'] = 'OK'

    return jsonify(response_data), 200


if __name__ == '__main__': # 如果当前脚本被直接运行
    app.run(debug=True) # 启动 Flask 应用
