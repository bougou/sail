# sail 命令

## 变量

Ansible 和 Helm 的核心操作是根据变量去渲染出配置文件。`sail` 做的最多的一项工作就是针对「环境」的变量的处理。

`sail` 会把下面一些变量传给 Ansible 和 Helm，Ansible 和 Helm 可以使用这些变量去执行渲染动作。

1. 环境特有的变量信息

    - `targets/<target_name>/<zone_name>/vars.yaml` 文件中的变量，包含所有组件变量以及非组件变量。
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
