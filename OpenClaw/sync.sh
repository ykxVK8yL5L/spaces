#!/bin/sh
# 文件名: sync.sh
# 用法: ./sync.sh backup 或 ./sync.sh restore
# 兼容 /bin/sh，OPENCLAW_PATHS 可换行声明

# ----------------------------
# 配置部分
# ----------------------------

# 远程备份目录（Dockerfile 或环境变量中可声明）
# 示例: REMOTE_FOLDER="huggingface:/openclaw"
: "${REMOTE_FOLDER:?请在环境变量中设置 REMOTE_FOLDER}"

# 定义要备份的路径，用空格分隔，换行使用反斜杠
OPENCLAW_PATHS="\
/root/.openclaw/sessions/ \
/root/.openclaw/workspace/ \
/root/.openclaw/agents/main/sessions/ \
/root/.openclaw/openclaw.json"

BACKUP_FILENAME="openclaw_backup.tar.gz"
TMP_BACKUP="/tmp/$BACKUP_FILENAME"

# ----------------------------
# 方法：备份
# ----------------------------
backup() {
    # echo "=== 创建备份文件 ==="
    # tar -czf "$TMP_BACKUP" $OPENCLAW_PATHS || {
    #     echo "失败: 创建备份失败!"
    #     return 1
    # }

    # echo "=== 上传备份文件到远程 ==="
    # rclone copy "$TMP_BACKUP" "$REMOTE_FOLDER/" --progress || {
    #     echo "失败: 上传备份失败!"
    #     return 1
    # }
    echo "=== 开始备份 ==="
    # 遍历每个路径
    for path in $OPENCLAW_PATHS; do
        name=$(basename "$path")
        # 判断是文件还是目录
        if [ -d "$path" ]; then
            # 文件夹 → 使用 sync 保证完全一致
            rclone sync "$path" "$REMOTE_FOLDER/$name" --checksum --progress --create-empty-src-dirs
        elif [ -f "$path" ]; then
            # 文件 → 用 copy 覆盖
            rclone copy "$path" "$REMOTE_FOLDER/" --checksum --progress
        else
            echo "⚠️ 路径不存在: $path"
        fi
    done

    echo "✅ 备份上传成功!"
}

# ----------------------------
# 方法：还原
# ----------------------------
restore() {
    # echo "=== 检查远程备份文件是否存在 ==="
    # if rclone ls "$REMOTE_FOLDER/$BACKUP_FILENAME" >/dev/null 2>&1; then
    #     echo "=== 下载远程备份到本地 ==="
    #     rclone copy "$REMOTE_FOLDER/$BACKUP_FILENAME" "/tmp/" --progress
    #     echo "=== 解压备份文件 ==="
    #     tar -xzf "$TMP_BACKUP" -C / || {
    #         echo "失败: 解压备份文件失败!"
    #         return 1
    #     }
    #     echo "✅ 还原完成!"
    # else
    #     echo "⚠️ 没有发现远程备份文件!"
    #     return 1
    # fi
    echo "=== 开始还原备份 ==="
    for path in $OPENCLAW_PATHS; do
        name=$(basename "$path")
        echo "🔄 还原 $name ..."
        if [ "${path%/}" != "$path" ]; then
            # 👉 文件：复制到父目录
            target_dir=$(dirname "$path")
            mkdir -p "$target_dir"
            rclone copy "$REMOTE_FOLDER/$name" "$target_dir" \
                --checksum --progress
        else
            # 👉 目录：正常处理
            mkdir -p "$path"
            rclone copy "$REMOTE_FOLDER/$name" "$path" \
                --checksum --progress --create-empty-src-dirs
        fi
    done
}

# ----------------------------
# 命令行参数处理
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
