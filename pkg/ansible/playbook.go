package ansible

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Playbook []Play

func NewPlaybookFromFile(file string) (*Playbook, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read file failed, err: %s", err)
	}

	playbook := &Playbook{}
	if err := yaml.Unmarshal(b, playbook); err != nil {
		return nil, fmt.Errorf("yaml unmarshal failed, err: %s", err)
	}

	return playbook, nil
}

func (p *Playbook) PlaysTags() []string {
	out := []string{}
	for _, play := range *p {
		out = append(out, play.Tags...)
	}

	return out
}

func (p *Playbook) PlaysTagsStartAt(tag string) []string {
	if !strings.HasPrefix(tag, "play-") {
		tag = "play-" + tag
	}

	out := []string{}

	playsTags := p.PlaysTags()
	for i, v := range playsTags {
		if v == tag {
			out = playsTags[i:]
			break
		}
	}

	return out
}
