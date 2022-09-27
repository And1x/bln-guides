# bln--h
- Web Page where Users can Login with LNURL-Auth
- Post Guides in Markdown Language
- WebLN enabeld when looking at single-guides
- -> Get Upvotes in Sats

## Tech-Stack
- ServerSide: Go
- Frontend: Plain HTML/CSS/Javascript
- Database: Postgresql

### Roadmap
#### 1. Post Guides
- Build bare minimum to post something 
- Create a DB to post guides -> Title, Content, Created_at, Creator
- -> html forms -> content with markdown enabled -> save content as plain markdown in db -> render it to html as html.ResponseWriter to display at page
#### 2. Create User
- Login with Oauth -> see goth or self implementation
- lnurl-auth later... 
- User creates wallet within LNbits to get upvotes

### Current ToDo List:
- [x] Go trough handler and add html.StatusCodes for Error returns
- [x] Build Basic CSS 
- [ ] Build TestCases
- [x] Day/Night mode
- [ ] Check/investigate about Dependency Injection with LNbits (type, config ...)
### Possible Todo:
- [ ] 

### Data stored in session cookie:
1. temporary session k/v 
- FlashMsg on: (Guide created, Guide deleted, Registerd, Login, settings changed, password changed, Logout)
2. normal as long as session is valid
- userID, userName

 # Routes:

| method                          | route                   | handler            | description                            |
|---------------------------------|-------------------------|--------------------|----------------------------------------|
| GET                             | /                       | homeSite           | Default landing page                   |
| GET                             | /allguides              | allGuides          | lists all guides                       |
| GET                             | /guide/$id              | singleGuide        | list specified Guide                   |
| GET                             | /user/register          | registerUserForm   | Form for registration                  |
| POST                            | /user/register          | registerUser       | Creates new User in DB with a wallet   |
| GET                             | /user/login             | loginUserForm      | Form for login                         |
| POST                            | /user/login             | loginUserHandler   | Authenticates a user                   |
| Routes for authenticated users: |                         |                    |                                        |
| GET                             | /createguide            | createGuideForm    | Form to create a new guide             |
| POST                            | /createguide            | createGuide        | Creates a new guide in DB              |
| POST                            | /deleteguide            | deleteGuide        | Deletes a Guide by ID                  |
| GET                             | /editguide/$id          | editGuideForm      | Form to edit an existing guide by id   |
| POST                            | /editguide              | editGuide          | Updates guide in DB                    |
| GET                             | /user/profile           | profileUser        | Display user Profile                   |
| GET                             | /user/settings          | settingsUserForm   | Form to change user settings           |
| POST                            | /user/settings          | settingUser        | Update edited user settings in DB      |
| GET                             | /user/settings/password | settingsUserPwForm | Form to change user password           |
| POST                            | /user/settings/password | settingsUserPw     | Update edited user password in DB      |
| POST                            | /user/logout            | logoutUser         | Logs user out by invalidating session  |
| POST                            | /allguides              | upvoteAllGuides    | Upvotes a Guide from allguides page    |
| POST                            | /guide/$id              | upvoteSingleGuide  | Upvotes a Guide at specific guide page |

--> Rethink about upvote-handlers - I currently use 2 bc. it's more readable and cleaner flow, however 1 could also be use similar to the delete-function...

### More on lnurl-auth
- [fiatjaf blog:](https://fiatjaf.com/e0a35204.html)
- [lnurl-auth explained github:](https://github.com/fiatjaf/lnurl-rfc/blob/legacy/lnurl-auth.md)
- [lnurl-auth expample implementation in go:](https://github.com/xplorfin/lnurlauth)

### Thougth's and Questions...
- LNbits Keys also unique entrys in DB? (+check serverside?)
- Use less precise DB queries eg. Pull whole Guide instead just guides.UserId ?
  - This would result in less Queries(code) but we also pull always way more Data than needed



