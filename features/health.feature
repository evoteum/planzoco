Feature: Health check
  As an operator
  I want a health endpoint
  So that Kubernetes can tell the app is alive

  Scenario: The health endpoint reports OK
    When I request the health endpoint
    Then the response status should be 200
    And the response body should be "OK"
