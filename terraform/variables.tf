variable "instance_name" {
  default = "hmh-demo-example-instance"
  type = string 
}

variable "bucket_name" {
  default="hmh-demo-example-bucket"
  type = string
}

variable "project_id" {
  default = "caseytin-sandbox"
  type = string
}

variable "region" {
  default = "us-east1"
  type = string
}

variable "zone" {
  default = "us-east1-b"
  type = string
}

variable "machine_type" {
  default = "f1-micro"
  type = string
}
