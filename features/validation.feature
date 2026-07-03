Feature: Input length limits
  As the service operator
  I want event, question, option, and voter name text capped at 255 characters
  So that one participant can't flood the page with a huge block of text

  Scenario: An event name over 255 characters is rejected
    When I create an event with a 256 character name
    Then the event should be rejected

  Scenario: An event name of exactly 255 characters is accepted
    When I create an event with a 255 character name
    Then the event should exist

  Scenario: A question over 255 characters is rejected
    Given an event named "Saturday Hike"
    When I add a question with 256 characters to the event
    Then the question should be rejected

  Scenario: An option over 255 characters is rejected
    Given an event named "Saturday Hike"
    And the event has a question "Where should we go?"
    When I suggest an option with 256 characters for the question
    Then the option should be rejected
