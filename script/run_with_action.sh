#!/bin/bash

# 从GitHub Action获取输入参数
TARGET_PATH="$1"
CODE_TYPE="$2"
OUTPUT_REPORT="$3"
TARGET_CODE="$4"


# 移动递归 rule 到根目录 /rule
mv /app/rule /


# 确定使用的规则文件
case "$CODE_TYPE" in
    "trpc")
        RULE_FILE="/rule/find_trpc_api.yaml"
        ;;
    "gin")
        RULE_FILE="/rule/find_gin_api.yaml"
        ;;
    "go_swagger")
      RULE_FILE="/rule/find_go_swagger_api.yaml"
      ;;
    *)
        RULE_FILE="/rule/find_go_swagger_api.yaml"
        ;;
esac

# 创建输出目录（如果不存在）
mkdir -p "$(dirname "$OUTPUT_REPORT")"

# 处理 target_path
# 如果是 ./ 开头 就将这个路径添加去除 ./ 之后 添加到  /github/workspace。如果是 / 开头不处理，否则就添加 /github/workspace/
if [[ $TARGET_PATH == "./"* ]]; then
  TARGET_PATH="/github/workspace/${TARGET_PATH#./}"
  echo "处理 target_path: $TARGET_PATH"
elif [[ $TARGET_PATH != "/"* ]]; then
  echo "处理 target_path: $TARGET_PATH"
else
  echo "target_path 不处理: $TARGET_PATH"
fi

echo "当前目录："
ls -al

echo "/app 目录"
ls -al /app

# 制作FindCodeCommand 到 全局可调用的地方
ln -s /app/FindCodeCommand /usr/local/bin/FindCodeCommand

# 执行FindCodeCommand扫描
echo "执行FindCodeCommand扫描..."
echo "目标代码路径: $TARGET_PATH"
echo "代码类型: $CODE_TYPE"
echo "使用规则文件: $RULE_FILE"
echo "输出报告路径: $OUTPUT_REPORT"

FindCodeCommand -l "$TARGET_PATH" -r "$RULE_FILE" -o "$OUTPUT_REPORT" -go_target "$TARGET_CODE" -c "/app/etc/config.yaml"


# 检查执行状态
if [ $? -eq 0 ]; then
    cat "$OUTPUT_REPORT"
    ls /home/runner/work/augeu/augeu/report/ -al

    echo "FindCodeCommand扫描完成，报告已生成在: $OUTPUT_REPORT"
    exit 0
else
    echo "FindCodeCommand扫描失败"
    exit 1
fi

