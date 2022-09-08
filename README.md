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
### Possible Todo:
- [ ] 

 # Routes:

| method | route          | handler                | description                                       |
|--------|----------------|------------------------|---------------------------------------------------|
| GET    | /              | homeSiteHandler        | default home page                                 |
| GET    | /createguide   | createGuideFormHandler | empty Form to create guides                       |
| POST   | /createguide   | createGuideHandler     | insert new guide in DB → redirect to guide/id     |
| GET    | /editguide/id  | editGuideFormHandler   | Form with values from guide by ID                 |
| POST   | /editguide     | editGuideHandler       | Updated edited guide in DB → redirect to guide/id |
| GET    | /allguides     | allGuideHandler        | lists all guides                                  |
| POST   | /deleteguide   | deleteGuideHandler     | deletes Guide by ID                               |
| GET    | /guide/id      | singleGuideHandler     | shows specific guide by ID                        |
| GET    | /user/register | registerUserFormHandler| Form to register new users                        |
| POST   | /user/register | registerUserHandler    | Create a new user in DB                           |
| GET    | /user/login      | loginUserFormHandler | Form for the Login                                |
| POST   | /user/login      | loginUserHandler     | Authentication + Login                            |
| POST   | /user/logout     | logoutUserHandler    | Logout                                            |


### More on lnurl-auth
- [fiatjaf blog:](https://fiatjaf.com/e0a35204.html)
- [lnurl-auth explained github:](https://github.com/fiatjaf/lnurl-rfc/blob/legacy/lnurl-auth.md)
- [lnurl-auth expample implementation in go:](https://github.com/xplorfin/lnurlauth)


