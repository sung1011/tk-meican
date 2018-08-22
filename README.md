# tk-meican
美餐相关

### INSTALL
1. git clone ...
2. 安装依赖 ```dep ensure -v```

### Usage
##### 编译
- go build -o $HOME/go/bin/tk-meican $GOPATH/src/sung/tk-meican/main.go # 编译到环境变量所在的目录

##### regi 注册用户
tk-meican regi -j {工号} -n {名字}

##### join 参加当天随机取餐
tk-meican join -j {工号} -m {今天美餐号}

##### roll 随机取餐
tk-meican roll

### TODO
统计
修改注册信息
钉钉通知
