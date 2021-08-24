# 产品

products 目录下的每一个 `<product-name>` 目录就是该产品的「部署代码」。里面封装了该产品的所有部署及运维逻辑。

理想情况下，products 目录的内容应该维护在一个 Git 仓库中（或者按产品维护在多个仓库里）。
当产品的运维逻辑需要更新时，你都应该去修改仓库中对应的产品运维代码。

## Product 的目录结构

`products/<product-name>` 目录下的结构：

```bash
## sail 规范
components.yml      # 组件声明
components/         # 组件声明
vars.yml            # 非组件变量，通常定义一些所有组件都通用的变量，如数据目录 data_dir，时区 timezone
# 所有的「组件变量」以及「非组件变量」最终会合并在一起，称为「环境配置」
# 把产品部署到不同环境中时，每个环境都拥有自己独立的一份「环境配置」
# 每个环境的「环境配置」位于环境自己的目录下面（targets/<target>/<zone>/)，可以按照环境实际情况修改

## Ansible 相关
sail.yml            # 必须存在的一个 Ansible Playbook 文件
<playbook>.yml      # 其它的 Playbook 文件
roles/              # Ansible Roles，实现各个组件的实际的、具体的安装逻辑

## Helm 相关
templates/          # Helm templates
Chart.yml           # Helm Chart definition
values.yml          # Helm values 变量定义

## 通用
resources/          # 其它默认资源，如默认的证书文件、试用 License 文件、默认的 icon 图标等

## 其它文件
README.md           # 产品的介绍文档

```
