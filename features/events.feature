Feature: Events
  As a group organiser
  I want to create an event with a single shareable link
  So that my group can decide on everything from one place

  Scenario: Creating an event
    When I create an event named "Saturday Hike"
    Then the event should exist
    And the event page should show "Saturday Hike"

  Scenario: Renaming an event
    Given an event named "Saturday Hike"
    When I rename the event to "Sunday Hike"
    Then the event page should show "Sunday Hike"

  Scenario: Deleting an event removes its questions and options
    Given an event named "Saturday Hike"
    And the event has a question "Where should we go?"
    And the question has an option "Snowdon"
    When I delete the event
    Then the event should no longer exist
    And the question should no longer exist
    And the option should no longer exist
