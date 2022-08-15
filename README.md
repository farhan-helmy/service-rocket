
## Pre requisites 

> - Before running the app, make sure that you have aws cli installed in your machine [AWS-CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)

> - Create 1 IAM user that have permission S3fullAccess

> - Type aws configure on your local machine terminal and put the access key and secret key of the IAM user that you have just created

## Step by step

```shell
git clone git@github.com:farhan-helmy/service-rocket.git

cd service-rocket
```
```go 
go mod tidy
go run main.go
```
2 endpoints will be generated 
localhost:3001/api/v1/image/multiple-upload
localhost:3001/api/v1/image/upload

For me I use insomnia to test the end point 

Testing using CURL

```shell

#Endpoint for upload zip file
curl --request POST \
  --url http://localhost:3001/api/v1/image/multiple-upload \
  --header 'Content-Type: multipart/form-data' \
  --header 'content-type: multipart/form-data; boundary=---011000010111000001101001' \
  --form 'file="filepath to zip file"'

#Endpoint for upload single image
  curl --request POST \
  --url http://localhost:3001/api/v1/image/upload \
  --header 'Content-Type: multipart/form-data' \
  --header 'content-type: multipart/form-data; boundary=---011000010111000001101001' \
  --form image=@/Users/farhanhelmy/Downloads/download.png
```
## TODOS

- [x] upload service 
- [x] unzip file 
- [x] upload multiple image 
- [x] check if zip file -> go to unzip function 
- [x] check if unzipped file png or jpeg? send to function that generate .png extension or .jpeg 
- [x] resize image 
- [x] generate two type of resized image 
- [x] clear folder after upload to S3 
- [x] connect to s3 account 
- [x] AWS SDK to upload data to S3 


## Challenges

- Lack of knowledge in Go
- Pointer and memory segment error (Solved by changing code writing)
- Return value, confuse to use array or struct
- Reduce time and space complexity challenge
- AWS SDK s3 when upload object dont have public read access (Already solved, need to modify IAM JSON)

## Please describe in the README the difference for the solution to handle average daily call count 5, 5K and 100K.

- First is in term of architecture, we will containerize this application to make it scalable, can be deployed anywhere.
- The plan to hold a large request is to deploy it inside Kubernetes cluster so that auto replica can be done based on treshold that we set inside deployment configuration in the K8s cluster.
- Also planning to use EKS (AWS), GKE(Google Cloud) or AKS(Azure) to deploy the Kubernetes engine easily, and also make use of other service provided by the Cloud provider such as Load Balancer, Content Delivery Network, also storage service 
- In terms of code, this application uses Go lang, in the future we can also utilize the use of Concurrency in Go to handle such a large request
