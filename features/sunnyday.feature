
Feature: Sunny Day
  This is a sunny day example

Scenario: The HASD web site is still available
  Given the HASD has a covid site
  When I consume the url
  Then the http status should be "ok"
