package io.paketo.demo;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class IntegrationTestsController {

	@GetMapping("/runtime")
	public String index() {
		StringBuilder sb = new StringBuilder();
		System.getProperties().forEach((key, value) -> {
			sb.append(key + "=" + value + "\n");
		});

		return sb.toString();
	}
}
