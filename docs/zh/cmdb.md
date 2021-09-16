# cmdb

cmdb 中存放部署产品的各个组件所用到的服务器信息或者平台类的信息。根据组件的部署形态，服务器部署形态的组件需要用到主机清单， Pod 部署形态的组件需要用到容器平台的信息。

服务器的信息存放在 cmdb.Inventory 中。平台类的信息存放在 cmdb.Platforms 中。

## cmdb.Inventory

在环境的目录下（`targets/<target>/<zone>`）的 `hosts.yaml` 文件，就是文本格式的 Inventory 信息。`hosts.yaml` 本质上就是 Ansible 能够识别的主机清单。

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

如果主机的 SSH 连接地址并不是 Inventory 的 `hosts` 字段下的地址，可以使用 `ansible_host` 变量指定 SSH 的连接地址，ansible 的 `inventory_hostname` 变量依然是 `hosts` 字段下的地址。

```yaml
somegroupname:
  hosts:
    10.0.0.1:
      ansible_user: root
      ansible_port: 2222
      ansible_host: 127.0.0.1
```

## cmdb.Platforms

在环境的目录下（`targets/<target>/<zone>`）的 `platforms.yaml` 文件中配置组件部署的 K8S 平台的信息。

通常一个产品的所有组件都部署到一个 K8S 集群中的一个 namespace 下。这种情况，`platforms.yaml` 中只需要配置一条 `all` 的记录，如下。

```yaml
all:
  k8s:
    kubeConfig: ~/.kube/config
    kubeContext: ""
    namespace: default
```

如果 Pod 组件需要部署到不同的 K8S 集群或者是不同的 namespace 下，你可以分别为特定的组件配置 K8S 平台信息。如下：

```yaml
all:
  k8s:
    kubeConfig: ~/.kube/config
    kubeContext: ""
    namespace: default

component-name1:
  k8s:
    kubeConfig: ~/.kube/config
    kubeContext: ""
    namespace: another

component-name2:
  k8s:
    kubeConfig: ~/.kube/config
    kubeContext: ""
    namespace: another
```

`all` 记录用于所有没有明确配置平台信息的组件。
