# cmdb

cmdb 中存放部署产品的各个组件所用到的服务器信息或者平台类的信息。根据组件的部署形态，服务器部署形态的组件需要用到主机清单， Pod 部署形态的组件需要用到容器平台的信息。

服务器的信息存放在 cmdb.Inventory 中。平台类的信息存放在 cmdb.Platforms 中。

## cmdb.Inventory

在实际环境的目录下（`targets/<target>/<zone>`）的 `hosts.yml` 文件，就是文本格式的 Inventory 信息。`hosts.yml` 本质上就是 Ansible 能够识别的主机清单。

### 如何更改一个主机的 SSH 连接信息

`ansible_user` 和 `ansible_port` 变量可以设置在 Host、Group 或者 `all` Group 下面。

```yaml
all:
  vars:
    ansible_user: root
    ansible_port: 22
somegroupname:
  vars:
    ansible_user: root
    ansible_port: 22
  hosts:
    10.0.0.1:
      ansible_user: root
      ansible_port: 2222
```

## 主机的 SSH 连接地址并不是 Inventory 的 hosts 信息

使用 `ansible_host` 变量指定 SSH 的连接地址。

```yaml
somegroupname:
  hosts:
    10.0.0.1:
      ansible_user: root
      ansible_port: 2222
      ansible_host: 127.0.0.1
```
