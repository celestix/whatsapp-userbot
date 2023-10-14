package helputil

import (
	"fmt"
	"strings"
)

const _HELP_SEP = "\n\n"

// map["cmd name"][]string{"short description", "explanation"}
type Help map[string][]string

func (h Help) Add(name, desc string) {
	h[strings.ToLower(name)] = strings.SplitN(desc, _HELP_SEP, 2)
}

func (h Help) GetOne(name string) string {
	return strings.Join(h[strings.ToLower(name)], _HELP_SEP)
}

func (h Help) Get() string {
	var ghelp strings.Builder
	var n int
	for _, e := range h {
		if len(e) == 0 || e[0] == "" {
			continue
		}
		n++
		ghelp.WriteString(
			fmt.Sprintf(
				"\n %d.) %s",
				n,
				e[0],
			),
		)
	}
	return ghelp.String()
}
