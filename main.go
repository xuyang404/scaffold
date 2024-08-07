package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	var (
		localTemplatePath  string
		remoteTemplatePath string
		projectName        string
		branch             string
	)

	flag.StringVar(&localTemplatePath, "local", "", "本地模板路径")
	flag.StringVar(&remoteTemplatePath, "remote", "", "远程仓库url")
	flag.StringVar(&projectName, "name", "", "项目名称")
	flag.StringVar(&branch, "branch", "main", "要使用的分支（仅当模板是远程仓库时）")

	flag.Parse()
	if localTemplatePath == "" && remoteTemplatePath == "" {
		log.Println("请指定本地模板路径或远程仓库url")
		os.Exit(1)
	}

	if projectName == "" {
		log.Println("请指定项目名称")
		os.Exit(1)
	}

	templatePath := ""
	if localTemplatePath != "" {
		templatePath = localTemplatePath
	}

	if remoteTemplatePath != "" {
		templatePath = remoteTemplatePath
	}

	replacements := getReplacements()
	if _, ok := replacements["{{PROJECT_NAME}}"]; !ok {
		replacements["{{PROJECT_NAME}}"] = projectName
	}

	var err error

	if remoteTemplatePath != "" {
		err = handleRemoteTemplate(templatePath, branch, projectName, replacements)
	} else {
		err = copyTemplate(templatePath, projectName, replacements)
	}

	if err != nil {
		log.Printf("Error creating project: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Project %s created successfully!\n", projectName)
}

func handleRemoteTemplate(templateRepo, branch, projectName string, replacements map[string]string) (err error) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "template-*")
	if err != nil {
		return fmt.Errorf("error creating temporary directory: %s", err)
	}

	// 清理临时目录
	defer os.RemoveAll(tempDir)

	// 克隆模板仓库
	cloneCmd := exec.Command("git", "clone", "-b", branch, templateRepo, tempDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr

	if err = cloneCmd.Run(); err != nil {
		return fmt.Errorf("error cloning template repository: %s", err)
	}

	return copyTemplate(tempDir, projectName, replacements)
}

func getReplacements() map[string]string {
	scanner := bufio.NewScanner(os.Stdin)
	replacements := make(map[string]string)
	fmt.Println("输入替换值(key=value)，空行结束: ")
	for {
		fmt.Print("> ")
		scanner.Scan()
		line := scanner.Text()
		if line == "" {
			break
		}

		splits := strings.Split(line, "=")
		if len(splits) != 2 {
			fmt.Println("无效输入，请以key=value格式输入")
			continue
		}

		replacements[splits[0]] = splits[1]
	}

	return replacements
}

func copyTemplate(src, dist string, replacements map[string]string) (err error) {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// 获取目标路径
		targetPath := filepath.Join(dist, relPath)

		if info.IsDir() {
			// 创建目录
			return os.MkdirAll(targetPath, info.Mode())
		}

		return copyAndReplaceFile(path, targetPath, info.Mode(), replacements)
	})
}

func copyAndReplaceFile(src, dist string, mode os.FileMode, replacements map[string]string) (err error) {

	// 读取源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	content, err := io.ReadAll(sourceFile)
	if err != nil {
		return err
	}

	newContent := string(content)
	for key, value := range replacements {
		newContent = strings.ReplaceAll(newContent, key, value)
	}

	// 读取目标文件
	targetFile, err := os.OpenFile(dist, os.O_CREATE|os.O_RDWR|os.O_TRUNC, mode)
	if err != nil {
		return err
	}

	defer targetFile.Close()

	_, err = targetFile.WriteString(newContent)
	return
}
