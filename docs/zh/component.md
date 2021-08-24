# 组件

一个产品是由多个组件组成。你可以使用 `components.yml` 文件或和 `components` 目录来定产品的组件。

## 组件声明

「组件声明」也被称作「组件变量」。

```yaml
# 组件名称
redis:
  # 组件的版本
  version: 3.2.10

  # 是否部署该组件
  enabled: false

  # 该组件是否有外部系统提供
  external: false

  # 该组件对外提供的服务（端口）
  services:
    # key can be any string
    default:
      port: 6379
    sentinel:
      port: 7379
    check:
      port: 6479

  # 组件的变量
  vars:
    database: 6
    pass: ePxoMflue6jhx9oYvYV3
```

基本上，「组件声明」中除了字段名称外的所有字段都是可选的。一个最简单的「组件声明」如下：

```yaml
theComponentName: {}
```

在 `products/<productName>/` 的 `components.yml` 或者 `comopnents/*.yml` 文件中定义的组件变量的值可以看做是默认值。
在一个实际的部署环境中，所有组件变量和 `<productName>` 目录下的 `vars.yml` 中的变量会被合并，并保存到对应的环境的目录的 `vars.yml` 文件中。
所有的「组件变量」都可以根据实际的部署环境而变化，比如组件的 `version` 变量随着对应组件的升级而相应更改。

## Example

假设你的产品名称叫做 `foobar`，该产品由 `foobar-web`，`foobar-api`，`foobar-backend`, `foobar-db` 和 `foobar-cache` 五个组件构成。

你可以把它们都定义到 `components.yml` 文件中，如下：

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

你也可以使用 `components` 目录并用不用的文件来定义，甚至使用嵌套的子目录，如下：

```yml

# components/foobar-web.yml
foobar-web:
  version: "v0.0.1"
  # ... more other fields

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
组件的名称是由 `.yml` 文件的顶层字段决定的。


## Roles 目录

通常情况下，你会为每一个组件在 `products/<productName>/roles` 目录下开发一个对应的组件 Role。

`sail` 要求的组件 Role 的目录结构以 Ansible 的 Role 目录作为基准。

```bash
<roleName>/
  # 下面 8 个子目录是 Ansible 的 Role 标准目录结构，使用这几个目录去开发标准的 Ansible Role
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
  Dockerfile

  # 比如介绍组件的 README.md
  README.md

  # 比如文档，学习资料等
  docs

  # 任何文件或目录
  # any-other-files-or-dirs
```

## Ansible 模板文件可以读到的变量

`sail` 会把下面一些变量传给 Ansible，Ansible Roles 中的模板文件可以使用这些变量去渲染。

1. 环境特有的变量信息

    `targets/<target>/<zone>/vars.yml` 文件中的变量，包含一个产品的所有组件变量以及非组件变量。

2. 三个 `sail` 全局变量

    - `packages_dir`
    - `targets_dir`
    - `products_dir`

3. Ansible Role 本身能够识别的其它变量
