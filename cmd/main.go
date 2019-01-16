package main

import (
	"bytes"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
}

func main() {
	var component string
	var packageName string
	var version string
	flag.StringVar(&component, "component", "", "Component name")
	flag.StringVar(&version, "version", "", "Component version")
	flag.StringVar(&packageName, "package", "", "Package name")
	flag.Parse()
	if len(component) == 0 {
		log.Println("Please specify component name")
		return
	}
	if len(packageName) == 0 {
		currentDir, _ := os.Getwd()
		packageName = findPackageNameInPath(currentDir)
		log.Println("Deduced package name:", packageName)
	}
	if len(version) == 0 {
		version = NodeVersion
	}
	log.Println("Creating component:", component, "with version:", version)
	t := template.New("component")
	t, err := t.Parse(NodeTemplate)
	if err != nil {
		log.Println("Unable to parse component template", err)
		return
	}
	config := TemplateConfig{
		PackageName:    packageName,
		Component:      component,
		ComponentLower: strings.ToLower(component),
		Version:        version,
	}
	buffer := bytes.Buffer{}
	t.Execute(&buffer, config)

	lowerName := config.ComponentLower + ".go"
	err = ioutil.WriteFile(lowerName, buffer.Bytes(), 0644)
	if err != nil {
		log.Println("Unable to write generated file", err)
	}
}

func findPackageNameInPath(currentPath string) string {
	files, err := ioutil.ReadDir(currentPath)
	if err != nil {
		log.Fatal(err)
	}

	var absoluteFilePath string
	for _, f := range files {
		if f.IsDir() || !strings.Contains(f.Name(), ".go") {
			continue
		}
		absoluteFilePath = currentPath + "/" + f.Name()
		data, err := ioutil.ReadFile(absoluteFilePath)
		if err != nil {
			continue
		}
		var re = regexp.MustCompile(`.*package (.*)`)
		matched := re.FindAllStringSubmatch(string(data), -1)
		if len(matched) == 0 {
			continue
		}
		return matched[0][1]
	}
	// package name will be the dir name
	return path.Base(currentPath)
}
