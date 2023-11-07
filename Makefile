# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test

# Mockery parameters
MOCKERYCMD = mockery
MOCKERYFLAGS = --all --dir=. --recursive --inpackage --inpackage-suffix --case=snake

# Generate mocks using mockery
generate_mocks:
	$(MOCKERYCMD) $(MOCKERYFLAGS)

# Run tests
test:
	$(GOTEST) -v ./...


.PHONY: generate_mocks test
