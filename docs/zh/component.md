# 组件

## 组件声明

「组件声明」也被称作「组件变量」。

在 `products/<productName>/` 的 `components.yml` 或者 `comopnents/*.yml` 文件中定义的组件变量的值可以看做是默认值。

在一个实际的部署环境中，所有组件变量和 `<productName>` 目录下的 `vars.yml` 中的变量会被合并，并保存到对应的环境的目录的 `vars.yml` 文件中。

所有的「组件变量」都可以根据实际的部署环境而变化，比如组件的 `version` 变量随着对应组件的升级而相应更改。

```yaml
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

## Roles 目录

你应该为每一个组件在 `products/<productName>/roles` 目录下开发一个对应的组件 Role。

组件 Role 的目录结构以 Ansible 的 Role 目录作为基准。

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

  # 比如用于构建组件镜像的Dockerfile
  Dockerfile

  # 比如组件介绍的 README.md
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
