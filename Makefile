clean :
	if [ -f "terraform-provider-onelogin" ]; then rm terraform-provider-onelogin; fi

clean-terraform :
	rm terraform.*

compile :
	make clean && cd ./cmd && go build -o terraform-provider-onelogin && mv terraform-provider-onelogin ../.

ti:
	terraform init

tp:
	terraform plan

ta:
	terraform apply
