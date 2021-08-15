# product

products 目录下的每一个 `<product-name>` 目录就是该产品的「部署代码」。里面封装了该产品的所有部署及运维逻辑。

理想情况下，products 下面的内容应该维护在一个或一些独立的仓库中。
当产品的部署逻辑需要更新时，你都应该在这些独立的仓库中去修改。

## Product 的目录结构

`products/<product-name>` 目录下的结构：

```bash
## Ansible 相关
roles/              # Ansible Roles
vars.yml            # 默认变量
components.yml      # 定义组件
components/         # 定义组件
sail.yml            # Ansible 必须存在的一个 Ansible Playbook 文件
<playbook>.yml      # 其它的 Playbook 文件

## Helm 相关
templates/          # Helm templates
Chart.yml           # Helm Chart definition
values.yml          # Helm values 变量定义

## 通用
resources/          # 其它默认资源，如默认的证书文件、试用License 文件、默认的 icon 图标等

## 其它
README.md           # 介绍文档
```

## 定义产品的组件

你可以使用 `components.yml` 或和 `components` 目录来定产品的组件。

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
