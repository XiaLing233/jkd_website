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

    return list(unique_majors)

# 以下是 Flask 的路由函数

# 用于获取所有可检索的字段
@app.route('/api/get_searchable_fields', methods=['GET'])
def get_searchable_fields():
    # 示例字段，可以从数据库获取
    fields = ["课程序号", "课程名称", "授课教师", "教师工号", "课程性质", "校区", "开课学院", "排课信息", "听课专业"]

    response_data['data'] = fields
    response_data['status'] = 200
    response_data['message'] = "获取字段成功！"

    return jsonify(response_data)

# 有一些字段，提供给用户下拉菜单选择
@app.route('/api/get_field_options', methods=['POST'])
def get_field_options():
    data = request.json # 获取请求的 JSON 数据
    field_name = data.get('field_name') # 获取字段名
    select_term = data.get('select_term') # 获取学期

    # 连接数据库，查询字段的可选值
    conn = mysql.connector.connect(**db_config) # 连接数据库
    cursor = conn.cursor() # 创建游标

    # 构建学期的 LIKE 语句
    term_like = "(学期 LIKE "
    for term in select_term:
        # 第一个
        if term == select_term[0]:
            term_like += f"'{term}'"
        else:
            term_like += f" OR 学期 LIKE '{term}'"
    
    term_like += ")"

    # 根据字段名查询不同的表
    if field_name == '校区': # 如果字段名是“校区”
        query = "SELECT DISTINCT 校区 FROM course_all" # 查询 course_all 表的校区字段
    elif field_name == '开课学院':
        query = "SELECT DISTINCT 开课学院 FROM course_all"
    elif field_name == '课程性质':
        query = f"SELECT DISTINCT 课程性质 FROM course_all WHERE {term_like}"
    elif field_name == '听课专业':
        query = f"SELECT DISTINCT 听课专业 FROM course_all WHERE {term_like} ORDER BY SUBSTRING(听课专业, 1, 4) DESC"
    else:
        response_data['message'] = "空"
        response_data['status'] = 200
        response_data['data'] = []
        return jsonify(response_data) # 返回空数组，如果没有匹配的字段

    print(query) # 打印查询语句

    cursor.execute(query) # 执行查询
    options = [row[0] for row in cursor.fetchall()] # 获取查询结果

    # print(options)

    if field_name == '听课专业':
        options = get_unique_majors(options)
    
    cursor.close() # 关闭游标
    conn.close() # 关闭数据库连接

    # 排除空的 options
    options = [option for option in options if option]

    # print(options)

    response_data['data'] = options
    response_data['status'] = 200
    response_data['message'] = "获取字段成功！"

    return jsonify(response_data)

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
    response_data['status'] = 200
    response_data['message'] = "获取学期成功！"

    return jsonify(response_data)


