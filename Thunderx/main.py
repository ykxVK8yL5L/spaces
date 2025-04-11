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
    Body,
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


@front_router.get("/", response_class=HTMLResponse)
async def home(request: Request):
    return templates.TemplateResponse("index.html", {"request": request})


@api_router.post("/files")
async def get_files(item: FileRequest):
    return await THUNDERX_CLIENT.file_list(
        item.size, item.parent_id, item.next_page_token, item.additional_filters
    )


@api_router.get("/files/{file_id}")
async def get_file_info(file_id: str):
    return await THUNDERX_CLIENT.get_download_url(file_id)


@api_router.post("/offline")
async def offline(item: OfflineRequest):
    return await THUNDERX_CLIENT.offline_download(
        item.file_url, item.parent_id, item.name
    )


@api_router.get("/userinfo")
async def userinfo():
    return THUNDERX_CLIENT.get_user_info()


app.include_router(front_router)
app.include_router(api_router)
