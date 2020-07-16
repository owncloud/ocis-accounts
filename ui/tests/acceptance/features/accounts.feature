Feature: Accounts

	Scenario: list accounts
		Given user "user1" has been created with default attributes
		And user "user1" has logged in using the webUI
		When the user browses to the accounts page
		Then user "einstein" should be displayed in the accounts list on the WebUI
		And user "konnectd" should be displayed in the accounts list on the WebUI
		And user "marie" should be displayed in the accounts list on the WebUI
		And user "reva" should be displayed in the accounts list on the WebUI
		And user "richard" should be displayed in the accounts list on the WebUI

