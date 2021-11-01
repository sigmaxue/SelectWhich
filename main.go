package main

import (
	"os"
	"io"
	"fmt"
	"log"
	"regexp"
	"strings"
	 "os/exec"

	"github.com/manifoldco/promptui"
)

var (
	Info *log.Logger
	Warning *log.Logger
	Error * log.Logger
)

func init(){
	errFile,err:=os.OpenFile("errors.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	if err!=nil{
		log.Fatalln("打开日志文件失败：",err)
	}

	Info = log.New(os.Stdout,"Info:",log.Ldate | log.Ltime | log.Lshortfile)
	Warning = log.New(errFile,"Warning:",log.Ldate | log.Ltime | log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr,errFile),"Error:",log.Ldate | log.Ltime | log.Lshortfile)

}

func main() {
	items := []Snippet{
		{Name: "Help", Content: "ls", Comment: "add custom cmdline to snippets.yaml", CmdType: "shell"},
	}

	info := BaseInfo{};
	snippets, err := info.GetConf("snippets.yml");
	if err != nil {
		Warning.Println(err.Error())
	} else {
		for _, v := range snippets.Snippets {
			items = append(items, Snippet{Name: v.Name, Content: v.Content, Comment: v.Comment, CmdType: v.CmdType })
		}
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Content | cyan }} ({{ .Name | red }})",
		Inactive: "  {{ .Content | blue }} ({{ .Name | blue }})",
		Selected: " {{ .Content | red | cyan }}",
		Details: `
--------- Detail ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Content:" | faint }}	{{ .Content }}
{{ "CmdType:" | faint }}	{{ .CmdType }}
{{ "Comment:" | faint }}	{{ .Comment }}`,
	}

	searcher := func(input string, index int) bool {
		item := items[index]
		ctx := strings.ToLower(item.Content)
		input = strings.Replace(strings.ToLower(input), " ", ".*", -1)

		match, _ := regexp.MatchString(input, ctx)  
		return match
	}

	prompt := promptui.Select{
		Label:     "Select Which",
		Items:     items,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v, %d\n", err, i)
		return
	}

	command := items[i].Content
	if strings.Contains(items[i].CmdType, "snippet") {
		command = fmt.Sprintf("echo '%s' | vipe ", items[i].Content)
		cmd := exec.Command("/bin/bash", "-c", command)
		output, err := cmd.Output()
		if err != nil {
			Warning.Printf("Cmd error choose number %d: %v\n", i+1, items[i])
			return
		}
		command = string(output)
	}
	if strings.Contains(items[i].CmdType, "shell") {
		cmd := exec.Command("/bin/bash", "-c", command)
		output, err := cmd.Output()
		if err != nil {
			Warning.Printf("Cmd run error: %v\n", command)
			return
		}
		command = string(output)
	}

	fmt.Printf("%s", command)
}
