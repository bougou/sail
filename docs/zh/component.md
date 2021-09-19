# 组件

一个产品是由多个组件组成。你可以使用 `components.yaml` 文件或和 `components` 目录来定义（声明）产品的组件。

## 什么是组件

- 一个组件实现了一个产品的某个功能，或提供产品运行时需要的某个能力。
- 一个组件可以以多个副本运行（比如组成多节点集群，或在 K8S 中运行多个 Pods）
- 一个组件的任何「一个副本」只能部署到一个「计算单元」中，即一个服务器中或一个 Pod 中。

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
    pass: 123456

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

- 在 `products/<productName>/` 的 `components.yaml` 或者 `components/*.yaml` 文件中定义的组件，各个字段的值可以看做是默认值。
- 组件的各字段的值都可以根据实际的部署环境而变化，比如组件的 `version` 变量随着对应组件的升级会相应更改。
- `<productName>` 目录下所有的 `components.yaml` 或者 `components/*.yaml` 会和 `vars.yaml` 中的变量合并，并保存到对应的环境目录 `targets/<target_name>/<zone_name>` 的 `vars.yaml` 文件中。

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

# components/storage/cache.yaml
foobar-cache:
  version: "v0.0.5"
  # ... more other fields
```

`components.yaml` 或和 `components` 目录下的任何 `.yaml` 文件（文件名没有实际的意义）都可以用来定义产品的组件。
组件的名称是由 `.yaml` 文件的「顶层字段」决定的。

在给组件命名时，请确保组件名称不要和 `vars.yaml` 中的顶层变量名冲突。

## 组件的实现

`products/<productName>/roles` 目录下存放每一个组件的实际运维代码。

通常情况下，你会为每一个组件在 `products/<productName>/roles` 目录下开发一个对应的组件 Role。

> 但是请注意，components 与 roles 下面的目录并不是一对一的关系，并且 role 的名字也不一定需要和 component 的名字相同。

Sail 制定的组件 Role 的目录结构是以 Ansible 的 Role 目录作为基准，并扩展了其它功能。

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

  # 存放该组件 Helm Chart
  helm/

  # 除了上面几个有特殊用途的目录外，
  # 你可以在 Role 目录下放置任何文件或子目录

  # 比如用于构建组件镜像的 Dockerfile
  # Dockerfile 中的 COPY，ADD 指令可以直接引用 <roleName> 目录内其它文件，如 files/xxx，或者 templates/xxx
  Dockerfile

  # 比如文档，学习资料等
  README.md
  docs

  # 任何文件或目录
  # any-other-files-or-dirs
```

## 组件的声明详解

### 组件的部署形式 Form

组件可能以 Pod 中的形式部署到 K8S 平台上，或者以 Server 的形式部署到常规的服务器上。
按照部署形式，我们把组件分别称作「容器组件」和「常规组件」。


```yaml
componentName:
  # "server" 或 "pod", 如果为空，默认为 "server"
  form: "server"
```


Sail 目前通过 Ansible Role 来部署常规组件，通过 Helm Chart 来部署 Pod 组件。

- [开发 Ansible Role](./ansible.md)
- [开发 Helm Chart](./helm.md)

### 组件 Enabled/External

- `enabled` 表示在该环境中是否部署该组件。
- `external` 表示在该环境中，该组件的能力是否是由外部系统或平台提供的。

比如你的产品中用到了 MySQL，在产品的组件声明中你必须声明该组件，假设你把该组件命名为 "mysql" （可以是任意名字）。

在实际部署该产品时，不同的环境对于 "mysql" 组件就有不同的选择。

- 公有云环境可能直接购买云产商提供的数据库实例，如 AWS RDS
- 有些环境中可能选择直接使用几台服务器部署出一套 MySQL 集群
- 有些环境中，部署运维该产品的团队可能直接使用由另外的专门的数据库团队在维护并对外提供 MySQL 实例。

当你决定使用外部系统提供的能力而非自己去部署时，就把该组件的 `external` 设置为 `true`。

当你决定自己去部署该组件时，请把 `enabled` 设置为 `true`。

在有些环境中，产品的某些组件甚至不用部署，请把 `external` 和 `enabled` 都设置为 `false`。

`enabled` 和 `external` 有三种使用组合：

```yaml
# case1: 自己部署该组件
mysql:
  enabled: true
  external: false

# case2: 使用外部系统提供的能力
# 当同时配置为 true 时，Sail 自动把 enabled 改成 false
mysql:
  enabled: false
  external: true

# case3: 禁用该组件（不部署）
mysql:
  enabled: false
  external: false
```


### 组件 Services

- 一个组件可以向其它组件提供 0 个，1 个或多个服务。
- 一个服务对应一个端口。

如果一个组件对外暴露了多个端口，那么必须定义为多个 Service。
端口名（或服务名）可以是任意有意义的 alphanum 字符串，比如 "web", "http", "rpc", "cluster", "default" 等等。

