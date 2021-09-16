package options

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintColorHeader(targetName string, zoneName string) {
	// d.Printf("👉 target: (%s), zone: (%s)\n", o.TargetName, o.ZoneName)
	d := color.New(color.FgWhite, color.Bold, color.BgBlue)
	s := fmt.Sprintf("👉 target: (%s), zone: (%s)", d.Sprint(targetName), d.Sprint(zoneName))
	fmt.Println(s)
}
