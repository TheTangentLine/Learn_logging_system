.PHONY: test check

test:
	$(MAKE) -C services/api-server test
	$(MAKE) -C services/relay-worker test
	$(MAKE) -C services/rmq-consumer test

check:
	$(MAKE) -C services/api-server check
	$(MAKE) -C services/relay-worker check
	$(MAKE) -C services/rmq-consumer check
