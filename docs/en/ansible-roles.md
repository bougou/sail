# ansible roles

Sail uses [Ansible Roles](https://docs.ansible.com/ansible/latest/user_guide/playbooks_reuse_roles.html#role-directory-structure) as the building unit for sail server components.

You should develop an ansible role for the component that may be deployed to plain servers.

## Write ansible role for component

The files of `products/<productName>/{components.yaml,components/*.yaml}` hold ONLY the definition of the components that make up the product. Those declaration files DO NOT describe how the components are installed and configured on the servers. The actual installation code for the components should be developed in ansible role.

I will demonstrate an ansible role I wrote. This role will show some common used patterns when writing ansible role for Sail. I will assume that you are already familar with Ansible Role, at least understand ansible role directory structure.

This role is used for deploying `etcd` (and well-known key-value store) cluster on three nodes.

### Preparation

1. Prepare the installation package of `etcd`. Place it under the `packages-dir` dir you specified for `sail`.

    ```yaml
    $ cat ~/.sailrc.yaml
    products-dir: /path/to/products
    targets-dir: /path/to/targets
    packages-dir: /path/to/packages
    ```

    > You can download it from [etcd-v3.5.0-linux-amd64.tar.gz](https://github.com/etcd-io/etcd/releases/download/v3.5.0/etcd-v3.5.0-linux-amd64.tar.gz)
    >
    > You can place it under `packages-dir` directly like `/path/to/packages/etcd-v3.5.0-linux-amd64.tar.gz` or arbitrary nested subdirs like `/path/to/packages/files/etcd-v3.5.0-linux-amd64.tar.gz`.

2. Declare a component for `etcd`.

    ```yaml
    # products/<productName>/components.yaml

    etcd:
      version: v3.5.0
      enabled: true
      external: false
      services:
        client:
          port: 2379
        peer:
          port: 2380
        metric:
          port: 2381
    ```

### Write Role

Create a role dir for `etcd` component, `products/<productName>/roles/etcd`.

In ansible role, you will usually uses at least the following files.

```bash
defaults/main.yaml
tasks/main.yaml
templates/
```

1. In `defaults/main.yaml`, I usually defines some commonly used variables. These variables will only be accessed within the role.

    ```bash
    # defaults/main.yml

    # I use sail's component variables etcd['version'] here
    etcd_file: etcd-{{ etcd['version'] }}-linux-amd64.tar.gz

    # the data_dir variable is sail's general variables defined in  products/<productName>/vars.yaml
    etcd_base_dir: "{{ data_dir }}/etcd"

    etcd_data_dir: "{{ etcd_base_dir }}/data"
    etcd_cert_dir: "{{ etcd_base_dir }}/ssl"
    ```

2. In `tasks/main.yaml`, I will write the installation tasks.

    You can write all ansible tasks in one file, that is `tasks/main.yaml`. But for clarity, readability and maintainability, I always split the tasks into several files.

    ```bash
    tasks/main.yaml
    tasks/install.yaml
    tasks/config.yaml
    tasks/start.yaml
    ```

    In `tasks/main.yaml`, compose them together.

    ```yaml
    - import_tasks: install.yml
      tags:
        - install
        - install-etcd
        - update
        - update-etcd

    - import_tasks: config.yml
      tags:
        - config
        - config-etcd
        - update
        - update-etcd

    - import_tasks: start.yml
      tags:
        - start
        - start-etcd
        - update
        - update-etcd
     ```

     > Note the tags I added to each tasks.
     > Sail largely uses [ansible tags](https://docs.ansible.com/ansible/latest/user_guide/playbooks_tags.html) in its core logic.

3. `tasks/install.yaml`

    In `install.yaml`, I fistly make sure the user/group to be created, then make sure required directories exist, and finally decompress the package file to proper path on the target machine.

    ```yaml
    - name: "create etcd group"
      group:
        name: "etcd"
        state: present

    - name: "create etcd user"
      user:
        name: "etcd"
        group: "etcd"
        shell: /bin/false
        create_home: no

    - name: prepare etcd dirs
      file:
        name: "{{ item }}"
        state: directory
        owner: "etcd"
        group: "etcd"
      with_items:
        - "{{ etcd_cert_dir }}"
        - "{{ etcd_data_dir }}"

    # only extract etcd and etcdctl binaries, ignore all other files in the tarball
    - name: unarchive etcd tar file
      unarchive:
        src: "{{ remote_packages_dir }}/files/{{ etcd_file }}"
        dest: "/usr/local/bin/"
        exclude:
          - "*/Documentation*"
          - "*/README*"
        mode: 0755
        owner: "root"
        group: "root"
        extra_opts:
          - "--strip-components=1"
    ```

4. `tasks/config.yaml`

    In `config.yaml`, I usually complete the configuration tasks.


    ```yaml
    - name: render etcd service file
      template:
        src: etcd.service.j2
        dest: /usr/lib/systemd/system/etcd.service
      tags:
        - config-systemd-service

    - name: reload systemd config
      systemd:
        daemon_reload: yes
      tags:
        - config-systemd-service

    # You can do any other necessary things
    ```

    Note, we use a ansible template file. This template file should be put at `templates/etcd.service.j2`.

    ```ini
    [Unit]
    Description=Etcd Server
    After=network.target
    After=network-online.target
    Wants=network-online.target
    Documentation=https://github.com/coreos

    [Service]
    Type=notify

    User=etcd
    Group=etcd
    WorkingDirectory={{ etcd_data_dir }}

    ExecStart=/usr/bin/etcd \
      --name {{ inventory_hostname }} \
      # omited
      # ...
      # ...
      # omitted
      --listen-client-urls=https://0.0.0.0:{{ etcd['services']['client']['port'] }} \
      --listen-peer-urls=https://0.0.0.0:{{ etcd['services']['peer']['port'] }} \
      --advertise-client-urls=https://{{ inventory_hostname }}:{{ etcd['services']['client']['port'] }} \
      --initial-advertise-peer-urls=https://{{ inventory_hostname }}:{{ etcd['services']['peer']['port'] }} \
      --listen-metrics-urls=http://127.0.0.1:{{ etcd['services']['metric']['port'] }} \
      --data-dir={{ etcd_data_dir }}

    # Run ExecStartPre with root-permissions
    PermissionsStartOnly=true
    ExecStartPre=/usr/bin/chown -R etcd.etcd {{ etcd_cert_dir }}
    ExecStartPre=/usr/bin/chown -R etcd.etcd {{ etcd_data_dir }}

    Restart=always
    RestartSec=10
    LimitNOFILE=65536

    [Install]
    WantedBy=multi-user.target
    ```

    Note the template file uses **sail component variables** `{ etcd['services']['client']['port'] }}` to configure the listened port for etcd.

5.  `tasks/start.yaml`

    ```yaml
    - name: start etcd
      service:
        name: etcd
        state: restarted
        enabled: yes
    ```

## Tips for Write Ansible Tasks

- [Ansible Tips](./../ansible-tips/)
