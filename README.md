# WEP
A Workbook Exchange Platform 一个作业本流转平台

## 由来
学校每学期都会分发大量作业本，而由于各位同学使用方式存在差异，出现了“有人没本子用，有本子没人用”的现象，其中又以后者居多，造成极大的浪费。

制作本站的目的正是为了平衡现有需求量以减少浪费。后期还可能与学校合作，借助本平台减少作业本的制作量、改善作业本分配。

## 部署
### 1 安装环境
到[Golang中文网](https://studygolang.com/dl) 下载最新的Go安装包并配置好环境变量

输入以下指令安装数据库驱动

    go get -u github.com/go-sql-driver/mysql
### 2 配置数据库
创建并选择数据库

    CREATE DATABASE q5xy DEFAULT CHARSET utf8;
    USE q5xy;
为什么是q5xy？因为本来这个平台叫泉五闲鱼，暂时懒得改后端名称。

创建新用户，password替换为自己设定的密码

    CREATE USER app IDENTIFIED BY "password";
    GRANT INSERT,SELECT,UPDATE,DELETE ON q5xy.* TO 'app'@'localhost' IDENTIFIED BY 'password';
    FLUSH PRIVILEGES;

创建用户表

    CREATE TABLE users ( id INT, name VARCHAR(20), class VARCHAR(40));
创建订单表

    CREATE TABLE orders ( id INT AUTO_INCREMENT, creator INT, matcher INT, item INT, amount INT, kind INT, date DATETIME, status INT, PRIMARY KEY (id));
创建日志表

    CREATE TABLE logs ( time DATETIME，action VARCHAR(20), user VARCHAR(20), ip VARCHAR(20), status VARCHAR(20));
### 3 下载并编译
克隆

    git clone https://github.com/Q5CS/WEP.git
将src目录下文件复制到你设置的gopath的src文件夹，然后在gopath中打开命令窗口，编译

    go build ./src/wep.go
将生成的文件与front目录放在一起，运行。client_secret要向电研社[五中人开放授权平台](https://doc.qz5z.ren/oauth2/)申请，mysql_password是前面设置的数据库密码。

    ./wep -s client_secret -p mysql_password
打开浏览器输入`localhost:9090`就可以看到效果了。
## TODO
以下是未来可能实现的东西
- [ ] 将dashboard与marketplace的数据改为由AJAX拉取
- [x] 通过命令行指定数据库地址与账密
- [x] 实现自己的SessionID而不是套用官网接口返回的ID
- [ ] 添加更多的注释供诸公参考
