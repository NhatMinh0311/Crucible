package dev.crucible.learning.infrastructure.runtime;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClientManager;
import org.springframework.security.oauth2.client.web.client.support.OAuth2RestClientHttpServiceGroupConfigurer;
import org.springframework.web.client.support.RestClientHttpServiceGroupConfigurer;
import org.springframework.web.service.registry.ImportHttpServices;

/**
 * Client gọi runtime-api sinh từ contracts/runtime-api.yaml.
 * Token client credentials được gắn tự động theo {@code @ClientRegistrationId} trên interface sinh ra.
 */
@Configuration
@ImportHttpServices(group = "runtime", basePackages = "dev.crucible.learning.infrastructure.runtime.generated")
class RuntimeClientConfig {

    @Bean
    RestClientHttpServiceGroupConfigurer runtimeBaseUrl(@Value("${crucible.services.runtime-api}") String baseUrl) {
        return groups -> groups.filterByName("runtime")
                .forEachClient((group, builder) -> builder.baseUrl(baseUrl));
    }

    @Bean
    OAuth2RestClientHttpServiceGroupConfigurer runtimeOAuth2(OAuth2AuthorizedClientManager manager) {
        return OAuth2RestClientHttpServiceGroupConfigurer.from(manager);
    }
}
