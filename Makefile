VERSION?=$(shell git describe --tags --dirty=-dirty)
BP_UNDER_TEST?="../java-buildpack.cnb"

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
	pack buildpack package ${BP_UNDER_TEST} --format=file --config ./target/package.toml

integration:
	BP_UNDER_TEST=${BP_UNDER_TEST} INCLUDE_INTEGRATION_TESTS=true go test -v -count=1 -timeout=20m ./integration --ginkgo.label-filter integration --ginkgo.v

.PHONY: integration pre-integration package create-package clean
 