import os
import asyncio
import json
import logging

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

THUNDERX_CLIENT = None


async def log_token(THUNDERX_CLIENT, extra_data):
    logging.info(f"Token: {THUNDERX_CLIENT.encoded_token}, Extra Data: {extra_data}")


@app.on_event("startup")
async def init_client():
    global THUNDERX_CLIENT
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


@front_router.get(
    "/",
    response_class=HTMLResponse,
    summary="前台页面",
    description="前台管理页面，需要在设置里设置SECRET_TOKEN才能正常请求",
    tags=["前端"],
)
async def home(request: Request):
    return templates.TemplateResponse("index.html", {"request": request})


@api_router.post("/files", summary="文件列表", description="获取文件列表", tags=["文件"])
async def get_files(item: FileRequest):
    return await THUNDERX_CLIENT.file_list(
        item.size, item.parent_id, item.next_page_token, item.additional_filters
    )


@api_router.post("/file_star_list", summary="加星文件列表", description="加星文件列表", tags=["文件"])
async def file_star_list( size: int = Query(default=100, title="显示数量", description="显示数量"), next_page_token: str | None = Query(default=None, title="分页Token", description="分页Token")):
    return await THUNDERX_CLIENT.file_star_list(
        size,next_page_token
    )

@api_router.get(
    "/files/{file_id}", summary="文件信息", description="获取文件信息", tags=["文件"]
)
async def get_file_info(file_id: str = Path(..., title="文件ID", description="文件ID") ):
    return await THUNDERX_CLIENT.get_download_url(file_id)


@api_router.delete(
    "/files/{file_id}", summary="删除文件", description="删除文件", tags=["文件"]
)
async def delete_file_info(file_id: str = Path(..., title="文件ID", description="文件ID") ):
    return await THUNDERX_CLIENT.delete_forever([file_id])

@api_router.post(
    "/file_rename", summary="重命名文件", description="重命名文件", tags=["文件"]
)
async def file_rename(
    file_id: str = Query(title="文件ID", description="文件ID"),
    new_file_name: str = Query(title="新文件名", description="新文件名"),
    ):
    return await THUNDERX_CLIENT.file_rename(file_id,new_file_name)


@api_router.post(
    "/file_batch_copy", summary="批量复制文件", description="批量复制文件", tags=["文件"]
)
async def file_batch_copy(
    ids: List[str] = Body(title="文件ID列表", description="文件ID列表"),
    to_parent_id: str = Query(title="复制到的文件夹id, 默认为根目录", description="复制到的文件夹id, 默认为根目录"),
    ):
    return await THUNDERX_CLIENT.file_batch_copy(ids,to_parent_id)


@api_router.post(
    "/file_batch_move", summary="批量移动文件", description="批量移动文件", tags=["文件"]
)
async def file_batch_move(
    ids: List[str] = Body(title="文件ID列表", description="文件ID列表"),
    to_parent_id: str = Query(title="移动到的文件夹id, 默认为根目录", description="移动到的文件夹id, 默认为根目录"),
    ):
    return await THUNDERX_CLIENT.file_batch_move(ids,to_parent_id)


@api_router.post(
    "/create_folder", summary="新建文件夹", description="新建文件夹", tags=["文件"]
)
async def create_folder(
    name: str = Query(title="文件夹名称", description="文件夹名称"),
    parent_id: str = Query(title="父文件夹id, 默认创建到根目录", description="父文件夹id, 默认创建到根目录"),
    ):
    return await THUNDERX_CLIENT.create_folder(name,parent_id)


@api_router.post(
    "/delete_to_trash", summary="将文件夹、文件移动到回收站", description="将文件夹、文件移动到回收站", tags=["文件"]
)
async def delete_to_trash(ids: List[str] = Body(title="文件ID列表", description="文件ID列表") ):
    return await THUNDERX_CLIENT.delete_to_trash(ids)


@api_router.post(
    "/delete_forever", summary="将文件夹、文件彻底删除", description="将文件夹、文件彻底删除", tags=["文件"]
)
async def delete_forever(ids: List[str] = Body(title="文件ID列表", description="文件ID列表") ):
    return await THUNDERX_CLIENT.delete_forever(ids)


@api_router.post(
    "/untrash", summary="将文件夹、文件移出回收站", description="将文件夹、文件移出回收站", tags=["文件"]
)
async def untrash(ids: List[str] = Body(title="文件ID列表", description="文件ID列表") ):
    return await THUNDERX_CLIENT.untrash(ids)


@api_router.post(
    "/file_batch_star", summary="批量给文件加星标", description="批量给文件加星标", tags=["文件"]
)
async def file_batch_star(ids: List[str] = Body(title="文件ID列表", description="文件ID列表") ):
    return await THUNDERX_CLIENT.file_batch_star(ids)


@api_router.post(
    "/file_batch_unstar", summary="批量给文件加星标", description="批量给文件加星标", tags=["文件"]
)
async def file_batch_unstar(ids: List[str] = Body(title="文件ID列表", description="文件ID列表") ):
    return await THUNDERX_CLIENT.file_batch_unstar(ids)



@api_router.post(
    "/emptytrash", summary="清空回收站", description="清空回收站【慎用】", tags=["文件"]
)
async def emptytrash():
    return await THUNDERX_CLIENT.emptytrash()


##############  分享 ################
@api_router.post(
    "/get_share_list", summary="获取账号分享列表", description="获取账号分享列表", tags=["分享"]
)
async def get_share_list(page_token: str | None = Query(default=None, title="分页Token", description="分页Token") ):
    return await THUNDERX_CLIENT.get_share_list(page_token)


@api_router.post(
    "/file_batch_share", summary="创建分享", description="创建分享", tags=["分享"]
)
async def file_batch_share(
    ids: List[str] = Body(default=None, title="文件ID列表", description="文件ID列表"),
    need_password: bool | None = Query(default=False, title="是否需要密码", description="是否需要密码"), 
    expiration_days:int | None = Query(default=-1, title="过期时间", description="过期时间【天数，默认永远】"), 
    ):
    return await THUNDERX_CLIENT.file_batch_share(ids,need_password,expiration_days)


@api_router.post(
    "/share_batch_delete", summary="取消分享", description="取消分享", tags=["分享"]
)
async def share_batch_delete(ids: List[str] = Body(title="文件ID列表", description="文件ID列表")):
    return await THUNDERX_CLIENT.share_batch_delete(ids)

@api_router.post(
    "/get_share_folder", summary="获取分享信息", description="获取分享信息", tags=["分享"]
)
async def get_share_folder(
    share_id: str = Query(title="分享ID", description="分享ID"), 
    pass_code_token: str | None = Query(default=None,title="密码", description="密码"),
    parent_id:str | None=Query(default=None,title="父ID", description="父ID")):
    return await THUNDERX_CLIENT.get_share_folder(share_id,pass_code_token,parent_id)


@api_router.post(
    "/restore", summary="转存分享文件", description="转存分享文件", tags=["分享"]
)
async def restore(share_id: str, pass_code_token: str | None = None,file_ids:List[str] | None=None):
    return await THUNDERX_CLIENT.restore(share_id,pass_code_token,file_ids)




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
    "/delete_tasks", summary="删除离线任务", description="删除离线任务", tags=["离线任务"]
)
async def delete_tasks(task_ids: List[str], delete_files: bool = False):
    return await THUNDERX_CLIENT.delete_tasks(
        task_ids,delete_files
    )


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
