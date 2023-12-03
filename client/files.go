package client

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type FileClient struct {
	root    string
	dayPath string
}

func NewFileClient(rootDir string, year int, day int) (FileClient, error) {
	if _, err := os.Stat(rootDir); err != nil {
		return FileClient{}, fmt.Errorf("Error when stating root directory: %s", err)
	}

	dayDirPath := filepath.Join(rootDir, strconv.Itoa(year), strconv.Itoa(day))

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

func createFileIfNotExists(filePath string, contents []byte) error {
	if exists, err := pathExists(filePath); err == nil && !exists {
		fmt.Printf("File %s did not exist, creating..\n", filePath)
		os.WriteFile(filePath, contents, 0755)
	} else {
		return err
	}

	fmt.Printf("File %s already exists, not creating..\n", filePath)

	return nil
}

func (fc FileClient) ScaffoldDay(year int, day int) error {
	if err := os.MkdirAll(fc.dayPath, 0755); err != nil {
		return err
	}

	s1FilePath := filepath.Join(fc.dayPath, "s1.py")
	s2FilePath := filepath.Join(fc.dayPath, "s2.py")

	solFileStr := fmt.Sprintf(`import sys

with open(f"%d/%d/{sys.argv[1]}", "r") as file:
    lines = [l.strip() for l in file.readlines()]

print(lines[-1], file=sys.stderr)
`, year, day)

	if err := createFileIfNotExists(s1FilePath, []byte(solFileStr)); err != nil {
		return err
	}

	if err := createFileIfNotExists(s2FilePath, []byte(solFileStr)); err != nil {
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
