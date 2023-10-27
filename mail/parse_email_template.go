package mail

import (
	"github.com/flosch/pongo2"
)

func ParseMailTemplate(username, link string) (string, error) {
	tpl, err := pongo2.FromFile("templates/layouts/main.django")
	if err != nil {
		return "", err
	}

	context := pongo2.Context{
		"Name": username,
		"Link": link,
	}

	rendered, err := tpl.Execute(context)
	if err != nil {
		return "", err
	}

	return string(rendered), nil
}