package ansible

import (
	"fmt"
	"testing"

	"github.com/kr/pretty"
	"gopkg.in/yaml.v3"
)

func Test_Task(t *testing.T) {

	type A []string

	var a A = []string{"{{ a['b'] }}", "b"}
	s, _ := yaml.Marshal(a)
	fmt.Printf("---%s---\n", string(s))

	var c string = "'a'"
	fmt.Println(c)

	var d A = []string{"{{ a }}", "'b'"}
	sss, _ := yaml.Marshal(d)
	fmt.Printf("---%s---\n", string(sss))

	tt := `name: block ip
iptables:
  chain: INPUT
  table: filter
  source: "{{ item }}"
  jump: DROP
  action: insert
  comment: "block security vulnerabilities scan"
with_items:
  - "{{ block_ips }}"`

	task := &Task{}
	if err := yaml.Unmarshal([]byte(tt), task); err != nil {
		t.Error(err)
	}

	pretty.Println(task)

	b, err := yaml.Marshal(task)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("marshal result")
	fmt.Println(string(b))
}
