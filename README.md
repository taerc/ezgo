# easygo

常用的一些应用库,用于开发一些项目的基础模板，这样用起来会方便一些。

## 使用过程中的一些方便的地方

- 今天在处理参数保持时，确实会遇到统一个项目需要把参数保存在不同的库里面

## TODO

- 配置文件热加载
- 配置文件按照不同组件独立出来
- 配置文件检查，去掉对应的空格，转换成为系统绝对路径
- 测试和验证 gorm 相关操作
- 连接池相关操作
- controller-service-dao

# 通用模块

```go

"github.com/sirupsen/logrus"
"gopkg.in/ini.v1"

```

## licence

- https://github.com/google/flatbuffers

## RESTfull

- URI 本意只有资源
- 资源和操作独立开来
- 关系也要被抽象为一种资源

### 参考链接

- https://martinfowler.com/articles/richardsonMaturityModel.html
- http://www.ruanyifeng.com/blog/2011/09/restful.html