---
title: Qinglong
emoji: 🏃
colorFrom: yellow
colorTo: purple
sdk: docker
pinned: false
app_port: 5700
---
演示视频： https://youtu.be/fGxS3GOmI3Y?si=8LtETCU_EeSpbp8U  



Space secrets(Private)说明:   
ADMIN_USERNAME:登陆用户名   
ADMIN_PASSWORD:登陆密码   
RCLONE_CONF:rclone配置内容，可选   
NOTIFY_CONFIG:通知配置内容，可选，需要到通知设置里抓包

## 20250308更新   
加入coder-server 密码同青龙面板 管理地址为:/coder/

## 20241020更新   
加入services.json文件，可自行参照修改，添加更多的邮箱服务支持：serv00仅供参考，需要修改邮箱服务器为注册时收到的地址   
新添加邮箱:   
cock.li-Cock   
serv00.com-Serv00   
mail.com-Mail.com   

## 注意 
docker-entrypoint-rclone.sh为可以在容器重置时恢复数据的配置，需要设置上面的RCLONE_CONF  
其中名称为huggingface,文件夹为:/qinglong   
rclone远程文件夹为huggingface:/qinglong   
进入后台后需要新建个同步的计划任务代码为:  
```
rclone delete huggingface:/qinglong/db/database.sqlite && rclone sync /ql/data huggingface:/qinglong
```
备份时间可以根据自己的情况设置我这里设置每天1点更新: 0 1 1 * * *    

如需安装linux依赖需要修改Dockerfile手动安装：搜索sshpass在下面继续添加即可。  

Check out the configuration reference at https://huggingface.co/docs/hub/spaces-config-reference
