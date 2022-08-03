# bln--h
- Web Page where Users can Login with LNURL-Auth
- Post Guides in Markdown Language
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
- [ ] Go trough handler and add html.StatusCodes for Error returns
- [x] Build Basic CSS 
- [ ] Build TestCases
- [x] Day/Night mode


### More on lnurl-auth
- [fiatjaf blog:](https://fiatjaf.com/e0a35204.html)
- [lnurl-auth explained github:](https://github.com/fiatjaf/lnurl-rfc/blob/legacy/lnurl-auth.md)
- [lnurl-auth expample implementation in go:](https://github.com/xplorfin/lnurlauth)


