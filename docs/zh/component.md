# 组件

一个产品是由多个组件组成。你可以使用 `components.yml` 文件或和 `components` 目录来定义（声明）产品的组件。

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

- 在 `products/<productName>/` 的 `components.yml` 或者 `comopnents/*.yml` 文件中定义的组件，其各个字段的值可以看做是默认值。
- 组件的各个字段的值都可以根据实际的部署环境而变化，比如组件的 `version` 变量随着对应组件的升级而相应更改。
- `<productName>` 目录下所有的 `components.yml` 或者 `comopnents/*.yml` 会和 `vars.yml` 中的变量合并，并保存到对应的环境目录 `targets/<target>/<zone>` 的 `vars.yml` 文件中。

## Example

假设你的产品名称叫做 `foobar`，该产品由 `foobar-web`，`foobar-api`，`foobar-backend`, `foobar-db` 和 `foobar-cache` 五个组件构成。

你可以把组件都定义到 `components.yml` 文件中，如下：

```yaml
# all components are defined in components.yml

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
# components/foobar-web.yml
foobar-web:
  version: "v0.0.1"
  # ... more other fields

# 一个文件中也可以定义多个组件
# components/api.yml
foobar-api:
  version: "v0.0.2"
  # ... more other fields
foobar-backend:
  version: "v0.0.3"
  # ... more other fields

# components/storage/db.yml
foobar-db:
  version: "v0.0.4"
  # ... more other fields

# components/storeage/cache.yml
foobar-cache:
  version: "v0.0.5"
  # ... more other fields
```

`components.yml` 或和 `components` 目录下的任何 `.yml` 文件（文件名没有实际的意义）都可以用来定义产品的组件。
组件的名称是由 `.yml` 文件的「顶层字段」决定的。

## Roles 目录

组件的实际运维代码是由组件的 Role 实现的。

通常情况下，你会为每一个组件在 `products/<productName>/roles` 目录下开发一个对应的组件 Role。

`sail` 制定的组件 Role 的目录结构以 Ansible 的 Role 目录作为基准。

```bash
<roleName>/
  # 下面 8 个子目录是 Ansible 的 Role 标准目录结构，请使用这几个目录去开发标准的 Ansible Role
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

`sail` 会把下面一些变量传给 Ansible 和 Helm，Ansible 和 Helm 可以使用这些变量去执行渲染动作。

1. 环境特有的变量信息

    `targets/<target>/<zone>/vars.yml` 文件中的变量，包含所有组件变量以及非组件变量。

2. 三个 `sail` 全局变量

    - `packages_dir`
    - `targets_dir`
    - `products_dir`

3. Ansible 或 Helm 本身能够识别的其它变量
