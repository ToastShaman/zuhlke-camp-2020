package com.zuhlke.todo.client;

import au.com.dius.pact.consumer.MockServer;
import au.com.dius.pact.consumer.dsl.PactDslWithProvider;
import au.com.dius.pact.consumer.junit5.PactConsumerTestExt;
import au.com.dius.pact.consumer.junit5.PactTestFor;
import au.com.dius.pact.core.model.RequestResponsePact;
import au.com.dius.pact.core.model.annotations.Pact;
import com.zuhlke.todo.client.http.HttpTodoClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;

import java.util.List;

import static org.assertj.core.api.AssertionsForClassTypes.assertThat;

//@formatter:off
@ExtendWith(PactConsumerTestExt.class)
@PactTestFor(providerName = "TodoApi", port = "1234")
public class TodoApiPactConsumerTest {

    @BeforeEach
    public void setUp(MockServer mockServer) {
        assertThat(mockServer).isNotNull();
    }

    @Pact(consumer = "TodoClient")
    public RequestResponsePact createTodo(PactDslWithProvider builder) {
        return builder
                .given("An empty repository")
                .uponReceiving("a new todo")
                    .method("POST")
                    .path("/todo")
                    .body("""
                        {
                            "text": "Don't forget the milk",
                            "status": "TODO",
                            "category": "shopping",
                            "tags": [
                                "groceries",
                                "food"
                            ]
                        }
                        """, "application/json")
                    .matchHeader("Content-Type", "application/json", "application/json")
                .willRespondWith()
                    .status(201)
                    .body("""
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
                            """)
                .toPact();
    }

    @Test
    @PactTestFor(pactMethod = "createTodo")
    void testCreateTodo(MockServer mockServer)  {
        HttpTodoClient todoClient = TodoClient.builder()
                .setHost("http://localhost:%d".formatted(1234))
                .build();

        CreateTodoRequest createTodoRequest = CreateTodoRequest.builder()
                .withText("Don't forget the milk")
                .withStatus("TODO")
                .withCategory("shopping")
                .withTags(List.of("groceries", "food"))
                .build();

        CreateTodoResponse response = todoClient.create(createTodoRequest);

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
