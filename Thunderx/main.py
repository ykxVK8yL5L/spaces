import os
import asyncio
import json
import logging
from telegram import Update, InlineKeyboardButton, InlineKeyboardMarkup
from telegram.ext import (
    Application,
    CommandHandler,
    ConversationHandler,
    MessageHandler,
    CallbackQueryHandler,
    CallbackContext,
    ContextTypes,
    filters,
)
import httpx

from pikpakapi import PikPakApi

from typing import Union, Any, Dict, List, Optional
from fastapi import (
    FastAPI,
    APIRouter,
    Depends,
    Request,
    Query,
    Body,
    Path,
    Response,
    HTTPException,
    status,
    Request,
)
from fastapi.responses import StreamingResponse, HTMLResponse, JSONResponse
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from fastapi.templating import Jinja2Templates
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel, Extra


class PostRequest(BaseModel):
    class Config:
        extra = Extra.allow


class FileRequest(BaseModel):
    size: int = 100
    parent_id: str | None = ""
    next_page_token: str | None = ""
    additional_filters: Dict | None = {}

    class Config:
        extra = Extra.allow


class OfflineRequest(BaseModel):
    file_url: str = ""
    parent_id: str | None = ""
    name: str | None = ""

    class Config:
        extra = Extra.allow


security = HTTPBearer()
# SECRET_TOKEN = "SECRET_TOKEN"
SECRET_TOKEN = os.getenv("SECRET_TOKEN")
if SECRET_TOKEN is None:
    raise ValueError("请在环境变量中设置SECRET_TOKEN，确保安全!")

THUNDERX_USERNAME = os.getenv("THUNDERX_USERNAME")
if THUNDERX_USERNAME is None:
    raise ValueError("请在环境变量中设置THUNDERX_USERNAME，用户名【邮箱】用来登陆!")


THUNDERX_PASSWORD = os.getenv("THUNDERX_PASSWORD")
if THUNDERX_PASSWORD is None:
    raise ValueError("请在环境变量中设置THUNDERX_PASSWORD，密码用来登陆!")

PROXY_URL = os.getenv("PROXY_URL")
TG_BOT_TOKEN = os.getenv("TG_BOT_TOKEN")
TG_WEBHOOK_URL = os.getenv("TG_WEBHOOK_URL")


async def verify_token(
    request: Request, credentials: HTTPAuthorizationCredentials = Depends(security)
):
    # excluded_paths = ["/"]  # 需要排除的路径列表
    # if request.url.path in excluded_paths:
    #     return  # 直接跳过验证

    # 验证Bearer格式
    if credentials.scheme != "Bearer":
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid authentication scheme",
        )

    # 验证令牌内容
    if credentials.credentials != SECRET_TOKEN:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid or expired token"
        )


def format_bytes(size: int) -> str:
    # 预设单位
    units = ["B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"]

    # 确保字节数是正数
    if size < 0:
        raise ValueError("字节大小不能为负数")

    # 选择合适的单位
    unit_index = 0
    while size >= 1024 and unit_index < len(units) - 1:
        size /= 1024.0
        unit_index += 1

    # 格式化输出，保留两位小数
    return f"{size:.2f} {units[unit_index]}"


app = FastAPI()


app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

api_router = APIRouter(dependencies=[Depends(verify_token)])
front_router = APIRouter()

templates = Jinja2Templates(
    directory="templates", variable_start_string="{[", variable_end_string="]}"
)


async def log_token(THUNDERX_CLIENT, extra_data):
    logging.info(f"Token: {THUNDERX_CLIENT.encoded_token}, Extra Data: {extra_data}")


THUNDERX_CLIENT = None
TG_BOT_APPLICATION = None
TG_BASE_URL = "https://tg.alist.dpdns.org/bot"


###################TG机器人功能区###################
# ❗❗❗❗❗❗❗❗❗注意TG机器人callbackdata不能超过64位，否则会报无效按钮的错误
# 定义命令处理函数
async def start(update: Update, context):
    commands = (
        "🚀欢迎使用我的机器人！\n\n"
        "📋可用命令:\n"
        "•直接发送magent:开头的磁力将直接离线下载\n"
        "•直接发送分享码:开头的分享ID将直接离线下载\n"
        "•/tasks - 查看下载任务\n"
        "•/files - 查看文件列表\n"
        "•/shares - 查看分享列表\n"
        "•/quota - 查看存储空间\n"
        "•/emptytrash - 清空回收站\n"
        "•/help - 获取帮助信息\n"
    )
    await update.message.reply_text(commands)


