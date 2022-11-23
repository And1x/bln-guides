# About bln-guides:
  BLN Guides is a fullstack project to learn and improve my skills in Web development with a focus on Go, PostgreSQL and Bitcoin/Lightning. So if you encounter any bugs or have any suggestions for improvement, I would be happy to hear from you.

---
- Write Guides, preferably about Bitcoin or Ligthning and get bitcoin(sats) in form of upvotes. Use Markdown to format your guides. 
- Registration creates automatically a Wallet within LNbits where your bitcoin balance is managed.

## Features: 
- Get upvotes in Sats
- WebLN enabled -> no account needed to send Sats to an author (LN address is also visible)
- Format guides in Markdown after [CommonMark](https://commonmark.org/) specification

## Tech-Stack
- Frontend: HTML/CSS/Javascript
- Backend: Go, PostgreSQL and LNbits

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
  - [x] Deposit/Withdraw in app
- [ ] Add Unit tests
  - [x] Handlers 
  - [x] middleware
  - [ ] helpers
  - [x] render
- [x] Day/Night mode based on system settings
- [ ] Login with LNURL-Auth
- [ ] Improve withdraw UX
- [ ] Reset Password per Mail




## How to run:
- go to [LNbits](https://legend.lnbits.com) (or run your own Instance), add a new wallet and activate User Manager
- setup PostgreSQL (see notes.md for schema) 
- rename 'example.env' to '.env' and fill it out with your config
- then build and run it
