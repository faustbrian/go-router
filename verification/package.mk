.PHONY: conformance docs integration safety

conformance:
	./scripts/check-conformance.sh

docs:
	./scripts/check-docs.sh

integration:
	./scripts/check-integrations.sh

safety:
	./scripts/check-safety.sh
