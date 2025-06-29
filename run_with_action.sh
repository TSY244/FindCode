#!/bin/bash

# 从GitHub Action获取输入参数
TARGET_PATH="$1"
CODE_TYPE="$2"
OUTPUT_REPORT="$3"
TARGET_CODE="$4"

# 移动递归rule到根目录/rule（使用sudo）
sudo mv /app/rule /

# 确定使用的规则文件
case "$CODE_TYPE" in
    "trpc") RULE_FILE="/rule/find_trpc_api.yaml" ;;
    "gin") RULE_FILE="/rule/find_gin_api.yaml" ;;
    "go_swagger") RULE_FILE="/rule/find_go_swagger_api.yaml" ;;
    *) RULE_FILE="/rule/find_go_swagger_api.yaml" ;;
esac

# 创建输出目录（确保有权限）
REPORT_DIR=$(dirname "$OUTPUT_REPORT")
mkdir -p "$REPORT_DIR"
sudo chown runner:runner "$REPORT_DIR"

# 创建符号链接（使用绝对路径）
sudo ln -s /app/FindCode /usr/local/bin/FindCode

# 执行FindCode扫描
echo "执行FindCode扫描..."
echo "目标路径: $TARGET_PATH"
echo "代码类型: $CODE_TYPE"
echo "规则文件: $RULE_FILE"
echo "输出报告: $OUTPUT_REPORT"

FindCode -l "$TARGET_PATH" -r "$RULE_FILE" -o "$OUTPUT_REPORT" -go_target "$TARGET_CODE" -c "/app/etc/config.yaml"

# 检查执行状态并输出日志
if [ $? -eq 0 ]; then
    echo "=== 报告内容 ==="
    cat "$OUTPUT_REPORT"
    echo "=== 报告目录 ==="
    ls -al "$REPORT_DIR"
    echo "FindCode扫描完成，报告已生成在: $OUTPUT_REPORT"
    exit 0
else
    echo "FindCode扫描失败"
    exit 1
fi