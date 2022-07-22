//go:build generate
// +build generate

package main

// package prlabels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"
)

const (
	filename = `../../../.github/labeler-pr-test.yml`
)

type templateData struct {
	Names []Label
}

func main() {
	fmt.Printf("Generating %s\n", strings.TrimPrefix(filename, "../../../"))

	lbs := getLabels()

	td := templateData{}

	td.Names = append(td.Names, lbs...)

	sort.SliceStable(td.Names, func(i, j int) bool {
		return td.Names[i].WithUnderscore < td.Names[j].WithUnderscore
	})

	writeTemplate(tmpl, "prlabeler", td)
}

type Label struct {
	LowerCase      string
	WithUnderscore string
}

func getLabels() []Label {
	url := "https://developer.atlassian.com/cloud/jira/platform/swagger-v3.v3.json"
	c := http.Client{Timeout: time.Duration(2) * time.Second}
	resp, err := c.Get(url)
	if err != nil {
		fmt.Printf("Error %s", err)
		return nil
	}
	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	resp.Body.Close()

	var tags = result["tags"].([]interface{})

	var labels []Label
	r := regexp.MustCompile(`\(apps\)`)
	for _, value := range tags {
		// Each value is an interface{} type, that is type asserted as map[string]interface{}
		// due to nested objects in the original JSON response
		m := value.(map[string]interface{})
		name := m["name"]
		ok := r.MatchString(name.(string))
		if ok {
			continue
		}
		l := Label{}
		l.LowerCase = strings.TrimSuffix(strings.ToLower(strings.ReplaceAll(name.(string), " ", "")), "s")
		l.WithUnderscore = strings.TrimSuffix(strings.ToLower(strings.ReplaceAll(name.(string), " ", "_")), "s")
		labels = append(labels, l)
	}

	return labels
}

func writeTemplate(body string, templateName string, td templateData) {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening file (%s): %s", filename, err)
	}

	tp, err := template.New(templateName).Parse(body)
	if err != nil {
		log.Fatalf("error parsing template: %s", err)
	}

	var buffer bytes.Buffer
	err = tp.Execute(&buffer, td)
	if err != nil {
		log.Fatalf("error executing template: %s", err)
	}

	if _, err := f.Write(buffer.Bytes()); err != nil {
		f.Close()
		log.Fatalf("error writing to file (%s): %s", filename, err)
	}

	if err := f.Close(); err != nil {
		log.Fatalf("error closing file (%s): %s", filename, err)
	}
}

var tmpl = `# YAML generated by internal/generate/prlabels/main.go; DO NOT EDIT.
dependencies:
  - '.github/dependabot.yml'
documentation:
  - '**/*.md'
  - 'docs/**/*'
  - 'templates/**/*'
examples:
  - 'examples/**/*'
github_actions:
  - '.github/*.yml'
  - '.github/workflows/*.yml'
linter:
  - 'scripts/*'
  - '.github/workflows/terraform-provider-check.yml'
  - '.github/workflows/workflow-lint.yml'
provider:
  - '.gitignore'
  - '*.md'
  - 'internal/provider/**/*'
  - 'main.go'
repository:
  - '.github/**/*'
  - 'GNUmakefile'
  - 'infra'
tests:
  - '**/*_test.go'
{{- range .Names }}
jira/{{ .LowerCase }}:
  - 'internal/provider/*{{ .WithUnderscore }}.go'
  - 'internal/provider/*{{ .WithUnderscore }}_test.go'
{{- end }}
`
