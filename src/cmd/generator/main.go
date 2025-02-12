package main

import (
	"encoding/json"
	"log"
	"os"
	"text/template"

	"github.com/italiaremote/awesome-italia-remote/pkg/cmp"
)

const homePath = "../../../"
const dataPath = homePath + "data/"

type TemplateParameters struct {
	Len       int
	Companies map[string][]cmp.Company
}

func main() {
	companies := TemplateParameters{}
	companies.Companies = make(map[string][]cmp.Company)
	companiesFile, err := os.ReadDir(dataPath)
	if err != nil {
		log.Fatalln(err)
	}
	companies.Len = len(companiesFile)
	var allCompanies []cmp.Company

	for _, companyFiles := range companiesFile {
		var company cmp.Company
		file, err := os.ReadFile(dataPath + companyFiles.Name())
		if err != nil {
			log.Println(companyFiles.Name())
			log.Fatalln(err)
		}

		err = json.Unmarshal(file, &company)
		if err != nil {
			log.Println(companyFiles.Name())
			log.Fatalln(err)
		}
		err = company.Validate()
		if err != nil {
			log.Println(companyFiles.Name())
			log.Fatalln(err)
		}
		company.Fix()

		allCompanies = append(allCompanies, company)
		for _, category := range company.Categories {
			companies.Companies[category] = append(companies.Companies[category], company)
		}
	}

	templ, err := template.ParseFiles("./template.md")
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.Create(homePath + "README.md")
	if err != nil {
		log.Fatalln(err)
	}

	err = templ.Execute(f, companies)
	if err != nil {
		log.Fatalln(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalln(err)
	}

	jsonByte, err := json.MarshalIndent(allCompanies, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile(homePath+"outputs.json", jsonByte, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
