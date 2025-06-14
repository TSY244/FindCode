#!/bin/bash


url="${GIT_URL:-$1}"


# 从URL提取仓库名称
repo_name=$(basename $url .git)

# 克隆仓库
echo "正在克隆仓库: $url"
git clone $url

if [ $? -ne 0 ]; then
    echo "错误：克隆仓库失败"
    exit 1
fi

# 记录仓库路径
path=$repo_name
echo "仓库路径: $path"

# 运行FindCode程序
echo "正在运行FindCode..."
./FindCode -l $path -r rule/find_go_swagger_api.yaml

if [ $? -ne 0 ]; then
    echo "警告：FindCode运行可能出错"
fi

# 启动HTTP服务器
echo "在result目录下启动HTTP服务器..."
cd result
python3 -m http.server 8000 --bind 0.0.0.0