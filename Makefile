all:
	goxc go-test xc archive rmbin

clean:
	rm -rf dist
