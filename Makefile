title:
	@echo "********************************************************************************************************************************************"

lint: title
	golangci-lint version
	golangci-lint run -v --color always --out-format colored-line-number
