from io import BytesIO
import os
import random
import string
import sys
import argparse
from huggingface_hub import HfApi

parser = argparse.ArgumentParser(description="抱脸SDK创建N8n")
parser.add_argument(
    "--token",
    type=str,
    required=True,
    help="抱脸的Token,需要写权限",
    default="",
)
parser.add_argument("--image", help="Docker镜像地址", default="")
parser.add_argument("--key", help="N8N的N8N_ENCRYPTION_KEY", default="")
parser.add_argument("--password", help="Code Server的管理密码", default="")
parser.add_argument("--rclone_conf_path", help="Rclone配置", default="")

args = parser.parse_args()

def generate_random_string(length=10):
    if length < 1:
        return ""

    chars = string.ascii_letters + string.digits  # 包含字母和数字
    # 1. 先强制加入一个随机字母
    mandatory_letter = random.choice(string.ascii_letters)

    # 2. 生成剩余的字符
    remaining_chars = random.choices(chars, k=length - 1)

    # 3. 将强制字母加入随机位置
    full_chars = remaining_chars + [mandatory_letter]
    random.shuffle(full_chars)

    return "".join(full_chars)


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
    api = HfApi(token=token)
    user_info = api.whoami()
    if not user_info.get("name"):
        print("未获取到用户名信息，程序退出。")
        sys.exit(1)
    userid = user_info.get("name")
    image = "ghcr.io/ykxvk8yl5l/spaces/n8n:latest"
    if len(args.image) > 0:
        image = args.image

    key = "952a63942ad0c68a"
    if len(args.key) > 0:
        key = args.key
    rclone_conf_path = "~/.config/rclone/rclone.conf"
    rclone_conf = ""
    if len(args.rclone_conf_path) > 0:
        rclone_conf_path = args.rclone_conf_path
    password = ""
    if len(args.password) > 0:
        password = args.password
    rclone_conf = read_file_if_not_empty(rclone_conf_path)
    #space_name = generate_random_string(2)
    space_name = "n8n"
    repoid = f"{userid}/{space_name}"

    # readme.md的字符串内容
    readme_content = f"""
---
title: {space_name}
emoji: ⚡
colorFrom: red
colorTo: gray
sdk: docker
app_port: 5700
pinned: false
---
Check out the configuration reference at https://huggingface.co/docs/hub/spaces-config-reference
    """
    # 转成 file-like object（以字节形式）
    readme_obj = BytesIO(readme_content.encode("utf-8"))

    api.create_repo(
        repo_id=repoid,
        repo_type="space",
        space_sdk="docker",
        space_secrets=[
            {"key": "N8N_ENCRYPTION_KEY", "value": key},
            {"key": "RCLONE_CONF", "value": rclone_conf},
            {"key": "ADMIN_PASSWORD", "value": password},
        ],
        space_variables=[
            {"key": "GENERIC_TIMEZONE", "value": "Asia/Shanghai"},
            {"key": "TZ", "value": "Asia/Shanghai"},
            {"key": "NODE_FUNCTION_ALLOW_BUILTIN", "value": "*"},
            {"key": "N8N_ENFORCE_SETTINGS_FILE_PERMISSIONS", "value": "true"},
            {"key": "WEBHOOK_URL", "value": f"https://{userid}-{space_name}.hf.space"},
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
