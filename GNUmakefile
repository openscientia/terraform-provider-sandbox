default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

gen: 
	rm -f .github/labeler-pr-labels.yml
	rm -f .github/labeler-issue-labels.yml
	rm -f infrastructure/repository/labels-resource.tf
	go generate ./...