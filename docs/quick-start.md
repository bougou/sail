# Quick Start

Suppose, your boss assigns you a task to deploy a software product called `foobar`.

0. Firstly, you should get the operation code (installation script) of product, and prepare some servers.

    ```bash
    # use the examples of this project to get started.
    sail-examples/products/foobar/

    # suppose you put the sail-examples at your home dir
    $ ls ~/sail-examples/products/foobar/

    $ cd ~/sail-examples
    $ mkdir targets
    $ mkdir packages


    # create a ~/.sailrc.yaml to set global options, or else
    # you have to explicitly specify these options everytime when running sail command.
    $ vim ~/.sailrc.yaml
    products-dir: /path/to/home/dir/sail-examples/products
    targets-dir: /path/to/home/dir/sail-examples/targets
    packages-dir: /path/to/home/dir/sail-examples/packages


    # place the sail exectuable binary file to some place.
    # like /usr/bin/sail or any place, and chmod it with executable permission.
    ```

1. Create a environment (target/zone) and specify the product to be deployed.

    ```bash
    # the target and zone name can be any string
    $ ./sail conf-create --target myfoobarapp --zone default --product foobar --hosts 10.0.0.1
    ```

2. Now start to deploy

    ```bash
    $ ./sail apply --target myfoobarapp --zone default
    ```
