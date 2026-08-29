.PHONY: conformance interoperability

conformance:
	./scripts/check-conformance.sh

interoperability:
	./scripts/check-integrations.sh
