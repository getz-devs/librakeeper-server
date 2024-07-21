## API Documentation

### Authentication

All endpoints under `/api` require authentication using Firebase. Include the `Authorization` header in your request
with a valid Firebase ID token:

`Authorization: Bearer <token>`

### Book Endpoints

| Method   | Endpoint                   | Description                                     | Query Params                                             | Path Params     | Data Structures         |
|----------|----------------------------|-------------------------------------------------|----------------------------------------------------------|-----------------|-------------------------|
| `POST`   | `/api/books/add`           | Create a new book in the user's library.        | None                                                     | None            | `Book`                  |
| `POST`   | `/api/books/add/advanced`  | Add book from advanced search to user's library | `isbn` (string), `index` (number)                        | None            | None                    |
| `GET`    | `/api/books/`              | Retrieve books for the authenticated user.      | `page` (number, default 1), `limit` (number, default 10) | None            | `PaginatedBookResponse` |
| `GET`    | `/api/books/:id`           | Retrieve a book by ID.                          | None                                                     | `id` (string)   | `Book`                  |
| `GET`    | `/api/books/isbn/:isbn`    | Retrieve a book by ISBN.                        | None                                                     | `isbn` (string) | `Book`                  |
| `GET`    | `/api/books/bookshelf/:id` | Retrieve books from a specific bookshelf.       | `page` (number, default 1), `limit` (number, default 10) | `id` (string)   | `PaginatedBookResponse` |
| `PUT`    | `/api/books/:id`           | Update a book.                                  | None                                                     | `id` (string)   | `BookUpdate`            |
| `DELETE` | `/api/books/:id`           | Delete a book.                                  | None                                                     | `id` (string)   | None                    |

#### Data Structures

**`Book`:**

```typescript
interface Book {
    id: string;
    userId: string;
    bookshelfId: string;
    isbn: string;
    title: string;
    author: string;
    publishing: string;
    description: string;
    coverImage: string;
    shopName: string;
    createdAt: Date;
    updatedAt: Date;
}
```

**`BookUpdate`:**

```typescript
interface BookUpdate {
    isbn?: string;
    bookshelfId?: string;
    title?: string;
    author?: string;
    publishing?: string;
    description?: string;
    coverImage?: string;
    shopName?: string;
    updatedAt: Date;
}
```

**`PaginatedBookResponse`:**

```typescript
interface PaginatedBookResponse {
    books: Book[];
    total: number;
    page: number;
    limit: number;
}
```

### Bookshelf Endpoints

| Method   | Endpoint               | Description                                      | Query Params                                             | Path Params   | Data Structures              |
|----------|------------------------|--------------------------------------------------|----------------------------------------------------------|---------------|------------------------------|
| `POST`   | `/api/bookshelves/add` | Create a new bookshelf.                          | None                                                     | None          | `Bookshelf`                  |
| `GET`    | `/api/bookshelves/`    | Retrieve bookshelves for the authenticated user. | `page` (number, default 1), `limit` (number, default 10) | None          | `PaginatedBookshelfResponse` |
| `GET`    | `/api/bookshelves/:id` | Retrieve a bookshelf by ID.                      | None                                                     | `id` (string) | `Bookshelf`                  |
| `PUT`    | `/api/bookshelves/:id` | Update a bookshelf.                              | None                                                     | `id` (string) | `BookshelfUpdate`            |
| `DELETE` | `/api/bookshelves/:id` | Delete a bookshelf.                              | None                                                     | `id` (string) | None                         |

#### Data Structures

**`Bookshelf`:**

```typescript
interface Bookshelf {
    id: string;
    userId: string;
    name: string;
    createdAt: Date;
    updatedAt: Date;
}
```

**`BookshelfUpdate`:**

```typescript
interface BookshelfUpdate {
    name?: string;
    updatedAt: Date;
}
```

**`PaginatedBookshelfResponse`:**

```typescript
interface PaginatedBookshelfResponse {
    bookshelves: Bookshelf[];
    total: number;
    page: number;
    limit: number;
}
```

### Search Endpoints

| Method | Endpoint               | Description                                               | Query Params    | Path Params | Data Structures  |
|--------|------------------------|-----------------------------------------------------------|-----------------|-------------|------------------|
| `GET`  | `/api/search/simple`   | Search for a book by ISBN in the local database.          | `isbn` (string) | None        | `SearchResponse` |
| `GET`  | `/api/search/advanced` | Perform an advanced search for a book by ISBN using gRPC. | `isbn` (string) | None        | `SearchResponse` |

#### Data Structures

**`SearchResponse`:**

```typescript
enum SearchStatus {
    PROCESSING = 0,
    SUCCESS = 1,
    FAILED = 2,
}

interface SearchResponse {
    status: SearchStatus;
    books: Book[];
}
```

### Health Check Endpoint

| Method | Endpoint      | Description                                           | Query Params | Path Params | Data Structures |
|--------|---------------|-------------------------------------------------------|--------------|-------------|-----------------|
| `GET`  | `/api/health` | Checks the health of the server and its dependencies. | None         | None        | None            |
