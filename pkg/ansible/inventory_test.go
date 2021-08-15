package ansible

import (
	"encoding/json"
	"fmt"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

func TestInventory_Json(t *testing.T) {
	data := `{
		"a": {
			"hosts": {
				"192.168.28.71": {
					"hello": "world",
					"xyz": 100
				},
				"192.168.28.72": {
					"hello": "test",
					"xyz": 89
				}
			},
			"vars": {
				"ansible_ssh_user": "johndoe",
				"ansible_ssh_private_key_file": "~/.ssh/mykey",
				"example_variable": "value"
			},
			"children": {
				"b": {
					"hosts": {
						"10.0.0.1": {
							"abc": "xyz",
							"test": 10
						},
						"10.0.0.2": {
							"abc": "xyz",
							"test": 10
						}
					},
					"vars": {
						"b_v": "test"
					},
					"children": {}
				}
			}
		}
	}`

	i := NewInventory()
	err := json.Unmarshal([]byte(data), i)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(i)

	b, _ := json.MarshalIndent(i, "", "  ")
	fmt.Println(string(b))

	y, _ := yaml.Marshal(i)
	fmt.Println(string(y))
	t.Error()
}

func TestInventoryYaml(t *testing.T) {
	data := `
a:
    hosts:
        192.168.28.71:
            hello: world
            xyz: "100"
        192.168.28.72:
            hello: test
            xyz: "89"
    vars:
        ansible_ssh_private_key_file: ~/.ssh/mykey
        ansible_ssh_user: johndoe
        example_variable: value
    children:
        b:
            hosts:
                10.0.0.1:
                    abc: xyz
                    test: "10"
                10.0.0.2:
                    abc: xyz
                    test: "10"
            vars:
                b_v: test
            children: {}`

	i := NewInventory()
	err := yaml.Unmarshal([]byte(data), i)
	if err != nil {
		fmt.Println(err)
	}

	b, _ := yaml.Marshal(i)
	fmt.Println(string(b))

	j, _ := json.MarshalIndent(i, "", "  ")
	fmt.Println(string(j))

	t.Error()
}
