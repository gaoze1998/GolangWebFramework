package Cli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gaoze1998/GolangWebFramework/Helper"
)

// createApiProject 创建API项目
func createAPIProject(name string) {
	cmd := exec.Command("go", "env", "GOPATH")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	godir := out.String()
	godir = godir[:len(godir)-1]
	exampleAPIProjectZipPath := godir + "/src/" + "github.com/gaoze1998/GolangWebFramework/Cli/exampleApiProject.zip"
	currentWorkDirctory, err := os.Getwd()
	//fmt.Println(godir)
	//fmt.Println(exampleApiProjectZipPath)
	//fmt.Println(currentWorkDirctory)
	if err != nil {
		fmt.Printf("错误发生了: %s\n", err)
		return
	}
	err = Helper.Unzip(exampleAPIProjectZipPath, filepath.Join(currentWorkDirctory, name))
	if err != nil {
		fmt.Printf("grest不完整，请查看文档后重新下载")
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("创建了API项目: %s\n", name)
	os.Chdir(filepath.Join(currentWorkDirctory, name))
	cmd = exec.Command("go", "mod", "init", name)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

// createRegistryProject 创建Rgistry项目
func createRegistryProject(name string) {
	cmd := exec.Command("go", "env", "GOPATH")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	godir := out.String()
	godir = godir[:len(godir)-1]
	exampleRegistryProjectZipPath := godir + "/src/" + "github.com/gaoze1998/GolangWebFramework/Cli/exampleRegistryProject.zip"
	currentWorkDirctory, err := os.Getwd()
	//fmt.Println(godir)
	//fmt.Println(exampleApiProjectZipPath)
	//fmt.Println(currentWorkDirctory)
	if err != nil {
		fmt.Printf("错误发生了: %s\n", err)
		return
	}
	err = Helper.Unzip(exampleRegistryProjectZipPath, filepath.Join(currentWorkDirctory, name))
	if err != nil {
		fmt.Printf("grest不完整，请查看文档后重新下载")
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("创建了Registry项目: %s\n", name)
	os.Chdir(filepath.Join(currentWorkDirctory, name))
	cmd = exec.Command("go", "mod", "init", name)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

// Create 创建项目
func Create(args []string) {
	switch args[2] {
	case "api":
		createAPIProject(args[3])
	case "registry":
		createRegistryProject(args[3])
	}
}
