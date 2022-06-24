# Acme Books project

1. Start the emulator. Changes to the datastore are not persisted between sessions
```
docker compose up
```

&nbsp;&nbsp;&nbsp;&nbsp;**Note:** For Mac M1, if you have trouble with the amd64 `singularities/datastore-emulator` docker image, running the datastore emulator locally could be an option:
```
gcloud components install cloud-datastore-emulator
gcloud beta emulators datastore start --project=acme-books --host-port=localhost:3031
```

2. Start the server
```
go run .
```

3. See the current book list
```
http://localhost:3030/books
```

4. Check a single book's details using its id/key
```
http://localhost:3030/books/4
```

# Exercise
1. Refactor the datastore client usage to avoid duplicate code (see main.go and library.go)
2. Order the results for the book list by id
3. Allow filtering to be applied to the book list (via a query parameter)
4. Add a new end point to rent or return a book
```
request:
PUT: http://localhost:3030/:id/borrow
PUT: http://localhost:3030/:id/return

response:
204 if ok
400 if invalid id
400 if already borrowed
400 if not borrowed when returned
```
5. Add a new end point to add a book to the library
```
request:
POST: http://localhost:3030/book
Body: Book in JSON format

response:
200 + Book object including generated key
```
6. Raise a PR and have it peer reviewed by a fellow Gohorter
