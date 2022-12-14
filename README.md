## 部署 （默认是在linux环境下）

#### 修改配置

修改config.yaml

```yaml
db:
  host: mysql host
  port: mysql port
  user: user
  password: password
  charset: utf8mb4
  database: database


redis:
  host: localhost
  port: 6379
  password:
  db: 8

server:
  secret: qwe123 # secret  
  addr: :8081 # 对外暴露的域名和端口，如example.com
  domain: http://127.0.0.1:8081 # 当前对外暴露的域名

jwt:
  secret: 12313
  expire: 86400
```

#### 安装golang环境

* 版本要求: go1.18+

```shell
wget https://studygolang.com/dl/golang/go1.19.1.linux-amd64.tar.gz
```

解压下载的压缩包，将解压后的文件夹下的bin文件夹加入到环境变量即可

检查是否安装完成，直接在控制台输入 go 回车即可

#### 安装依赖并编译

在项目根目录下

```shell
export GOPROXY=https://goproxy.cn 
go mod tidy
go build -o app
```

编译完成之后得到一个名为app的可执行文件

在config.yaml同级目录下执行 `./app`

可配合`nohup`或`supervisor`来对进程进行管理