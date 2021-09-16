# 组件

一个产品是由多个组件组成。你可以使用 `components.yaml` 文件或和 `components` 目录来定义（声明）产品的组件。

## 什么是组件

一个组件实现了一个产品的某个功能，或提供产品运行时的需要的某个能力。
一个组件可以以多个副本运行（比如组成多节点集群，K8S Pods）
一个组件的任何「一个副本」只能部署到一个「计算单元」中，即一个服务器中或一个 Pod 中。

如果多个程序在部署时，必须部署在一起的话，那就可能意味着应该把它们包含在一个组件中。

## 组件声明

`sail` 的组件声明有具体的规范要求。

```yaml
# 组件名称
redis:
  # 组件的版本
  version: 3.2.10

  # 是否部署该组件
  enabled: false

  # 该组件是否由外部系统提供（如使用公有云提供的服务）
  external: false

  # 该组件对外提供的服务（端口）
  services:
    # each key represents a serviceName, it can be any string
    default:
      # 服务监听端口（或者是服务对外暴露的端口）
      port: 6379
      # other fields
    sentinel:
      port: 7379
    check:
      port: 6479

  # 组件相关的变量
  # 任意自定义
  vars:
    database: 6
    pass: ePxoMflue6jhx9oYvYV3

  # 组件相关的 Tag
  # 任意自定义
  # 目前与 `vars` 字段没有任何区别
  tags: {}

```

基本上，「组件声明」中除了组件名称外的所有字段都是可选的。一个最简单的「组件声明」如下：

```yaml
theComponentName: {}
```

请注意：

- 在 `products/<productName>/` 的 `components.yaml` 或者 `components/*.yaml` 文件中定义的组件，其各个字段的值可以看做是默认值。
- 组件的各个字段的值都可以根据实际的部署环境而变化，比如组件的 `version` 变量随着对应组件的升级而相应更改。
- `<productName>` 目录下所有的 `components.yaml` 或者 `components/*.yaml` 会和 `vars.yaml` 中的变量合并，并保存到对应的环境目录 `targets/<target>/<zone>` 的 `vars.yaml` 文件中。

## Example

假设你的产品名称叫做 `foobar`，该产品由 `foobar-web`，`foobar-api`，`foobar-backend`, `foobar-db` 和 `foobar-cache` 五个组件构成。

你可以把组件都定义到 `components.yaml` 文件中，如下：

```yaml
# all components are defined in components.yaml

foobar-web:
  version: "v0.0.1"
  # ... more other fields

foobar-api:
  version: "v0.0.2"
  # ... more other fields

foobar-backend:
  version: "v0.0.3"
  # ... more other fields

foobar-db:
  version: "v0.0.4"
  # ... more other fields

foobar-cache:
  version: "v0.0.5"
  # ... more other fields
```

你也可以使用 `components` 目录甚至使用嵌套的子目录，并使用不用的文件来定义，如下：

```yml
# components/foobar-web.yaml
foobar-web:
  version: "v0.0.1"
  # ... more other fields

# 一个文件中也可以定义多个组件
# components/api.yaml
foobar-api:
  version: "v0.0.2"
  # ... more other fields
foobar-backend:
  version: "v0.0.3"
  # ... more other fields

# components/storage/db.yaml
foobar-db:
  version: "v0.0.4"
  # ... more other fields

# components/storeage/cache.yaml
foobar-cache:
  version: "v0.0.5"
  # ... more other fields
```

`components.yaml` 或和 `components` 目录下的任何 `.yaml` 文件（文件名没有实际的意义）都可以用来定义产品的组件。
组件的名称是由 `.yaml` 文件的「顶层字段」决定的。

## 组件的实现

`products/<productName>/roles` 目录下存放每一个组件的实际运维代码。

通常情况下，你会为每一个组件在 `products/<productName>/roles` 目录下开发一个对应的组件 Role。

> 但是请注意，components 与 roles 下面的目录并不是一对一的关系，并且 role 的名字也不一定需要和 component 的名字相同。

`sail` 制定的组件 Role 的目录结构以 Ansible 的 Role 目录作为基准，并扩展了其它功能。

```bash
<roleName>/
  # 下面 8 个子目录是 Ansible 的 Role 标准目录结构，请使用这几个目录去开发标准的 Ansible Role
  # 通常情况下，如果组件需要使用常规形式（非容器化）部署的话，你肯定需要编写 ansible 格式的 role
  defaults/
  tasks/
  files/
  templates/
  handlers/
  vars/
  meta/
  library/

  # 存放该组件的 helm 模板
  helm/templates

  # 除了上面几个有特殊用途的目录外，
  # 你可以在 Role 目录下放置任何文件或子目录

  # 比如用于构建组件镜像的 Dockerfile
  # Dockerfile 中的 COPY，ADD 指令可以直接引用 <roleName> 目录内其它文件，如 files/xxx，或者 templates/xxx
  Dockerfile

  # 比如介绍组件的 README.md
  README.md

  # 比如文档，学习资料等
  docs

  # 任何文件或目录
  # any-other-files-or-dirs
```

## 变量

Ansible 和 Helm 的核心操作是根据变量去渲染出配置文件。`sail` 做的最多的一项工作就是针对「环境」的变量的处理。

`sail` 会把下面一些变量传给 Ansible 和 Helm，Ansible 和 Helm 可以使用这些变量去执行渲染动作。

1. 环境特有的变量信息

    `targets/<target>/<zone>/vars.yaml` 文件中的变量，包含所有组件变量以及非组件变量。

2. 几个 `sail` 相关变量

    - `sail_packages_dir`
    - `sail_targets_dir`
    - `sail_products_dir`
    - `sail_target_dir`
    - `sail_zone_dir`
    - `sail_target_name`
    - `sail_zone_name`

除了 `sail` 提供的变量，Ansible 或 Helm 本身也能够识别很多其它变量。

## 组件的部署形式

组件可能以 Pod 中的形式部署到 K8S 平台上，或者以 Server 的形式部署到常规的服务器上。
按照部署形式，我们把组件分别称作「容器组件」和「常规组件」。

Sail 目前通过 Ansible Role 来部署常规组件，通过 Helm Chart 来部署 Pod 组件。

- [开发 Ansible Role](./ansible.md)
- [开发 Helm Chart](./helm.md)
