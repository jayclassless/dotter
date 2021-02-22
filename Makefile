test::
	@ginkgo -r -cover -coverprofile=coverage.out

coverage:
	@find . -name coverage.out | xargs gocovmerge > combined_coverage.out
	@go tool cover -html=combined_coverage.out
