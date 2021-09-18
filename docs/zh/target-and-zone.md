# target and zone

Sail 有三个主要概念，Products （产品），Targets (环境)，Packages （包）。

通过 Sail 来进行运维操作，第一步就是创建出一个环境目标，并指定在这个目标中部署的是什么产品。

```bash
$ sail conf-create -t <target_name> -z <zone_name> -p <product_name> [--hosts ...]
```

Sail 使用了两个层级来表示环境目标。

第一个层级时 Target，第二个层级是 Zone。Target 是完全独立的。Target 下面可以创建多个 Zone(s)。

`sail` 命令在运行时必须同时通过 `-t` 和 `-z` 指明 target 和 zone，从而指定一个具体的操作目标。`sail` 在一个命令行中只能操作一个具体的 Zone。

虽然 `sail` 一次只能操作一个 Zone，但是 `sail` 会把该 Zone 所属的 Target 下面的所有 Zones 的变量整合起来传递给正在操作中的 Zone。这也是 Sail 设计两个层级的目的所在。

`sail` 在执行 `ansible-playbook` 或 `helm` 命令时，会把 `<target_name>/<zone_name>/vars.yaml`
和 `<target_name>/<zone_name>/_computed.yaml` 两个文件传递过去。

重点看一下 `_computed.yaml` 文件中的变量。

- `inventory` 的值是当前 Zone 的 `hosts.yaml` 的文件内容。
-  `platforms` 的值是当前 Zone 的 `platforms.yaml` 的文件内容。

> 一个 Zone 「自身变量」则是由 `vars.yaml` 中的变量，以及 `_computed.yaml` 中 `inventory` 和 `platforms` 变量组成。

`_computed.yaml` 中的 `targetvars` 的值则保存了当前 Zone 所在 Target 下面所有 Zones 的「自身变量」。

在某些部署场景下，通过 `targetvars` 变量可以访问到 Target 下所有 Zones 的信息。

```yaml
inventory:
  ...
platforms:
  ...


targetvars:
  zones:
    zone1:
      ...
    zone2:
      ...
```

## `targetvars` 的使用场景举例

Todo
