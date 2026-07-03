Feature: Voting
  As a group member
  I want my vote to count once and be changeable
  So that the group can trust the leading answer shown for each question

  Background:
    Given an event named "Saturday Hike"
    And the event has a question "Where should we go?"
    And the question has an option "Snowdon"
    And the question has an option "Ben Nevis"

  Scenario: Voting for the first time asks for a name
    When I vote for "Snowdon" without giving a name
    Then I should be asked who I am

  Scenario: A named vote counts
    When I vote for "Snowdon" as "Alice"
    Then "Snowdon" should have 1 vote
    And the event page should show "Snowdon" as the leading answer

  Scenario: The same person voting again changes their vote instead of adding one
    Given I have voted for "Snowdon" as "Alice"
    When I vote for "Ben Nevis" as "Alice"
    Then "Snowdon" should have 0 votes
    And "Ben Nevis" should have 1 vote

  Scenario: Different people voting for the same option each count
    Given I have voted for "Snowdon" as "Alice"
    When I vote for "Snowdon" as "Bob"
    Then "Snowdon" should have 2 votes

  Scenario: Everyone can see who voted for what
    Given I have voted for "Snowdon" as "Alice"
    When I vote for "Snowdon" as "Bob"
    Then the question page should show "Alice" voted for "Snowdon"
    And the question page should show "Bob" voted for "Snowdon"
