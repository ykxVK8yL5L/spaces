from io import BytesIO
import os
import random
import string
import sys
import argparse
from huggingface_hub import HfApi


parser = argparse.ArgumentParser(
    description="抱脸SDK创建Directus【无头CMS】,目前仅支持postgres"
)
parser.add_argument(
    "--token",
    type=str,
    required=True,
    help="抱脸的Token,需要写权限",
    default="",
)
parser.add_argument("--userid", type=str, required=True, help="抱脸用户名", default="")
parser.add_argument("--image", help="Docker镜像地址", default="")
parser.add_argument("--rclone_conf_path", help="Rclone配置", default="")
parser.add_argument(
    "--db_connect_string", help="数据库链接；目前仅支持postgres", default=""
)


args = parser.parse_args()


def generate_random_string(length=10):
    chars = string.ascii_letters + string.digits  # 包含大小写字母和数字
    return "".join(random.choices(chars, k=length))


def read_file_if_not_empty(file_path):
    if not os.path.exists(file_path):
        print("Rclone配置文件不存在。")
        return ""

    if os.path.getsize(file_path) == 0:
        print("Rclone配置文件为空。")
        return ""

    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()

    if content.strip():
        return content
    else:
        print("Rclone配置文件内容全是空白。")
        return ""


if __name__ == "__main__":
    token = ""
    if len(args.token) > 0:
        token = args.token
    else:
        print("Token 不能为空")
        sys.exit(1)
        # raise ValueError("字符串不能为空！")
    userid = ""
    if len(args.userid) > 0:
        userid = args.userid
    else:
        print("userid 不能为空")
        sys.exit(1)
    image = "FROM directus/directus"
    if len(args.image) > 0:
        image = args.image
    rclone_conf_path = "~/.config/rclone/rclone.conf"
    rclone_conf = ""
    if len(args.rclone_conf_path) > 0:
        rclone_conf_path = args.rclone_conf_path
    rclone_conf = read_file_if_not_empty(rclone_conf_path)
    space_name = "directus"
    # space_name = generate_random_string(2)
    repoid = f"{userid}/{space_name}"
    db_connect_string = ""
    if len(args.db_connect_string) > 0:
        db_connect_string = args.db_connect_string

    # readme.md的字符串内容
    readme_content = f"""
---
title: {space_name}
emoji: ⚡
colorFrom: red
colorTo: indigo
sdk: docker
pinned: false
app_port: 8055
---
Check out the configuration reference at https://huggingface.co/docs/hub/spaces-config-reference
    """

    # 转成 file-like object（以字节形式）
    readme_obj = BytesIO(readme_content.encode("utf-8"))

    api = HfApi(token=token)
    api.create_repo(
        repo_id=repoid,
        repo_type="space",
        space_sdk="docker",
        space_secrets=[
            {"key": "DB_CLIENT", "value": "pg"},
            {"key": "DB_CONNECTION_STRING", "value": db_connect_string},
            {"key": "NODE_TLS_REJECT_UNAUTHORIZED", "value": 0},
        ],
    )

    api.upload_file(
        repo_id=repoid,
        path_in_repo="README.md",
        path_or_fileobj=readme_obj,
        repo_type="space",
    )
    dockerfile_content = f"FROM {image}"
    api.upload_file(
        repo_id=repoid,
        path_in_repo="Dockerfile",
        path_or_fileobj=BytesIO(dockerfile_content.encode("utf-8")),
        repo_type="space",
    )
