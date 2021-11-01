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

type pepper struct {
	Name     string
	Content  string
	Comment  string
	CmdType  string
}

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
	peppers := []pepper{
		{Name: "Help", Content: "Help", Comment: "add custom cmdline to snippets.yaml"},
	}

	info := BaseInfo{};
	snippets, err := info.GetConf("snippets.yml");
	if err != nil {
		Warning.Println(err.Error())
	} else {
		for _, v := range snippets.Snippet {
			peppers = append(peppers, pepper{Name: v.Name, Content: v.Content, Comment: v.Comment })
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
{{ "Comment:" | faint }}	{{ .Comment }}`,
	}

	searcher := func(input string, index int) bool {
		pepper := peppers[index]
		ctx := strings.ToLower(pepper.Content)
		input = strings.Replace(strings.ToLower(input), " ", ".*", -1)

		match, _ := regexp.MatchString(input, ctx)  
		return match
	}

	prompt := promptui.Select{
		Label:     "Select Which",
		Items:     peppers,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v, %d\n", err, i)
		return
	}

	command := fmt.Sprintf("echo '%s' | vipe ", peppers[i].Content)
	cmd := exec.Command("/bin/bash", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		Warning.Printf("Cmd error choose number %d: %v\n", i+1, peppers[i])
		return
	}
	fmt.Printf("%s", output)
}