async def help(update: Update, context):
    commands = (
        "🚀欢迎使用我的机器人！\n\n"
        "📋可用命令:\n"
        "•直接发送magent:开头的磁力将直接离线下载\n"
        "•直接发送分享码:开头的分享ID将直接离线下载\n"
        "•/tasks - 查看下载任务\n"
        "•/files - 查看文件列表\n"
        "•/shares - 查看分享列表\n"
        "•/quota - 查看存储空间\n"
        "•/emptytrash - 清空回收站\n"
        "•/help - 获取帮助信息\n"
    )
    await update.message.reply_text(commands)


async def quota(update: Update, context):
    """
    返回信息
    {
      "kind": "drive#about",
      "quota": {
        "kind": "drive#quota",
        "limit": "72057604737418240",
        "usage": "18700975438",
        "usage_in_trash": "0",
        "play_times_limit": "2",
        "play_times_usage": "0",
        "is_unlimited": true
      },
      "expires_at": "2026-04-08T21:47:59.000+08:00",
      "quotas": {}
    }
    """
    quota_info = await THUNDERX_CLIENT.get_quota_info()
    if quota_info["quota"]["usage"] is None:
        await update.message.reply_text("❌未找到使用信息，请稍后再试！")
    else:
        await update.message.reply_text(
            f"✅使用信息:\n{format_bytes(int(quota_info['quota']['usage']))}/{format_bytes(int(quota_info['quota']['limit']))}\n⏰到期时间:\n{quota_info['expires_at']}"
        )


async def tg_emptytrash(update: Update, context):
    """
    返回信息
    """
    result = await THUNDERX_CLIENT.emptytrash()
    if result["task_id"] is None:
        await update.message.reply_text("❌未成功创建任务，请稍后重试!!")
    else:
        await update.message.reply_text(f"✅操作成功")


# 消息处理
async def handle_message(update: Update, context: ContextTypes.DEFAULT_TYPE):
    text = update.message.text
    if text.lower().startswith("magnet:"):
        result = await THUNDERX_CLIENT.offline_download(text, "", "")
        if result["task"]["id"] is not None:
            await update.message.reply_text(f"✅操作成功")
        else:
            await update.message.reply_text(f"❌未成功创建任务，请稍后重试!!")
    elif text.lower().startswith("分享码:"):
        share_id = text.split(":")[1]
        result = await THUNDERX_CLIENT.restore(share_id, None, None)
        if isinstance(result, str):
            await update.message.reply_text(f"❌未成功创建任务:{result}，请稍后重试!!")
        else:
            await update.message.reply_text(f"操作结果:{result['share_status_text']}")

    else:
        await update.message.reply_text(f"收到不支持的消息:{text}")


# 消息处理
async def handle_copy_text(update: Update, context: CallbackContext):
    query = update.callback_query
    await query.answer()

    # 获取操作类型和文件 ID
    action, text = (query.data.split(":")[0], query.data.split(":")[1])
    await query.edit_message_text(f"{text}")


#################### 分享操作 #############################
async def tg_show_shares(update: Update, context: CallbackContext):
    shares = await THUNDERX_CLIENT.get_share_list("")
    keyboard = []

    if shares["data"] is None:
        await update.message.reply_text("❌未找到分享!!")
    else:
        # 为每个文件创建按钮和操作选项
        for share in shares["data"]:
            keyboard.append(
                [
                    InlineKeyboardButton(
                        f"{share['title']}",
                        callback_data=f"copy_text:{share['share_id']}",
                    ),
                    InlineKeyboardButton(
                        f"{share['share_id']}",
                        callback_data=f"copy_text:{share['share_id']}",
                    ),
                    InlineKeyboardButton(
                        f"取消",
                        callback_data=f"del_s:{share['share_id']}",
                    ),
                ]
            )
        reply_markup = InlineKeyboardMarkup(keyboard)
        await update.message.reply_text(f"📋分享列表:", reply_markup=reply_markup)


