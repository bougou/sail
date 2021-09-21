# Linux 运维反模式

## 不要在 `/etc/rc.local` 文件中放置任何与特定组件相关的任务

这些特定组件相关的任务包括但不限于下面几项：

- 在 `PATH` 变量中追加该组件提供的目录
- 系统启动后要执行的命令或脚本

可能有人要质疑，`/etc/rc.local` 文件本来就是做这些事情的，为什么要禁止呢？

因为，`/etc/rc.local` 文件中应该放置系统级通用的任务操作。不管该系统上面部署了什么组件，这个文件中的内容都应该与这个组件没有关系。

假设 componentA 和 componentB 都部署在了一台机器上，两个组件都要在 PATH 变量中追加一个目录 `/opt/component-a/bin` 和 `/opt/component-b/bin` ？

如果使用 /etc/rc.local 文件，两个组件中的部署脚本都需要去修改这个文件，通常会使用 Bash 的 `sed` 命令去追加。这样做有一些问题：

1. 难度较高
2. 不容易实现幂等

那与组件相关的这些任务应该存放在什么地方呢？

使用 `/etc/profile.d/` 目录。把特定组件相关的脚本放到这个目录下，并使用组件相关的名字命名脚本。如 `/etc/profile.d/component-a.sh`。

这样，你就可以做：

```bash
# /etc/profile.d/component-a.sh
export PATH=/opt/component-a/bin:$PATH

# /etc/profile.d/component-b.sh
export PATH=/opt/component-b/bin:$PATH
```

组件的脚本中不会出现其它组件的信息。并且可以比较轻松地实现脚本的幂等。

## 不要在 cron 中使用 ntpdate 去同步时间

- [ntpdate from cron -- DON'T DO THAT!
](https://lists.debian.org/debian-user/2002/12/msg04091.html)
