
## Step by step

```go 
go mod tidy
```
## TODOS

- upload service [x]
- unzip file [x]
- upload multiple image [x]
- check if zip file -> go to unzip function [x]
- check if unzipped file png or jpeg? send to function that generate .png extension or .jpeg [x]
- resize image [x]
- generate two type of resized image [x]
- clear folder after upload to S3 [x]
- connect to s3 account [x]
- AWS SDK to upload data to S3 [x]


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
