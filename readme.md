# Tig
Tig is a command-line interface (CLI) app that enables you to push your commit message to a spreadsheet using the Google Sheets API and the Golang programming language. I use this app because my Tribe needs to track what people (read: developers) do every day, and sometimes I forget to do that 😭✊.

## Installation

To use this app, you'll need to have Golang installed on your machine. Once you have Golang installed, you can download the app's source code and build it using the following steps:

1.  Clone the repository: `git clone https://github.com/zakyyudha/tig.git`
2.  Navigate to the project directory: `cd tig`
3.  Build the app: `go install ./cmd/tig`

## Usage

### Enable the Google API

Before using Google APIs, you need to turn them on in a Google Cloud project. You can turn on one or more APIs in a single Google Cloud project.

-   In the Google Cloud console, enable the Google Sheets API.

    [Enable the API](https://console.cloud.google.com/flows/enableapi?apiid=sheets.googleapis.com)

### Authorize Google credentials for a desktop application

To authenticate as an end-user and access user data in your app, you need to create one or more OAuth 2.0 Client IDs. A client ID is used to identify a single app to Google's OAuth servers. If your app runs on multiple platforms, you must create a separate client ID for each platform.

1.  In the Google Cloud console, go to Menu > **APIs & Services** > **Credentials**.

    [Go to Credentials](https://console.cloud.google.com/apis/credentials)

2.  Click **Create Credentials** > **OAuth client ID**.
3.  Click **Application type** > **Desktop app**.
4.  In the **Name** field, type a name for the credential. This name is only shown in the Google Cloud console.
5.  Click **Create**. The OAuth client created screen appears, showing your new Client ID and Client secret.
6.  Click **OK**. The newly created credential appears under **OAuth 2.0 Client IDs.**
7.  Save the downloaded JSON file as `credentials.json`, and move the file to `$HOME/.config/tig/credentials.json`.

Once you've set up your Google API credentials, you can run the app using the following command:

`$ tig [commit] [-m] "some commit message" [-tJC] "some jira number" [-tN] "some note"`

On the first run, you will be prompted by a message to complete the installation. Here's an example:

```bash
$ tig commit -m "TJC-123 | Update the folder structure and delete the built binary"
Go to the following link in your browser then type the authorization code:
https://some-long-url.com
Input authorization code: some-long-code
Saving credential file to: /path/to/home/.config/tig/token.json
If your input contains a space, please add double quotes (e.g., "Some Input").
Input Spreadsheet ID: 1ESHDjQalu30gtOzH3x18twauRHHyBdDV3x7eplZR7xg
Input Sheet Name: "Example Sheet"
Input Squad Name: "Your Squad" 
Input Your Name: "Your Name"
```

The authorization code is from the callback URL on your browser. Here's an example URL:

`http://localhost/?state=state-token&code=some-long-code&scope=https://www.googleapis.com/auth/spreadsheets`

Note that Tig supports all standard commit flags and arguments based on official Git. Tig only extracts the commit message from Git and converts it into an activity in the spreadsheet. In addition to the standard Git flags, Tig also supports its own flags, `-tJC` and `-tN`, for additional details.

Here is an example of how to use Tig to push a commit message to a spreadsheet:

```bash
$ tig commit -m "TJC-123 | Update the folder structure and delete the built binary" \
    -tJC "TJC-123" \
    -tN "Updated the project structure to improve maintainability"
```

In this example, the `-m` flag is used to specify the commit message, the `-tJC` flag is used to specify the Jira code associated with the commit, and the `-tN` flag is used to specify additional notes about the commit.

Once the commit has been made, the activity will be added to the Google Spreadsheet specified in the `spreadsheetId` parameter. 
Example: [Click Here](https://docs.google.com/spreadsheets/d/1ESHDjQalu30gtOzH3x18twauRHHyBdDV3x7eplZR7xg/edit#gid=0)