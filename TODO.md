# Java Integration Test Follow-ups

1. **Integrate `jvm-vendors` into the `java` composite**
   - Replace `bellsoft-liberica` with `jvm-vendors` while preserving BellSoft
     as the default vendor, or add vendor alternatives to the order.
   - Remove the direct Adoptium/Amazon Corretto buildpack arguments and image
     pre-pulls from `integration/init_test.go`.
   - Switch the vendor tests to `BP_JVM_VENDOR=adoptium` and
     `BP_JVM_VENDOR=corretto`.

2. **Split the quick and long test sets**
   - Keep a fast, high-value set on pull requests.
   - Add a scheduled/cron workflow for the longer-running tests.
   - Add a test-set selection mechanism, for example
     `INTEGRATION_TEST_SET=quick|long`.

3. **Scatter stack coverage**
   - Add noble and/or resolute builders to `integration.json`.
   - Register only the tests assigned to each stack in
     `integration/init_test.go`.
   - Coordinate this with the quick/long split so the PR workflow does not run
     every test on every stack.

4. **Optional Maven + npm sample**
   - The Java+Node test currently uses `java/java-node/gradle-node` (npm).
   - Add a `java/java-node/maven-npm` sample only if Maven+Node coverage is
     needed.

5. **Optional long-tail coverage**
   - Kotlin, Clojure (`deps`, Leiningen, tools-build), AspectJ, OpenTelemetry,
     and Azure Application Insights.
   - Only if priorities change; the current suite is intentionally focused on
     the core supported paths.

## Notes

- Local integration runs need a JDK compatible with the sample Gradle wrapper
  for the host-side WAR precompile. The CI runner provides a compatible JDK.
