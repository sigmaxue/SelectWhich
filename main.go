package main

import (
	"os"
	"io"
	"fmt"
	"log"
	"bufio"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

type pepper struct {
	Name     string
	HeatUnit int
	Peppers  int
	CmdType  int
}

var (
	Info *log.Logger
	Warning *log.Logger
	Error * log.Logger
)

func init(){
	errFile,err:=os.OpenFile("~/errors.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	if err!=nil{
		log.Fatalln("打开日志文件失败：",err)
	}

	Info = log.New(os.Stdout,"Info:",log.Ldate | log.Ltime | log.Lshortfile)
	Warning = log.New(errFile,"Warning:",log.Ldate | log.Ltime | log.Lshortfile)
	Error = log.New(io.MultiWriter(os.Stderr,errFile),"Error:",log.Ldate | log.Ltime | log.Lshortfile)

}

func main() {
	peppers := []pepper{
		{Name: "Bell Pepper", HeatUnit: 0, Peppers: 0},
		{Name: "Banana Pepper", HeatUnit: 100, Peppers: 1},
	}

	fp,err := os.Open("~/shell.txt")
	if err!=nil{
		fmt.Println(err) //打开文件错误
		return 
	}
	buf := bufio.NewScanner(fp)
	for {
		if !buf.Scan() {
			break //文件读完了,退出for
		}
		line := buf.Text() //获取每一行
		peppers = append(peppers, pepper{Name: line, HeatUnit: 0, Peppers: 0, CmdType: 1})
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Name | cyan }} ({{ .HeatUnit | red }})",
		Inactive: "  {{ .Name | blue }} ({{ .HeatUnit | blue }})",
		Selected: " {{ .Name | red | cyan }}",
		Details: `
--------- Pepper ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Heat Unit:" | faint }}	{{ .HeatUnit }}
{{ "Peppers:" | faint }}	{{ .Peppers }}`,
	}

	searcher := func(input string, index int) bool {
		pepper := peppers[index]
		name := strings.ToLower(pepper.Name)
		input = strings.Replace(strings.ToLower(input), " ", ".*", -1)

		match, _ := regexp.MatchString(input, name)  
		//Warning.Println("searcher: ", input, index)
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
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	// Warning.Printf("You choose number %d: %v\n", i+1, peppers[i])

	fmt.Println(peppers[i].Name)
}
