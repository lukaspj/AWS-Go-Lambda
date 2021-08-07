terraform {
  backend "s3" {
    bucket = "bryrupteaterterraform"
    key = "sample.tfstate"
    region = "eu-north-1"
  }

  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}

provider "aws" {
  region = "eu-north-1"
}

module "testtype_lambda" {
  source = "./modules/lambda"

  archive = "test_lambda.zip"
  source_dir = "${path.root}/../out"
  lambdas = {
    get_items = {
      route = "GET /"
    }

    get_item = {
      route = "GET /{id}"
    }

    create_item = {
      route = "POST /"
    }
  }


  lambda_name = "sample_lambda"
}