# 处理任务操作的回调
async def handle_share_operation(update: Update, context: CallbackContext):
    query = update.callback_query
    await query.answer()

    # 获取操作类型和文件 ID
    action, share_id = (query.data.split(":")[0], query.data.split(":")[1])

    # 需要确认的操作
    if action in ["del_s"]:
        # 生成确认消息
        keyboard = [
            [InlineKeyboardButton("确认", callback_data=f"yes_s_{action}:{share_id}")],
            [InlineKeyboardButton("取消", callback_data=f"no_s_{action}:{share_id}")],
        ]
        reply_markup = InlineKeyboardMarkup(keyboard)
        await query.edit_message_text(
            f"你确定要{action}分享 {share_id} 吗？", reply_markup=reply_markup
        )


async def handle_share_confirmation(update: Update, context: CallbackContext):
    query = update.callback_query
    await query.answer()

    # 获取确认操作的类型和文件 ID
    action, share_id = (query.data.split(":")[0], query.data.split(":")[1])

    if action == "yes_s_del_s":
        await THUNDERX_CLIENT.share_batch_delete([share_id])
        await query.edit_message_text(f"✅分享 {share_id} 已取消。")


async def handle_share_cancel(update: Update, context: CallbackContext):
    query = update.callback_query
    await query.answer()
    await query.edit_message_text(f"操作已取消")


#################### 文件操作 #############################


async def tg_show_files(update: Update, context: CallbackContext):
    files = await THUNDERX_CLIENT.file_list(100, "", "", {})
    keyboard = []

    if files["files"] is None:
        await update.message.reply_text("❌未找到文件!!")
    else:
        # 为每个文件创建按钮和操作选项
        for file in files["files"]:
            if file["kind"].lower() == "drive#folder":
                keyboard.append(
                    [
                        InlineKeyboardButton(
                            f"查看📁: {file['name']}",
                            callback_data=f"ls_f:{file['id']}:{file['parent_id']}",
                        ),
                        InlineKeyboardButton(
                            f"删除",
                            callback_data=f"del_f:{file['id']}:{file['parent_id']}",
                        ),
                        InlineKeyboardButton(
                            f"分享",
                            callback_data=f"sh_f:{file['id']}:{file['parent_id']}",
                        ),
                    ]
                )
            else:
                keyboard.append(
                    [
                        InlineKeyboardButton(
                            f"下载📄: {file['name']}",
                            callback_data=f"dw_f:{file['id']}:{file['parent_id']}",
                        ),
                        InlineKeyboardButton(
                            f"删除",
                            callback_data=f"del_f:{file['id']}:{file['parent_id']}",
                        ),
                        InlineKeyboardButton(
                            f"分享",
                            callback_data=f"sh_f:{file['id']}:{file['parent_id']}",
                        ),
                    ]
                )

        reply_markup = InlineKeyboardMarkup(keyboard)
        await update.message.reply_text(f"📋文件列表:", reply_markup=reply_markup)


async def handle_file_confirmation(update: Update, context: CallbackContext):
    query = update.callback_query
    await query.answer()

    # 获取确认操作的类型和文件 ID
    action, file_id = (query.data.split(":")[0], query.data.split(":")[1])

    if action == "yes_f_del_f":
        await THUNDERX_CLIENT.delete_forever([file_id])
        await query.edit_message_text(f"✅文件 {file_id} 已删除。")


async def handle_file_cancel(update: Update, context: CallbackContext):
    query = update.callback_query
    await query.answer()
    # 获取取消操作的类型和文件 ID
    # action, file_id, parent_id = (
    #     query.data.split(":")[0],
    #     query.data.split(":")[1],
    #     query.data.split(":")[2],
    # )
    # 返回文件夹列表
    await query.edit_message_text(f"操作已取消")


# 处理任务操作的回调
async def handle_file_operation(update: Update, context: CallbackContext):
    query = update.callback_query
    await query.answer()

    # 获取操作类型和文件 ID
    action, file_id, parent_id = (
        query.data.split(":")[0],
        query.data.split(":")[1],
        query.data.split(":")[2],
    )

    # 需要确认的操作
    if action in ["del_f"]:
        # 生成确认消息
        keyboard = [
            [InlineKeyboardButton("确认", callback_data=f"yes_f_{action}:{file_id}")],
            [InlineKeyboardButton("取消", callback_data=f"no_f_{action}:{file_id}")],
        ]
        reply_markup = InlineKeyboardMarkup(keyboard)
        await query.edit_message_text(
            f"你确定要{action}文件 {file_id} 吗？", reply_markup=reply_markup
        )
    else:
        # 不需要确认的操作，直接处理
        await perform_file_action(update, context, action, file_id, parent_id)


