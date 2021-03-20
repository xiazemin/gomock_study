https://github.com/golang/mock


go mod init gomock_study
go get github.com/golang/mock/gomock
go get github.com/golang/mock/mockgen


反射模式
通过构建一个程序用反射理解接口生成一个mock类文件，它通过两个非标志参数生效：导入路径和用逗号分隔的符号列表（多个interface）。
mockgen -destination mock_sql_driver.go database/sql/driver Conn,Driver

源码模式
通过一个包含interface定义的文件生成mock类文件，它通过 -source 标识生效，-imports 和 -aux_files 标识在这种模式下也是有用的。

 mockgen -source=exp1/foo.go -destination=exp1/mock/mock_foo.go


mock控制器
mock控制器通过NewController接口生成，是mock生态系统的顶层控制，它定义了mock对象的作用域和生命周期，以及它们的期望。多个协程同时调用控制器的方法是安全的。
当用例结束后，控制器会检查所有剩余期望的调用是否满足条件。

mock对象的行为注入
对于mock对象的行为注入，控制器是通过map来维护的，一个方法对应map的一项。因为一个方法在一个用例中可能调用多次，所以map的值类型是数组切片。当mock对象进行行为注入时，控制器会将行为Add。当该方法被调用时，控制器会将该行为Remove。


行为调用的保序
默认情况下，行为调用顺序可以和mock对象行为注入顺序不一致，即不保序。如果要保序，有两种方法：

通过After关键字来实现保序
通过InOrder关键字来实现保序


关键字InOrder是After的语法糖，源码如下：
// InOrder declares that the given calls should occur in order.
func InOrder(calls ...*Call) {
    for i := 1; i < len(calls); i++ {
        calls[i].After(calls[i-1])
    }
}


当测试用例执行完成后，并没有回滚interface到真实对象，有可能会影响其它测试用例的执行。所以，笔者强烈建议大家使用GoStub框架完成mock对象的注入。

stubs := StubFunc(&redisrepo.GetInstance, mockDb)
defer stubs.Reset()


GoConvey + GoStub + GoMock组合使用
```
 Convey("create obj", func() {
            ctrl := NewController(t)
            defer ctrl.Finish()
            mockRepo := mock_db.NewMockRepository(ctrl)
            mockRepo.EXPECT().Retrieve(Any()).Return(nil, ErrAny)
            mockRepo.EXPECT().Create(Any(), Any()).Return(nil)
            mockRepo.EXPECT().Retrieve(Any()).Return(objBytes, nil)
            stubs := StubFunc(&redisrepo.GetInstance, mockRepo)
            defer stubs.Reset()
            ...
        })
```

全局变量可通过GoStub框架打桩
过程可通过GoStub框架打桩
函数可通过GoStub框架打桩
interface可通过GoMock框架打桩
https://www.jianshu.com/p/f4e773a1b11f

https://github.com/prashantv/gostub

go get github.com/prashantv/gostub

https://www.jianshu.com/p/44355571888d?from=timeline&isappinstalled=0

https://www.sohu.com/a/210573755_99930294

GoMock非常优秀，但是对于普通的函数打桩来说也有一些缺点：

必须引入额外的抽象(interface)
打桩过程比较重
既有代码必须适配新增的抽象
我们知道，Golang支持闭包，这使得函数可以作为另一个函数的参数或返回值，而且可以赋值给一个变量。

func Exec(cmd string, args ...string) (string, error) {}

这种函数没法直接用gostub的
需要改成
var Exec = func(cmd string, args ...string) (string, error) {}





其实GoStub框架专门提供了StubFunc函数用于函数打桩，我们重构打桩代码：

stubs := StubFunc(&Exec,"xxx-vethName100-yyy", nil)

https://cloud.tencent.com/developer/article/1076111



stubs := Stub(&Exec, func(cmd string, args ...string) (string, error) { return "xxx-vethName100-yyy", nil})
defer stubs.Reset()

GoStub框架专门提供了StubFunc函数用于函数打桩，我们重构打桩代码：

stubs := StubFunc(&Exec,"xxx-vethName100-yyy", nil)
defer stubs.Reset()

Golang的库函数或第三方的库函数
定义库函数的变量：

package adaptervar 
Stat = os.Statvar 
Marshal = json.Marshalvar 
UnMarshal = json.Unmarshal...

源码解析

gostub 源码很简单，只有gostub.go一个源文件

核心原理就是利用反射进行值的替换

以函数打桩的接口函数​为例：
func StubFunc(funcVarToStub interface{}, stubVal ...interface{}) *Stubs {
	return New().StubFunc(funcVarToStub, stubVal...)
}
其实底层和​变量替换的实现是差不多的。

核心的数据如下，存储了最初的原始变量，便于在stub结束后数据的恢复，用stubs 变量存储了被替换的变量和替换后变量的映射。
type Stubs struct {
	// stubs is a map from the variable pointer (being stubbed) to the original value.
	stubs   map[reflect.Value]reflect.Value
	origEnv map[string]envVal
}

核心的stub函数如下：
1，先通过反射获取被替换变量的值，和即将替代的值
2，把原始的值用刚刚讲的map存起来
3，修改被替换变量的值
4，用到了反射的核心函数有
reflect.ValueOf //获取interface的值
v.Elem().Interface()//获取interface 包含的值的接口
v.Elem().Set(stub) //修改interface 包含的值

func (s *Stubs) Stub(varToStub interface{}, stubVal interface{}) *Stubs {
	v := reflect.ValueOf(varToStub)
	stub := reflect.ValueOf(stubVal)

    if _, ok := s.stubs[v]; !ok {
		// Store the original value if this is the first time varPtr is being stubbed.
		s.stubs[v] = reflect.ValueOf(v.Elem().Interface())
	}

	// *varToStub = stubVal
	v.Elem().Set(stub)
	return s
}

恢复最初定义变量或者函数

func (s *Stubs) Reset() {
  for v, originalVal := range s.stubs {
    v.Elem().Set(originalVal)
  }
}

对于函数替换过程稍微有点复杂，需要对函数的内容进行重新构造，具体代码如下
func FuncReturning(funcType reflect.Type, results ...interface{}) reflect.Value {
	return reflect.MakeFunc(funcType, func(_ []reflect.Value) []reflect.Value {
		return resultValues
	})
}
构造函数的过程用到了反射的MakeFunc方法

reflect.MakeFunc //用给定函数来构造出funcType类型的函数，底层和c函数的转换类似，​就是函数指针的转换。

