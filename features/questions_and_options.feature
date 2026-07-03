Feature: Questions and options
  As a group member
  I want to add, edit, and remove questions and options
  So the group can decide together, all from the one event link

  Background:
    Given an event named "Saturday Hike"

  Scenario: Adding a question to an event
    When I add a question "Where should we go?" to the event
    Then the event page should list the question "Where should we go?"

  Scenario: Adding a question takes you straight to it, ready to add options
    When I add a question "Where should we go?" to the event
    Then I should land on the question's own page

  Scenario: Changing a question
    Given the event has a question "Where should we go?"
    When I change the question to "Where should we eat?"
    Then the question should now read "Where should we eat?"

  Scenario: Removing a question
    Given the event has a question "Where should we go?"
    When I delete the question
    Then the question should no longer exist

  Scenario: Suggesting an option for a question
    Given the event has a question "Where should we go?"
    When I suggest the option "Snowdon" for the question
    Then the question page should list the option "Snowdon"

  Scenario: Changing an option
    Given the event has a question "Where should we go?"
    And the question has an option "Snowdon"
    When I change the option "Snowdon" to "Ben Nevis"
    Then the question page should list the option "Ben Nevis"

  Scenario: Removing an option
    Given the event has a question "Where should we go?"
    And the question has an option "Snowdon"
    When I delete the option "Snowdon"
    Then the option "Snowdon" should no longer exist

  Scenario: A question with no votes shows no leading answer
    Given the event has a question "Where should we go?"
    Then the event page should show "No votes yet" for the question

  Scenario: A question with no options yet invites you to add one, not vote
    Given the event has a question "Where should we go?"
    Then the question should invite adding an option

  Scenario: Once a question has an option, the event page invites you to vote
    Given the event has a question "Where should we go?"
    And the question has an option "Snowdon"
    Then the question should invite voting

  Scenario: Everyone can see who asked a question
    When I add a question "Where should we go?" to the event as "Alice"
    Then the question page should show that "Alice" asked the question
    And the event page should show that "Alice" asked the question
