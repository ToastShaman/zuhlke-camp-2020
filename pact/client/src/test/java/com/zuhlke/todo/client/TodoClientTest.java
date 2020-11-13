package com.zuhlke.todo.client;

import com.github.tomakehurst.wiremock.WireMockServer;
import com.github.tomakehurst.wiremock.client.WireMock;
import com.zuhlke.todo.client.http.HttpTodoClient;
import com.zuhlke.todo.client.model.CreateTodoRequest;
import com.zuhlke.todo.client.model.CreateTodoResponse;
import com.zuhlke.todo.client.model.Todo;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.util.List;

import static com.github.tomakehurst.wiremock.client.WireMock.*;
import static com.github.tomakehurst.wiremock.core.WireMockConfiguration.options;
import static org.assertj.core.api.Assertions.assertThat;

class TodoClientTest {

    private WireMockServer server;

    @BeforeEach
    public void setup() {
        server = new WireMockServer(options().dynamicPort());
        server.start();
        WireMock.configureFor(server.port());
    }

    @AfterEach
    public void teardown() {
        server.stop();
    }

    private String formatWithPort(String host) {
        return host.formatted(server.port());
    }

    @Test
    @DisplayName("can create a new createTodoRequest")
    void can_create_todo() {
        stubFor(post("/todo")
                .withHeader("Accept", equalTo("application/json"))
                .withHeader("Content-Type", equalTo("application/json"))
                .withRequestBody(equalToJson("""
                        {
                            "text": "Don't forget the milk",
                            "status": "TODO",
                            "category": "shopping",
                            "tags": [
                                "groceries",
                                "food"
                            ]
                        }
                        """))
                .willReturn(aResponse()
                        .withStatus(201)
                        .withHeader("Content-Type", "application/json")
                        .withBody("""
                                {
                                    "id": "-MLqrG6LkLkkKc1iMLBt",
                                    "rev": "-MLivp1BrS59mMbSN7Jr",
                                    "text": "Don't forget the milk",
                                    "status": "TODO",
                                    "category": "shopping",
                                    "tags": [
                                        "groceries",
                                        "food"
                                    ]
                                }
                                """)));

        HttpTodoClient client = TodoClient.builder()
                .setHost(formatWithPort("http://localhost:%d"))
                .build();

        CreateTodoRequest createTodoRequest = CreateTodoRequest.builder()
                .withText("Don't forget the milk")
                .withStatus("TODO")
                .withCategory("shopping")
                .withTags(List.of("groceries", "food"))
                .build();

        CreateTodoResponse response = client.create(createTodoRequest);

        Todo actual = response.getTodo();

        Todo expected = Todo.builder()
                .withId("-MLqrG6LkLkkKc1iMLBt")
                .withRev("-MLivp1BrS59mMbSN7Jr")
                .withText("Don't forget the milk")
                .withStatus("TODO")
                .withCategory("shopping")
                .withTags(List.of("groceries", "food"))
                .build();

        assertThat(actual).isEqualTo(expected);
    }
}