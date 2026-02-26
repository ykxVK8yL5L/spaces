#!/bin/sh

#获取登陆cookie:需要设置用户名和密码
USERNAME="admin"
echo "开始登陆......"
BHToken=$(
curl -c cookies.txt -s -D - -o /dev/null \
'http://localhost:8052/api/v1/auth/login' \
-H 'content-type: application/json' \
--data-raw "{\"username\":\"admin\",\"password\":\"$ADMIN_PASSWORD\"}" \
| awk -F'[=;]' '/Set-Cookie: BHToken=/{print $2}'
)
echo "登陆成功Token为:$BHToken,开始请求备份"

#用cookie请求生成备份
BACKUP_RESPON=$(
    curl -b cookies.txt -X POST "http://localhost:8052/api/v1/settings/backup" -H 'content-type: application/json' 
)
echo "请求备份完成，等待10秒后备份"

#休眠一段时间后再复制，具体根据文件内容多少来定
sleep 10
#获取备份目录下的最新文件
latest_file=$(ls -t /app/data/backups | head -n 1)
rclone delete $REMOTE_FOLDER  #删除之前的备份，如需要保留之前的备份可以注释掉
rclone copy /app/data/backups/$latest_file $REMOTE_FOLDER
rm /app/data/backups/$latest_file #删除生成的备份文件，如需保留可注释
rm -f cookies.txt
