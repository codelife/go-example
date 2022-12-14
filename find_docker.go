package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

// Labels is a map of label names and values.
type Labels map[string]string

// DockerfileInfo holds the path and labels of a Dockerfile.
type DockerfileInfo struct {
	Path   string
	Labels Labels
}

// FindDockerfiles finds all Dockerfiles in the specified directory and subdirectories.
// It returns a map of Dockerfile paths and their labels.
func FindDockerfiles(dir string) (map[string]Labels, error) {
	result := make(map[string]Labels)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.EqualFold(filepath.Base(path), "Dockerfile") {
			return nil
		}
		fmt.Println(path)

		// labels, err := getLabelsFromDockerfile(path)
		labels, err := readDockerfileLabels(path)
		fmt.Println(labels)
		if err != nil {
			return err
		}

		result[path] = labels
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// getLabelsFromDockerfile parses the specified Dockerfile and returns a map of the labels defined in it.
func getLabelsFromDockerfile(path string) (Labels, error) {
	// Read the Dockerfile.
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse the Dockerfile.
	parsed, err := parser.Parse(strings.NewReader(string(content)))
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v", parsed.AST)

	// Extract the labels from the parsed Dockerfile.
	// labels := parsed.AST.Labels

	// Convert the labels slice to a map.
	labelsMap := make(Labels)
	// for _, label := range labels {
	// 	labelsMap[label.Key] = label.Value
	// }

	return labelsMap, nil
}

func readDockerfileLabels(dockerfilePath string) (map[string]string, error) {
	f, err := os.Open(dockerfilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	labels := make(map[string]string)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "LABEL") {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid LABEL line: %q", line)
			}
			kv := strings.SplitN(parts[1], "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("invalid LABEL format: %q", parts[1])
			}
			labels[kv[0]] = kv[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return labels, nil
}

func GitUpdateDir(path string) {
	// 获取给定目录的路径
	repoDir := "/path/to/repo"

	// 调用Git，并在给定目录中执行"git pull"命令
	cmd := exec.Command("git", "pull")
	cmd.Dir = repoDir

	// 执行命令并处理输出
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(output))
}

func main() {
	FindDockerfiles(".")
}
