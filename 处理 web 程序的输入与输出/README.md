
@[TOC]
## 一、概述
设计一个 web 小应用，展示静态文件服务、js 请求支持、模板输出、表单处理、Filter 中间件设计等方面的能力。（不需要数据库支持）

## 二、任务要求
编程 web 应用程序 cloudgo-io。 请在项目 README.MD 给出完成任务的证据！
### 基本要求
支持静态文件服务
支持简单 js 访问
提交表单，并输出一个表格
对 /unknown 给出开发中的提示，返回码 5xx


## 三、参考博客
详细的介绍和入门操作见潘老师的这篇博客，https://blog.csdn.net/pmlpml/article/details/78539261

## 四、实验过程
### 1. 配置所需环境
我的实验环境：win10，vscode，go语言
首先使用`git clone https://github.com/.../...`命令下载以下三个库（在终端中转到GOPATH下运行即可）：
```cpp
"github.com/codegangsta/negroni" 
"github.com/gorilla/mux"        
"github.com/unrolled/render"
```
所以，具体运行一下三个指令即可：
```cpp
go get -v github.com/codegangsta/negroni
go get -v github.com/gorilla/mux
go get -v github.com/unrolled/render
```

![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113171656392.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L0xMTF9TaWNpbHk=,size_16,color_FFFFFF,t_70)
### 2. 实验过程
 - 文件结构
assets放网页素材和js文件，service放处理的逻辑实现，templates中是表格模板
 - 代码实现
代码实现参考老师的教程，对给出的代码模块进行适当的修改即可，这里便不详细说明了。

### 3. 实验结果
在main.go文件同级目录下执行文件（注意**这里必须使终端跳转到该文件目录下，不然会报404错误**，真的是血一样的教训，一开始就是直接用的VSCode的run code，结果一直404到怀疑人生。。。。）：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113213145156.png)
#### 1）静态文件服务
在浏览器中输入http://localhost:8080/static/可以查看结果：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113215513485.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L0xMTF9TaWNpbHk=,size_16,color_FFFFFF,t_70)
当没有放入html文件时候显示的是当前文件夹下的所有文件。
直接在地址后输入文件名则访问文件内容：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113215456561.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L0xMTF9TaWNpbHk=,size_16,color_FFFFFF,t_70)
#### 2）简单 js 访问
在浏览器中输入http://localhost:8080/api/Test：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113220223223.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L0xMTF9TaWNpbHk=,size_16,color_FFFFFF,t_70)
#### 3）提交与输出表格
在浏览器中输入http://localhost:8080/table：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113221155754.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L0xMTF9TaWNpbHk=,size_16,color_FFFFFF,t_70)
点击submit：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113221230804.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L0xMTF9TaWNpbHk=,size_16,color_FFFFFF,t_70)
#### 4）对 /unknown 给出开发中的提示，返回码 5xx
在浏览器中输入http://localhost:8080/unknown，返回505错误。
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191113221510420.png)