`services` 字段是一个由端口名（或者称为服务名）作为键的字典结构。

```yaml
<componentName>:
  services:
    <portName>:  # 我们把下面的对象称为 Service
      description: ""
      scheme: ""
      host: ""
      ipv4: ""
      ipv6: ""
      port: 80
      addr: ""
      endpoints: ""
      urls: ""
      lbPort: 10080
      pubPort: 80
```

Component 中的 Service(s) 用于向其它组件声明应该如何访问本组件的服务。
Service 的字段该如何配置，通常取决于是否使用了外部组件。

#### 使用外部组件 `external: true`

当使用了外部组件时，Service 的字段完全用来表明其它组件该如何访问这个组件的这个服务。

比如，你从云服务商购买了一个 MySQL 实例，你可以直接把实例的访问地址配置到该组件的 Service 下。

```yaml
mysql:
  external: true
  services:
    default:
      port: 3306
      host: "some-addr"
  # 在组件的 vars 下面配置任意其它的变量
  vars:
    dbuser: "someuser"
    dbpass: "somepass"
    authuser: "root"
    authpass: "rootpass"
```

如果服务的访问地址是多个，请使用 Service 的 `endpoints` 字段或者 `urls`。

```yaml
elasticsearch:
  external: true
  services:
    http:
      endpoints:
        - "192.168.1.10:9200"
        - "192.168.1.11:9200"
        - "192.168.1.12:9200"
```

#### `external: false`

当不使用外部组件时，通常只需要在 Service 中中定义 `port`。port 通常用于在安装配置组件时配置进程的监听端口。

```yaml
mysql:
  enabled: true
  external: false
  services:
    default:
      port: 3306
  vars:
    dbuser: "someuser"
    dbpass: "somepass"
    authuser: "root"
    authpass: "rootpass"

elasticsearch:
  enabled: true
  external: false
  services:
    http:
      port: 9200
    cluster:
      port: 9300
```

只配置 `port`， 不只配 `host` 或其它地段，这个组件该被如何访问呢？

在准备一个环境时，肯定都会为这些需要部署的组件（`enabled: true`）指定主机信息。组件的主机信息维护在主机清单中 `<target_name>/<zone_name>/hosts.yaml`。

所以，对于 `enabled: true` 的组件，不需要在组件的 Service 下面配置 `host` 信息。

并且不管是否使用了外部组件，Sail 都不建议你在实际的安装脚本中直接应用的 `services` 下的字段来作为连接地址，建议你使用 `computed` 下的字段。

### 组件 Computed

Computed 与 Services 相关。
Computed 的键名和 Services 的键名一一对应。

Sail 会对组件的 Services 下声明的所有 Service 进行计算，将计算后的 ServiceComputed 设置给 Computed 对应的服务。

永远不要直接修改 Computed 下面的字段。

根据组件 `enabled` 和 `external` 的不同配置，Sail 对服务的计算逻辑也不一样。

#### 当 `external: true`

当使用外部组件时，Computed 的计算逻辑比较简单，如下：

```yaml
mysql:
  enabled: false
  external: true
  services:
    default:  # 我们把下面的对象称为 Service
      scheme: ""
      host: "ipaddr-or-hostname"
      port: 3306
      addr: ""
      endpoints: []
      urls: []

  computed:
    default:
      # 如果 Service.scheme 非空，则取值 Service.scheme, 否则设置为 "tcp"
      scheme: "tcp"

      # 如果 Service.host 非空，则取值 Service.host, 否则设置为 "127.0.0.1"
      host: "ipaddr-or-hostname"

      # 直接取值 Service.port
      port: 3306

      # 如果 Service.addr 非空，则取值 Service.addr, 否则设置为 "<host>:<port>"
      # <host> 和 <port> 为上面计算出来的 host 和 port
      addr: "ipaddr-or-hostname:3306"

      # 如果 Service.path 非空且以 "/" 开头，则取值 Service.path
      # 如果 Service.path 非空且不以 "/" 开头，则取值 "/"+Service.path
      # 否则设置为 "/"
      path: "/"

      # 否则把上面计算出来的 host 作为唯一元素加到 hosts 中
      hosts:
        - "ipaddr-or-hostname"

      # 如果 Service.addrs 非空，则取值 Service.addrs
      # 否则把上面计算出来的 addr 作为唯一元素加到 endpoints 中
      addrs:
        - "ipaddr-or-hostname:3306"

      # 如果 Service.endpoints 非空，则取值 Service.endpoints
      # 否则把上面计算出来的 addrs 中的每个元素加上 scheme 前缀 "<scheme>://<addr>" 加入 urls 中
      endpoints:
        - "tcp://ipaddr-or-hostname:3306"

      # 如果 Service.urls 非空，则取值 Service.urls
      # 否则把上面计算出来的 endpoints 中的每个元素加上 path 后缀 "<endpoint><path>" 加入 urls 中
      urls:
        - "tcp://ipaddr-or-hostname:3306/"
```

