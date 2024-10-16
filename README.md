# Posts

Eventually a website to create and read posts, basically title and body in
markdown, but for now I'm experimenting with user management.

## TODO

- if refresh home page after login to fast, the session invalidates and exists
- if logout after token becomes invalid, then it doesn't logout
- if clicking on other links after token becomes invalid, it doesn't logout

## References

- https://www.alexedwards.net/blog/working-with-cookies-in-go
- https://stackoverflow.com/questions/54258233/do-i-have-to-store-tokens-in-cookies-or-localstorage-or-session
- https://stackoverflow.com/questions/53678019/laravel-5-6-passport-jwt-httponly-cookie-spa-authentication-for-self-consuming/54011649#54011649
- https://stackoverflow.com/questions/41496924/how-to-authenticate-spa-users-using-oauth2/53988717
- https://en.wikipedia.org/wiki/Salt_(cryptography)
- https://auth0.com/blog/adding-salt-to-hashing-a-better-way-to-store-passwords/
- https://stackoverflow.com/questions/7562675/proper-way-to-send-username-and-password-from-client-to-server
