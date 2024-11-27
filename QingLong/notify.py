#!/usr/bin/python
import sqlite3
import os

conn = sqlite3.connect('/ql/data/db/database.sqlite')
# 创建一个游标对象
cursor = conn.cursor()

# 从环境变量获取新的信息值
new_info_value = os.getenv('NOTIFY_CONFIG')
# 要更新的信息条件
update_condition = 'notification'

# 执行更新操作
try:
    cursor.execute("""
        UPDATE Auths 
        SET info = ? 
        WHERE id = (
            SELECT MIN(id) 
            FROM Auths 
            WHERE type = ?
        )
    """, (new_info_value, update_condition))

    # 提交事务
    conn.commit()

    # 检查影响的行数
    if cursor.rowcount > 0:
        print(f"成功更新通知设置")
    else:
        print("没有找到符合条件的记录。")

except sqlite3.Error as e:
    print(f"发生错误: {e}")

finally:
    # 关闭游标和连接
    cursor.close()
    conn.close()


