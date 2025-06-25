VERSION?=$(shell git describe --tags  --abbrev=0 | cut -c2-  )
export BP_UNDER_TEST ?= "paketobuildpacks/java:$(VERSION)"

clean:
	rm -fr target
	rm -f *.cnb
	mvn clean -f integration/pom.xml

create-package:
	create-package --include-dependencies --destination ./target --version "${VERSION}"
	cp package.toml ./target/package.toml
	echo "[buildpack]" >> ./target/package.toml
	echo "uri = \"./\"" >> ./target/package.toml
	echo "" >> ./target/package.toml
	echo "[platform]" >> ./target/package.toml
	echo "os = \"linux\"" >> ./target/package.toml

package: create-package
	pack buildpack package ${BP_UNDER_TEST} --format=image --config ./target/package.toml

integration: samples
	TESTCONTAINERS_RYUK_DISABLED=true TESTCONTAINERS_RYUK_CONTAINER_PRIVILEGED=true go test -v -count=1 -timeout=20m ./integration

samples:
	test -d integration/samples && git -C integration/samples pull || git clone https://github.com/paketo-buildpacks/samples integration/samples

.PHONY: integration pre-integration package create-package clean samples
