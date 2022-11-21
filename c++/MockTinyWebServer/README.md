# C++ 实现的 web 服务器

参考：https://github.com/qinguoyi/TinyWebServer

开发的平台是在 Ubuntu20.04 上面

# MySQL 安装
安装并启动 mysql
```bash
sudo apt-get update
sudo apt-get install mysql-server
sudo systemctl start mysql.service
```
配置 mysql 信息
```bash
sudo mysql_secure_installation
```
具体参考：https://www.digitalocean.com/community/tutorials/how-to-install-mysql-on-ubuntu-20-04

使用 `mysql_secure_installation` 方式配置 root 密码出错的话，需要使用 `sudo mysql` 进入 mysql 使用 `ALTER user` 进行密码修改，具体参考上面链接。

还需要安装 mysql 的开发头文件：
```bash
sudo apt-get install libmysqlclient-dev
```


