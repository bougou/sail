# Components

A **Product** is composed of **Component(s)**.
You can use  `components.yaml` file and `components` directory to declare the components of the product.

## Declare Components

```yaml
# component name
redis:
  # component version
  version: 3.2.10

  # the component is deployed or not
  enabled: false

  # the component is provided by external system or platform
  external: false

  # the exposed services (ports) by the component
  services:
    # each key represents a service name, it can be any alphanum string
    default:
      # listend port, or exposed port
      port: 6379
      # other fields
    sentinel:
      port: 7379
    check:
      port: 6479

  # any component related variables
  vars:
    database: 6
    pass: 123456

  # any component related tags
  tags: {}

```

Basically, all fields are optional except the component name.
So, the most simplest component declaration is:

```yaml
theComponentName: {}
```

Note:

- the values of the fields of the component declaration in the `components.yaml` and `components/*.yaml` can be seen as default values.
- the values of the fields of the component may be changed according to the deploying target environments。
- the variables in the  `components.yaml` and `components/*.yaml` will be merged with `vars.yaml` and stored to `targets/<target_name>/<zone_name>/vars.ayml` for a specific environment.

## Example

Suppose you manage a product called `foobar` which is composed of five components, `foobar-web`，`foobar-api`，`foobar-backend`, `foobar-db` and `foobar-cache`.

You can delcare them all in `components.yaml` file:

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

Or you can use `components` directory and even nested directories to define the components:

```yml
# components/foobar-web.yaml
foobar-web:
  version: "v0.0.1"
  # ... more other fields

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

The `components.yaml` file or any files with suffix `.yaml` under `components` directory can be used to declare components.
The files under `components` has no actual meanings.

The names of components are infered from the top-level keys in yaml files.

When naming components, you have to make sure they are not conflict with the variables defined in `vars.yaml`.

## Component Implementation

The component implementations are the real Ansible and/or Helm code for the component.

The `products/<productName>/roles` directory holds the implementations for each product. Generally, you will develop a corresponding role for each component of the product under `products/<productName>/roles` directory.

Sail's role directory of component is based on [Ansible Role](https://docs.ansible.com/ansible/latest/user_guide/playbooks_reuse_roles.html#role-directory-structure), but has some extensions.

```bash
<roleName>/
  # the following eight directories are standard ansible role structure
  # Generally, you have to define an ansible role for the component that
  # needs to be deployed to normal servers (non-k8s pods).
  defaults/
  tasks/
  files/
  templates/
  handlers/
  vars/
  meta/
  library/

  # store helm chart for the component.
  helm/

  # except above directories, you can put any files or dirs under the role dir.

  # like Dockerfile for building docker image
  # the COPY，ADD directives can directly reference files/xxx or templates/xxx, thus you can put all data about the component together
  Dockerfile

  # like documents or learning materials
  README.md
  docs

  # any-other-files-or-dirs
```

## Explanation for declaring components

### Component deploy formation

Component(s) can be deployed to normal servers as or deployed to K8S platforms.
According to deployment mechanism, we can call the components as "container/pod component" or "normal/server component".

```yaml
componentName:
  # "server" or "pod", defaults to "server" if empty
  form: "server"
```

Sail currently supports using Ansible Role to deploy normal components and using Helm Chart to deploy pod components.

- [Develop Ansible Role](./ansible.md)
- [Develop Helm Chart](./helm.md)

### Component Enabled/External

- `enabled` means whether to deploy the component in the target environment.
- `external` means whether the component is provided by external system or platform.

Suppose you product uses MySQL, so you have to declare one component for MySQL. Let's say you named this component as "mysql".

When deploying the product, there are different choices for "mysql" component as to the different target envionments.

- You can directly use mysql instance purchased from cloud provider on public cloud environment, like AWS RDS.
- You may want to deploy a mysql cluster using some plain servers.
- You may want to use the mysql instance served and maintained by other specical DBA teams from your company.

Anyway, when you want to use external services for the component, you have to set `external` to `true` in the `vars.yaml` of the environment `targets/<target_name>/<zone_name>/vars.yaml`.

When you want to deploy the component yourself, you have to set `enabled` to `true`。

For some cases, the component may even not used, you have to set both `external` and `enabled` to `false`。

In short, there are 3 combination cases to `enabled` and `external`:

```yaml
# case1: deploy the component yourself
mysql:
  enabled: true
  external: false

# case2: use external services
# Sail will automatically set enabled to false if both are set to true
mysql:
  enabled: false
  external: true

# case3: the component are totally disabled
mysql:
  enabled: false
  external: false
```

### Component Services

- one component can provides 0, 1, or multiple services to other components
- one service for one port

`services` field is a dict with port name (service name) as keys.

If the component exposed multiple ports, it must be defined as multiple services.
The port name (or called service name) can be any meaningful  alphanum string,
like "web", "http", "rpc", "cluster", "default" etc.

```yaml
<componentName>:
  services:
    <portName>:  # we call the following object as Service
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

How to configure the fields of Service are largely determined by whether `external` is true or false.

#### Use external component `external: true`

When using external component, the fields of Service are totally indicates
how to access the service for other components.

