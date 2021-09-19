# Product

The operation code for a specific product should be put under `products-dir` with `<product-name>` as the directory name.

The `products-dir` is refered as `products` in other places in the documents.

> Ideally, the contents of `products` should be maintained in a Git repo.
> Or, you can keep one separated Git repo for each product.
> If have to change the operation code for the product when the operation logic of the product is changed.

There are three special files or directories reserved by `sail` for special purposes.

> - `ansible.cfg` `sail` will auto generate one if it not exists
> - `shared_roles`
> - `shared_tasks`

## Product directory structure

The directory structure for `products/<product-name>`:

```bash
## Sail specification
components.yaml       # component declaration/variables
components/           # component declaration/variables
vars.yaml             # general variables, like data_dir or timezone etc
# all component variables and general variables are called product variables
# the product variables SHOULD be referenced by Ansible/Helm code
# each deploy target environment will have its own copied and product variables,
# which can be called environment configurations and the are kept under targets/<target>/<zone>/vars.yaml
# `sail` will automatically keep the environment configurations and product variables have a same structure,
# like adding new variables from product variables to environment configurations

## General
resources/            # other resources files，like https cert files, license files, or icon images

## others
README.md             # introduction

## Ansible specific
sail.yaml             # A must exist Ansible Playbook, you can use `sail gen-sail` to generate it automatically
<playbook>.yaml       # any other playbook file
                      # In most cases, one playbook file is sufficient.

## implementation code for components
roles/                # Roles


## Helm Chart
# If you want to develop the product as a whole helm chart,
# you can use `products/<product-name>` directory as standarad helm chart base dir.
Chart.yaml    # helm Chart.yaml
values.yaml   # helm values.yaml
templates     # helm templates
crds          # helm crds
charts        # helm charts
## Note
# Sail also support developing helm chart for each component of the product,
```

A product is composed of components, so you have to develop operation code for each [Component](./component.md).

## Product Variables

You have to declare all product variables which may used when deploying the product.
The product operation code needs to do and finish its job by these variables.

For different target envionrments, the "business operators" just need to maintain
these product variables matched with the real environment, and then execute `sail` command
to install or upgrade the product.

> The "business operator" does not have to care how the product is installed, because the
> installtaion details have been already coded and implemented in product operation code.
>
> If `sail` reported failure in executation, there are basically two reasons:
> 1. the variables of the environment is not configured properly.
> 2. there are bugs in the production operation code.

`sail` reserved the following variables, please don't use them when declaring product variables.

- `_sail*` Sail keeps all `_sail` prefixed variables reserved.
- `inventory`
- `platforms`
- `targetvars`
