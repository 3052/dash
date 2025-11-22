# dash

Go language, I need a package to parse DASH (MPD) files

1. standard library only
2. use a separate file for each type
3. package will be called "dash"
4. parse input will be byte slice
5. BaseURL is type "string", not a slice
6. omit ProgramInformation
7. include test Go file in same package
8. do not skip any tests
9. test will read from rakuten.mpd, which user will provide
10. only send new or updated files
