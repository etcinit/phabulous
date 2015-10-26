all:
	gox --os="linux darwin" --output="build/{{.Dir}}_{{.OS}}_{{.Arch}}" ./cmd/phabulous

clean:
	rm -r build
