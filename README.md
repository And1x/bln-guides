# About bln-guides:
  bln guides is my first real fullstack project to learn and improve my skills in Web development with focus on Go, PostgreSQL and Bitcoin/Lightning. So if you encounter any bugs or have any suggestions for improvement, I would be happy to hear from you.
  
---
- Write Guides, preferably about Bitcoin or Ligthning and get bitcoin(sats) in form of upvotes. Use Markdown to format your guides. 
- Registration creates automatically a Wallet within LNbits where your bitcoin balance is managed.

## Features: 
- Get upvotes in Sats
- WebLN enabled -> no account needed to send Sats to an author (LN address is also visible)
- Format guides in Markdown after [CommonMark](https://commonmark.org/) specification

## Tech-Stack
- ServerSide: Go
- Frontend: HTML/CSS/Javascript
- Database: Postgresql
- Wallet Management: LNbits

## Roadmap/ ToDo List:

- [x] Create DB and Forms to post guides 
- [x] Add markdown
- [x] Build basic CSS 
- [x] Create DB and Forms to register users
- [x] Add basic Authentication
  - currently username + password
- [x] Add LNbits 
  - [x] show balance
  - [x] send/receive payments in form of upvotes 
  - [ ] Deposit/Withdraw in app
- [ ] Add Unit tests
  - [x] Handlers 
  - [x] middleware
  - [ ] helpers
  - [x] render
- [x] Day/Night mode based on system settings
- [ ] Login with LNURL-Auth


 # Routes:

| method                          | route                   | handler            | description                            |
|---------------------------------|-------------------------|--------------------|----------------------------------------|
| GET                             | /                       | homeSite           | Default landing page                   |
| GET                             | /allguides              | allGuides          | Lists all guides                       |
| GET                             | /guide/{id}             | singleGuide        | List specified Guide                   |
| GET                             | /user/register          | registerUserForm   | Form for registration                  |
| POST                            | /user/register          | registerUser       | Creates new User in DB with a wallet   |
| GET                             | /user/login             | loginUserForm      | Form for login                         |
| POST                            | /user/login             | loginUserHandler   | Authenticates a user                   |
|                                 |                         |                    |                                        |
| **for authenticated users:**    |                         |                    |                                        |
| GET                             | /createguide            | createGuideForm    | Form to create a new guide             |
| POST                            | /createguide            | createGuide        | Creates a new guide in DB              |
| POST                            | /deleteguide            | deleteGuide        | Deletes a Guide by ID                  |
| GET                             | /editguide/{id}         | editGuideForm      | Form to edit an existing guide by id   |
| POST                            | /editguide              | editGuide          | Updates guide in DB                    |
| GET                             | /user/profile           | profileUser        | Display user Profile                   |
| GET                             | /user/settings          | settingsUserForm   | Form to change user settings           |
| POST                            | /user/settings          | settingUser        | Update edited user settings in DB      |
| GET                             | /user/settings/password | settingsUserPwForm | Form to change user password           |
| POST                            | /user/settings/password | settingsUserPw     | Update edited user password in DB      |
| POST                            | /user/logout            | logoutUser         | Logs user out by invalidating session  |
| POST                            | /allguides              | upvoteAllGuides    | Upvotes a Guide from allguides page    |
| POST                            | /guide/{id}             | upvoteSingleGuide  | Upvotes a Guide at specific guide page |

## How to run:
- go to [LNbits](https://legend.lnbits.com) (or run your own Instance), add a new wallet and activate User Manager
- setup PostgreSQL (see notes.md for schema) 
- rename 'example.env' to '.env' and fill it out with your config
- then build and run it
