参考链接：

https://blog.csdn.net/qq_35015497/category_9483904.html



# Go语言的一些知识总结：



指针：一般不应在函数中传入指针来修改值，Unmarshal 这类除外。但在结构体较大时，向函数传入指针的性能会比较好，传递指针大约1纳秒，而10M大小的数据需要耗时1毫秒。返回值则不同，1M 以下的数据结构比指针类型要快，如100字节数据花费10纳秒，而这一数据结构的指针耗时在10纳秒（i7-8700 32GB内存测试数据）。

## 小技巧

1、交叉编译

```bash
# macOS下编译Linux及Windows 64位可执行程序(-o用于指定路径及文件名，可省略)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o xxx main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o xxx main.go
 
# Linux下编译Mac及Windows 64位可执行程序
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
 
#Windows下编译Mac及Linux 64位可执行程序
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build main.go
 
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build main.go
```

2、查询出所有的库函数、结构体等，如 net/http

```bash
go doc net/http | grep "^func"
go doc net/http | grep "^type"|grep struct
```



3、使用Docker容器编译

```bash
docker pull golang:1.15.12-alpine3.13
docker run --rm -it \
-v /path/xxx:/app \
-w /app/src \          # 假设main.go及go.mod放在/path/xxx目录下
-e CGO_ENABLED=0 \     # 如在 alpine 中执行可不指定
-e GOPROXY=https://goproxy.cn \
golang:1.15.12-alpine3.13 \
go build -o ../path/xxx main.go
```



出现/lib/ld-musl-x86_64.so.1: bad ELF interpreter: No such file or directory报错可在编译时使用-e CGO_ENABLED=0，或依然在 alpine 内执行

4、Golang操作 Docker API

API文档地址：https://docs.docker.com/engine/api/v1.41/

```bash
# 配置文件/usr/lib/systemd/system/docker.service中的ExecStart最后面添加-H tcp://0.0.0.0:2345
ExecStart=xxxxx -H tcp://0.0.0.0:2345
sudo systemctl daemon-reload
sudo systemctl restart docker
# 验证命令
docker -H tcp://ip.address.xxx:2345 ps
 
go get github.com/docker/docker/client
 
# 示例代码 NewClient第2个参数为API 版本，可通过 docker version 进行查看
cli, err := client.NewClient("tcp://the.ip.addr:2345", "1.41", nil, nil)
if err != nil {
   log.Fatal(err)
}
images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
if err != nil {
   log.Fatal(err)
}
for _, image := range images{
   fmt.Println(image.RepoTags)
}
```



5、相关资源

