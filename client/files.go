package client

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileClient struct {
	root    string
	dayPath string
}

func NewFileClient(rootDir string, year string, day string) (FileClient, error) {
	if _, err := os.Stat(rootDir); err != nil {
		return FileClient{}, fmt.Errorf("Error when stating root directory: %s", err)
	}

	dayDirPath := filepath.Join(rootDir, year, day)

	return FileClient{
		root:    rootDir,
		dayPath: dayDirPath,
	}, nil
}

func pathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func createFileIfNotExists(filePath string, contents []byte) (bool, error) {
	exists, err := pathExists(filePath)
	if err != nil {
		return false, err
	}

	if exists {
		fmt.Printf("File %s already exists, not creating..\n", filePath)
		return true, nil
	}

	fmt.Printf("File %s did not exist, creating..\n", filePath)
	return false, os.WriteFile(filePath, contents, 0755)
}

func (fc FileClient) ScaffoldDay(year string, day string) error {
	if err := os.MkdirAll(fc.dayPath, 0755); err != nil {
		return err
	}

	s1FilePath := filepath.Join(fc.dayPath, "s1.py")

	solFileStr := fmt.Sprintf(`import sys
from pprint import pprint

with open(f"%s/%s/%s/{sys.argv[1]}", "r") as file:
    lines = [l.strip() for l in file.readlines()]

print(lines[-1], file=sys.stderr)
`, fc.root, year, day)

	if _, err := createFileIfNotExists(s1FilePath, []byte(solFileStr)); err != nil {
		return err
	}


	gitignorePath := filepath.Join(fc.root, year)

	gitignoreStr := `*

!/*
!/*/*.py
!.gitignore
`

	if _, err := createFileIfNotExists(gitignorePath, []byte(gitignoreStr)); err != nil {
		return err
	}

	return nil
}

func (fc FileClient) WriteInput(contents []byte) error {
	inputPath := filepath.Join(fc.dayPath, "input")
	if err := os.WriteFile(inputPath, contents, 0755); err != nil {
		return fmt.Errorf("Something went wrong when writing input: %s", err)
	}

	return nil
}

func (fc FileClient) WriteTestInput(contents []byte) error {
	inputPath := filepath.Join(fc.dayPath, "test")
	if err := os.WriteFile(inputPath, contents, 0755); err != nil {
		return fmt.Errorf("Something went wrong when writing test input: %s", err)
	}

	return nil
}

func (fc FileClient) SolutionFileExists(problem int) (bool, error) {
	solPath := filepath.Join(fc.dayPath, fmt.Sprintf("%d.sol", problem))
	return pathExists(solPath)
}

func (fc FileClient) ReadSolutionFile(problem int) (string, error) {
	solPath := filepath.Join(fc.dayPath, fmt.Sprintf("%d.sol", problem))
	bytes, err := os.ReadFile(solPath)
	return string(bytes), err
}

func (fc FileClient) ProblemSolved(problem int) (bool, error) {
	solPath := filepath.Join(fc.dayPath, fmt.Sprintf("%d.solved", problem))
	return pathExists(solPath)
}

func (fc FileClient) SetProblemSolved(problem int) error {
	solPath := filepath.Join(fc.dayPath, fmt.Sprintf("%d.sol", problem))
	solvedPath := filepath.Join(fc.dayPath, fmt.Sprintf("%d.solved", problem))
	return os.Rename(solPath, solvedPath)
}

func (fc FileClient) CreateSecondSolutionFile() (bool, error) {
	s1FilePath := filepath.Join(fc.dayPath, "s1.py")

	s1FileContents, err := os.ReadFile(s1FilePath)
	if err != nil {
		return false, err
	}

	s2FilePath := filepath.Join(fc.dayPath, "s2.py")

	return createFileIfNotExists(s2FilePath, s1FileContents)
}
