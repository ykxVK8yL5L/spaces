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
            echo "📁 备份目录: $path"
            
            # 先确保远程目录存在
            rclone mkdir "$REMOTE_FOLDER/$path" 2>/dev/null || true
            
            # 然后同步内容
            rclone sync --checksum --progress --create-empty-src-dirs \
                "$path" "$REMOTE_FOLDER/$path"
            echo "✅ 完成: $path"
        elif [ -f "$path" ]; then
            # 文件备份，保留父目录结构
            echo "📄 备份文件: $path"
            parent_dir=$(dirname "$path")
            
            # 确保远程父目录存在
            rclone mkdir "$REMOTE_FOLDER$parent_dir/" 2>/dev/null || true
            
            rclone copy --checksum --progress \
                "$path" "$REMOTE_FOLDER$parent_dir/"
            echo "✅ 完成: $path"
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
        if [ -d "$path" ] || [[ "$path" == */ ]]; then
            # 目录还原，使用 sync 确保完全匹配
            echo "📁 还原目录: $path"
            mkdir -p "$path"
            rclone sync --checksum --progress --create-empty-src-dirs \
                "$REMOTE_FOLDER/$path" "$path"
            echo "✅ 完成: $path"
        else
            # 文件还原
            echo "📄 还原文件: $path"
            target_dir=$(dirname "$path")
            mkdir -p "$target_dir"
            
            # 从远程父目录复制文件
            parent_dir=$(dirname "$path")
            filename=$(basename "$path")
            rclone copy --checksum --progress \
                "$REMOTE_FOLDER$parent_dir/$filename" "$target_dir/"
            echo "✅ 完成: $path"
        fi
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