你应该在 Ansible 和 Helm 中引用组件的 computed 下的字段用于访问相关的服务。

比如在 Ansible 的 Jinja2 模板中，你可以使用下面的格式：

```j2
db_host: "{{ mysql['computed']['default']['host'] }}"
db_port: {{ mysql['computed']['default']['port'] }}
```

#### 当 `external: false`

当不使用外部组件时，自己部署组件时，Computed 的计算逻辑稍微复杂一点，如下：

```yaml
mysql:
  enabled: true
  external: false
  services:
    default:  # 我们把下面的对象称为 Service
      port: 3306

  computed:
    default:
      # 如果 Service.scheme 非空，则取值 Service.scheme, 否则设置为 "tcp"
      scheme: "tcp"

      # 如果 Service.host 非空，则取值 Service.host （自己部署组件时，Service.host 通常都应该为空）
      # 否则去查询主机清单 "targets/<target_name>/<zone_name>/hosts.yaml"
      # 如果主机清单中存在该组件，并且对应的主机列表不为空，则从主机列表中取出第一条作为 host 值 (主机列表取自 hosts 字典的键值列表，所以列表顺序具有不确定性)
      # 所有其它情况，设置为 "127.0.0.1"。
      host: <computed>

      # 如果 Service.pubPort 非空（不为 0)，则取值 Service.pubPort
      # 如果 Service.lbPort 非空（不为 0)，则取值 Service.lbPort
      # 所有其它情况，直接取值 Service.port
      port: 3306

      # 如果 Service.addr 非空，则取值 Service.addr, 否则设置为 "<host>:<port>"
      # <host> 和 <port> 为上面计算出来的 host 和 port
      addr: <computed>

      # 如果 Service.path 非空且以 "/" 开头，则取值 Service.path
      # 如果 Service.path 非空且不以 "/" 开头，则取值 "/"+Service.path
      # 否则设置为 "/"
      path: "/"

      # 去查询主机清单 "targets/<target_name>/<zone_name>/hosts.yaml"
      # 如果主机清单中存在该组件，则取值为主机清单中该组件的主机列表
      # 所有其它情况，把上面计算出来的 host 作为唯一元素加到 hosts 中
      hosts:
        - <computed>

      # 如果 Service.addrs 非空，则取值 Service.addrs
      # 否则把上面计算出来的 hosts 中的每个元素加上 port 后缀 "<host>:<port>" 加入 addrs 中
      addrs:
        - <computed>

      # 如果 Service.endpoints 非空，则取值 Service.endpoints
      # 否则把上面计算出来的 addrs 中的每个元素加上 scheme 前缀 "<scheme>://<addr>" 加入 endpoints 中
      endpoints:
        - <computed>

      # 如果 Service.urls 非空，则取值 Service.urls
      # 否则把上面计算出来的 endpoints 中的每个元素加上 path 后缀 "<endpoint>:<path>" 加入 urls 中
      urls:
        - <computed>
```

假设，产品中需要部署 elasticsearch 组件，在主机清单中为 elasticsearch 分配了 3 台机器。

```yaml
# <target_name>/<zone_name>/hosts.yaml
elasticsearch:
  hosts:
    192.168.1.10: {}
    192.168.1.11: {}
    192.168.1.12: {}
  vars: {}
  children: {}
```

该组件配置为不使用外部组件。

```yaml
# <target_name>/<zone_name>/vars.yaml
elasticsearch:
  enabled: true
  external: false
  services:
    default:
      scheme: "http"
      port: 9200
    cluster:
      scheme: "http"
      port: 9300
```

那么自动计算出的 Computed 字段如下：

```yaml
elasticsearch:
  computed:
    cluster:
      scheme: tcp
      host: 192.168.1.10
      port: 9300
      addr: 192.168.1.10:9300
      path: /
      hosts:
        - 192.168.1.10
        - 192.168.1.11
        - 192.168.1.12
      addrs:
        - 192.168.1.10:9300
        - 192.168.1.11:9300
        - 192.168.1.12:9300
      endpoints:
        - tcp://192.168.1.10:9300
        - tcp://192.168.1.11:9300
        - tcp://192.168.1.12:9300
      urls:
        - tcp://192.168.1.10:9300/
        - tcp://192.168.1.11:9300/
        - tcp://192.168.1.12:9300/
    default:
      scheme: tcp
      host: 192.168.1.11
      port: 9200
      addr: 192.168.1.11:9200
      path: /
      hosts:
        - 192.168.1.10
        - 192.168.1.11
        - 192.168.1.12
      addrs:
        - 192.168.1.10:9200
        - 192.168.1.11:9200
        - 192.168.1.12:9200
      endpoints:
        - tcp://192.168.1.10:9200
        - tcp://192.168.1.11:9200
        - tcp://192.168.1.12:9200
      urls:
        - tcp://192.168.1.10:9200/
        - tcp://192.168.1.11:9200/
        - tcp://192.168.1.12:9200/

```
