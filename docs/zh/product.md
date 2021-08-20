# 产品

假设我们负责一个产品的运维，产品名称为 `foobar`。首先以产品的名字创建一个文件夹，该文件夹里面包含该产品的运维代码。

这个文件夹的结构如下：

```bash
components/     # 该产品的组件声明
component.yml   # 该产品的组件声明

vars.yml        # 该产品在不同环境中部署时，可能需要改变的变量

sail.yml        # 一个Ansible Playbook 文件，使用命令自动生成

roles/          # Ansible Roles，实现各个组件的实际的、具体的安装逻辑

README.md
```
