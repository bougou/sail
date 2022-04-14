# sail 命令

## 变量

Ansible 和 Helm 的核心操作是根据变量去渲染出配置文件。`sail` 做的最多的一项工作就是针对「环境」的变量的处理。

`sail` 会把下面一些变量传给 Ansible 和 Helm，Ansible 和 Helm 可以使用这些变量去执行渲染动作。

1. 环境特有的变量信息

    - `targets/<target_name>/<zone_name>/vars.yaml` 文件中的变量，包含所有组件变量以及通用变量。
    - `targets/<target_name>/<zone_name>/_computed.yaml`

    在使用 Helm 部署容器组件时，还会传递一些其它的文件。见文档 [Helm](./helm.md)。

2. 几个 `sail` 相关变量

    - `sail_packages_dir`
    - `sail_targets_dir`
    - `sail_products_dir`
    - `sail_target_dir`
    - `sail_zone_dir`
    - `sail_target_name`
    - `sail_zone_name`

除了 `sail` 传递的变量，Ansible 或 Helm 本身也能够识别很多内置的变量。

- [Ansible Special Variables](https://docs.ansible.com/ansible/latest/reference_appendices/special_variables.html)
- [Helm Builtin Objects](https://helm.sh/docs/chart_template_guide/builtin_objects/)

## sail list-components

列出一个产品中所有的组件。

`list-components` 是一个无须制定环境信息 (target 和 zone 选项）的命令。

```bash
$ ./sail list-components -p <productName>
the product <productName> contains (<number>) components:
- <componentName1>
- <componentName2>
- ...
```

## sail conf-create

创建一个全新的部署环境。

在 `conf-create` 时，部署人员必须做的事情是：

- 决定环境的名称
- 决定该环境要部署的产品名称
- 决定该产品中那些组件需要部署
- 决定每个组件分别部署到哪些服务器上

1. 环境名称通过`-t <targetName> -z <zoneName>` 指定。
2. 要部署的产品通过 `-p <productName>` 指定。产品名称会持久化在 Zone 的环境信息变量中。这样，在运行其它 sail 命令时，只需要制定环境名称，该环境部署的产品就自动解析出来了。
3. 哪些组件需要部署和部署到哪些服务器上，这两步往往和环境以及部署需求相关。灵活性非常高。

组件是否部署的计算原则：

1. 一个组件是否部署的默认值取自「产品运维代码」中的组件声明中的 `enabled` 默认值。
2. 使用 `--hosts componentName1/ip1,ip2,ip3` 选项明确制定的组件会部署，且对应的主机组设置为相应的值。

```bash
$ sail conf-create -t <targetName> -z <zoneName> \
    -p <productName> \
    --hosts ip1,ip2,ip3 \                   # will be used for all components which are not explicityly specified by --hosts options
    --hosts componentName1/ip1,ip2,ip3 \
    --hosts componentName2/ip11,ip12,ip13
```

### `--hosts` 选项格式

`--hosts` 选项可以多次指定，其参数值有两种格式：

- `--hosts <ip1>[,<ip2>,<ip3>,...]` 由任意数量的 IP 地址组成（由逗号连接），这种形式必须且只能使用一次，指定多次时，后指定的有效。
- `--hosts <componentName1>[,<componentName2>,<componentName3>,...]/<ip1>[,<ip2>,<ip3>,...]` 由任意数量的组件名（由逗号连接）和任意数量的 IP 地址（由逗号连接）组成，组件名列表和 IP 地址列表之间通过斜杠分隔。这种形式可以多次指定。

> 提示：使用 `list-components` 查看一个产品的组件列表。

## sail conf-update

用于更新一个部署环境的配置文件。

自动同步「产品运维代码」中的变动（主要是新增的变量）到特定的环境中。

事实上，`sail conf-create` 做的事情，`sail apply` 和 `sail upgrade` 命令也会在后台自动地执行。
所以，你很少需要手动去执行 `sail conf-create` 命令。

## sail apply

执行部署操作。

通过 `ansible-playbook` 命令部署常规组件，通过 `helm` 部署容器组件。

```bash
$ sail apply -t <targetName> -z <zoneName> [-c <componentName>]
```

### 透传 `ansible-playbook` 的选项

```bash
$ sail apply -t <targetName> -z <zoneName> -- <put-any-ansible-playbook-options-here>

# eg
$ sail apply -t <targetName> -z <zoneName> -- --skip-tags sometag1,sometag2
```

## sail upgrade

用于升级特定的组件，比如组件版本更新，组件配置变化等。

```bash
$ sail upgrade -t <targetName> -z <zoneName> [-c <componentName>]
```

事实上，`sail apply` 和 `sail upgrade` 之间的区别很小。

对于使用 `helm` 部署容器组件，`apply` 和 `upgrade` 没有任何区别。

对于使用 `ansible-playbook` 部署常规组件：

1. 如果没有使用 `--component` 选项指定了特定的组件，`apply` 和 `upgrade` 没有任何区别。
2. 如果使用 `--component` 选项指定了特定的组件，`apply` 会传给  `--tags play-<componentName>` 选项给 `ansible-playbook`，`upgrade` 则会传递 `--tags update-<componentName>` 选项给 `ansible-playbook`。
