---
title: OpenClaw
emoji: 🦀
colorFrom: green
colorTo: blue
sdk: docker
pinned: false
short_description: 麻辣小龙虾
---


# 请先确认你Rclone的文件中的openclaw文件夹存在
如需要设置TG机器人可将start-openclaw.sh中的示例代码进行修改并使用
注意Telegarm的设置为pairing即配对模式。需要向机器人发送/pair获取配对代码，然后在龙虾中执行配对代码即可   


Check out the configuration reference at https://huggingface.co/docs/hub/spaces-config-reference   
Secrets：   
RCLONE_CONF:rclone配置内容   
OPENAI_API_BASE:AI的BaseURL通常是https://xxxxx.com/v1   
OPENAI_API_KEY:AI的key目前代码中用的是OpenAI的   
MODEL:AI调用的模型   
OPENCLAW_GATEWAY_PASSWORD:龙虾前台管理密码   
