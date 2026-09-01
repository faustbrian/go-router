.PHONY: conformance docs integration

conformance:
	./scripts/check-conformance.sh

docs:
	./scripts/check-docs.sh

integration:
	./scripts/check-integrations.sh
