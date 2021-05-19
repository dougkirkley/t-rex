variable "api_name" {
  type = string
}

variable "profile" {
  type = string
  default = "default"
}

variable "region" {
  type = string
  default = "us-east-1"
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

variable "module_provider" {
  type = object({
    name    = string
    version = string
    region  = string
  })
  default = null
}