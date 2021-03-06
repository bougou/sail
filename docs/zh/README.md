# sail

`sail` 是一个运维框架，用于实现产品的安装部署，变更升级等运维工作。

> `sail` 尤其适合软件产品的私有化交付工作。

## 主要概念

`sail` 有三个主要概念：产品（Product)，环境 (Target)，和包（Package）

### Product 产品

「产品」表示一个具体的软件产品。从运维的角度理解，一个软件「产品」是由一些「组件」组成的。一个「产品」可以简单到由几个「组件」组成，也可以复杂到由成百上千个「组件」组成。

当你负责运维一个产品时，你应该实现（或由其他人实现）该产品的运维代码。「产品运维代码」，Product Operation Code (POC) 并不是那些由「产品开发工程师」实现的产品功能代码，而是用来指导「产品运维工程师」去安装部署以及管理该产品的运维操作代码。

### Target 环境

「环境」表示用来部署和运行一个产品的实际服务器环境。

在 `sail`, 用了两个层级来表达环境，分别是 `target` 和 `zone`。

`target(s)` 之间是独立的，`target` 下面可以有多个 `zone` 组成。

### Package 包

「包」表示具体的软件制品。常见的各种软件安装包文件，如 `.rpm`，`.tar.gz`, `.gzip` 等，或者容器镜像文件。
有些软件除了安装包外，可能还有一些体积比较大的资料包，文档包，镜像包等，这些也可以算作 Package。

## 为什么使用 `sail`

运维工作的核心就是把一个「产品」以「正确的方式」部署到一个具体的「环境」中。

`sail` 的目的就是把这个操作变简单，简单到只需要执行一个命令 `sail apply <特定的环境>`。

为了实现这个目标，`sail` 制定了一套为「产品」编写 「产品运维代码」（Product Operation Code）的规范。
在 `sail` 中，一个「产品」是由「组件」构成，所以开发「产品运维代码」的主要工作就是开发各个组件的运维代码。

`sail` 定义了组件的运维代码规范，用户可以使用一种统一的方式去编写组件的运维代码。

需要声明的是，`sail` 的组件规范并不会把运维工程师已经拥有的运维技能知识变过时，而是在这些技能知识的基础上，把组件的运维变的更有条理，更简单，更清晰。

- [如何开发产品运维代码](./product.md)
- [如何开发组件](./component.md)
- [如何维护环境的 CMDB](./cmdb.md)
- [Sail 命令使用](./sail-commands.md)