async def perform_file_action(
    update: Update, context: CallbackContext, action: str, file_id: str, parent_id: str
):

    if action == "ls_f":
        files = await THUNDERX_CLIENT.file_list(100, file_id, "", {})
        keyboard = []

        if files["files"] is None:
            await update.message.reply_text("❌未找到文件!!")
        else:
            keyboard.append(
                [
                    InlineKeyboardButton(
                        f"↩️返回上级",
                        callback_data=f"ls_f:{parent_id}:{parent_id}",
                    ),
                ]
            )
            # 为每个文件创建按钮和操作选项
            for file in files["files"]:
                if file["kind"].lower() == "drive#folder":
                    keyboard.append(
                        [
                            InlineKeyboardButton(
                                f"查看📁: {file['name']}",
                                callback_data=f"ls_f:{file['id']}:{file['parent_id']}",
                            ),
                            InlineKeyboardButton(
                                f"删除",
                                callback_data=f"del_f:{file['id']}:{file['parent_id']}",
                            ),
                            InlineKeyboardButton(
                                f"分享",
                                callback_data=f"sh_f:{file['id']}:{file['parent_id']}",
                            ),
                        ]
                    )
                else:
                    keyboard.append(
                        [
                            InlineKeyboardButton(
                                f"下载📄: {file['name']}",
                                callback_data=f"dw_f:{file['id']}:{file['parent_id']}",
                            ),
                            InlineKeyboardButton(
                                f"删除",
                                callback_data=f"del_f:{file['id']}:{file['parent_id']}",
                            ),
                            InlineKeyboardButton(
                                f"分享",
                                callback_data=f"sh_f:{file['id']}:{file['parent_id']}",
                            ),
                        ]
                    )

            reply_markup = InlineKeyboardMarkup(keyboard)
            # await update.message.reply_text(f"📋文件列表:", reply_markup=reply_markup)
            await update.callback_query.edit_message_text(
                f"📋文件列表:", reply_markup=reply_markup
            )
    elif action == "dw_f":
        result = await THUNDERX_CLIENT.get_download_url(file_id)
        download_url = result["web_content_link"]
        for media in result["medias"]:
            if media["link"]["url"] is not None:
                download_url = media["link"]["url"]
                break
        if download_url is not None:
            await update.callback_query.edit_message_text(
                f"📋文件下载地址:{download_url}"
            )
        else:
            await update.callback_query.edit_message_text(f"❌未找到文件下载地址!!")
    elif action == "sh_f":
        result = await THUNDERX_CLIENT.file_batch_share([file_id], False, -1)
        share_id = result["share_id"]
        if share_id is not None:
            await update.callback_query.edit_message_text(f"分享码:{share_id}")
        else:
            await update.callback_query.edit_message_text(f"❌分享失败!!")


#################### 离线任务处理 ##########################
# 确认操作的回调
async def handle_task_confirmation(update: Update, context: CallbackContext):
    query = update.callback_query
    await query.answer()

    # 获取确认操作的类型和文件 ID
    action, task_id = query.data.split(":")[0], query.data.split(":")[1]

    if action == "confirm_task_delete_task":
        await THUNDERX_CLIENT.delete_tasks([task_id])
        await query.edit_message_text(f"✅任务 {task_id} 已删除。")


async def handle_task_cancel(update: Update, context: CallbackContext):
    query = update.callback_query
    await query.answer()
    # 获取取消操作的类型和文件 ID
    action, file_id = query.data.split(":")[0], query.data.split(":")[1]
    # 返回文件夹列表
    await query.edit_message_text(f"操作已取消")


