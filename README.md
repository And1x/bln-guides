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
- connect to DB for authenication?

### Current ToDo List:
- [x] Go trough handler and add html.StatusCodes for Error returns
- [x] Build Basic CSS 
- [ ] Build TestCases
- [x] Day/Night mode
### Possible Todo:
- [ ] Use own middleware Chains insead of justinas/alice same - seems not that difficult to implement
- [ ] 

 # Routes:

| method | route          | handler                | description                                       |
|--------|----------------|------------------------|---------------------------------------------------|
| get    | /              | homeSiteHandler        | default home page                                 |
| get    | /createguide   | createGuideFormHandler | empty Form to create guides                       |
| post   | /createguide   | createGuideHandler     | insert new guide in DB → redirect to guide/id     |
| get    | /editguide/id  | editGuideFormHandler   | Form with values from guide by ID                 |
| post   | /editguide     | editGuideHandler       | Updated edited guide in DB → redirect to guide/id |
| get    | /allguides     | allGuideHandler        | lists all guides                                  |
| post   | /deleteguide   | deleteGuideHandler     | deletes Guide by ID                               |
| get    | /guide/id      | singleGuideHandler     | shows specific guide by ID                        |

### More on lnurl-auth
- [fiatjaf blog:](https://fiatjaf.com/e0a35204.html)
- [lnurl-auth explained github:](https://github.com/fiatjaf/lnurl-rfc/blob/legacy/lnurl-auth.md)
- [lnurl-auth expample implementation in go:](https://github.com/xplorfin/lnurlauth)


