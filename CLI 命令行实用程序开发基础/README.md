@[TOC](目录)

# 开发实践要求
使用 golang 开发 开发 Linux 命令行实用程序 中的 selpg
提示：

 - 请按文档 使用 selpg 章节要求测试你的程序 
 - 请使用 pflag 替代 goflag 以满足 Unix 命令行规范，参考：Golang之使用Flag和Pflag 
 - golang 文件读写、读环境变量，请自己查 os 包 
 - “-dXXX” 实现，请自己查os/exec 库，例如案例 Command，管理子进程的标准输入和输出通常使用 io.Pipe，具体案例见 Pipe

# 设计说明
## 安装包
使用了pflag包（flag的升级版本）。安装地址为github.com/spf13/pflag
```cpp
go get github.com/spf13/pflag

```

## 代码实现
### 整体思路
参考C++程序：[selpg.c源码](https://www.ibm.com/developerworks/cn/linux/shell/clutil/selpg.c)
具体步骤为：读取指令 → 分析指令并设置参数 → 判断参数合法 → 执行指令

### 结构
```go
//构造结构，设置变量
type args struct {
	start int  //起始位置
	end int  //终止位置
	len int  //行数
	cal_type string  //是否按页结束符计算(l’代表按照行数计算页；‘f’为按照换页符计算页) ???
	dest string  //定向位置
	in_type string  //输入方式（文件输入还是键盘输入）
}
```

### 参数绑定与合法性检查
```go
func process(sa args) {
	//参数绑定：将所有的参数使用pflag绑定到变量上。  //之后改为单独的函数
	flag.IntVar(&sa.start, "s", 0, "the start page")  //起始位置，默认为0
	flag.IntVar(&sa.end, "e", 0, "the end page")  //终止位置，默认也是0
	flag.IntVar(&sa.len, "l", 72, "the length of the page")  //每页行数，默认为72
	flag.StringVar(&sa.dest, "d", "", "the destination of the out")  //输出位置，默认为空字符？？

	//查找是否为f
	forl := flag.Bool("f", false, "")  
	flag.Parse()

	//f和l的两种情况
	if *forl {
		sa.cal_type = "f"
		sa.len = -1
	} else {
		sa.cal_type = "l"
	}

	//输入方式（文件输入还是键盘输入）
	//若使用了文件输入，将方式为文件名？？
	sa.in_type = ""
	if flag.NArg() == 1 {
		sa.in_type = flag.Arg(0)
	}

	//判断参数合法 → 主要是剩余参数个数、l/f是否同时出现、起始页与终止页是否冲突
	//参数个数
	n := flag.NArg()
	if n!= 1 && n != 0 {
		usage()
		os.Exit(1)
	}

	//起始终止位置
	if sa.start > sa.end || sa.start < 1 {
		usage()
		os.Exit(1)
	}

	//l、f不能同时出现
	if sa.cal_type == "f" && sa.len != -1 {
		usage()
		os.Exit(1)
	}
}
```

## 根据参数执行指令

```go
//根据参数执行指令，步骤：判断输入方式，绑定输入流 → 绑定管道（如果有管道） → l/f读取
func run(sa args) {
	//初始化
	fin := os.Stdin  //输入
	fout := os.Stdout  //输出
	nowline := 0  //当前行
	nowpage := 1  //当前页
	var inpipe io.WriteCloser  //管道
	var err error  //错误

	//判断输入方式，绑定输入流
	if sa.in_type != "" {  //如果不是？？？
		fin, err = os.Open(sa.in_type)
		if err != nil {
			fmt.Println(err)
			usage()
			os.Exit(1)
		}
		defer fin.Close()  //全部结束后关闭
	}

	//用管道接通grep模拟打印机测试，结果输出到屏幕
	if sa.dest != "" {
		//Command(name string, arg ...string) *Cmd
		cmd := exec.Command("grep", "-nf", "keyword")
		inpipe, err = cmd.StdinPipe()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer inpipe.Close()  //最后执行
		cmd.Stdout = fout
		cmd.Start()
	}

	//分页方式
	//
	if sa.cal_type == "l" {  //若按行读取
		line := bufio.NewScanner(fin)
		for line.Scan() {
			if nowpage >= sa.start && nowpage <= sa.end {
				//输出到窗口
				fout.Write([]byte(line.Text() + "\n"))
				if sa.dest != "" {
					//定向输出到文件管道
					inpipe.Write([]byte(line.Text() + "\n"))
				}
			}
			nowline ++
			//翻页
			if nowline >= sa.len {
				nowline = 0
				nowpage ++
			}
		}
	} else {  //若按换页符
		//rite(p []byte) (n int, err error)
		rd := bufio.NewReader(fin)
		for {
			page, ferr := rd.ReadString('\f')
			if ferr != nil || ferr == io.EOF {
				if ferr == io.EOF {
					if nowpage >= sa.start && nowpage <= sa.end {
						fmt.Fprintf(fout, "%s", page)
					}
				}
				break
			}
			//'\f'翻页
			page = strings.Replace(page, "\f", "", -1)
			nowpage++
			if nowpage >= sa.start && nowpage <= sa.end {
				fmt.Fprintf(fout, "%s", page)
			}
		}

	}
	if nowpage < sa.end {  //比较输出的页数与期望输出的数量
		fmt.Fprintf(os.Stderr, "./selpg: end (%d) greater than total pages (%d), less output than expected\n", sa.end, nowpage)
	}

}
```

## 输入错误提示信息
```go
func usage() {
	fmt.Fprintf(os.Stderr, "\nUSAGE: ./selpg [--s start] [--e end] [--l lines | --f ] [ --d dest ] [ in_filename ]\n")
	fmt.Fprintf(os.Stderr, "\n selpg --s start    : start page")
	fmt.Fprintf(os.Stderr, "\n selpg --e end      : end page")
	fmt.Fprintf(os.Stderr, "\n selpg --l lines    : lines/page")
	fmt.Fprintf(os.Stderr, "\n selpg --f          : check page with '\\f'")
	fmt.Fprintf(os.Stderr, "\n selpg --d dest     : pipe destination\n")
}

```

## 主函数

```go
func main() {
	sa := new(selpg_Args)

	process(*sa)

	run(*sa)  //执行指令
}
```


# 测试结果
命令1
```cpp
selpg
```
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191008111433718.png)
命令2
```cpp
selpg -s 1 -e 1 cs.txt
```
该指令会将第一页的内容输出到屏幕上
测试文件：我的代码做成的txt文件
结果：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191008113038354.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L0xMTF9TaWNpbHk=,size_16,color_FFFFFF,t_70)
命令3

```cpp
more input.txt | ./Selpg -s1 -e2
```
该指令将more指令的输出，可以看到会出现中文乱码
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191008113422280.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L0xMTF9TaWNpbHk=,size_16,color_FFFFFF,t_70)
命令4

```cpp
Selpg -s1 -e2 cs.txt > out.txt
```
将cs.txt的1-2页输入到out.txt中
结果：
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191008113558767.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L0xMTF9TaWNpbHk=,size_16,color_FFFFFF,t_70)
命令5

```cpp
Selpg -s10 -e20 input.txt 2>error.txt
```
没有报错信息

命令6

```cpp
Selpg -s1 -e2 -l40 cs.txt
```

通过-l指定每一页的行数，能看到最后一行是代码的第79行
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191008114146583.png)
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191008114203388.png)
