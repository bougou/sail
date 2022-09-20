# sail commands

```yaml
$ cat ~/.sailrc.yaml
products-dir: /path/to/products
targets-dir: /path/to/targets
packages-dir: /path/to/packages
```

Most `sail` subcommands require explicitly specifying the environment you want to operate by `--target <targetname>` and `--zone <zonename>`, which corresponds to the `<targets-dir>/<targetname>/<zonename>`.

Sometimes, if you ONLY have one target and one zone to operate, specify `--target` and `--zone` options would seems tedious each time you run `sail`.

You can configure `default-target` and `default-zone` variables in `~/.sairc.yaml` in advance to shorten the typing.
And explicitly specified `--target` and `--zone` arguments would override them.

```yaml
products-dir: /path/to/products
targets-dir: /path/to/targets
packages-dir: /path/to/packages

default-target: targetname
default-zone: zonename
```

## sail list-components

List all components for the specified product.

```bash
$ ./sail list-components -p <productName>
the product <productName> contains (<number>) components:
- <componentName1>
- <componentName2>
- ...
```

## sail conf-create

Create a new deploy target environments.
You have to specify a target name and a zone name, and the name of the product which will be deployed in this environment.

The product name is persisted under zone's variables. So you won't need to specify the product when running other `sail` commands.
When running other `sail` commands, ONLY the envionment (`-t <targetName> -z <zoneName>`) may need to be specified.

```bash
$ sail conf-create -t <targetName> -z <zoneName> \
    -p <productName> \
    --hosts ip1,ip2,ip3 \                   # will be used for all components which are not explicityly specified by --hosts options
    --hosts componentName1/ip1,ip2,ip3 \
    --hosts componentName2/ip11,ip12,ip13
```

## sail conf-update

Syncs, updates, and computes the vars for the zone and dumps them into files.
Thus `ansible-playbook` and/or `helm` can load these variables files from the command line.

Actually, all the things did by `sail conf-update` are also automatically executed in other `sail` commands,
like `sail apply` and `sail upgrade`.
So there is no need to manually run `sail conf-update` by yourself before you run other `sail` commands.

Only on the cases that you want to check and see what the computed variables looks like,
you can run `sail conf-update` manually.

```bash
$ sail conf-update -t <targetName> -z <zoneName>
```

### add or remove hosts for specific components

```bash
# enable k8s-node component and add a node for k8s-node component
$ sail conf-update -t k8s -c k8s-node --hosts k8s-node/192.168.2.204

# add two nodes for k8s-node component
# note: use '+' sign before component name
$ sail conf-update -t k8s --hosts +k8s-node/192.168.2.205,192.168.2.206
```

## sail apply

`sail apply` will execute `ansible-playbook` for the server components, and execute `helm` for the pod components.

```bash
$ sail apply -t <targetName> -z <zoneName> [-c <componentName>]
```

## sail upgrade

`sail upgrade` will execute `ansible-playbook` for the server components, and execute `helm` for the pod components.

```bash
$ sail upgrade -t <targetName> -z <zoneName> [-c <componentName>]
```

> Actually, there's very little differences between `sail apply` and `sail upgrade`.
> When execute `helm` for pod components, there are no differeneces at all.
> When execute `ansible-playbook` for server components, if no server components are explicitly specified on the command line,
> then `sail apply` and `sail upgrade` are totally same.
> If there are server components explicitly specified, then
> `sail apply` pass `--tags play-<componentName>` options to `ansible-palybook` and
> `sail upgrade` pass `--tags update-<componentName>` options to `ansible-playbook`.
