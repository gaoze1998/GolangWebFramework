package Cli

import (
	"fmt"
	"github.com/gaoze1998/GolangWebFramwork/Helper"
	"os"
)

func createApiProject(name string) {
	godir := os.Getenv("GOPATH")
	exampleApiProjectZipPath := godir + "/src/" + "github.com/gaoze1998/GolangWebFramwork/Cli/exampleApiProject.zip"
	currentWorkDirctory, err := os.Getwd()
	//fmt.Println(godir)
	//fmt.Println(exampleApiProjectZipPath)
	//fmt.Println(currentWorkDirctory)
	if err != nil {
		fmt.Printf("错误发生了: %s\n", err)
		return
	}
	err = Helper.Unzip(exampleApiProjectZipPath, currentWorkDirctory)
	if err != nil {
		fmt.Printf("grest不完整，请查看文档后重新下载")
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("在\"%s\"创建了\"%s\"API项目\n", currentWorkDirctory, name)
}

func Create(args []string) {
	switch args[2] {
	case "api":
		createApiProject(args[3])

	}
}
