# FAQ

## How to modify the hosts for a component

1. You can directly edit `hosts.yaml` file

```bash
$ vim <targets-dir>/<target>/<zone>/hosts.yaml

your-component-name:
  hosts:
    10.0.0.1: {}
    10.0.0.2: {}
    10.0.0.3: {}
```

2. You can also use `--hosts` option of `conf-update` subcommand

```bash
# the hosts will be finally updated into <target>/<zone>/hosts.yaml file.
$ ./sail conf-update -t <target> -z <zone> --hosts <componentName>/10.0.0.1,10.0.0.2,10.0.0.3
```

> The format of `--hosts` options is as follow:
> - `--hosts A,B/10.0.0.1,10.0.0.2`
> - `--hosts +C/10.0.0.3,10.0.0.4`
> - `--hosts -C,D,E/10.0.0.4`
> - `--hosts 10.0.0.1`
>
> Whether `+` or `-` is used as the first character will have different behaviors.
>
> - `+` means to add new hosts into existing hosts
> - `-` means to remove hosts from existing hosts
> - neither `+` nor `-` means to update(override) the existing hosts
>
> Using `--hosts` option, normally you would specify the component names,
> and separate the components and the hosts by using `/`.
>
> If no components is specified, it means for **a special component** name `_cluster`.
