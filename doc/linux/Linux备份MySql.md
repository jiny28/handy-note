```shell
编写脚本
db_user="root"
db_passwd="123456"
db_name="hlhz_duoji"
backup_dir="/myfile/backfile"  
time="$(date +"%Y_%m_%d_%H:%M:%S")"
mysqldump -h192.168.31.66  -P33066 -u$db_user  -p$db_passwd $db_name > $(backup_dir)/$(db_name)_$(time).sql

文件执行权限
 chmod +x  （filePath）
 编写定时
修改/etc/crontab

vim crontab

45 22 * * * root /home/mysql_data/mysql_databak.sh #表示每天22点45分执行备份
* * * * * root /home/mysql_data/mysql_databak.sh #表示每分执行一次

```