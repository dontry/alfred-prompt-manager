package service

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	aw "github.com/deanishe/awgo"
)

type Prompt struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

type Service struct {
	wf                         *aw.Workflow
	customPromptsFileName      string
	awesomePromptsDownloadLink string
	awesomePromptsFileName     string
}

func NewService(
	customPromptsFileName string,
	awesomePromptsDownloadLink string,
	awesomePromptsFileName string,
) *Service {
	return &Service{
		wf:                         aw.New(),
		customPromptsFileName:      customPromptsFileName,
		awesomePromptsDownloadLink: awesomePromptsDownloadLink,
		awesomePromptsFileName:     awesomePromptsFileName,
	}
}

func (s *Service) Run(action func()) {
	s.wf.Run(action)
}

func (s *Service) Query(input string, all bool, action string) func() {
	return func() {
		var awesomePrompts, customPrompts []Prompt
		var err error

		if all {
			awesomePrompts, err = s.readPrompts(s.awesomePromptsFileName)
			s.check(err)
		}

		customPrompts, err = s.readPrompts(s.customPromptsFileName)
		s.check(err)

		prompts := append(awesomePrompts, customPrompts...)

		for _, p := range prompts {
			var arg string
			if action == "delete" {
				arg = p.Title
			} else {
				arg = p.Subtitle
			}

			s.wf.
				NewItem(p.Title).
				Subtitle(p.Subtitle).
				Arg(arg).Valid(true).
				Copytext(p.Subtitle).
				Var("action", action)
		}
		s.wf.Filter(input)
		s.wf.SendFeedback()
	}
}

func (s *Service) Add(title string, subtitle string) func() {
	return func() {
		if title == "" || subtitle == "" {
			s.wf.Fatal("Name or content cannot be empty")
			os.Exit(1)
		}
		// read a json file
		file, err := os.OpenFile(s.customPromptsFileName, os.O_RDWR|os.O_CREATE, 0644)
		s.check(err)
		defer file.Close()

		prompts, err := s.readPrompts(s.customPromptsFileName)
		s.check(err)

		for _, p := range prompts {
			if p.Title == title {
				s.wf.Fatal("Prompt already exists")
				os.Exit(1)
			}
		}

		prompts = append(prompts, Prompt{title, subtitle})

		err = s.writePrompts(prompts, s.customPromptsFileName)
		s.check(err)
	}
}

func (s *Service) Delete(input string) func() {
	return func() {
		prompts, err := s.readPrompts(s.customPromptsFileName)
		s.check(err)

		if len(prompts) == 0 {
			s.wf.Warn("Warning", "No prompts to delete")
			os.Exit(1)
		}

		newPrompts := prompts[:0]
		for _, p := range prompts {
			if strings.ToLower(p.Title) != strings.ToLower(input) {
				newPrompts = append(newPrompts, p)
			}
		}
		err = s.writePrompts(newPrompts, s.customPromptsFileName)
		s.check(err)
	}
}

func (s *Service) Download() func() {
	return func() {
		const tempFile = "tmp_awesome_prompts.csv"
		s.downloadFile(s.awesomePromptsDownloadLink, tempFile)
		data, err := s.readCSV(tempFile)
		s.check(err)
		err = s.writeJsonFile(data, s.awesomePromptsFileName)
		s.check(err)
		err = s.deleteFile(tempFile)
		s.check(err)
	}
}

func (s *Service) readPrompts(fileName string) ([]Prompt, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%s': %v", fileName, err)
	}
	if len(data) == 0 {
		return []Prompt{}, nil
	}
	var prompts []Prompt
	err = json.Unmarshal(data, &prompts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %v", err)
	}
	return prompts, nil
}

func (s *Service) writePrompts(prompts []Prompt, fileName string) error {
	data, err := json.Marshal(prompts)
	if err != nil {
		return fmt.Errorf("failed to marshal prompts: %v", err)
	}
	err = ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file '%s': %v", fileName, err)
	}
	return nil
}

func (s *Service) downloadFile(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.Status)
		panic("Failed to download file")
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy file: %v", err)
	}
	return nil
}

func (s *Service) check(e error) {
	if e != nil {
		s.wf.FatalError(e)
		os.Exit(1)
	}
}

func (s *Service) readCSV(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file '%s': %v", filePath, err)
	}
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv: %v", err)
	}
	return records, nil
}

func (s *Service) writeJsonFile(raw [][]string, fileName string) error {
	data := make([]Prompt, len(raw))
	for i, r := range raw {
		data[i] = Prompt{Title: r[0], Subtitle: r[1]}
	}
	err := s.writePrompts(data, fileName)
	if err != nil {
		return fmt.Errorf("failed to write json file: %v", err)
	}
	return nil
}

func (s *Service) deleteFile(fileName string) error {
	err := os.Remove(fileName)
	if err != nil {
		return fmt.Errorf("failed to delete file '%s': %v", fileName, err)
	}
	return nil
}
