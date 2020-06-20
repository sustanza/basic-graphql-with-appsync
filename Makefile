#############################################################################################
#								Build Configuration     									#
#############################################################################################

# the project name used throughout the process to define created resources
PROJECT_NAME ?= basic-graphql-with-appsync

# AWS S# bucket created to store the project artifacts
AWS_BUCKET_NAME ?= $(PROJECT_NAME)-artifacts

# AWS S3 bucket created to store the Cloudformation stack artifacts
AWS_STACK_NAME ?= $(PROJECT_NAME)-stack

# AWS region to deploy resources to
AWS_REGION ?= us-east-1

# Cloudformation template for project
FILE_TEMPLATE = template.yml

# Cloudformation template package name to be generated from 
FILE_PACKAGE = package.yml

# operating system target for go build https://golang.org/pkg/runtime/
GOOS ?= linux

# dir for funcs used by build process
PATH_FUNCTIONS := ./functions/

# used to dive into the above funcs dir to get all funcs
LIST_FUNCTIONS := $(subst $(PATH_FUNCTIONS),,$(wildcard $(PATH_FUNCTIONS)*))

#############################################################################################
#								Build Job / Task Definitions
#############################################################################################

# cleans the distrubution directory
clean:
	@ rm -rdf dist/

# installs go mods
install:
	@ go mod download

# runs go tests
test:
	@ go test ./... -v

# instructions for a standard build of a go func
build-%:
	@ env GOOS=linux \
		go build \
		-gcflags "all=-N -l"  \
		-o ./dist/$*/handler ./functions/$*

# instructions for a debug build of a go func
debug-%:
	@ env GOOS=linux \
		go build \
		-a -installsuffix cgo -ldflags="-s -w " \
		-o ./dist/$*/handler ./functions/$*

# build step
build:
	$(info Building: $(LIST_FUNCTIONS))
	@ $(MAKE) clean
	@ $(MAKE) $(foreach FUNCTION,$(LIST_FUNCTIONS),build-$(FUNCTION))

# debug step
debug:
	$(info Building: $(LIST_FUNCTIONS))
	@ $(MAKE) clean
	@ $(MAKE) $(foreach FUNCTION,$(LIST_FUNCTIONS),build-$(FUNCTION))
	@ $(MAKE) dlv

# generates .env file for local development
env:
	echo \
	"COFFEE_TABLE_NAME=${PROJECT_NAME}-coffee\n"\
	"AWS_REGION=${AWS_REGION}\n"\
	"AWS_SDK_LOAD_CONFIG=1\n"\
	> .env

# helper to grab dlv for local debugging
dlv:
	$(info Building Task: Dlv)
	@ env GOARCH=amd64 GOOS=linux go build -o dist/dlv github.com/go-delve/delve/cmd/dlv

# creates s3 buckets for project artifacts
s3:
	@ aws s3 mb s3://$(AWS_BUCKET_NAME) --region $(AWS_REGION)
	@ aws s3 mb s3://$(AWS_STACK_NAME) --region $(AWS_REGION)

# packages the resources for cloudformation deploy step
package:
	@ aws cloudformation package \
		--template-file $(FILE_TEMPLATE) \
		--s3-bucket $(AWS_BUCKET_NAME) \
		--region $(AWS_REGION) \
		--output-template-file $(FILE_PACKAGE)

# deploys projects cloudformation resources
deploy:
	@ aws cloudformation deploy \
		--template-file $(FILE_PACKAGE) \
		--region $(AWS_REGION) \
		--capabilities CAPABILITY_NAMED_IAM \
		--stack-name $(AWS_STACK_NAME) \
		--force-upload \
		--s3-bucket $(AWS_BUCKET_NAME) \
		--parameter-overrides \
			ProjectName=$(PROJECT_NAME) \

# describes the cloudformation stack
describe:
	@ aws cloudformation describe-stacks \
		--region $(AWS_REGION) \
		--stack-name $(AWS_STACK_NAME)

# cleans up resources not associated with cloudformation stack, ex S3
cleanup:
	@ aws s3 rb s3://$(AWS_BUCKET_NAME) --region $(AWS_REGION) --force
	@ aws s3 rb s3://$(AWS_STACK_NAME) --region $(AWS_REGION) --force

# describes output of stack
outputs:
	@ make describe \
		| jq -r '.Stacks[0].Outputs'

.PHONY: clean install test build build-% debug debug-% env dlv s3 package deploy describe cleanup output
