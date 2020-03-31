package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"text/template"

	"github.com/prmsrswt/gophercises/adventure"
)

var (
	storyPath string
	tpl       *template.Template
)

func init() {
	flag.StringVar(&storyPath, "story", "gopher.json", "path of the JSON file containing story")
	tpl = template.Must(template.New("").Parse(storyArcTemplate))
}

func main() {
	file, err := os.Open(storyPath)
	if err != nil {
		panic("Error opening json file")
	}
	defer file.Close()

	story, err := adventure.ParseStory(file)
	if err != nil {
		panic("Error parsing json story")
	}

	play(story, "intro")
}

var storyArcTemplate = `
TITLE:		{{ .Title }}

{{ range .Story }}
{{ . }}
{{ end }}
{{ if .Options }}
What would you do?
{{ range $i, $v := .Options }}
{{ $i }}) {{ $v.Text }}
{{ end }}
{{ else }}
The End.
{{ end }}

`

func askWithRetries(maxValue int) int {
	var ans int
	fmt.Print("Choose an option: ")
	_, err := fmt.Scan(&ans)

	if err != nil {
		fmt.Print("Invalid input.\nHere, try again. ")
		return askWithRetries(maxValue)
	}

	if ans < 0 || ans > maxValue {
		fmt.Print("Invalid input.\nHere, try again. ")
		return askWithRetries(maxValue)
	}

	return ans
}

func play(story adventure.Story, arc string) {
	storyArc, ok := story[arc]
	if !ok {
		panic("Error: cannot find given chapter")
	}

	clear()
	tpl.Execute(os.Stdout, storyArc)

	if len(storyArc.Options) > 0 {
		ans := askWithRetries(len(storyArc.Options) - 1)

		option := storyArc.Options[ans]

		play(story, option.Arc)
	}
}

func clear() {
	switch runtime.GOOS {
	case "linux":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
