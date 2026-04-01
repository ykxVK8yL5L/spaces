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

Check out the configuration reference at https://huggingface.co/docs/hub/spaces-config-reference   
Secrets：   
RCLONE_CONF:rclone配置内容   
OPENAI_API_BASE:AI的BaseURL通常是https://xxxxx.com/v1   
OPENAI_API_KEY:AI的key目前代码中用的是OpenAI的   
MODEL:AI调用的模型   
OPENCLAW_GATEWAY_PASSWORD:龙虾前台管理密码   

# 关于备份可以让AI建立定时任务执行以下代码
```
sh /app/sync.sh backup
```
## Telegram机器人注意事项：
- 抱脸不支持TG的api需要设置Telegram API Root URL
- WebHook地址为：https://抱脸用户名-抱脸空间名.hf.space/telegram/webhook
- WebHook端口为：8787 【固定的】
示例配置：

```
 "channels": {
    "telegram": {
      "enabled": true,
      "botToken": "机器人Token",
      "dmPolicy": "pairing",
      "apiRoot": "https://xxxxx.com",
      "groups": { "*": { "requireMention": true } },
      "webhookUrl": "https://抱脸用户名-抱脸空间名.hf.space/telegram/webhook",
      "webhookSecret": "$OPENCLAW_GATEWAY_PASSWORD",
      "webhookPath": "/telegram/webhook",
      "webhookHost": "0.0.0.0",
      "webhookPort": 8787,
    }
  }
```
