DIRECTORIES = $(wildcard */)

test: $(patsubst %,test-%, $(DIRECTORIES))

test-%:
	$(MAKE) -C $* test-example