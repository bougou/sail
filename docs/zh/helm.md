# Helm

在部署一个产品时，产品的所有组件或部分组件可能以 Pod 形式部署。
Sail 目前支持以 Helm Chart 的形式部署 Pod 组件。
Sail 支持两种模式来使用 Helm Chart。

1. 将整个产品作为一个 Helm Chart 来开发。 `_sail_helm_mode: "product"`
2. 将产品中的每一个组件作为独立的 Helm Chart 来开发。`_sail_helm_mode: "component"`

> 注：两个模式不能同时使用。

在 `<target_name>/<zone_name>/vars.yaml` 文件中使用 `_sail_helm_mode` 来指定模式。


## 将整个产品作为一个 Helm Chart 来开发

将整个产品作为一个 Helm Chart 来开发，你可以直接把 `products/<productName>` 目录当做 Helm Chart 的目录来使用。

如下：

```bash
# 不管使用 Ansible 还是 Helm 来部署，
# 都必须要在 components.yaml 和 vars.yaml 中声明该产品的组件和变量
components.yaml # sail 组件声明
vars.yaml       # sail 非组件变量

Chart.yaml      # helm Chart.yaml
values.yaml     # helm values.yaml
templates/      # helm templates
crds/           # helm crds
charts/         # helm charts
```

可以把为产品组件开发的 K8S manifest 模板文件放到上面的 `templates` 目录下。如：

```bash
templates/foobar-api-service.yaml
templates/foobar-api-deployment.yaml
templates/foobar-api-configmap.yaml
templates/foobar-web-service.yaml
templates/foobar-web-deployment.yaml
templates/foobar-web-configmap.yaml
templates/...
```

## 将产品中的每一个组件作为独立的 Helm Chart 来开发

将每一个组件作为独立的 Helm Chart 来开发，需要在组件的 role 目录下 `products/<productName>/roles/<roleName>` 存放对应组件的 Helm Chart。

组件的 Chart 目录：`products/<productName>/roles/<roleName>/helm/<componentName>/`

```bash
Chart.yaml
values.yaml
templates/
crds/
charts/
```

每一个组件作为独立的 Chart 时，可以定义独立的 `values.yaml`。

另外，即使在 `_sail_helm_mode: "component"` 模式下，你依然可以定义一个全局的 `products/<productName>/values.yaml` 文件 （可选的）。

## Helm 运行时

Sail 在使用 Helm 部署 Pod 组件时，会在环境的 `<target_name>/<zone_name>/helm` 目录中构造出最终的 Chart 目录来执行 `helm` 命令。

### 将整个产品作为一个 Helm Chart 来运行

`_sail_helm_mode: "product"` 模式

Sail 会以 `<productName>` 构造出一个 `<target_name>/<zone_name>/helm/<productName>` 目录作为标准的 Helm Chart 目录。

然后把产品代码中的 `templates/`,  `crds`, `charts/`，`Chart.yaml` 文件拷贝到 `<target_name>/<zone_name>/helm/<productName>` Chart 目录下。

然后把产品代码中的 `values.yaml` 与 `<target_name>/<zone_name>/values.yaml` 进行合并并保存（如果有的话）。

然后在 `<target_name>/<zone_name>/helm/<productName>` Chart 目录下创建一个 `resources` 软链文件，链接到 `<target_name>/<zone_name>/resources` 目录上。

> Sail 使用 `<target_name>/<zone_name>/resources` 目录来统一管理部署时的其它资源文件，如证书，秘钥等等。

在执行 `helm` 命令时，Sail 会把 Zone 目录下的以下几个文件作为 values 文件以 `--values <valuesFile>` 参数形式「依次」传给 `helm` 命令。

- `<target_name>/<zone_name>/vars.yaml`
- `<target_name>/<zone_name>/_computed.yaml`
- `<target_name>/<zone_name>/values.yaml`

> 1. values 文件传递顺序很重要。`helm` 命令会合并变量，后边的覆盖前面的。
> 2. `<target_name>/<zone_name>/values.yaml` 中的变量值可以根据实际环境进行修改，并且会被持久化。

### 将产品中的每一个组件作为独立的 Helm Chart 来运行

`_sail_helm_mode: "component"` 模式

首先处理全局 `values.yaml` 文件，把产品代码中的 `values.yaml` 与 `<target_name>/<zone_name>/values.yaml` 进行合并并保存（如果有的话）。

然后处理每一个需要部署的 Pod 组件的 Helm Chart 目录。每一个组件都会按照下面步骤操作。

1. Sail 会以组件的名字 `<componentName>` 构造出一个 `<target_name>/<zone_name>/helm/<componentName>` 目录作为该组件的 Helm Chart 目录。
2. 然后把组件代码中的 `templates/`,  `crds/`, `charts/`，`Chart.yaml` 拷贝到 `<target_name>/<zone_name>/helm/<componentName>` Chart 目录下。
3. 然后把组件代码中的 `values.yaml` 与 `<target_name>/<zone_name>/helm/<componentName>/values.yaml` 进行合并并保存。
4. 然后在 `<target_name>/<zone_name>/helm/<componentName>` Chart 目录下创建一个 `resources` 软链文件，链接到 `<target_name>/<zone_name>/resources` 目录上。

在执行 `helm` 命令时，Sail 为会每一个 Pod 组件分别执行 `helm` 命令。

Sail 会把 Zone 目录下的以下几个文件作为 values 文件以 `--values <valuesFile>` 参数形式「依次」传给 `helm` 命令。

- `<target_name>/<zone_name>/vars.yaml`
- `<target_name>/<zone_name>/_computed.yaml`
- `<target_name>/<zone_name>/values.yaml`
- `<target_name>/<zone_name>/helm/<componentName>/values.yaml`

> 1. values 文件传递顺序很重要。`helm` 命令会合并变量，后边的覆盖前面的。
> 2. `<target_name>/<zone_name>/values.yaml` 变量值可以根据实际环境进行修改，并且会被持久化。
> 3. `<target_name>/<zone_name>/helm/<componentName>/values.yaml` 变量值可以根据实际环境进行修改，并且会被持久化。
