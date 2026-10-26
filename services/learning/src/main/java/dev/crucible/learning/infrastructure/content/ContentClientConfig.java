package dev.crucible.learning.infrastructure.content;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClientManager;
import org.springframework.security.oauth2.client.web.client.support.OAuth2RestClientHttpServiceGroupConfigurer;
import org.springframework.web.client.support.RestClientHttpServiceGroupConfigurer;
import org.springframework.web.service.registry.ImportHttpServices;

/**
 * Client gọi content-service sinh từ contracts/content-api.yaml.
 * Token client credentials được gắn tự động theo {@code @ClientRegistrationId} trên interface sinh ra.
 */
@Configuration
@ImportHttpServices(group = "content", basePackages = "dev.crucible.learning.infrastructure.content.generated")
class ContentClientConfig {

    @Bean
    RestClientHttpServiceGroupConfigurer contentBaseUrl(@Value("${crucible.services.content-service}") String baseUrl) {
        return groups -> groups.filterByName("content")
                .forEachClient((group, builder) -> builder.baseUrl(baseUrl));
    }

    @Bean
    OAuth2RestClientHttpServiceGroupConfigurer contentOAuth2(OAuth2AuthorizedClientManager manager) {
        return OAuth2RestClientHttpServiceGroupConfigurer.from(manager);
    }
}