async def tg_show_task(update: Update, context: CallbackContext):
    """
    {
      "tasks": [
        {
          "kind": "drive#task",
          "id": "VONrJ4Skj4Qs7ALhxXlFudfJAA",
          "name": "Billy Elliot (2000) 1080p (Deep61)[TGx]",
          "type": "offline",
          "user_id": "2000403406",
          "statuses": [],
          "status_size": 2,
          "params": {
            "folder_type": "",
            "predict_type": "1",
            "url": "magnet:?xt=urn:btih:96451E6F1ADBC8827B43621B74EDB30DF45012D6"
          },
          "file_id": "VONrJ4dZ8zf9KVWQuVEKmW8sTT",
          "file_name": "Billy Elliot (2000) 1080p (Deep61)[TGx]",
          "file_size": "3748030421",
          "message": "Task timeout",
          "created_time": "2025-04-15T10:38:54.320+08:00",
          "updated_time": "2025-04-17T10:39:12.189+08:00",
          "third_task_id": "",
          "phase": "PHASE_TYPE_ERROR",
          "progress": 0,
          "icon_link": "https://backstage-img.xunleix.com/65d616355857aef8af40b89f187a8cf2770cb0ce",
          "callback": "",
          "reference_resource": {
            "@type": "type.googleapis.com/drive.ReferenceFile",
            "kind": "drive#folder",
            "id": "VONrJ4dZ8zf9KVWQuVEKmW8sTT",
            "parent_id": "VONS0fwXf3FNvt-g_IlMVKPxAA",
            "name": "Billy Elliot (2000) 1080p (Deep61)[TGx]",
            "size": "3748030421",
            "mime_type": "",
            "icon_link": "https://backstage-img.xunleix.com/65d616355857aef8af40b89f187a8cf2770cb0ce",
            "hash": "",
            "phase": "PHASE_TYPE_ERROR",
            "audit": null,
            "thumbnail_link": "",
            "params": {},
            "space": "",
            "medias": [],
            "starred": false,
            "tags": []
          },
          "space": ""
        }
      ],
      "next_page_token": "",
      "expires_in": 60,
      "expires_in_ms": 60000
    }
    """
    tasks = await THUNDERX_CLIENT.offline_list(
        size=100,
        next_page_token=None,
        phase=None,
    )
    keyboard = []

    if tasks["tasks"] is None:
        await update.message.reply_text("❌未找到任务!!")
    else:
        # 为每个文件创建按钮和操作选项
        for task in tasks["tasks"]:
            # 为每个文件添加操作按钮：删除
            keyboard.append(
                [
                    InlineKeyboardButton(
                        f"取消任务: {task['name']}",
                        callback_data=f"delete_task:{task['id']}",
                    ),
                ]
            )

        reply_markup = InlineKeyboardMarkup(keyboard)
        await update.message.reply_text(f"📋任务列表:", reply_markup=reply_markup)


# 处理任务操作的回调
async def handle_tasks_operation(update: Update, context: CallbackContext):
    query = update.callback_query
    await query.answer()

    # 获取操作类型和文件 ID
    action, task_id = query.data.split(":")

    # 需要确认的操作
    if action in ["delete_task"]:
        # 生成确认消息
        keyboard = [
            [
                InlineKeyboardButton(
                    "确认", callback_data=f"confirm_task_{action}:{task_id}"
                )
            ],
            [
                InlineKeyboardButton(
                    "取消", callback_data=f"cancel_task_{action}:{task_id}"
                )
            ],
        ]
        reply_markup = InlineKeyboardMarkup(keyboard)
        await query.edit_message_text(
            f"你确定要{action}任务 {task_id} 吗？", reply_markup=reply_markup
        )
    else:
        # 不需要确认的操作，直接处理
        await perform_task_action(update, context, action, task_id)


async def perform_task_action(
    update: Update, context: CallbackContext, action: str, file_id: str
):
    if action == "cancel_task":
        await update.callback_query.edit_message_text(f"你选择了取消任务：{file_id}")


