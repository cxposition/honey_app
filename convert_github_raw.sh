#!/bin/bash

# ==========================================================
# 脚本名称: replace_self_in_readme.sh
# 描述: 转换用户输入的 URL 为 Raw URL，并在 README.md 中将
#       输入的 URL 替换为转换后的 Raw URL。
# ==========================================================

# --- 1. 定义变量 ---
README_FILE="README.md"  # 要操作的文件名 (请确保脚本和此文件在同一目录)

# --- 2. 函数：URL 转换 ---
convert_to_raw() {
    local input_url="$1"

    # 替换域名和路径段
    # 注意: 这个转换函数会失败，如果输入的 URL 本身就不是 /blob/ 格式。
    local converted_url=$(echo "$input_url" | \
        sed 's/github.com/raw.githubusercontent.com/g' | \
        sed 's/\/blob\//\//g')

    echo "$converted_url"
}

# --- 3. 脚本主执行部分 ---

echo "--- README.md 链接自我替换工具 ---"

# 检查 README.md 文件是否存在
if [ ! -f "$README_FILE" ]; then
    echo "错误: 文件 $README_FILE 不存在于当前目录。请确保您在正确的目录下运行脚本。"
    exit 1
fi

# 提示用户输入链接 (该链接既是旧链接，也是新链接的源头)
read -p "请输入需要被替换的 GitHub 文件 URL (例如: https://github.com/.../blob/...): " OLD_GITHUB_URL

# 检查输入是否为空
if [ -z "$OLD_GITHUB_URL" ]; then
    echo "错误: URL 输入不能为空。"
    exit 1
fi

echo "-----------------------------------"

# 转换 URL
NEW_RAW_URL=$(convert_to_raw "$OLD_GITHUB_URL")

# OLD_GITHUB_URL 现在是我们需要查找并替换的【旧文本】
# NEW_RAW_URL 是我们转换后的【新文本】

echo "✅ 旧链接 (查找目标):  $OLD_GITHUB_URL"
echo "✅ 新链接 (替换目标): $NEW_RAW_URL"
echo "🔍 正在 ${README_FILE} 中执行替换..."

# --- 4. 执行替换操作 ---
# 替换前，先备份原始文件 (强烈推荐)
cp "$README_FILE" "${README_FILE}.bak"

# 由于 URL 中包含大量特殊字符 (如 / . -)，我们必须使用一个安全的 sed 分隔符，例如 |
# sed 命令格式: 's|旧文本|新文本|g'

# 注意: 使用 #!/bin/bash 的 Shell，通常使用 GNU sed。 macOS/BSD sed 略有不同。
# 增加转义以确保特殊字符被正确处理
ESCAPED_OLD_URL=$(printf "%s\n" "$OLD_GITHUB_URL" | sed 's/[][\/.^$*]/\\&/g')
ESCAPED_NEW_URL=$(printf "%s\n" "$NEW_RAW_URL" | sed 's/[][\/.^$*]/\\&/g')


# 使用 sed 替换。这里使用 | 作为分隔符，并使用转义后的 URL
if sed -i.bak "s|$ESCAPED_OLD_URL|$ESCAPED_NEW_URL|g" "$README_FILE"; then

    # 检查是否有内容被替换
    if grep -q "$NEW_RAW_URL" "$README_FILE"; then
        echo "✅ 替换成功! ${README_FILE} 中的旧链接已被更新。"
        rm -f "${README_FILE}.bak" # 替换成功后删除备份
    else
        echo "❌ 警告: sed 命令执行成功，但未在新文件中找到替换后的文本。"
        echo "请检查 ${README_FILE} 中是否包含旧链接 $OLD_GITHUB_URL。"
        echo "原始文件已备份为 ${README_FILE}.bak"
    fi

else
    echo "❌ 替换失败。请检查权限或 sed 语法。"
    echo "原始文件已备份为 ${README_FILE}.bak"
fi

echo "-----------------------------------"