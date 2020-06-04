
.PHONY: clean build ti tp ta

PKG_NAME=onelogin
WEBSITE_REPO=github.com/hashicorp/terraform-website
DIST_DIR=./dist
BIN_NAME=terraform-provider-onelogin
BIN_PATH=${DIST_DIR}/${BIN_NAME}

GO111MODULE=on

PLUGINS_DIR=~/.terraform.d/plugins

clean:
	rm -r ${DIST_DIR}
	rm -r ${PLUGINS_DIR}

clean-terraform:
	rm terraform.*

build:
	mkdir -p ${DIST_DIR}
	go build -o ${DIST_DIR} ./...
	mv ${DIST_DIR}/cmd ${BIN_PATH}

sideload: build
	mkdir -p ${PLUGINS_DIR}
	cp ${BIN_PATH} ${PLUGINS_DIR}/${BIN_NAME}

testacc:
	TF_ACC=1 go test ./... -v -timeout 120m

ti:
	terraform init

tp:
	terraform plan

ta:
	terraform apply -auto-approve

test:
	go get -u github.com/dcaponi/gopherbadger@v2.2.1
	gopherbadger -root="./ol_schema" -md="readme.md" -png=false

secure:
	# or install it into ./bin/
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s
	./bin/gosec -exclude=G104 ./...
