package terraform

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/dougkirkley/trex/config"
)

var tplFile = `
package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Hello struct {
    Response string 
}

// The input type and the output type are defined by the API Gateway.
func handleRequest(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var hello Hello
	err := json.Unmarshal([]byte(event.Body), &hello)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{Body: hello.Response, StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}`

var gomod = `
module lambda

go 1.14

require (
	github.com/aws/aws-lambda-go v1.19.0 // indirect
)
`

var terraformtmpl = `
variable "profile" {
	type = string
	default = "default"
}
  
variable "region" {
	type = string
	default = "us-east-1"
}

variable "module_provider" {
	type = object({
	  version = string
	  region  = string
	  profile = string
	})
	default = null
}
  
variable "api_name" {
    type = string
}
  
variable "functions" {
	type = set(object({
	  root_path      = string
	  path           = string
	  method         = string
	  authorization  = string
	  lambda_name    = string
	  lambda_runtime = string
	  iam_policy     = string
	}))
}

provider "{{ .Provider.Name }}" {
	version = "{{ .Provider.Version }}"
	region  = "{{ .Provider.Region }}"
	profile = "{{ .Provider.Profile }}"
}

module "{{ .Provider.Name }}" {
  source = "git::https://gitlab.com/dougkirkley/trex.git//terraform/modules/{{ .Provider.Name }}"
  region = var.region
  api_name = var.api_name
  functions = var.functions
}

output "test" {
	value = module.aws.test
  }
`

// CreateGoFiles creates go files
func CreateGoFiles(path string) {
	filePath, err := filepath.Abs(path)
	if err != nil {
		log.Print(err.Error())
	}
	f, err := os.Create(filePath + "/main.go")
	if err != nil {
		log.Print(err.Error())
	}
	_, err = f.WriteString(tplFile)
	if err != nil {
		log.Print(err.Error())
	}
	f, err = os.Create(filePath + "/go.mod")
	if err != nil {
		log.Print(err.Error())
	}
	_, err = f.WriteString(gomod)
	if err != nil {
		log.Print(err.Error())
	}
}

// CreateMainFile creates a terraform main.tf file from template
func CreateMainFile(cfg config.TrexConfig) ([]byte, error) {
	tmpl, err := template.New("test").Parse(terraformtmpl)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer

	if err = tmpl.Execute(&buffer, cfg); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
