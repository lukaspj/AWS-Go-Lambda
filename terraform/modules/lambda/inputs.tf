variable "lambdas" {
  type = map(object({route=string}))
}

variable "archive" {
  type = string
}

variable "source_dir" {
  type = string
}

variable "lambda_name" {
  type = string
}