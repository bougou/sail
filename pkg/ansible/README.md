# ansible

## Inventory vs group

```go
type Inventory struct {
	items map[string]*Group
	lock  sync.Mutex
}


type Group struct {
	name GroupName

	// for normal group
	Hosts    *[]string               `json:"hosts,omitempty"`
	Vars     *map[string]interface{} `json:"vars,omitempty"`
	Children Inventory               `json:"children,omitempty"`

	// only for "_meta" group
	Hostvars *map[string]map[string]interface{} `json:"hostvars,omitempty"`
}
```

1. Inventory 虽然是一个 struct，但本质上是一个里面的 items map
2. Group 是一个 struct，

```json
{
  "group-a": {
    "hosts": [
      "192.168.28.71",
      "192.168.28.72"
    ],
    "vars": {
      "ansible_ssh_user": "johndoe",
      "ansible_ssh_private_key_file": "~/.ssh/mykey",
      "example_variable": "value"
    },
    "children": {
      "group-B": {
        "hosts": [
          "192.168.28.71",
          "192.168.28.72"
        ],
        "vars": {
          "ansible_ssh_user": "johndoe",
          "ansible_ssh_private_key_file": "~/.ssh/mykey",
          "example_variable": "value"
        },
        "children": {}
      }
    }
  },
  "this-is-group-name": {
    "hosts": [
      "192.168.28.71",
      "192.168.28.72"
    ],
    "vars": {
      "ansible_ssh_user": "johndoe",
      "ansible_ssh_private_key_file": "~/.ssh/mykey",
      "example_variable": "value"
    }
  },
  "_meta": {
    "hostvars": {
      "192.168.28.71": {
        "host_specific_var": "bar"
      },
      "192.168.28.72": {
        "host_specific_var": "foo"
      }
    }
  }
}
```

以上面的 JSON 内容为例：
1. 整个 JSON 可以映射到一个 Inventory 结构体中（实际上是 inventory 结构体 items 字段里面， items 是一个 map）
2. Inventory 的 items map 里面 key 表示 group Name，value 表示一个 group，可以映射到一个 Group 结构体中。
3. Group 里面有 Hosts, Vars, Children 字段。
