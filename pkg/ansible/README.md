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

## Inventory Content

```json
{
  "group-a": {
    "hosts": [
      "192.168.28.71",
      "192.168.28.72"
    ],
    "vars": {
      "ansible_ssh_user": "someuser",
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
          "ansible_ssh_user": "someuser",
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
      "ansible_ssh_user": "someuser",
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

- The above JSON content can be decoded into a `Inventory` struct, actually the `items` map of `Inventory`.
- The key of `items` map represents a group name, and the value can be decoded into `Group` struct.