- [Go社区推荐的目录结构](https://github.com/golang-standards/project-layout)(官方未进行认可)
- [Uber Go语言编码规范](https://github.com/uber-go/guide)
- [Style guideline for Go packages](https://rakyll.org/style-packages/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](http://github.com/golang/go/wiki/CodeReviewComments)

6、单元测试覆盖率

```bash
go test -race -cover  -coverprofile=./coverage.out -timeout=10m -short -v ./...
go tool cover -func ./coverage.out
```



7、利用工具检测变量遮蔽问题(不能定位所有问题)

```bash
$go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest
$go vet -vettool=$(which shadow) -strict xxx.go
```



![img](https://cdn.nlark.com/yuque/0/2022/png/1207203/1649371068579-2e9a3600-113e-45d7-ab58-76b720f02292.png)

Go 一些工具和功能

![img](https://cdn.nlark.com/yuque/0/2022/png/1207203/1649371066716-0397e547-3dca-4ff4-b503-5a7211641c99.png)

## Go语言的安装与开发环境

下载：

国内：http://studygolang.com/dl

https://golang.org/dl/

```bash
# 设置国内镜像
go env -w GOPROXY=https://goproxy.cn,direct
# 开启 Go Module
go env -w GO111MODULE=on
# goimports
go get -v golang.org/x/tools/cmd/goimports
```



开发环境：vi, emacs, idea, eclipse, vs, sublime … + go 插件

IDE：Goland, liteIDE

本课程使用 idea + go 插件

多版本：https://golang.org/doc/manage-install

本文GitHub仓库：https://github.com/alanhou/learning-go

VS Code：使用快捷键：command+shift+P，然后键入：go:install/update tools，将所有 16 个插件都勾选上，然后点击 OK 即开始安装

学习资料：

- [Effective go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go by Example](https://gobyexample.com/)
- [GoCN 社区](https://gocn.vip/wiki)

## 01 基础语法

### 变量

#### 变量定义

- - var a, b, c bool
  - var s1, s2 string = “hello”, “world”
  - 可放在函数内，或直接放在包内
  - 使用 var()集中定义变量
  - 让编译器自动决定类型
  - var a, b, i, s1, s2 = true, false, 3, “hello”, “world”

- 使用 **:=** 定义变量

- - a, b, i, s1, s2 := true, false, 3, “hello”, “world”
  - 只能在函数内使用

#### 内建变量类型

- bool, string
- 整数类型：(u)int, (u)int8, (u)int16, (u)int32, (u)int64, **uintptr（指针）**

- - (u)int ：int类型，加上u表示无符号int类型，不规定长度则int长度根据操作系统决定，32位系统中为32位，64位系统中为64位。
  - 不加符号的可以指定长度，未指定长度时根据操作系统是多少位来决定

- byte, **rune**（长度32位，相当于 char，解决多国语言问题）

- - **rune补充**

- float32, float64, complex64, complex128

- - complex64：复数，实部和虚部都为32位
  - complex128：复数，实部和虚部都为64位

##### 复数回顾

![img](https://cdn.nlark.com/yuque/0/2022/png/1207203/1649379677623-192fda5a-368a-40a8-81ef-f0f4bbc2f1fa.png)

- i = −1−−−√−1
- 复数：3 + 4i （实部+虚部）
- 斜线的长度：|3+4i|=32+42−−−−−−√=5|3+4i|=32+42=5
- i2=−1,i3=−i,i4=1,…i2=−1,i3=−i,i4=1,…



![img](https://cdn.nlark.com/yuque/0/2022/png/1207203/1649379837870-25af0586-4b5e-4377-8e8d-d76f8765c635.png)

- eiϕ=cosϕ+isinϕeiϕ=cosϕ+isinϕ
- |eiϕ|=cos2ϕ+sin2ϕ−−−−−−−−−−−√=1|eiϕ|=cos2ϕ+sin2ϕ=1
- e0=1,eiπ2=ie0=1,eiπ2=i
- eiπ=−1,ei32π=−i,ei2π=1eiπ=−1,ei32π=−i,ei2π=1

##### 最美公式 – 欧拉公式

![img](https://cdn.nlark.com/yuque/0/2022/png/1207203/1649379856904-f86eb847-a14d-4782-a90d-c53d27c70297.png)

#### 强制类型转换

![img](https://cdn.nlark.com/yuque/0/2022/png/1207203/1649380808495-ae2f1e55-3c3b-48a3-b037-22c80b14141d.png)

- 类型转换是强制的
- var a, b int = 3, 4
- var c int = math.Sqrt(a*a + b*b)   ✕
- var c int = int(math.Sqrt(float64(a*a + b*b)))  ?(尚有浮点数所带来的偏差问题需解决)

### 常量与枚举

#### 常量的定义

- const filename = “abc.txt”
- const 数值可作为各种类型使用
- const a,b = 3,4
- var c = int(math.Sqrt(a * a + b * b)) // a,b 未指定类型无需转换为 float

const关键字：表示常量，常量可定义在包内部，放法外面，放法内部可直接使用，**可直接指定常量的类型，也可不指定类型**，常量的数值可以当作任何类型使用，当使用这个常量时会自动转换，常量定义也可以使用括号括起来。

#### 使用常量定义枚举类型

- 普通枚举类型
  - go没有指定的枚举关键字，可以使用const使用枚举
- 自增值枚举类型
  - go语言提供一种简单写法：关键字---》iota 表示一个变量是自增值

```go
const(
	cpp = iota // 用于自增，无需再为下面的项赋值
	java
	python
	golang
)
// iota高级用法：可以参与运算
const (
  b = 1 << (10 * iota)
  kb
  mb
  gb
  tb
  pb
)
```

注：iota 字母中第9个字母：Ι, ι，英文释义为极小值

#### 变量定义要点回顾

- 变量类型写在变量名之后
- 编译器可推测变量类型
- 没有 char，只有 rune
- 原生支持复数类型

### 条件语句

#### if

```go
func bounded(v int) int {
	if v > 100 {
		return 100
	} else if v < 0 {
		return 0
	} else {
		return v
	}
}
```

- if 的条件里不需要括号

  ```go
  if contents, err := ioutil.ReadFile(filename); err != nil {
  	fmt.Println(err)
  }else{
  	fmt.Printf("%s\n", contents)
  }
  ```

- if 的条件里可以赋值

- if 的条件里赋值的变量作用域就在这个 if 语句里

#### switch

```go
func eval(a, b int, op string) int {
	var result int
	switch op {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		result = a / b
	default:
		panic("unsupported operator:" + op)
	}
	return result
}
```

- **switch 会自动 break，除非使用fallthrough**
- switch 后可以没有表达式

### 循环

#### for

```go
sum := 0
for i := 1; i <= 100; i++ {
	sum += i
}
```

- for 的条件里不需要括号

- for 的条件里**可以省略初始条件，结束条件，递增表达式**

  ```go
  // 整数转二进制
  func convertToBin(n int) string {
  	result := ""
  	for ; n > 0; n /= 2 {
  		lsb := n % 2
  		result = strconv.Itoa(lsb) + result
  	}
  	return result
  }
  ```

- 省略初始条件，相当于 while

- **省略所有条件，无限循环/死循环**

  ```go
  for {
  	fmt.Println("abc")
  }
  ```

#### 基本语法要点回顾

- for，if后面的条件没有括号
- if 条件里也可定义变量
- 没有 while
- switch 不需要 break，也可直接 switch 多个条件

### 函数

- **func** eval(a, b int, op string) int

- 函数可返回多个值

  ```go
  func div(a, b int) (int, int) {
  	return a / b, a % b
  }
  ```

- **函数返回多个值时可以起名字，但仅用于非常简单的函数，对于调用者而言没有区别**

  ```go
  func div(a, b int) (q, r int) {
  	q = a / b
  	r = a % b
  	return
  }
  ```

- **函数可以作为参数，函数式编程**

  ```go
  func apply(op func(int, int) int, a, b int) int{
  	fmt.Printf("Calling %s with %d, %d\n",
  		runtime.FuncForPC(reflect.ValueOf(op).Pointer()).Name(),
  		a, b)
  	return op(a, b)
  }
  ```

- 可变参数列表

  ```go
  func sumArgs(values ...int) int {
  	sum := 0
  	for i := range values{
  		sum += values[i]
  	}
  	return sum
  }
  ```

  

#### 函数语法要点回顾

- 返回值类型写在最后面
- 可返回多个值
- 函数可作为参数
- 没有默认参数、可选参数

### 指针

- 指针不能运算

#### 参数传递

值传递？引用传递？

![image-20220409080330770](https://cdn.jsdelivr.net/gh/Fly0905/note-picture@main/img/202204091045743.png)

![image-20220409080401059](https://cdn.jsdelivr.net/gh/Fly0905/note-picture@main/img/202204091045563.png)

![image-20220409080538144](https://cdn.jsdelivr.net/gh/Fly0905/note-picture@main/img/202204091045512.png)

- Go 语言只有值传递一种方式

  ```go
  // 通过指针来交换值
  func swap(a, b *int){
  	*b, *a = *a, *b
  }
  ```

  

## 02内建容器

### 数组、切片和容器

#### 数组

```go
var arr1 [5]int // 声明数组
arr2 := [3]int{1, 3, 5}  // 声明数组并赋值
arr3 := [...]int{2, 4, 6, 8, 10} // 不输入数组长度，让编译器来计算长度
var grid [4][5]int // 二维数组
```

- 数量写在类型前

- 可通过 _ 来省略变量，不仅仅是 range，任何地方都可通过 _ 来省略变量

  ```go
  sum := 0
  for _, v := range numbers {
      sum += v
  }
  ```

- 如果只要下标 i，可写成for i := range numbers

**为什么要用 range?**

- 意义明确、美观
- c++：没有类似能力
- Java/Python：只能 for each value，不能同时获取 i, v

**数组是值类型**

- [10]int 和[20]int 是不同类型
- 调用 func f(arr [10]int)会 **拷贝** 数组
- 在 go 语言中一般不直接使用数组（指针），使用切片

#### 切片（Slice）

```go
arr := [...]int{0, 1, 2, 3, 4, 5, 6, 7}
s := arr[2:6]
```

- s 就是一个切片，值为[2 3 4 5]

- Slice本身没有数据，是对底层 array 的一个 view

  ```go
  arr := [...]int{0, 1, 2, 3, 4, 5, 6, 7}
  s := arr[2:6]
   
  s[0] = 10
  ```

- arr 的值变为 [0 1 10 3 4 5 6 7]

**Reslice**

```go
s := arr[2:6]
s = s[:3]
s = s[1:]
s = arr[:]
```



**Slice 的扩展**

```go
arr := [...]int{0, 1, 2, 3, 4, 5, 6, 7}
s1 := arr[2:6]
s2 := s1[3:5]
```

- s1的值为？
- s2的值为？

[![img](https://cdn.nlark.com/yuque/0/2022/jpeg/1207203/1649371066190-e7c41174-7522-4a55-a994-6078906e2234.jpeg)](http://alanhou.org/homepage/wp-content/uploads/2019/04/2019040510193038.jpg)

- s1的值下为[2 3 4 5]，s2的值为[5 6]
- slice 可以向后扩展，不可以向前扩展
- s[i]不可以超越 len(s)，向后扩展不可以超越底层数组 cap(s)

**Slice 的实现**

[![img](https://cdn.nlark.com/yuque/0/2022/jpeg/1207203/1649371068301-6b47262c-fa23-4984-8b30-3245f17a0233.jpeg)](http://alanhou.org/homepage/wp-content/uploads/2019/04/2019040510202984.jpg)

**向 Slice 添加元素**

```go
arr := [...]int{0, 1, 2, 3, 4, 5, 6, 7}
s1 := arr[2:6]
s2 := s1[3:5]
s3 := append(s2, 10)
s4 := append(s3, 11)
s5 := append(s4, 12)
```

- s3, s4, s5的值为？arr 的值为？
- 添加元素时如果超越 cap，系统会重新分配更大的底层数组
- 由于值传递的关系，必须接收 append 的返回值
- s = append(s, val)
- s := make([]int, 10, 32)，其中10和32分别是 length 和 capacity

### Map

```go
m := map[string] string {
	"name": "ccmouse",
	"course": "golang",
	"site": "imooc",
	"quality": "notbad",
}
```

- map[K]V, map[K1]map[K2]V（复合 map）

#### map 的操作

- 创建：make(map[string] int)
- 获取元素：m[key]
- key 不存在时，获得Value 类型的初始值（Zero value）
- 用 value, ok := m[key]来判断是否存在 key
- 用 delete 删除一个 key

#### map 的遍历

- 使用 range 遍历 key，或者遍历 key, value 对
- 不保证遍历顺序，如需顺序，需手动对 key 排序
- 使用 len 获取元素个数

#### map 的 key

- map 使用哈希表，必须可以比较相等
- 除 slice, map, function 外的内建类型都可以作为 key
- Struct 类型不包含上述字段，也可作为 key

#### map例题

寻找最长不含有重复字符的子串

- https://leetcode.com/problems/longest-substring-without-repeating-characters/
- abcabcbb → abc
- bbbbb → b
- pwwkew → wke

**解题思路：**

对于一个字母 x

- lastOccurred[x]不存在，或者<start → 无需操作
- lastOccurred[x] >= start → 更新 start
- 更新 lastOccurred[x]，更新 maxLength

解决中文等国际化字符的问题：

**rune相当于 go 语言的 char**

- 使用 range 遍历 pos, rune 对
- 使用 utf8.RuneCountInString 获得字符数量
- 使用 len 获得字节长度
- 使用[]byte 获得字节

### 其它字符串操作

- Fields, Split, Join
- Contains, Index
- ToLower, ToUpper
- Trim, TrimRight, TrimLeft

strings 包下有更多可进行查看

## 03 面向“对象”

### 结构体和方法

#### 面向对象

- go 语言仅支持封装，不支持继承和多态
- go语言没有 class，只有 struct

**结构的创建**

```go
root.left = &treeNode{}
root.right = &treeNode{5, nil, nil}
root.right.left = new(treeNode)
```

- 不论地址还是结构本身，一律使用 . 来访问成员



| **1****2****3****4****5** | **func** **createNode****(****value** **int****)*********treeNode****{**	**return****&treeNode****{****value****:****value****}****}** **root****.****left****.****right****=****createNode****(****2****)** |
| ------------------------- | ------------------------------------------------------------ |
|                           |                                                              |

- 使用自定义工厂函数
- 注意返回了局部变量的地址！

**结构创建在堆上还是栈上？**

- 不需要知道

**为结构定义方法**

```go
func (node treeNode) print(){
	fmt.Print(node.value, " ")
}
```

- 显示定义和命名方法接收者

**使用指针作为方法接收者**

```go
func (node *treeNode) setValue(value int) {
	node.value = value
}
```

- 只有使用指针才可以改变结构内容
- nil 指针也可以调用方法！

**值接收者 vs 指针接收者**

- 要改变内容必须使用指针接收者
- 结构过大也考虑使用指针接收者（性能考虑）
- 一致性：如有指针接收者，最好都是指针接收者
- 值接收者是 go语言特有
- 值/指针接收者均可接收值/指针

### 包和封装

#### 封装

- 名字一般使用 CamelCase
- 首字母大写：public
- 首字母小写：private

#### 包

- 每个目录一个包
- main 包包含可执行入口
- 为结构定义的方法必须放在同一包内
- 可以是不同文件

**go 语言中如何扩充系统类型或者别人的类型**

- 定义别名
- 使用组合

### GOPATH以及目录结构

#### GOPATH 环境变量

- 默认在~/go(Unix, Linux), %USERPROFILE%\go(Windows)
- 官方推荐：所有项目和第三方库都放在同一个 GOPATH 下
- 也可以将每个项目放在不同的 GOPATH

以 Mac 为例(~/.bash_profile)

```bash
export GOPATH=/Users/alan/go
export PATH="$GOPATH/bin:$PATH"
```

**go get 获取第三方库**

- go get 命令演示

- 使用 gopm 来获取无法下载的包(如官网 golang.org 下的包)
  **注：**1.13以上无需安装 gopm，可借助 [Go Modules Proxy](https://github.com/golang/go/wiki/Modules#are-there-always-on-module-repositories-and-enterprise-proxies)直接使用 go get 安装

  通过Preferences 的如下配置可以在保存时自动整理所导入的包，如删除未使用或错误的包导入（IDEA 原来通过 On Save 来实现，已淘汰并即将删除）
  
  ```bash
  go get github.com/gpmgo/gopm
  gopm -g -v -u golang.org/x/tools/cmd/goimports # 未 build 执行下面的命令 build 到 bin(PATH) 目录
  go install golang.org/x/tools/cmd/goimports
  ```

[![img](https://cdn.nlark.com/yuque/0/2022/jpeg/1207203/1649371069680-4601af9a-113e-4181-a784-83eeb9eb2a72.jpeg)](http://alanhou.org/homepage/wp-content/uploads/2019/04/2019040613314795.jpg)

- go build 来编译
- go install 产生 pkg 文件和可执行文件
- go run 直接编译运行

**GOPATH下目录结构**

- src

- - git repository 1
  - git repository 2

- pkg

- - git repository 1
  - git repository 2

- bin

- - 执行文件1, 2, 3…

## 04 面向接口

```go
type Traversal interface {
    Traverse()
}
 
func main() {
    traversal := getTraversal()
    traversal.Traverse()
}
```



### duck typing的概念

[![img](https://cdn.nlark.com/yuque/0/2022/jpeg/1207203/1649371071805-3c3c1d24-6c46-4473-bf5f-80d5bfe1b5c9.jpeg)](http://alanhou.org/homepage/wp-content/uploads/2019/04/2019040615094493.jpg)

**大黄鸭是鸭子吗？**

- 传统类型系统：脊索动物门，脊椎动物亚门，鸟纲雁形目…  不是
- duck typing：是鸭子
- 概念：“像鸭子走路，像鸭子叫（长得像鸭子），那么就是鸭子”
- 描述事物的外部行为而非内部结构
- 严格说 go 属于结构化类型系统，类似 duck typing

**Python 中的 duck typing**

```python
def download(retriever):
    return retriever.get("www.baidu.com")

```

- 运行时才知道传入的 retriever 有没有 get
- 需要注释来说明接口

**C++中的 duck typing**

```c++
template <class R>
string download(const R& retriever) {
    return retriever.get("www.baidu.com")
}
```

- 编译时才知道传入的 retriever 有没有 get
- 需要注释来说明接口

**Java 中的类似代码**

```java
<R extends Retriever>
String download(R r) {
    return r.get("www.baidu.com")
}
```

- 传入的参数必须实现 Retriever 接口
- 不是 duck typing
- 同时需要 Readable, Appendable 怎么办？（apache polygene）

**Go 语言的 duck typing**

- 同时具有 Python, C++的 duck typing 的灵活性
- 又具有 Java 的类型检查

### 接口的定义和实现

#### 接口的定义

使用者（download）→实现者（retriever）

- 接口由使用者定义

```go
type Retriever interface {
	Get(url string) string
}
 
func download(r Retriever) string {
	return r.Get("http://www.baidu.com")
}
```

#### 接口的实现

- 接口的实现是隐式的
- 只要实现接口里的方法

### 接口的值类型

#### 接口变量里面有什么

**接口变量**

- 实现者的类型
- 实现者的值或实现者的指针

**接口变量里面有什么**

- 接口变量自带指针
- 接口变量同样采用值传递，几乎不需要使用接口的指针
- 指针接收者实现只能以指针方式使用；值接收者都可

#### 查看接口变量

- 表示任何类型：interface{}
- Type Assertion
- Type Switch

### 接口的组合

```go
type ReaderWriter interface {
	Reader
	Writer
}
```

### 常用系统接口

- Stringer
- Reader/Writer

## 05 函数式编程

### 函数与闭包

**函数式编程 vs. 函数指针**

- 函数是一等公民：参数、变量、返回值都可以是函数
- 高阶函数
- 函数 ￫ 闭包

**“正统”函数式编程**

- 不可变性：不能有状态，只有常量和函数
- 函数只能有一个参数
- 本课程不作上述规定

### 闭包

[![img](https://cdn.nlark.com/yuque/0/2022/jpeg/1207203/1649371072508-d5045c86-c42c-4452-aa2d-a25a16b7264e.jpeg)](https://alanhou.org/homepage/wp-content/uploads/2019/04/2019041009400420.jpg)

#### Python 中闭包



| **1****2****3****4****5****6****7****8** | **def** **adder****(****)****:**	**sum****=****0** 	**def****f****(****value****)****:**		**nonlocal** **sum**		**sum****+=****value**		**return****sum**	**return****f** |
| ---------------------------------------- | ------------------------------------------------------------ |
|                                          |                                                              |

- Python 原生支持闭包
- 使用__closure__来查看闭包的内容

#### C++中的闭包



| **1****2****3****4****5****6****7** | **auto** **adder****(****)****{**	**auto** **sum****=****0****;**	**return****[****=****]****(****int****value****)****mutable****{**		**sum****+=****value****;**		**return****sum****;**	**}****;****}** |
| ----------------------------------- | ------------------------------------------------------------ |
|                                     |                                                              |

- 过去：stl 或者 boost 带有类似库
- C++11及以后：支持闭包

#### Java 中的闭包



| **1****2****3****4****5****6****7** | **Function****<****Integer****,****Integer****>****adder****(****)****{**	**final****Holder****<****Integer****>****sum****=****new****Holder****<****>****(****0****)****;**	**return****(****Integer****value****)****-****>****{**		**sum****.****value****+=****value****;**		**return****sum****.****value****;**	**}****;****}** |
| ----------------------------------- | ------------------------------------------------------------ |
|                                     |                                                              |

- 1.8以后：使用 Function 接口和 Lambda表达式来创建函数对象
- 匿名类或 Lambda 表达式均支持闭包

### Go语言闭包的应用

- 例一：斐波那契数列
- 例二：为函数实现接口
- 例三：使用函数来遍历二叉树
  [![img](https://cdn.nlark.com/yuque/0/2022/jpeg/1207203/1649371072627-2fb0ca06-54cc-4824-b4a4-f2dd7e2dfde7.jpeg)](https://alanhou.org/homepage/wp-content/uploads/2019/04/201904200052096.jpg)

#### 总结

- 更为自然，不需要修饰如何访问自由变量
- 没有 Lambda 表达式，但是有匿名函数

## 06 错误处理和资源管理

### defer调用

- 确保调用在函数结束时发生
- 参数在 defer 语句时计算
- defer 列表为后进先出

#### 何时使用 defer 调用

- Open/Close
- Lock/Unlock
- PrintHeader/PrintFooter

### 错误处理概念

错误处理



| **1****2****3****4****5****6****7****8** | **file****,****err****:****=****os****.****Open****(****"abc.txt"****)****if****err****!****=****nil****{**	**if****pathError****,****ok****:****=****err****.****(*********os****.****PathError****)****;****ok****{**		**fmt****.****Println****(****pathError****.****Err****)**	**}****else****{**		**fmt****.****Println****(****"unknown error"****,****err****)**	**}****}** |
| ---------------------------------------- | ------------------------------------------------------------ |
|                                          |                                                              |

### 服务器统一出错处理

- 如何实现统一错误处理逻辑

### panic和recover

#### panic

- 停止当前函数执行
- 一直向上返回，执行每一层的 defer
- 如果没有遇见 recover，程序退出

#### recover

- 仅在 defer 调用中使用
- 获取 panic 的值
- 如果无法处理，可重新 panic

### 服务器统一出错处理2

#### error vs panic

- 意料之中的：使用 error。如：文件打不开
- 意料之外的：使用 panic。如：数组越界

#### 错误处理综合示例

- defer + panic + recover
- Type Assertion
- 函数式编程的应用（errWrapper）

## 07 测试与性能调优

### 测试

Debugging Sucks! Testing Rocks!

#### 传统测试 vs 表格驱动测试

**传统测试**



| **1****2****3****4****5****6****7****8** | **@****Test** **public****void****testAdd****(****)****{**	**asserEquals****(****3****,****add****(****1****,****2****)****)****;**	**asserEquals****(****2****,****add****(****0****,****2****)****)****;**	**asserEquals****(****0****,****add****(****0****,****0****)****)****;**	**asserEquals****(****0****,****add****(****-****1****,****1****)****)****;**	**asserEquals****(****Integer****.****MIN_VALUE****,****add****(****1****,****Integer****.****MAX_VALUE****)****)****;**	**}** |
| ---------------------------------------- | ------------------------------------------------------------ |
|                                          |                                                              |

- 测试数据和测试逻辑混在一起
- 出错信息不明确
- 一旦一个数据出错测试全部结束

**表格驱动测试**

- 分离的测试数据和测试逻辑
- 明确的出错信息
- 可以部分失败
- go 语言的语法使得我们更易实践表格驱动测试

### 代码覆盖率和性能测试

#### 代码覆盖率



| **1****2** | **go** **test****-****coverprofile****=****c****.****out** **go** **tool** **cover****-****html****=****c****.****out** |
| ---------- | ------------------------------------------------------------ |
|            |                                                              |

#### 性能测试(Benchmark)



| **1** | **go** **test****-****bench****.** |
| ----- | ---------------------------------- |
|       |                                    |

### 使用pprof进行性能调优



| **1****2****3****4****5****6** | **go** **test****-****bench****.****-****cpuprofile** **cpu****.****out****go** **tool** **pprof** **cpu****.****out****# 进入 pprof 交互式命令行，最简单的为输入 web生成图形** **#  Failed to execute dot. Is Graphviz installed? Error: exec: "dot": executable file not found in $PATH****brew** **install** **graphviz** |
| ------------------------------ | ------------------------------------------------------------ |
|                                |                                                              |


[![img](https://cdn.nlark.com/yuque/0/2022/jpeg/1207203/1649371075814-f458dfcb-fe51-40bc-b894-c6cb26fc3634.jpeg)](https://alanhou.org/homepage/wp-content/uploads/2019/04/2019062315471254.jpg)

### http 测试

- 通过使用假的 Request/Response（TestErrWrapper）
- 通过起服务器（TestErrWrapperInServer）

### 生成文档和示例代码

#### 文档

- 用注释写文档
- 在测试中加入 Example
- 使用 go doc/godoc 来查看/生成文档

| **1****2****3** | **go** **doc****<****pkg****>****# 网页文档****godoc****-****http****:****6060** |
| --------------- | ------------------------------------------------------------ |
|                 |                                                              |

## 08 Goroutine

#### 协程 Coroutine

- 轻量级“线程”
- 非抢占式多任务处理，由协程主动交出控制权
- 编译器/解释器/虚拟机层面的多任务
- 多个协程可能在一个或多个线程上运行



| **1** | **go** **run****-****race** **xxx****.****go****// 检测数据访问冲突** |
| ----- | ------------------------------------------------------------ |
|       |                                                              |

Subroutines are special cases of more general program components, called *coroutines*.

子进程是协程的一个特例

普通函数：线程  main ➝ doWork

协程：线程（可能） main ⟺ doWork

#### 其它语言中的协程

- C++：Boost.Corouting
- Java：不支持
- Python：使用 yield 关键字实现协程，Python 3.5加入了 async def 对协程原生支持

#### goroutine 的定义

- 任何函数只需加上 go 就能送给调试器运行
- 不需要在定义时区分是否是异步函数
- 调度器在合适的点进行切的
- 使用-race 来检测数据访问冲突

#### go routine 可能的切换点

- I/0, select
- channel
- 等待锁
- 函数调用（有时）
- runtime.Gosched()
- 只是参考，不能保证切换，不能保证在其他地方不切换

## 09 管道Channel

[![img](https://cdn.nlark.com/yuque/0/2022/jpeg/1207203/1649371078990-656c6346-abc8-4bd7-93a5-3418e6c22a6d.jpeg)](https://alanhou.org/homepage/wp-content/uploads/2019/04/2019062513340472.jpg)

- channel
- buffered channel
- range
- 理论基础：Communication Sequential Process (CSP)
- Don’t communicate by sharing memory; share memory by communicating.
- 不要通过共享内存来通信；通过通信来共享内存

例一：使用 Channel 来等待 goroutine 结束

- 以及 WaitGroup的使用

例二：使用 Channel 来实现树的遍历

例三：使用 Select 来进行调度

- Select 的使用
- 定时器的使用
- 在 Select 中使用 Nil Channel

传统同步机制

- WaitGroup
- Mutex
- Cond

## 10 http及其他标准库

### http

- 使用 http 客户端发送请求
- 使用 http.Client 控制请求头部等
- 使用 httputil 简化工作

### http 服务器的性能分析

- import _ “net/http/pprof”
- 访问/debug/pprof
- 使用 go tool pprof 分析性能

| **1** | **# 内存使用情况** |
| ----- | ------------------ |
|       |                    |

**2**

**3**

**4**

**go** **tool** **pprof** **http****:****//localhost:6060/debug/pprof/heap**

**# 30秒 CPU使用情况**

**go** **tool** **pprof** **http****:****//localhost:6060/debug/pprof/profile?seconds=30**

### 其它标准库

- bufio
- log
- encoding/json
- regexp
- time
- strings/math/rand

**文档**

- godoc -http :8888
- https://studygolang.com/pkgdoc

第三方 http 框架

- gin-gonic

- - middleware的使用
  - context的使用

## 11 迷宫的广度优先搜索

### 广度优先算法

- 为爬虫实战项目做发准备
- 应用广泛，综合性强
- 面试常见

[![img](https://cdn.nlark.com/yuque/0/2022/jpeg/1207203/1649371080235-b66d8792-f2e7-4a5f-ad12-220792357e1e.jpeg)](https://alanhou.org/homepage/wp-content/uploads/2019/06/201906290647333.jpg)

例：广度优先搜索走迷宫

- 用循环创建二维 slice
- 使用 slice 来实现队列
- 用 Fscanf 读取文件
- 对 Point 的抽象

## 12开始实战项目

### 爬虫项目介绍

#### 为什么做爬虫项目

- 有一定的复杂性
- 可以灵活调整项目的复杂性
- 平衡语言/爬虫之间的比重

#### 网络爬虫分类

- 通用爬虫，如 baidu, google
- 聚焦爬虫，从互联网获取结构化数据

#### go语言的爬虫库/框架

- henrylee2cn/pholcus
- gocrawl
- colly
- hu17889/go_spider

#### 本课程爬虫项目

- 将不使用现成爬虫库/框架
- 使用 ElasticSearch 作为数据存储
- 使用 Go 语言标准模板库实现 http 数据展示部分

#### 爬虫的主题

爬取内容

- 内容：如新闻，博客，社区…

爬取人

- QQ 空间，人人网，微博，微信，facebook？
- 相亲网站，求职网站
- 出于隐私和趣味性考虑，本课程将爬取相亲网站

[![img](https://cdn.nlark.com/yuque/0/2022/jpeg/1207203/1649371082159-36040aec-da1a-40b4-95e4-2d8df3113af2.jpeg)](https://alanhou.org/homepage/wp-content/uploads/2019/04/2019062907200454.jpg)

### 总体算法

[![img](https://cdn.nlark.com/yuque/0/2022/jpeg/1207203/1649371082313-6d6751a7-5304-4c3c-8b33-603cdf69003f.jpeg)](https://alanhou.org/homepage/wp-content/uploads/2019/04/2019062907203058.jpg)

## 13 单任务版爬虫



## 14 并发版爬虫

### 数据存储和展示

## 15 分布式爬虫

## 16 课程总结

更新中…

其它示例代码：https://github.com/e421083458/gateway_demo

## 17 常用命令汇总



| **1****2****3****4****5****6****7****8****9****10****11****12****13****14****15****16****17****18** | **go** **run** **xxx****go** **build** **xxx****go** **doc** **fmt****.****Printf****godoc****-****http****:****6060****goimports****-****l****-****w****.****golint****.****/****.****.****.****go** **vet****.****/****.****.****.****# 获取外部包****go** **get****-****v****github****.****com****/****mactsouk****/****go****/****simpleGitHub****# 清理外部包****$****go** **clean****-****i****-****v****-****x****github****.****com****/****mactsouk****/****go****/****simpleGitHub****$****rm****-****rf****~****/****go****/****src****/****github****.****com****/****mactsouk****/****go****/****simpleGitHub****# 临时安装其它指定版本****go** **get** **golang****.****org****/****dl****/****go****.****1.15.6****go1****.****15.6****download****# 删除该临时安装版本****rm****-****rf****$****(****go1****.****15.6****env** **GOROOT****)****rm****$****(****go** **env** **GOPATH****)****/****bin****/****go1****.****15.6** |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
|                                                              |                                                              |

## 18 其它

下载安装

https://golang.org/dl/
https://golangtc.com/download



| **1****2****3****4****5****6****7****8****9****10****11****12****13****14****15****16****17****18****19****20****21****22****23****24****25****26****27** | **# 源码安装环境变量示例配置****export** **GOROOT****=/****usr****/****local****/****go****# 通常可自动推导，无需单独配置****export** **GOPATH****=/****Users****/****goRoot****:****/****Users****/****go****export** **PATH****=/****usr****/****local****/****go****/****bin****:****/****Users****/****goRoot****/****bin****:****$****PATH** **go** **doc** **fmt****.****Prinf****# 文档查看****godoc****-****http****=****:****8001****# 通过启动web服务器来浏览器端查看****go** **build** **xxx****.****go****# 编译****go** **run** **xxx****.****go****# 直接运行****go** **get****-****v****github****.****com****/****mactsouk****/****go****/****simpleGitHub****# 下载外部包，所在位置~/go/src；编译后包所在位置~/go/pkg/darwin_amd64/****go** **clean****-****i****-****v****-****x****github****.****com****/****mactsouk****/****go****/****simpleGitHub****# 清理外部包中间文件** **# 安装 Gin框架****go** **get****-****u****github****.****com****/****gin****-****gonic****/****gin** **# 安装 Beego 框架****go** **get****-****u****github****.****com****/****astaxie****/****beego****go** **get****-****u****github****.****com****/****beego****/****bee** **#创建****bee** **api** **xxx****# 打包****bee** **pack****-****be** **GOOS****=****linux** **# 使用 Gomod****export** **GO111MODULE****=****on****# Windows： go env -w GO111MODULE=on****go** **get****-****u****github****.****com****/****go****-****kratos****/****kratos****/****tool****/****kratos** |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
|                                                              |                                                              |

http://gorm.book.jasperxu.com/

http://www.topgoer.com/



| **1****2****3** | **# beego****github****.****com****/****astaxie****/****beego****/****client****/****orm****github****.****com****/****go****-****sql****-****driver****/****mysql** |
| --------------- | ------------------------------------------------------------ |
|                 |                                                              |

- [Go 语言规范文档](https://golang.google.cn/ref/spec)
- [Golang 基本项目结构](https://github.com/golang-standards/project-layout)
- [go-micro](https://github.com/asim/go-micro)

| **1** | **go** **get** **github****.****com****/****micro****/****go****-****micro** |
| ----- | ------------------------------------------------------------ |
|       |                                                              |

- [微服务-服务发现 etcd](https://github.com/etcd-io/etcd)

**2**

**3**

**4**

**5**

**6**

**7**

**8**

**9**

**go** **get** **github****.****com****/****micro****/****protobuf****/****{****proto****,****proto****-****gen****-****go****}**

**go** **get** **github****.****com****/****micro****/****protoc****-****gen****-****micro**

**protoc****-****proto_path****=****xx****-****go_out****=****xxx****-****micro_out****=****xxx** **xxx****.****proto**

 

**go** **get** **github****.****com****/****micro****/****micro**

**# export MICRO_REGISTRY=etcd**

**micro** **api****--****handler****=****api****--****address****=****0.0.0.0****:****8085**

**micro** **web**

## 19 常见问题

1、dyld: malformed mach-o image: segment __DWARF has vmsize < filesize

这是 macOS 升级到 Catalina 之后出现的问题，使用go build -ldflags “-w”来代替 go build：