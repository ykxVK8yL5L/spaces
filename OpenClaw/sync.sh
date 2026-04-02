#!/bin/sh

# ----------------------------
# 配置部分
# ----------------------------

# 本地需要备份/还原的路径
# 目录末尾加 "/" 表示目录
OPENCLAW_PATHS="
/root/.openclaw/sessions/
/root/.openclaw/workspace/
/root/.openclaw/agents/main/sessions/
/root/.openclaw/openclaw.json
"

# ----------------------------
# 备份函数
# ----------------------------
backup() {
    echo "=== 开始备份 ==="

    for path in $OPENCLAW_PATHS; do
        if [ -d "$path" ]; then
            # 目录备份，保持相对路径，保留空目录
            rclone sync --checksum --progress --create-empty-src-dirs \
                "$path" "$REMOTE_FOLDER/$path"
        elif [ -f "$path" ]; then
            # 文件备份，保留父目录结构
            rclone copy --checksum --progress --parents \
                "$path" "$REMOTE_FOLDER"
        else
            echo "⚠️ 路径不存在: $path"
        fi
    done

    echo "=== 备份完成 ==="
}

# ----------------------------
# 还原函数
# ----------------------------
restore() {
    echo "=== 开始还原备份 ==="

    for path in $OPENCLAW_PATHS; do
        case "$path" in
            */)  # 目录
                mkdir -p "$path"
                rclone copy --checksum --progress --create-empty-src-dirs \
                    "$REMOTE_FOLDER/$path" "$path"
                ;;
            *)  # 文件
                target_dir=$(dirname "$path")
                mkdir -p "$target_dir"
                rclone copy --checksum --progress --parents \
                    "$REMOTE_FOLDER/$path" "$target_dir"
                ;;
        esac
    done

    echo "=== 还原完成 ==="
}

# ----------------------------
# 主入口
# ----------------------------
case "$1" in
    backup)
        backup
        ;;
    restore)
        restore
        ;;
    *)
        echo "Usage: $0 {backup|restore}"
        exit 1
        ;;
esac