# 搜索功能
@app.route('/api/search', methods=['POST'])
def search():
    conditions = request.json # 获取请求的 JSON 数据
    # 构建查询条件
    query_conditions = []
    current_field = None
    field_conditions = []
    first_connector = "Empty"  # 第一个连接词
    
    # 黑白名单
    allowed_fields = ["学期", "课程序号", "课程名称", "授课教师", "教师工号", "课程性质", "校区", "开课学院", "排课信息", "听课专业"]
    allowed_connectors = ["AND", "OR", "NOT"] # 允许的连接词
    blacklist_searchWord = [
    "SELECT",    # 查询数据
    "INSERT",    # 插入数据
    "UPDATE",    # 更新数据
    "DELETE",    # 删除数据
    "CREATE",    # 创建数据库或表
    "DROP",      # 删除数据库或表
    "ALTER",     # 修改数据库或表结构
    "TRUNCATE",  # 清空表数据
    "GRANT",     # 授权权限
    "REVOKE",    # 撤销权限
    "COMMIT",    # 提交事务
    "ROLLBACK",  # 回滚事务
    "USE",       # 选择数据库
    "FROM",      # 用于指定表名
    "WHERE",     # 用于条件过滤
    "JOIN",      # 用于连接多个表
    "ON",        # 用于连接条件
    "AND",       # 逻辑与
    "OR",        # 逻辑或
    "NOT",       # 逻辑非
    "LIKE",      # 模糊匹配
    "IN",        # 用于指定多个值
    "BETWEEN",   # 用于范围查询
    "HAVING",    # 用于聚合条件
    "ORDER BY",  # 用于排序
    "GROUP BY",  # 用于分组
    "LIMIT",     # 用于限制结果集
    "OFFSET",    # 用于分页
    "UNION",     # 用于合并多个结果集
    "EXCEPT",    # 用于除去重复的结果
    "INTERSECT", # 用于交集
    "WITH",      # 公共表表达式
    "CASE",      # 用于条件表达式
    "WHEN",      # 用于条件判断
    "THEN",      # 用于条件判断
    "ELSE",      # 用于条件判断
    "END",        # 结束条件表达式
    ";",         # 分号
    "'",         # 单引号
    "\"",        # 双引号
    "--",        # 注释
    "/*",        # 注释
    "%",         # 用于模糊匹配
    "_",         # 用于模糊匹配
]
    

    for condition in conditions: # 遍历所有查询条件
        field = condition['selectedItem'] # 获取字段名
        try:
            value = condition['searchWord'] # 获取搜索关键词
        except:
            response_data['message'] = "搜索关键词不能为空！"
            response_data['status'] = 400
            response_data['data'] = []
            return jsonify(response_data), 400 # 返回错误信息
        
        connector = condition.get('connector', '') # 获取连接词

        if field not in allowed_fields: # 如果字段名不在允许的字段列表中
            response_data['message'] = "非法字段！"
            response_data['status'] = 400
            response_data['data'] = []
            return jsonify(response_data), 400 # 返回错误信息
        
        if not value: # 如果搜索关键词为空
            response_data['message'] = "搜索关键词不能为空！"
            response_data['status'] = 400
            response_data['data'] = []
            return jsonify(response_data), 400 # 返回错误信息
        
        if connector and connector not in allowed_connectors: # 如果连接词存在，但不在允许的连接词列表中
            response_data['message'] = "非法连接词！"
            response_data['status'] = 400
            response_data['data'] = []
            return jsonify(response_data), 400 # 返回错误信息
        
        if not connector and field_conditions: # 如果连接词不存在，并且临时条件列表不为空
            response_data['message'] = "连接词不能为空！"
            response_data['status'] = 400
            response_data['data'] = []
            return jsonify(response_data), 400 # 返回错误信息

        value_upper = value.upper() # 转换为大写

        for word in blacklist_searchWord:
            if word in value_upper:
                response_data['message'] = "非法关键词！"
                response_data['status'] = 400
                response_data['data'] = []
                return jsonify(response_data) # 返回错误信息

        if current_field is None:
            current_field = field

        # 如果当前字段与前一个字段不同
        if field != current_field:
            # 将当前字段的条件括起来，并添加到查询条件列表中
            if field_conditions:
                 # 如果是学期，加 AND
                if current_field == "学期":
                    query_conditions.append(f"AND ({' '.join(field_conditions)})")
                    print("学期")
                else:
                    # 如果 first_connector 为 NOT
                    if first_connector == "NOT":
                        # 如果是第一个条件，不加 AND
                        if not query_conditions:
                            print("第一个条件")
                            query_conditions.append(f"({' '.join(field_conditions)})")
                        else:
                            print(first_connector)
                            query_conditions.append(f"AND ({' '.join(field_conditions)})")
                    else:
                        query_conditions.append(f"{first_connector} ({' '.join(field_conditions)})")
            # 更新当前字段，并清空临时条件列表
            current_field = field
            field_conditions = []
            first_connector = connector  # 更新新的字段的第一个连接词

        # 构建单个条件
        # 如果是 NOT
        if connector == "NOT":
            field_condition = f"AND {field} NOT LIKE '%{value}%'"
            # 如果是第一个条件，不加 AND
            if not field_conditions: # 区分好 field_conditions 和 query_conditions!
                field_condition = f"{field} NOT LIKE '%{value}%'"
        else:
            field_condition = f"{connector} {field} LIKE '%{value}%'"

        if not field_conditions:  # 第一个条件不加 connector
            # 如果是 NOT
            if connector == "NOT":
                field_condition = f"{field} NOT LIKE '%{value}%'"
            else:
                field_condition = f"{field} LIKE '%{value}%'"

        # 将条件添加到临时列表中
        field_conditions.append(field_condition)

        # 对于第一个条件，记录连接词并重置为下一次使用
        if first_connector == "Empty":
            first_connector = connector

    # 添加最后一个字段的条件
    print(first_connector)
    if field_conditions:
        # 如果是学期，加 AND
        if current_field == "学期":
            # 如果选择的学期超过两个，并且学期是第一个条件（说明没有选择条件），返回错误信息
            if not query_conditions:
                if len(field_conditions) > 2: # 说明没有选择条件
                    # 返回错误信息，中间有换行符

                    response_data['message'] = "至少选择1个检索条件！不允许在不选择条件的情况下查看超过两个学期的课程喔！"
                    response_data['status'] = 400
                    response_data['data'] = []
                    return jsonify(response_data), 400
                else:
                    query_conditions.append(f"({' '.join(field_conditions)})")
            else:
                query_conditions.append(f"AND ({' '.join(field_conditions)})")
            print("学期")
        else:
            # 如果 first_connector 为 NOT
            if first_connector == "NOT":
                # 如果是第一个条件，不加 AND，因为只有可能 OR 是第一个条件
                if not query_conditions:
                    print("第一个条件")
                    query_conditions.append(f"{' '.join(field_conditions)}")
                else:
                    print(first_connector)
                    query_conditions.append(f"AND ({' '.join(field_conditions)})")
            else:
                query_conditions.append(f"{first_connector} ({' '.join(field_conditions)})")

    # 生成 SQL 查询语句
    where_clause = ' '.join(query_conditions)  # 使用空格连接所有字段的查询条件
    query = f"SELECT * FROM course_all WHERE {where_clause} ORDER BY `学期` ASC, `课程序号` ASC"  # 查询 course_all 表


    print(query) # 打印查询语句

    conn = mysql.connector.connect(**db_config) # 连接数据库
    cursor = conn.cursor() # 创建游标
    try:
        cursor.execute(query) # 执行查询
        results = cursor.fetchall() # 获取查询结果
    except Exception as e:
        response_data['message'] = "检索出错啦！<br>生成的 SQL 语句为：<br>" + query
        response_data['status'] = 400
        response_data['data'] = []
        return jsonify(response_data), 400

    # 给返回的 results 添加字段名
    results = [dict(zip(cursor.column_names, row)) for row in results]

    cursor.close() # 关闭游标
    conn.close()    # 关闭数据库连接

    # 返回查询结果

    response_data['data'] = results
    response_data['status'] = 200
    response_data['message'] = "检索成功！"

    return jsonify(response_data) # 返回 JSON 格式的查询结果


if __name__ == '__main__': # 如果当前脚本被直接运行
    app.run(debug=True) # 启动 Flask 应用
