from io import BytesIO
import os
import random
import string
import sys
import argparse
from huggingface_hub import HfApi


parser = argparse.ArgumentParser(description="抱脸创建青龙面板")
parser.add_argument(
    "--token",
    type=str,
    required=True,
    help="抱脸的Token,需要写权限",
    default="",
)
parser.add_argument("--userid", type=str, required=True, help="抱脸用户名", default="")
parser.add_argument("--image", help="青龙docker镜像地址", default="")
parser.add_argument("--admin", help="青龙管理用户名", default="")
parser.add_argument("--password", help="青龙管理密码", default="")
parser.add_argument("--rclone_conf_path", help="Rclone配置", default="")


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
    image = "ghcr.io/ykxvk8yl5l/spaces/qinglong:latest"
    if len(args.image) > 0:
        image = args.image

    admin = "Hadmin123456"
    if len(args.admin) > 0:
        admin = args.admin
    password = "Hpassword654321"
    if len(args.password) > 0:
        password = args.password
    rclone_conf_path = "~/.config/rclone/rclone.conf"
    rclone_conf = ""
    if len(args.rclone_conf_path) > 0:
        rclone_conf_path = args.rclone_conf_path
    rclone_conf = read_file_if_not_empty(rclone_conf_path)
    space_name = generate_random_string(2)
    repoid = f"{userid}/{space_name}"

    # readme.md的字符串内容
    readme_content = f"""
---
title: {space_name}
emoji: 😻
colorFrom: red
colorTo: blue
sdk: docker
app_port: 5700
pinned: false
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
            {"key": "ADMIN_USERNAME", "value": admin},
            {"key": "ADMIN_PASSWORD", "value": password},
            {"key": "RCLONE_CONF", "value": rclone_conf},
        ],
    )
    # api.add_space_secret(repo_id=repoid, key="ADMIN_USERNAME", value=admin)
    # api.add_space_secret(repo_id=repoid, key="ADMIN_PASSWORD", value=password)
    # api.add_space_secret(repo_id=repoid, key="RCLONE_CONF", value=rclone_conf)
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
