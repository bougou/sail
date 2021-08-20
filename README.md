# sail

[中文文档](./dcos/readme_zh.md)

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

## Documents

- [QuickStart](./docs/quick-start.md)
