#############################################################################################
#								Build Configuration     									#
#############################################################################################
PROJECT_NAME ?= basic-graphql-with-appsync-test

AWS_BUCKET_NAME ?= $(PROJECT_NAME)-artifacts
AWS_STACK_NAME ?= $(PROJECT_NAME)-stack
AWS_REGION ?= us-east-1
GOOS ?= linux
FILE_TEMPLATE = template.yml
FILE_PACKAGE = package.yml

PATH_FUNCTIONS := ./functions/
LIST_FUNCTIONS := $(subst $(PATH_FUNCTIONS),,$(wildcard $(PATH_FUNCTIONS)*))

#############################################################################################
#								Build Job / Task Definitions
#############################################################################################

clean:
	@ rm -rdf dist/

install:
	@ go mod download

test:
	@ go test ./... -v

build-%:
	@ env GOOS=linux \
		go build \
		-gcflags "all=-N -l"  \
		-o ./dist/$*/handler ./functions/$*

debug-%:
	@ env GOOS=linux \
		go build \
		-a -installsuffix cgo -ldflags="-s -w " \
		-o ./dist/$*/handler ./functions/$*

build:
	$(info Building: $(LIST_FUNCTIONS))
	@ $(MAKE) clean
	@ $(MAKE) $(foreach FUNCTION,$(LIST_FUNCTIONS),build-$(FUNCTION))

debug:
	$(info Building: $(LIST_FUNCTIONS))
	@ $(MAKE) clean
	@ $(MAKE) $(foreach FUNCTION,$(LIST_FUNCTIONS),build-$(FUNCTION))
	@ $(MAKE) dlv

dlv:
	$(info Building Task: Dlv)
	@ env GOARCH=amd64 GOOS=linux go build -o dist/dlv github.com/go-delve/delve/cmd/dlv

configure:
	@ aws s3 mb s3://$(AWS_BUCKET_NAME) --region $(AWS_REGION)
	@ aws s3 mb s3://$(AWS_STACK_NAME) --region $(AWS_REGION)

package:
	@ aws cloudformation package \
		--template-file $(FILE_TEMPLATE) \
		--s3-bucket $(AWS_BUCKET_NAME) \
		--region $(AWS_REGION) \
		--output-template-file $(FILE_PACKAGE)

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

describe:
	@ aws cloudformation describe-stacks \
		--region $(AWS_REGION) \
		--stack-name $(AWS_STACK_NAME)

outputs:
	@ make describe \
		| jq -r '.Stacks[0].Outputs'

.PHONY: clean install build build-lambda configure package deploy describe output
