# goWorkSpaceAPI
小工具  
通过命令行创建Google Admin WorkSpace 账号，并分配到相应的组。  

## 启动前准备  
1. 将对应的key放置到项目的上级目录中  
2. 将生成随机密码的脚本放置到项目的上级目录中  

## 启动方法
1. git clone https://github.com/SidneyCao/goWorkSpaceAPI.git
2. cd goWorkSpaceAPI
3. go build 
4. 编译完成后可以查看一下二进制的使用方法
```
    ./gcpWorkSpaceAPI -h
    Usage of ./gcpWorkSpaceAPI
    -a string
        管理员账号 (默认为空)
    -d string
        域名 (默认为空)
    -f string
        全名 (默认为空)
    -g string
        分组 (默认为空)
    -l string
        姓氏 (默认为空)
    -m string
        方法名
        list 列出域下的所有用户
        upload 创建用户
         (default "list")
    -o string
        组织名 (默认为空)
    -p string
        主邮箱 (默认为空)
```