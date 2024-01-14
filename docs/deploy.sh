#!/usr/bin/env bash

# 确保脚本抛出遇到的错误
set -e

npm run build

pub_dir=".vitepress/dist/"
oss_bucket="luchen-fun"
ossutil="ossutil -c ~/opt/conf/aliyun/oss-gz.conf "

# $ossutil rm -rf oss://${oss_bucket}/ --exclude "*.pdf"
$ossutil cp -rf ${pub_dir} oss://${oss_bucket}/ --exclude ".DS_Store" --exclude "*.drawio" -u

echo "done"
