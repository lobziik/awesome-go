package main

import (
	"fmt"
	"github.com/russross/blackfriday/v2"
	"io/ioutil"
)

const README_PATH = "../README.md"

// is there a way to make it const?
var skipHeaders = &[]string{"Awesome Go", "Sponsorships", "Contributing", "Contents"}

func skipHeader(header string) bool {
	for _, b := range *skipHeaders {
		if b == header {
			return true
		}
	}
	return false
}

type awsomeEntry struct {
	title string
	link  string
}

type awsomeEntriesCategorizedMap map[string]map[string]awsomeEntry

func (m awsomeEntriesCategorizedMap) fill(ast *blackfriday.Node) {
	skipHeaderFlag := false
	var currentHeader string

	ast.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {

		switch nodeType := node.Type; nodeType {

		case blackfriday.Heading:
			currentHeader = string(node.LastChild.Literal)
			if skipHeader(currentHeader) {
				skipHeaderFlag = true
			} else {
				skipHeaderFlag = false
				m[currentHeader] = make(map[string]awsomeEntry)
			}

		case blackfriday.Link:
			if !skipHeaderFlag {
				title := string(node.LastChild.Literal)
				link := string(node.LinkData.Destination)
				m[currentHeader][title] = awsomeEntry{title: title, link: link}
			}
			return blackfriday.SkipChildren

		default:
			return blackfriday.GoToNext
		}
		return blackfriday.SkipChildren
	})
}

func getAst(mdFilePath string) *blackfriday.Node {
	readmeContent, err := ioutil.ReadFile(mdFilePath)
	if err != nil {
		panic(err)
	}

	optList := []blackfriday.Option{blackfriday.WithExtensions(blackfriday.CommonExtensions)}
	mdProcessor := blackfriday.New(optList...)
	return mdProcessor.Parse(readmeContent)
}

func main() {
	mdAst := getAst(README_PATH)
	awesomeEntriesMap := make(awsomeEntriesCategorizedMap)
	awesomeEntriesMap.fill(mdAst)
	fmt.Println(awesomeEntriesMap)
}
