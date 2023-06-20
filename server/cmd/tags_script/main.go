package main

import (
	quizApi "github.com/dimayasha7123/quiz_service/server/internal/quiz_party_api_client"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
	"time"
)

type yamlData struct {
	Tags []string `yaml:"tags"`
}

const filename = "tags.yaml"

func main() {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, os.ModeAppend)
	if err != nil {
		log.Fatalf("can't open file: %v", err)
	}
	bytes, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("can't read file: %v", err)
	}

	var oldData yamlData
	err = yaml.Unmarshal(bytes, &oldData)
	if err != nil {
		log.Fatalf("can't unmarshal file: %v", err)
	}

	tagMap := make(map[string]struct{})
	for _, tag := range oldData.Tags {
		tagMap[tag] = struct{}{}
	}

	qclient := quizApi.New("cCP7GHWCoe2d4VBysvmaAT1x1hToB54mm1rdEiYK")
	for i := 0; i < 500; i++ {
		party, err := qclient.GetPartyData("")
		if err != nil {
			log.Fatalf("can't get party: %v", err)
		}
		for _, question := range party {
			for _, tag := range question.Tags {
				tagMap[tag.Name] = struct{}{}
			}
		}
		log.Println(i)
		time.Sleep(time.Millisecond * 500)
	}

	tags := make([]string, 0, len(tagMap))
	for tag := range tagMap {
		tags = append(tags, tag)
	}
	data := yamlData{Tags: tags}

	yamlFile, err := yaml.Marshal(&data)
	if err != nil {
		log.Fatalf("can't marshal data: %v", err)
	}

	log.Println(string(yamlFile))

	f, err = os.Create("tags.yaml")
	if err != nil {
		log.Fatalf("can't create file: %v", err)
	}
	defer f.Close()

	_, err = io.WriteString(f, string(yamlFile))
	if err != nil {
		log.Fatalf("can't write to file: %v", err)
	}
}
