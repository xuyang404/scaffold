#### 简介
一个通用的脚手架工具，从远程 GitHub 仓库克隆模板并生成新项目，同时允许用户输入自定义的替换内容。

#### 安装：
    go install github.com/xuyang404/scaffold


#### 命令行参数：
    -branch 要使用的分支（仅当模板是远程仓库时） (default "master")
    -local 本地模板路径
    -name 项目名称
    -remote 远程仓库url

#### eg:
    scaffold -remote https://golib.gaore.com/GaoreGo/hertz_demo.git -name ../hertz_new

#### 输入上面的命令后，可以在命令行交互中进行参数替换，下面例子会将模板中的{{PROJECT_NAME}}替换为my_project，{{GO_VERSION}}替换为1.20：
    输入替换值(key=value)，空行回车结束:    
    > {{PROJECT_NAME}}=my_project
    > {{GO_VERSION}}=1.20
    >
    