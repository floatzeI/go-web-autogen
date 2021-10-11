# Progress

This file will temporarily exist until all major features are finished.

- ✔️ Automated routing through decorator tags
- ✔️ Automated response bodies and "200" code for swagger
- ✔️ Automated response bodies and codes through `@Response()` decorator
- ✔️ Access to raw context object through controller/endpoints
- Customizable middleware for handling bad requests (argument not provided, argument not valid, etc)
- Middleware decorator
- Small wrappers around commonly accessed objects so you aren't constrained to one web framework (e.g. a "Cookie" param, or a "Header" param, etc)
- Ability to use program through an executable (e.g. from command line), and from another go program
- Tests on generated routes to make sure the code is valid and should (at the very least) compile correctly, immediately after it is generated
- A decently-sized example project that:
    - follows best practices
    - uses multiple datastores (e.g. redis, postgresql, and elasticsearch)
    - uses docker


## Smaller Things

Features I don't really care about right now but might add anyway.

- Ability to use other http servers (so you aren't forced to use fiber)