@app.on_event("startup")
async def init_client():
    global THUNDERX_CLIENT
    global TG_BOT_APPLICATION
    if not os.path.exists("thunderx.txt"):
        THUNDERX_CLIENT = PikPakApi(
            username=THUNDERX_USERNAME,
            password=THUNDERX_PASSWORD,
            httpx_client_args=None,
            token_refresh_callback=log_token,
            token_refresh_callback_kwargs={"extra_data": "test"},
        )
        await THUNDERX_CLIENT.login()
        await THUNDERX_CLIENT.refresh_access_token()
        with open("thunderx.json", "w") as f:
            f.write(json.dumps(THUNDERX_CLIENT.to_dict(), indent=4))
    else:
        with open("thunderx.txt", "r") as f:
            data = json.load(f)
            THUNDERX_CLIENT = PikPakApi.from_dict(data)
            # await client.refresh_access_token()
            print(json.dumps(THUNDERX_CLIENT.get_user_info(), indent=4))

            print(
                json.dumps(
                    await THUNDERX_CLIENT.events(),
                    indent=4,
                )
            )

    if TG_BOT_TOKEN is None:
        print("未设置TG_BOT_TOKEN无法实现TG机器人功能！")
    else:
        TG_BOT_APPLICATION = (
            Application.builder().base_url(TG_BASE_URL).token(TG_BOT_TOKEN).build()
        )
        # await TG_BOT_APPLICATION.bot.delete_webhook()
        await TG_BOT_APPLICATION.bot.set_webhook(
            url=TG_WEBHOOK_URL, allowed_updates=Update.ALL_TYPES
        )
        TG_BOT_APPLICATION.add_handler(
            CallbackQueryHandler(handle_tasks_operation, pattern="^delete_task:")
        )
        # 处理取消任务操作
        TG_BOT_APPLICATION.add_handler(
            CallbackQueryHandler(handle_task_cancel, pattern="^cancel_task")
        )
        # 处理确认操作（确认删除、复制等）
        TG_BOT_APPLICATION.add_handler(
            CallbackQueryHandler(handle_task_confirmation, pattern="^confirm_task")
        )

        ########## 分享操作 ###############
        TG_BOT_APPLICATION.add_handler(
            CallbackQueryHandler(handle_share_operation, pattern="^del_s:")
        )
        # 处理取消任务操作
        TG_BOT_APPLICATION.add_handler(
            CallbackQueryHandler(handle_share_cancel, pattern="^no_s")
        )
        # 处理确认操作（确认删除、复制等）
        TG_BOT_APPLICATION.add_handler(
            CallbackQueryHandler(handle_share_confirmation, pattern="^yes_s")
        )

        ########## 文件操作 ###############

        TG_BOT_APPLICATION.add_handler(
            CallbackQueryHandler(
                handle_file_operation, pattern="^(del_f|ls_f|dw_f|sh_f):"
            )
        )
        # 处理取消任务操作
        TG_BOT_APPLICATION.add_handler(
            CallbackQueryHandler(handle_file_cancel, pattern="^no_f")
        )
        # 处理确认操作（确认删除、复制等）
        TG_BOT_APPLICATION.add_handler(
            CallbackQueryHandler(handle_file_confirmation, pattern="^yes_f")
        )

        TG_BOT_APPLICATION.add_handler(CommandHandler("start", start))
        TG_BOT_APPLICATION.add_handler(CommandHandler("help", help))
        TG_BOT_APPLICATION.add_handler(CommandHandler("quota", quota))
        TG_BOT_APPLICATION.add_handler(CommandHandler("emptytrash", tg_emptytrash))
        TG_BOT_APPLICATION.add_handler(CommandHandler("tasks", tg_show_task))
        TG_BOT_APPLICATION.add_handler(CommandHandler("files", tg_show_files))
        TG_BOT_APPLICATION.add_handler(CommandHandler("shares", tg_show_shares))
        # Message 消息处理相关的命令！
        TG_BOT_APPLICATION.add_handler(MessageHandler(filters.TEXT, handle_message))
        # 处理取消任务操作
        TG_BOT_APPLICATION.add_handler(
            CallbackQueryHandler(handle_copy_text, pattern="^copy_text")
        )
        await TG_BOT_APPLICATION.initialize()


# FastAPI 路由：接收来自 Telegram 的 Webhook 回调
@app.post("/webhook")
async def webhook(request: Request):
    # 从请求获取 JSON 数据
    data = await request.json()

    # 将 Telegram Update 转换为 Update 对象
    update = Update.de_json(data, TG_BOT_APPLICATION.bot)

    # 将 Update 对象传递给 Application 进行处理
    await TG_BOT_APPLICATION.process_update(update)

    return JSONResponse({"status": "ok"})


@front_router.get(
    "/",
    response_class=HTMLResponse,
    summary="前台页面",
    description="前台管理页面，需要在设置里设置SECRET_TOKEN才能正常请求",
    tags=["前端"],
)
async def home(request: Request):
    return templates.TemplateResponse("index.html", {"request": request})


