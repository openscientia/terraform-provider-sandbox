default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

gen: 
	rm -f .github/labeler-pr-labels.yml
	rm -f infrastructure/repository/labels-product.tf
	go generate ./...