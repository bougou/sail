# sail

[中文文档](./docs/readme_zh.md)

`sail` is an operation framework based on Ansible. `sail` follows the principles of ***Infrastructure as Code (IaC)***,  ***Operation as Code (OaC)***，and ***Everything as Code***. So it is a tool for ***DevOps***.

Although `sail` strongly utilizes `Ansible`, `sail` does not write the ansible tasks/roles for you. It's still your responsibility to develop ansible roles and tasks.

## Product and Target and Package

`sail` has three main concepts: `Product`, `Target`, and `Package`

### Product

Product is a specific software product. It is composed of components.

A `Product` can be simple and small with just several components, or can be big and complex with hunderds of components.

When you have to manage and operate a product, you should prepare the product operation code. The operation code of a product is not the functional code that written by product developers. It is code that direct the operators to how to install and manage the product.

### Target

Target represents a the environment where the product softwares are installed and run.

In `sail`, we used two hierarchies to arrange environments.

The two hierarchies are:

- `target`
- `zone`

A `zone` must be created under a specific `target` and there can be multiple `zone`(s) under a `target`.

### Package

It's real software artifacts. It's normally compressed in some format, like `.rpm`, `.tar.gz`, `.gzip` ....

To use `sail` to do operation, you need to understand how to declare `product` and how to define the `component(s)` that make up of the `product`.

## Use sail

The `sail` command uses three directories to do its job.

```yaml
# cat ~/.sailrc.yml
products-dir: /path/to/products    # Store the operation code of product(s).
targets-dir: /path/to/targets      # Store the environment informations.
packages-dir: /path/to/packages    # Store package files.
```

> Note, these three directories can be located at different places (be not under a same parent dir).

> The `products-dir` contains all operation code of products(s). You should put your specific `<productName>` dir under `products-dir` even if you only have one product. Normally, the content of products dir should be managed by git.

> The `targets-dir` keeps all configurations of environments.

> The `packages-dir` holds all package files. You can create recursive directories under it.


## Documents

- [QuickStart](./docs/quick-start.md)