@api_router.post(
    "/files", summary="文件列表", description="获取文件列表", tags=["文件"]
)
async def get_files(item: FileRequest):
    return await THUNDERX_CLIENT.file_list(
        item.size, item.parent_id, item.next_page_token, item.additional_filters
    )


@api_router.post(
    "/file_star_list", summary="加星文件列表", description="加星文件列表", tags=["文件"]
)
async def file_star_list(
    size: int = Query(default=100, title="显示数量", description="显示数量"),
    next_page_token: str | None = Query(
        default=None, title="分页Token", description="分页Token"
    ),
):
    return await THUNDERX_CLIENT.file_star_list(size, next_page_token)


@api_router.get(
    "/files/{file_id}", summary="文件信息", description="获取文件信息", tags=["文件"]
)
async def get_file_info(file_id: str = Path(..., title="文件ID", description="文件ID")):
    return await THUNDERX_CLIENT.get_download_url(file_id)


@api_router.delete(
    "/files/{file_id}", summary="删除文件", description="删除文件", tags=["文件"]
)
async def delete_file_info(
    file_id: str = Path(..., title="文件ID", description="文件ID")
):
    return await THUNDERX_CLIENT.delete_forever([file_id])


@api_router.post(
    "/file_rename", summary="重命名文件", description="重命名文件", tags=["文件"]
)
async def file_rename(
    file_id: str = Query(title="文件ID", description="文件ID"),
    new_file_name: str = Query(title="新文件名", description="新文件名"),
):
    return await THUNDERX_CLIENT.file_rename(file_id, new_file_name)


@api_router.post(
    "/file_batch_copy",
    summary="批量复制文件",
    description="批量复制文件",
    tags=["文件"],
)
async def file_batch_copy(
    ids: List[str] = Body(title="文件ID列表", description="文件ID列表"),
    to_parent_id: str = Query(
        title="复制到的文件夹id, 默认为根目录",
        description="复制到的文件夹id, 默认为根目录",
    ),
):
    return await THUNDERX_CLIENT.file_batch_copy(ids, to_parent_id)


@api_router.post(
    "/file_batch_move",
    summary="批量移动文件",
    description="批量移动文件",
    tags=["文件"],
)
async def file_batch_move(
    ids: List[str] = Body(title="文件ID列表", description="文件ID列表"),
    to_parent_id: str = Query(
        title="移动到的文件夹id, 默认为根目录",
        description="移动到的文件夹id, 默认为根目录",
    ),
):
    return await THUNDERX_CLIENT.file_batch_move(ids, to_parent_id)


@api_router.post(
    "/create_folder", summary="新建文件夹", description="新建文件夹", tags=["文件"]
)
async def create_folder(
    name: str = Query(title="文件夹名称", description="文件夹名称"),
    parent_id: str = Query(
        title="父文件夹id, 默认创建到根目录", description="父文件夹id, 默认创建到根目录"
    ),
):
    return await THUNDERX_CLIENT.create_folder(name, parent_id)


@api_router.post(
    "/delete_to_trash",
    summary="将文件夹、文件移动到回收站",
    description="将文件夹、文件移动到回收站",
    tags=["文件"],
)
async def delete_to_trash(
    ids: List[str] = Body(title="文件ID列表", description="文件ID列表")
):
    return await THUNDERX_CLIENT.delete_to_trash(ids)


@api_router.post(
    "/delete_forever",
    summary="将文件夹、文件彻底删除",
    description="将文件夹、文件彻底删除",
    tags=["文件"],
)
async def delete_forever(
    ids: List[str] = Body(title="文件ID列表", description="文件ID列表")
):
    return await THUNDERX_CLIENT.delete_forever(ids)


@api_router.post(
    "/untrash",
    summary="将文件夹、文件移出回收站",
    description="将文件夹、文件移出回收站",
    tags=["文件"],
)
async def untrash(ids: List[str] = Body(title="文件ID列表", description="文件ID列表")):
    return await THUNDERX_CLIENT.untrash(ids)


