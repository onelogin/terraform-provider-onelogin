.PHONY: clean build ti tp ta

PKG_NAME=onelogin
WEBSITE_REPO=github.com/hashicorp/terraform-website
DIST_DIR=./dist
BIN_NAME=terraform-provider-onelogin
BIN_PATH=${DIST_DIR}/${BIN_NAME}

GO111MODULE=on

PLUGINS_DIR=~/.terraform.d/plugins
PLUGIN_PATH=onelogin.com/onelogin/onelogin
VERSION=0.8.1

clean:
	rm -r ${DIST_DIR}
	rm -r ${PLUGINS_DIR}

clean-terraform:
	rm terraform.*
	rm .terraform.lock.hcl
	rm -r .terraform/

build:
	mkdir -p ${DIST_DIR}
	go build -o ${DIST_DIR} ./...

sideload: build
	# Terraform v0.12.x
	mkdir -p ${PLUGINS_DIR}
	cp ${BIN_PATH} ${PLUGINS_DIR}/${BIN_NAME}
	# Terraform >= v0.13.x
	# macOS
	mkdir -p ${PLUGINS_DIR}/${PLUGIN_PATH}/${VERSION}/darwin_amd64
	cp ${BIN_PATH} ${PLUGINS_DIR}/${PLUGIN_PATH}/${VERSION}/darwin_amd64/${BIN_NAME}
	# macOS Apple Silicon
	mkdir -p ${PLUGINS_DIR}/${PLUGIN_PATH}/${VERSION}/darwin_arm64
	cp ${BIN_PATH} ${PLUGINS_DIR}/${PLUGIN_PATH}/${VERSION}/darwin_arm64/${BIN_NAME}
	# Linux
	mkdir -p ${PLUGINS_DIR}/${PLUGIN_PATH}/${VERSION}/linux_amd64
	cp ${BIN_PATH} ${PLUGINS_DIR}/${PLUGIN_PATH}/${VERSION}/linux_amd64/${BIN_NAME}

testacc:
	TF_ACC=1 go test ./... -v -timeout 120m

ti:
	terraform init

tp:
	terraform plan

ta:
	terraform apply -auto-approve

test:
	go test -v -count=1 -short ./...

secure:
	# excludes G104 (unhandled go errors) - Approved by security team
	# excludes G109 (potential integer rollover) - These function calls were recommended by hashicorp for developing a provider
	# excludes G304 (potential file inclusion) - The file needs to be read in order to run acceptance tests by hashicorp
	# excludes G401 (Use of weak cryptographic primitive) - md5 is perfectly enough to create a unique data source ID
	# excludes G501 (Blocklisted import crypto/md5: weak cryptographic primitive) - md5 is perfectly enough to create a unique data source ID
	gosec -exclude=G104,G304,G109,G401,G501 ./...
