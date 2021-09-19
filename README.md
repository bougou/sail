# sail

[中文文档](./docs/readme_zh.md)

`sail` is an operation framework based on Ansible/Helm. `sail` follows the principles of ***Infrastructure as Code (IaC)***,  ***Operation as Code (OaC)***，and ***Everything as Code***. So it is a tool for ***DevOps***.

> `sail` is especially suitable for the privatization delivery of software products.

Although `sail` strongly utilizes `Ansible` and `Helm`, `sail` does not write the ansible tasks or helm chart templates for you. It's still your responsibility to develop ansible tasks and helm templates files.

> - `Ansible` is for deploying to normal servers.
> - `Helm` is for deploying to Kubernetes.

## Product and Target and Package

`sail` has three main concepts: `Product`, `Target`, and `Package`

### Product

Product is a specific software product. It is composed of components.

A `Product` can be simple and small with just several components, or can be big and complex with hunderds of components.

When you are responsible for managing and operating a software product, you should prepare the ***product operation code***. The operation code of a product is not the functional code that written by product developers. It is code that direct the operators to install and manage the product.

### Target

Target represents the environment where the product softwares are installed and run.

In `sail`, we used two hierarchies to arrange environments.

The two hierarchies are:

- `target`
- `zone`

A `zone` must be created under a specific `target` and there can be multiple `zone`(s) under a `target`. The `target(s)` are totally isolated from each other in `sail`.

### Package

Package(s) are real software artifacts. They are normally compressed in some format, like `.rpm`, `.tar.gz`, `.gzip` ....

## Use sail

Generally, you run `sail` on a centralized machine (deploy machine).

The `sail` command uses three directories to do its job.

```yaml
# cat ~/.sailrc.yaml
products-dir: /path/to/products    # Store the operation code of product(s).
targets-dir: /path/to/targets      # Store the environment informations.
packages-dir: /path/to/packages    # Store package files.
```

> Note, these three directories can be located at different places (that is not under a same parent dir).

> The `products-dir` contains all operation code of products(s). You should put your specific `<productName>` dir under `products-dir` even if you only have one product. Normally, the content of products dir should be managed by git.

> The `targets-dir` keeps all configurations of environments. This dir should exists on the deploy machine.

> The `packages-dir` holds all package files. You can create recursive directories under it to store any files.

To use `sail` to do operation, you need to understand how to declare `product` and how to define the `component(s)` that make up of the `product`.

## Why use sail

Using `sail`, you can change ***the operations of any products*** into the following commands:

```bash
# Create a target environment (target/zone)
$ sail conf-create -t <targetname> -z <zonename> \
  -p <productname> \
  --hosts <hosts-for-components> \
  --hosts <hosts-for-components> \
  ...

# Deploy
$ sail apply -t <targetname> -z <zonename>

# Upgrade specific components
$ sail upgrade -t <targetname> -z <zonename> \
  -c <componentName1>/<componentVersion1> \
  -c <componentName2>/<componentVersion2>
```

## Documents

- [How to develop product operation code](./docs/en/product.md)
- [How to develop component operation code](./docs/en/component.md)
- [How to maintain cmdb for deploy target](./docs/en/cmdb.md)
- [Sail command usage](./docs/en/sail-commands.md)
- [QuickStart](./docs/quick-start.md)
