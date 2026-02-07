# mysql 库表生成代码
# powershell 以order.sql为例

# 生成库
# mysql -u {USER} -p {PWD} -e "CREATE DATABASE IF NOT EXISTS looklook_order DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;"

# 生成表
# Get-Content deploy/sql/looklook_order.sql -Encoding UTF8 | mysql -uroot -pwzy256411 looklook_order
# 若sql文件有中文注释,则需要添加-Encoding UTF8参数