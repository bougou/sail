# 安装软件

## `yum`

```yaml
- name: install cassandra
  yum:
    name:
      - cassandra-{{cassandra['version'] }}
      - cassandra-tools-{{cassandra['version'] }}
    state: present
  tags:
    - yum
```


## tar.gz

例子 1：

```yaml
# 直接从部署机解压到目标机的指定目录下
# cassandra_filename: apache-cassandra-3.11.11-bin.tar.gz
# 使用 tar -tf apache-cassandra-3.11.11-bin.tar.gz 查看压缩包里面的文件路径为：
# apache-cassandra-3.11.11/bin/
# apache-cassandra-3.11.11/conf/
# apache-cassandra-3.11.11/conf/triggers/
# apache-cassandra-3.11.11/......

- name: install cassandra
  unarchive:
    src: "{{packages_dir}}/files/{{ cassandra_filename }}"
    dest: "/opt"
    owner: "root"
    group: "root"

- name: create a link
  file:
    src: "/opt/apache-cassandra-3.11.11"
    dest: "/opt/cassandra"
    owner: root
    group: root
    state: link

```

> 这种用法解压出来的结果是 `/opt/apache-cassandra-3.11.11`。
> 你可能需要创建一个额外的软链 `/opt/cassandra` 链接到 `/opt/apache-cassandra-3.11.11`。
> 这里的问题是你还需要使用某种方法取得压缩包文件带来的路径（`/opt/apache-cassandra-3.11.11`）。

例子 2：

```yaml
# 灵活使用 --strip-components 参数
- name: install cassandra
  unarchive:
    src: "{{packages_dir}}/files/{{ cassandra_filename }}"
    dest: "/opt/cassandra"
    owner: "root"
    group: "root"
    extra_opts: "--strip-components=1"
```

> 这种用法通过 `--strip-components=1` 参数直接把压缩包的第一段目录移除，并解压到 `/opt/cassandra` 目录下。
> 效果就像是把压缩包 `apache-cassandra-3.11.11/{fileanddirs}` 的路径解压到了 `/opt/cassandra/{fileanddirs}`。
> 省去了创建软链的需要。
