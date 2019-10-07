package main

//C语言：读取指令 → 分析指令 → 设置参数 → 判断参数合法 → 执行指令
//读取指令 → 分析指令并设置参数 → 判断参数合法 → 执行指令

import (
	"bufio"	//bufio io os用作读取，输入输出
	"fmt"
	flag "github.com/spf13/pflag"  //安装的pflag包
	"io"
	"os"
	"os/exec"
	"strings"
)

//构造结构，设置变量
type args struct {
	start int  //起始位置
	end int  //终止位置
	len int  //行数
	cal_type string  //是否按页结束符计算(l’代表按照行数计算页；‘f’为按照换页符计算页) ???
	dest string  //定向位置
	in_type string  //输入方式（文件输入还是键盘输入）
}

func main() {
	sa := new(selpg_Args)

	process(*sa)

	run(*sa)  //执行指令
}

//根据参数执行指令，步骤：判断输入方式，绑定输入流 → 绑定管道（如果有管道） → l/f读取
func run(sa selpg_Args) {
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

func process(sa selpg_Args) {
	//参数绑定：将所有的参数使用pflag绑定到变量上。  //之后改为单独的函数
	flag.IntVar(&sa.start, "s", 0, "the start page")  //起始位置，默认为0
	flag.IntVar(&sa.end, "e", 0, "the end page")  //终止位置，默认也是0
	flag.IntVar(&sa.len, "l", 72, "the length of the page")  //每页行数，默认为72
	flag.IntVar(&sa.dest, "d", "", "the destination of the out")  //输出位置，默认为空字符？？

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


func usage() {
	fmt.Fprintf(os.Stderr, "\nUSAGE: ./selpg [--s start] [--e end] [--l lines | --f ] [ --d dest ] [ in_filename ]\n")
	fmt.Fprintf(os.Stderr, "\n selpg --s start    : start page")
	fmt.Fprintf(os.Stderr, "\n selpg --e end      : end page")
	fmt.Fprintf(os.Stderr, "\n selpg --l lines    : lines/page")
	fmt.Fprintf(os.Stderr, "\n selpg --f          : check page with '\\f'")
	fmt.Fprintf(os.Stderr, "\n selpg --d dest     : pipe destination\n")
}
