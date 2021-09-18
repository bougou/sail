# 产品

一个产品的运维代码，需要以产品名作为文件夹名称，放到 `sail` 指定的 `products-dir` 目录下。
文档中其它部分都假设 `products-dir` 目录为 `products`。

`products/<product-name>` 目录下面就是该产品的「部署代码」，里面封装了该产品的所有部署及运维逻辑。

理想情况下，products 目录的内容应该维护在一个 Git 仓库中（或者按产品维护在多个仓库里）。
当产品的运维逻辑需要更新时，你都应该去修改仓库中对应的产品运维代码。

> `products` 目录下除了各个产品的子文件夹外，有三个「文件或目录名」被 `sail` 用作特殊的目的。
> - `ansible.cfg` Ansible 的配置文件，`sail` 在执行的时候会自动生成该文件（如果不存在）
> - `shared_roles` 存放被多个产品公用的 role(s)
> - `shared_tasks` 存放通用的 Ansible Tasks，可以看做工具库

## Product 的目录结构

`products/<product-name>` 的目录结构：

```bash
## sail 规范

components.yaml       # 组件声明
components/           # 组件声明
vars.yaml             # 通用变量，通常定义一些所有组件都通用的变量，如数据目录 data_dir，时区 timezone
# 所有的「组件变量」以及「通用变量」称为「产品变量」
# 把产品部署到不同环境中时，每个环境都拥有自己独立的一份「产品变量」，可以称作「环境配置」。
# 产品的组件代码中可以引用这里的「环境配置」变量，比如渲染模板等。
# 每个环境的「环境配置」位于环境自己的目录下面（targets/<target>/<zone>/)，可以按照环境实际情况修改。
# `sail` 会自动负责一个环境的「环境配置」变量与 `products/<product-name>` 下的「产品变量」的结构保持一致。
# 比如 「产品变量」中新增加的变量会自动同步到「环境配置中」。

## 通用
resources/            # 其它默认资源，如默认的证书文件、试用 License 文件、默认的 icon 图标等

## 其它文件
README.md             # 产品的介绍文档

## Ansible 相关
sail.yaml             # 必须存在的一个 Ansible Playbook 文件，可以使用 `sail gen-sail` 命令自动生成 sail.yaml playbook 文件
<playbook>.yaml       # 其它的 Playbook 文件


## 组件代码目录
roles/                # Roles，实现各个组件的实际的、具体的安装逻辑


## Helm Chart 文件
# 如果你把整个产品作为一个 Helm Chart 来开发，你可以直接把 `products/<product-name>` 目录作为 Helm 的 Chart 目录来使用。
Chart.yaml    # helm Chart.yaml
values.yaml   # helm values.yaml
templates     # helm templates
crds          # helm crds
charts        # helm charts
# Sail 除了支持把整个产品当做一个 Helm Chart，还支持把产品中的每一个组件当做一个 Helm Chart。
```

一个产品是由多个组件组成的，你需要为每一个 [组件](./component.md) 编写运维代码。

## 产品的变量

在定义产品时，你必须声明该产品在部署时会用到的「所有的变量」。products 下面的「产品运维代码」则需要按照「这套变量」去执行完成相应的动作。

对于不同的部署环境，「运维操作者」只需要维护「这套变量」与实际的环境相符，然后通过执行 `sail` 命令去安装和升级产品即可。

>「运维操作者」在执行 `sail` 命令时，不需要去关心产品或产品的组件是怎么安装部署的。**安装部署的细节** 已经由「产品运维代码」实现。这种「运维模式」就像一个普通的用户在自己的电脑上安装一款电脑软件一样，如果安装过程出错，原因无外乎两种：1，用户自己的电脑环境不对；2，软件的安装流程本身有 Bug。
>
> 当通过 `sail` 安装部署产品的过程中失败了，原因也无外乎两种：1，部署环境的变量没有配置对；2，「产品运维代码」本身有 Bug。

产品变量分为「组件变量」和「通用变量」。

- 使用 `components.yaml` 文件或和 `components/` 目录来声明「组件变量」。组件变量是以组件名称作为顶层变量的。见 [组件声明](./component.md)
- 使用 `vars.yaml` 文件定义与具体组件无关的「通用变量」，如数据目录 `data_dir`，时区 `timezone` 等。

你自己必须确保所有的「组件变量」即所有的组件名称不能和「通用变量」冲突。Sail 会把所有的「组件变量」和「通用变量」合并在一起保存在具体部署环境的 `targets/<target_name>/<zone_name>/vars.yaml` 中。

在定义产品的「通用变量」时，除了要避免和组件名称冲突外，`sail` 还保留了下面一些变量，请不要使用下面变量。

- `_sail` 开头的变量。Sail 保留使用所有以 `_sail` 开头的变量。
- `inventory`
- `platforms`
- `targetvars`
