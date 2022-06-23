# zuhlke-camp-2020

Consumer Driven Contract Testing demo using [Pact](https://docs.pact.io/)

## Prerequisites

1. Find a pair to do the exercise with

1. If you are doing this exercise remotely:

    1. Each of you needs to install [Tuple][4] and register with your **Zuhlke email address**!

    1. Ask Kevin for the magic team link

    1. Start a call

1. Read the [Etiquette for Pair Programming][1]

1. Choose a technology for the *Consumer* you are going to write

1. Choose a technology for the *Provider* you are going to write

## Instructions

You'll be practising writing API contracts by using [PACT][2] in the language/tech-stack of your choice.

1. Start implementing the *Consumer* **incrementally** by using [PACT][2] and once you have a working test/contract for one of the specifications share the contract with the *Provider*.

1. Run the contract you just been given from the *Consumer* against the *Provider* and try to make it pass.

1. Repeat the process as you implement more and more of the specification.

## Functionality

Inspired by [TODO-Backend API][3], [todomvc][5] and [TodoMVC Scenarios][6].

### Scenario: Empty list can have item added

        Given an empty Todo list
        When I add a Todo for ‘Buy cheese’
        Then only that item is listed
        And the list summary is ‘1 item left’
        And the list’s filter is set to ‘All’ with ‘Completed’ & ‘Active’ unset
        And ‘Clear completed’ is unavailable

### Scenario: Empty list can have two items added

        Given an empty Todo list
        When I add Todos for ‘Buy cheese’ & ‘Wash the car’
        Then only those items are listed
        And the list summary is ‘2 items left’
        And the list’s filter is set to ‘All’ with ‘Completed’ & ‘Active’ unset
        And ‘Clear completed’ is unavailable

### Scenario: Item completion changes the list

        Given a Todo list with items ‘File taxes’ & ‘Walk the dog’
        When the first item is marked as complete
        Then only the second item is listed as active
        And the list summary is ‘1 item left’
        And ‘Clear completed’ is available

### Scenario: Completed items should not be visible in active filter

        Given a Todo list with items ‘File taxes’ & ‘Walk the dog’
        And the first item is completed
        When the filter is set to ‘Active’
        Then only the second item is listed
        And the list summary is ‘1 item left’
        And ‘Clear completed’ is available
        And the route is /active

### Scenario: All completed items should not be visible in active filter

        Given a Todo list with items ‘File taxes’ & ‘Walk the dog’
        And the first item is completed
        And the second item is completed
        When the filter is set to ‘Active’
        Then nothing is listed
        And the list summary is ‘0 items left’
        And ‘Clear completed’ is available
        And the route is /active

### Scenario: Uncompleted items should not be visible in the completed filter

        Given a Todo list with items ‘File taxes’ & ‘Walk the dog’
        And the first item is completed
        When the filter is set to ‘Completed’
        Then only the first item is listed
        And the list summary is ‘1 item left’
        And ‘Clear completed’ is available
        And the route is /completed

### Scenario: Items can be cleared

        Given a Todo list with items ‘File taxes’ & ‘Walk the dog’
        And the first item is completed
        When ‘Clear completed’ is clicked
        Then only the second item is listed
        And the list summary is ‘1 item left’

### Scenario: clear all works when none completed

        Given a Todo list with items ‘File taxes’ & ‘Walk the dog’
        When the clear all items affordance is toggled
        Then both items are listed
        But nothing listed is active
        And the list summary is ‘0 items left’

### Scenario: clear all works when one of two completed

        Given a Todo list with items ‘File taxes’ & ‘Walk the dog’
        And the first item is completed
        When the clear all items affordance is toggled
        Then both items are listed
        But nothing listed is active
        And the list summary is ‘0 items left’

### Scenario: Clear all un-clears when two of two completed

        Given a Todo list with items ‘File taxes’ & ‘Walk the dog’
        And the first item is completed
        And the second item is completed
        When the clear all items affordance is toggled
        Then both items are listed
        And both are active
        And the list summary is ‘2 items left’

### Scenario: Items can be straight removed

        Given a Todo list with a single item ‘File taxes’
        When the item is removed
        Then nothing is listed

### Scenario: Initiation of editing takes away delete and complete affordances

        Given a Todo list with single item ‘File taxes’
        When the item is selected for edit
        Then the item cannot be completed
        And the item cannot be removed

### Scenario: Editing can change the text of an item

        Given a Todo list with a single item ‘File taxes’ When the item is selected for edit
        And the Todo is changed to ‘Apply for 6-month tax extension’
        Then only the revised item is listed

### Scenario: Editing can be abandoned

        Given a Todo list with a single item ‘File taxes’
        And that item is being edited
        When editing is abandoned
        Then only that item is listed
        And the list summary is ‘1 item left’

[1]: https://github.com/ToastShaman/zuhlke-camp-2022/blob/open-space/pair-programming/build/pair-programming-guide.pdf
[2]: https://pact.io/
[3]: https://todobackend.com/
[4]: https://tuple.app/
[5]: https://github.com/tastejs/todomvc/blob/master/app-spec.md
[6]: https://paulhammant.com/2017/05/14/todomvc-and-given-when-then-scenarios/