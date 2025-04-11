import hashlib
from uuid import uuid4
import time

CLIENT_ID = "ZQL_zwA4qhHcoe_2"
CLIENT_SECRET = "Og9Vr1L8Ee6bh0olFxFDRg"
CLIENT_VERSION = "1.06.0.2132"
PACKAG_ENAME = "com.thunder.downloader"
SDK_VERSION = "2.0.3.203100 "
APP_NAME = PACKAG_ENAME


def get_timestamp() -> int:
    """
    Get current timestamp.
    """
    return int(time.time() * 1000)


def device_id_generator() -> str:
    """
    Generate a random device id.
    """
    return str(uuid4()).replace("-", "")


SALTS = [
    "kVy0WbPhiE4v6oxXZ88DvoA3Q",
    "lON/AUoZKj8/nBtcE85mVbkOaVdVa",
    "rLGffQrfBKH0BgwQ33yZofvO3Or",
    "FO6HWqw",
    "GbgvyA2",
    "L1NU9QvIQIH7DTRt",
    "y7llk4Y8WfYflt6",
    "iuDp1WPbV3HRZudZtoXChxH4HNVBX5ZALe",
    "8C28RTXmVcco0",
    "X5Xh",
    "7xe25YUgfGgD0xW3ezFS",
    "",
    "CKCR",
    "8EmDjBo6h3eLaK7U6vU2Qys0NsMx",
    "t2TeZBXKqbdP09Arh9C3",
]


def captcha_sign(device_id: str, timestamp: str) -> str:
    """
    Generate a captcha sign.

    在网页端的js中, 搜索 captcha_sign, 可以找到对应的js代码

    """
    sign = CLIENT_ID + CLIENT_VERSION + PACKAG_ENAME + device_id + timestamp
    for salt in SALTS:
        sign = hashlib.md5((sign + salt).encode()).hexdigest()
    return f"1.{sign}"


def generate_device_sign(device_id, package_name):
    signature_base = f"{device_id}{package_name}1appkey"

    # 计算 SHA-1 哈希
    sha1_hash = hashlib.sha1()
    sha1_hash.update(signature_base.encode("utf-8"))
    sha1_result = sha1_hash.hexdigest()

    # 计算 MD5 哈希
    md5_hash = hashlib.md5()
    md5_hash.update(sha1_result.encode("utf-8"))
    md5_result = md5_hash.hexdigest()

    device_sign = f"div101.{device_id}{md5_result}"

    return device_sign


def build_custom_user_agent(device_id, user_id):
    device_sign = generate_device_sign(device_id, PACKAG_ENAME)

    user_agent_parts = [
        f"ANDROID-{APP_NAME}/{CLIENT_VERSION}",
        "protocolVersion/200",
        "accesstype/",
        f"clientid/{CLIENT_ID}",
        f"clientversion/{CLIENT_VERSION}",
        "action_type/",
        "networktype/WIFI",
        "sessionid/",
        f"deviceid/{device_id}",
        "providername/NONE",
        f"devicesign/{device_sign}",
        "refresh_token/",
        f"sdkversion/{SDK_VERSION}",
        f"datetime/{get_timestamp()}",
        f"usrno/{user_id}",
        f"appname/{APP_NAME}",
        "session_origin/",
        "grant_type/",
        "appid/",
        "clientip/",
        "devicename/Xiaomi_M2004j7ac",
        "osversion/13",
        "platformversion/10",
        "accessmode/",
        "devicemodel/M2004J7AC",
    ]

    return " ".join(user_agent_parts)