@api_router.post(
    "/file_batch_star",
    summary="批量给文件加星标",
    description="批量给文件加星标",
    tags=["文件"],
)
async def file_batch_star(
    ids: List[str] = Body(title="文件ID列表", description="文件ID列表")
):
    return await THUNDERX_CLIENT.file_batch_star(ids)


@api_router.post(
    "/file_batch_unstar",
    summary="批量给文件加星标",
    description="批量给文件加星标",
    tags=["文件"],
)
async def file_batch_unstar(
    ids: List[str] = Body(title="文件ID列表", description="文件ID列表")
):
    return await THUNDERX_CLIENT.file_batch_unstar(ids)


@api_router.post(
    "/emptytrash", summary="清空回收站", description="清空回收站【慎用】", tags=["文件"]
)
async def emptytrash():
    return await THUNDERX_CLIENT.emptytrash()


##############  分享 ################
@api_router.post(
    "/get_share_list",
    summary="获取账号分享列表",
    description="获取账号分享列表",
    tags=["分享"],
)
async def get_share_list(
    page_token: str | None = Query(
        default=None, title="分页Token", description="分页Token"
    )
):
    return await THUNDERX_CLIENT.get_share_list(page_token)


@api_router.post(
    "/file_batch_share", summary="创建分享", description="创建分享", tags=["分享"]
)
async def file_batch_share(
    ids: List[str] = Body(default=None, title="文件ID列表", description="文件ID列表"),
    need_password: bool | None = Query(
        default=False, title="是否需要密码", description="是否需要密码"
    ),
    expiration_days: int | None = Query(
        default=-1, title="过期时间", description="过期时间【天数，默认永远】"
    ),
):
    return await THUNDERX_CLIENT.file_batch_share(ids, need_password, expiration_days)


@api_router.post(
    "/share_batch_delete", summary="取消分享", description="取消分享", tags=["分享"]
)
async def share_batch_delete(
    ids: List[str] = Body(title="文件ID列表", description="文件ID列表")
):
    return await THUNDERX_CLIENT.share_batch_delete(ids)


@api_router.post(
    "/get_share_folder",
    summary="获取分享信息",
    description="获取分享信息",
    tags=["分享"],
)
async def get_share_folder(
    share_id: str = Query(title="分享ID", description="分享ID"),
    pass_code_token: str | None = Query(default=None, title="密码", description="密码"),
    parent_id: str | None = Query(default=None, title="父ID", description="父ID"),
):
    return await THUNDERX_CLIENT.get_share_folder(share_id, pass_code_token, parent_id)


@api_router.post(
    "/restore", summary="转存分享文件", description="转存分享文件", tags=["分享"]
)
async def restore(
    share_id: str, pass_code_token: str | None = None, file_ids: List[str] | None = None
):
    return await THUNDERX_CLIENT.restore(share_id, pass_code_token, file_ids)


##############  离线任务 ################


@api_router.get(
    "/offline", summary="离线任务列表", description="离线任务列表", tags=["离线任务"]
)
async def offline_list(size: int = 10000, next_page_token: str | None = None):
    return await THUNDERX_CLIENT.offline_list(
        size=size,
        next_page_token=next_page_token,
        phase=None,
    )


@api_router.post(
    "/offline", summary="添加离线任务", description="添加离线任务", tags=["离线任务"]
)
async def offline(item: OfflineRequest):
    return await THUNDERX_CLIENT.offline_download(
        item.file_url, item.parent_id, item.name
    )


@api_router.post(
    "/delete_tasks",
    summary="删除离线任务",
    description="删除离线任务",
    tags=["离线任务"],
)
async def delete_tasks(task_ids: List[str], delete_files: bool = False):
    return await THUNDERX_CLIENT.delete_tasks(task_ids, delete_files)


##############  账号 ################
@api_router.get(
    "/userinfo", summary="用户信息", description="获取用户登陆信息", tags=["账号"]
)
async def userinfo():
    return THUNDERX_CLIENT.get_user_info()


@api_router.get(
    "/quota", summary="空间使用信息", description="获取空间使用信息", tags=["账号"]
)
async def quota_info():
    return await THUNDERX_CLIENT.get_quota_info()


@api_router.get(
    "/invite_code", summary="查看邀请码", description="查看邀请码", tags=["账号"]
)
async def get_invite_code():
    return await THUNDERX_CLIENT.get_invite_code()


app.include_router(front_router)
app.include_router(api_router)