For example, if you buy one MySQL instance from a public cloud provider,
you can directly set mysql host and port to the fields of Service.

```yaml
mysql:
  external: true
  services:
    default:
      port: 3306
      host: "some-addr"
  # other mysql informations need to set under vars
  vars:
    dbuser: "someuser"
    dbpass: "somepass"
    authuser: "root"
    authpass: "rootpass"
```

If there are multiple accessing endpoints, you can set them by using the `endpoints` and/or `urls` fields.

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

#### Not using external component `external: false`

When not using external component, usually only `port` field needs to be set.
The `port` here generally used to configure the listened port of the program.

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

How to access this component if just `port` is configured?

As for those deployed components (`enabled: true`), their hosts information are kept in the hosts inventory file `<target_name>/<zone_name>/hosts.yaml`。
So, you don't have to configure the `host` field.

And no matther `external` is `true` or `false`, Sail does not recommend to directly reference the fields
under `services` in your actual ansible or helm code. Sail recommends to use `computed`.

### Component Computed

The keys of Computed is one-to-one for Services.

Sail will apply a computation prodedure for each service under `services`,
and set the computed object to `computed`.

Never directly to edit or change the values under `computed`.

The computation procedure is different according to value of `external`.

#### `external: true`

The computation is relatively simple:

```yaml
mysql:
  enabled: false
  external: true
  services:
    default:  # we call the following object as Service
      scheme: ""
      host: "ipaddr-or-hostname"
      port: 3306
      addr: ""
      endpoints: []
      urls: []

  computed:
    default:
      # set to Service.scheme if Service.scheme is not empty, or else set to "tcp"
      scheme: "tcp"

      # set to Service.host if Service.host is not empty, or else set to "127.0.0.1"
      host: "ipaddr-or-hostname"

      # set to Service.port
      port: 3306

      # set to Service.addr if Service.addr is not empty, or else set to "<host>:<port>"
      # the <host> and <port> are the above computed host and port
      addr: "ipaddr-or-hostname:3306"

      # set to Service.endpoints if Service.endpoints is not empty
      # or else use the above computed addr as the only element for endpoints
      endpoints:
        - "ipaddr-or-hostname:3306"

      # set to Service.urls if Service.urls is not empty
      # or else use set to the above computed endpoints with each item
      # prefixed by the above computed scheme ("<scheme>:<endpoint>")
      urls:
        - "tcp://ipaddr-or-hostname:3306"
```

You should always use fields under `computed` to access the services exposed by components.

Like, in Ansible Jinja2 templates:

```j2
db_host: "{{ mysql['computed']['default']['host'] }}"
db_port: {{ mysql['computed']['default']['port'] }}
```

#### `external: false`

The computation procedure is a little complicated when `external: false`:

```yaml
mysql:
  enabled: true
  external: false
  services:
    default:  # We call the following object as Service
      port: 3306

  computed:
    default:
      # set to Service.scheme if Service.scheme not empty, or else set to "tcp"
      scheme: "tcp"

      # set to Service.host if Service.host is not empty
      # or else to query hosts inventory "targets/<target_name>/<zone_name>/hosts.yaml"
      # if the component exists in the inventory and the hosts for the component is not emtpy,
      # then choose one host from hosts (the hosts list are fetched from dict keys, so the order is uncertainty)
      # for all other cases, set to 127.0.0.1
      host: <computed>

      # set to Service.pubPort if Service.pubPort is not zero
      # set to Service.lbPort if Service.lbPort is not zero
      # for all other cases, set to Service.port
      port: 3306

      # set to Service.addr if Service.addr is not empty, or else set to "<host>:<port>"
      # the <host> and <port> are the above computed host 和 port
      addr: <computed>

      # set to Service.endpoints if Service.endpoints not empty
      # or else to query hosts inventory "targets/<target_name>/<zone_name>/hosts.yaml"
      # if the component exists in the inventory and the hosts for the component is not emtpy,
      # then set to all hosts with each host suffixed with port
      # For all other cases, use the addr as the only element for endpoints.
      endpoints:
        - <computed>

      # set to Service.urls if Service.urls not empty
      # or else use set to the above computed endpoints with each item
      # prefixed by the above computed scheme ("<scheme>:<endpoint>")      urls:
        - <computed>
```

Suppose, you product needs the component elasticsearch, you assiged three servers for elasticsearch in hosts inventory.

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

If the elasticsearch component is configured not to using external component.

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

Then, the `computed` will be:

```yaml
elasticsearch:
  computed:
    default:
      scheme: http
      host: 192.168.1.10
      port: 9200
      addr: 192.168.1.10:9200
      endpoints:
        - 192.168.1.10:9200
        - 192.168.1.11:9200
        - 192.168.1.12:9200
      urls:
        - http://192.168.1.10:9200
        - http://192.168.1.11:9200
        - http://192.168.1.12:9200
    cluster:
      scheme: http
      host: 192.168.1.11
      port: 9300
      addr: 192.168.1.11:9300
      endpoints:
        - 192.168.1.12:9300
        - 192.168.1.11:9300
        - 192.168.1.10:9300
      urls:
        - http://192.168.1.12:9300
        - http://192.168.1.11:9300
        - http://192.168.1.10:9300
```
