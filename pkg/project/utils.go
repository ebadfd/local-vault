package project

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/manifoldco/promptui"
)

func projectSelector(projects []Project) (*Project, error) {
	if len(projects) == 0 {
		return nil, errors.New("no projects available or found matching the criteria")
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ .Name }}?",
		Active:   "\U0001F872 {{ .Name | cyan }} ({{ .ID | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .ID | red }})",
		Selected: "\U0001F872 {{ .Name | red | cyan }}",
		Details: `
--------- Project ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "ProjectId:" | faint }}	{{ .ID }}`,
	}

	searcher := func(input string, index int) bool {
		project := projects[index]
		name := strings.Replace(strings.ToLower(project.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "Select the project",
		Items:     projects,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return nil, errors.New("Failed to select a project")
	}

	return &projects[i], nil
}

func appsSelector(apps []string) (string, error) {
	if len(apps) == 0 {
		return "", errors.New("No available apps found for the specified criteria")
	}

	prompt := promptui.Select{
		Label: "Select App",
		Items: apps,
	}

	_, result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, nil
}

func listPrettyPrint(title string, content string, prefix string) {
	fmt.Printf("%s:\n", title)
	fmt.Println(strings.Repeat("-", len(title)+1))
	for _, line := range strings.Split(content, "\n") {
		fmt.Printf("%s%s\n", prefix, line)
	}
	fmt.Println()
}

func RenderList(data []string, title string) {
	l := list.NewWriter()

	for _, item := range data {
		l.AppendItem(item)
	}

	l.SetStyle(list.StyleBulletCircle)
	listPrettyPrint(title, l.Render(), "")
}
