// Package runoob 在线运行代码
package runoob

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/web"
)

var (
	// Templates ...
	Templates = map[string]string{
		"py2":        "print 'Hello World!'",
		"ruby":       "puts \"Hello World!\";",
		"rb":         "puts \"Hello World!\";",
		"php":        "<?php\n\techo 'Hello World!';\n?>",
		"javascript": "console.log(\"Hello World!\");",
		"js":         "console.log(\"Hello World!\");",
		"node.js":    "console.log(\"Hello World!\");",
		"scala":      "object Main {\n  def main(args:Array[String])\n  {\n    println(\"Hello World!\")\n  }\n\t\t\n}",
		"go":         "package main\n\nimport \"fmt\"\n\nfunc main() {\n   fmt.Println(\"Hello, World!\")\n}",
		"c":          "#include <stdio.h>\n\nint main()\n{\n   printf(\"Hello, World! \n\");\n   return 0;\n}",
		"c++":        "#include <iostream>\nusing namespace std;\n\nint main()\n{\n   cout << \"Hello World\";\n   return 0;\n}",
		"cpp":        "#include <iostream>\nusing namespace std;\n\nint main()\n{\n   cout << \"Hello World\";\n   return 0;\n}",
		"java":       "public class HelloWorld {\n    public static void main(String []args) {\n       System.out.println(\"Hello World!\");\n    }\n}",
		"rust":       "fn main() {\n    println!(\"Hello World!\");\n}",
		"rs":         "fn main() {\n    println!(\"Hello World!\");\n}",
		"c#":         "using System;\nnamespace HelloWorldApplication\n{\n   class HelloWorld\n   {\n      static void Main(string[] args)\n      {\n         Console.WriteLine(\"Hello World!\");\n      }\n   }\n}",
		"cs":         "using System;\nnamespace HelloWorldApplication\n{\n   class HelloWorld\n   {\n      static void Main(string[] args)\n      {\n         Console.WriteLine(\"Hello World!\");\n      }\n   }\n}",
		"csharp":     "using System;\nnamespace HelloWorldApplication\n{\n   class HelloWorld\n   {\n      static void Main(string[] args)\n      {\n         Console.WriteLine(\"Hello World!\");\n      }\n   }\n}",
		"shell":      "echo 'Hello World!'",
		"bash":       "echo 'Hello World!'",
		"erlang":     "% escript will ignore the first line\n\nmain(_) ->\n    io:format(\"Hello World!~n\").",
		"perl":       "print \"Hello, World!\n\";",
		"python":     "print(\"Hello, World!\")",
		"py":         "print(\"Hello, World!\")",
		"swift":      "var myString = \"Hello, World!\"\nprint(myString)",
		"lua":        "var myString = \"Hello, World!\"\nprint(myString)",
		"pascal":     "runcode Hello;\nbegin\n  writeln ('Hello, world!')\nend.",
		"kotlin":     "fun main(args : Array<String>){\n    println(\"Hello World!\")\n}",
		"kt":         "fun main(args : Array<String>){\n    println(\"Hello World!\")\n}",
		"r":          "myString <- \"Hello, World!\"\nprint ( myString)",
		"vb":         "Module Module1\n\n    Sub Main()\n        Console.WriteLine(\"Hello World!\")\n    End Sub\n\nEnd Module",
		"typescript": "const hello : string = \"Hello World!\"\nconsole.log(hello)",
		"ts":         "const hello : string = \"Hello World!\"\nconsole.log(hello)",
	}
	LangTable = map[string][2]string{
		"py2":        {"0", "py"},
		"ruby":       {"1", "rb"},
		"rb":         {"1", "rb"},
		"php":        {"3", "php"},
		"javascript": {"4", "js"},
		"js":         {"4", "js"},
		"node.js":    {"4", "js"},
		"scala":      {"5", "scala"},
		"go":         {"6", "go"},
		"c":          {"7", "c"},
		"c++":        {"7", "cpp"},
		"cpp":        {"7", "cpp"},
		"java":       {"8", "java"},
		"rust":       {"9", "rs"},
		"rs":         {"9", "rs"},
		"c#":         {"10", "cs"},
		"cs":         {"10", "cs"},
		"csharp":     {"10", "cs"},
		"shell":      {"10", "sh"},
		"bash":       {"10", "sh"},
		"erlang":     {"12", "erl"},
		"perl":       {"14", "pl"},
		"python":     {"15", "py3"},
		"py":         {"15", "py3"},
		"swift":      {"16", "swift"},
		"lua":        {"17", "lua"},
		"pascal":     {"18", "pas"},
		"kotlin":     {"19", "kt"},
		"kt":         {"19", "kt"},
		"r":          {"80", "r"},
		"vb":         {"84", "vb"},
		"typescript": {"1010", "ts"},
		"ts":         {"1010", "ts"},
	}
)

// RunOOB ...
type RunOOB string

// NewRunOOB ...
func NewRunOOB(token string) RunOOB {
	return RunOOB(token)
}

type result struct {
	Output string `json:"output"`
	Errors string `json:"errors"`
}

// Run ...
func (ro RunOOB) Run(code string, lang string, stdin string) (string, error) {
	// 对菜鸟api发送数据并返回结果
	api := "https://tool.runoob.com/compile2.php"
	runType, ok := LangTable[lang]
	if !ok {
		return "", errors.New("no such language")
	}

	header := http.Header{
		"Content-Type": []string{"application/x-www-form-urlencoded; charset=UTF-8"},
		"Origin":       []string{"https://m.runoob.com"},
		"Referer":      []string{"https://m.runoob.com/"},
		"User-Agent":   []string{web.RandUA()},
	}

	val := url.Values{
		"code":     []string{code},
		"token":    []string{string(ro)},
		"stdin":    []string{stdin},
		"language": []string{runType[0]},
		"fileext":  []string{runType[1]},
	}

	// 发送请求
	client := &http.Client{
		Timeout: time.Minute,
	}

	request, _ := http.NewRequest("POST", api, strings.NewReader(val.Encode()))
	request.Header = header
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("status code " + strconv.Itoa(resp.StatusCode))
	}
	var r result
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return "", err
	}
	// 结果处理
	e := strings.Trim(r.Errors, "\n")
	if e != "" {
		return "", errors.New(e)
	}
	return r.Output, nil
}
