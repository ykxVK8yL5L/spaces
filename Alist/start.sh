#!/bin/bash


if [ -n "$RCLONE_CONF" ]; then
  echo -e "##########同步备份############"
  # 指定远程文件夹路径，格式为 remote:path
  REMOTE_FOLDER="huggingface:/alist"

  # 使用 rclone ls 命令列出文件夹内容，将输出和错误分别捕获
  OUTPUT=$(rclone ls "$REMOTE_FOLDER" 2>&1)

  # 获取 rclone 命令的退出状态码
  EXIT_CODE=$?

  # 判断退出状态码
  if [ $EXIT_CODE -eq 0 ]; then
    # rclone 命令成功执行，检查文件夹是否为空
    if [ -z "$OUTPUT" ]; then
      #为空不处理
      #rclone sync --interactive /ql $REMOTE_FOLDER
      echo "初次安装,不备份"
    else
      #echo "文件夹不为空"
      mkdir /home/user/alist/data
      rclone sync $REMOTE_FOLDER /home/user/alist/data 
    fi
  elif [[ "$OUTPUT" == *"directory not found"* ]]; then
    echo "错误：文件夹不存在"
  else
    echo "错误：$OUTPUT"
  fi
else
    echo "没有检测到Rclone配置信息"
fi


cd $HOME/alist && alist server &

yacron &

# Wait for any process to exit
wait -n
  
# Exit with status of process that exited first
exit $?
