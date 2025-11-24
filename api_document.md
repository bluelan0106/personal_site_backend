## Posts / Comments / Reactions APIs

This section describes the post (article) system: creating and managing posts, nested comments, and reaction types (like, love, haha, wow, sad, angry, care).

Notes:
Authentication: most write endpoints require authentication (the server uses an HTTP-only `auth_token` cookie). Read endpoints are usually `AuthOptional`.
Visibility: posts support `public` and `private` visibility. Private posts are visible only to the author.

### GET /posts
**Description**: List posts with simple filtering and pagination. Authentication optional.

**Query Parameters**:
- `page` (int, optional): page number (default 1)
- `limit` (int, optional): items per page (default 20)
- `author_id` (uint, optional): filter by author
- `tag` (string, optional): filter by tag slug
- `status` (string, optional): `published|draft|archived` (only authors/admins can see drafts)

**Success Response (200)**:
```json
{
  "page": 1,
  "limit": 20,
  "total": 123,
  "posts": [
    {
      "id": 10,
      "author_id": 1,
      "title": "Hello world",
      "summary": "Short summary",
      "cover_image": "/uploads/cover.jpg",
      "status": "published",
      "visibility": "public",
      "view_count": 42,
      "published_at": "2025-11-10T12:00:00Z",
      "tags": ["news","update"]
    }
  ]
}
```

**Error Responses**:
- `400 Bad Request`: invalid query parameters

---

### POST /posts
**Description**: Create a new post. Requires authentication.

**Request Body**:
```json
{
  "title": "My post title",
  "content": "Full HTML or markdown content",
  "summary": "Optional short summary",
  "cover_image": "/uploads/cover.jpg",
  "status": "draft", // draft | published | archived
  "visibility": "public", // public | private
  "tags": ["go","programming"]
}
```

**Request Body Schema**:
- `title` (string, required)
- `content` (string, required)
- `summary` (string, optional)
- `cover_image` (string, optional)
- `status` (string, optional): `draft|published|archived` (default `draft`)
- `visibility` (string, optional): `public|private` (default `public`)
- `tags` (string[], optional)

**Success Response (201)**:
```json
{
  "id": 11,
  "author_id": 1,
  "title": "My post title",
  "status": "draft",
  "visibility": "public",
  "published_at": null
}
```

**Error Responses**:
- `400 Bad Request`: validation error
- `401 Unauthorized`: missing/invalid authentication
- `500 Internal Server Error`: DB error

---

### GET /posts/:id
**Description**: Get a single post by ID. Authentication optional. Private posts return 403 for non-author.

**Path Parameters**:
- `id` (uint, required): Post ID

**Success Response (200)**:
```json
{
  "id": 11,
  "author_id": 1,
  "title": "My post title",
  "content": "Full content",
  "summary": "Optional",
  "cover_image": "/uploads/cover.jpg",
  "status": "published",
  "visibility": "public",
  "view_count": 123,
  "published_at": "2025-11-10T12:00:00Z",
  "tags": ["go","programming"],
  "comments_count": 5,
  "reactions_summary": {"like": 10, "love": 2}
}
```

**Error Responses**:
- `401 Unauthorized` / `403 Forbidden`: access denied for private posts
- `404 Not Found`: post not found

---

### PUT /posts/:id
**Description**: Update a post. Requires authentication and must be the author (or admin).

**Path Parameters**:
- `id` (uint, required): Post ID

**Request Body**: same as `POST /posts` (partial updates allowed)

**Success Response (200)**:
```json
{
  "message": "Post updated",
  "id": 11
}
```

**Error Responses**:
- `400 Bad Request`: validation error
- `401 Unauthorized`: not logged in
- `403 Forbidden`: not the author
- `404 Not Found`: post not found

---

### DELETE /posts/:id
**Description**: Delete a post. Requires authentication and must be the author or an admin. This performs a soft delete.

**Path Parameters**:
- `id` (uint, required): Post ID

**Success Response (200)**:
```json
{
  "message": "Post deleted"
}
```

**Error Responses**:
- `401 Unauthorized`, `403 Forbidden`, `404 Not Found`

---

### GET /posts/:id/comments
**Description**: List comments for a post. Returns nested replies. Authentication optional.

**Path Parameters**:
- `id` (uint, required): Post ID

**Success Response (200)**:
```json
[
  {
    "id": 1,
    "post_id": 11,
    "author_id": 2,
    "content": "Top-level comment",
    "parent_id": null,
    "replies": [
      { "id": 2, "parent_id": 1, "content": "Reply" }
    ],
    "is_deleted": false
  }
]
```

---

### POST /posts/:id/comments
**Description**: Create a comment on a post. Requires authentication.

**Path Parameters**:
- `id` (uint, required): Post ID

**Request Body**:
```json
{
  "content": "Nice post!",
  "parent_id": null
}
```

**Success Response (201)**:
```json
{
  "id": 21,
  "post_id": 11,
  "author_id": 3,
  "content": "Nice post!",
  "parent_id": null
}
```

**Error Responses**:
- `400 Bad Request`, `401 Unauthorized`, `404 Not Found` (post not found)

---

### PUT /comments/:comment_id
**Description**: Update a comment. Requires authentication and comment owner.

**Request Body**:
```json
{
  "content": "Updated content"
}
```

**Success Response (200)**:
```json
{
  "message": "Comment updated",
  "id": 21
}
```

---

### DELETE /comments/:comment_id
**Description**: Soft-delete a comment. Requires authentication and comment owner or post author.

**Success Response (200)**:
```json
{
  "message": "Comment deleted"
}
```

---

### POST /posts/:id/reactions
**Description**: Add or toggle a reaction to a post. Requires authentication.

**Path Parameters**:
- `id` (uint, required): Post ID

**Request Body**:
```json
{ "type": "like" }
```

**Behavior**:
- If the user has no reaction on the target, the reaction is created.
- If the user has the same reaction, it is removed (toggle off).
- If the user has a different reaction, it is updated to the new type.

**Success Response (200)**:
```json
{ "message": "Reaction updated", "summary": {"like": 10, "love": 2} }
```

---

### POST /comments/:comment_id/reactions
**Description**: Add or toggle a reaction to a comment. Requires authentication. Same behavior as post reactions.

---

### GET /posts/:id/reactions
**Description**: Get reaction summary for a post. Authentication optional.

**Success Response (200)**:
```json
{ "like": 10, "love": 2, "haha": 0 }
```

---

### GET /comments/:comment_id/reactions
**Description**: Get reaction summary for a comment.

**Success Response (200)**:
```json
{ "like": 2, "love": 0 }
```

---

## Posts API Notes
- Reactions types: `like`, `love`, `haha`, `wow`, `sad`, `angry`, `care`.
- Comments use soft-delete: deleted comments keep place in thread but marked as deleted.
- Posts with `status: draft` are visible only to the author and admins.
- All write operations require authentication; read endpoints are often available without auth but hide private content.
# Personal Site Backend API Documentation

## Authentication APIs

### POST /auth/register
**Description**: Register a new user account

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "nickname": "username"
}
```

**Request Body Schema**:
- `email` (string, required): User's email address (must be valid email format)
- `password` (string, required): User's password (minimum 8 characters)
- `nickname` (string, required): User's display name

**Success Response (200)**:
```json
{
  "message": "User registered successfully",
  "user_id": 1
}
```

**Error Responses**:
- `400 Bad Request`: Invalid input data
  ```json
  {
    "error": "Key: 'registerRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"
  }
  ```
- `500 Internal Server Error`: Server error during registration
  ```json
  {
    "error": "Failed to create user"
  }
  ```

---

### POST /auth/login
**Description**: Login with email and password. Sets an `auth_token` HTTP-only cookie for the session.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Request Body Schema**:
- `email` (string, required): User's email address (must be valid email format)
- `password` (string, required): User's password (minimum 8 characters)

**Success Response (200)**:
```json
{
  "user_id": 1,
  "message": "Login successful",
  "role": "user",
  "nickname": "username",
}
```

**Response Schema**:
- `user_id` (uint): User's unique identifier
- `message` (string): Success message
- `role` (string): User's role (e.g., "user", "admin")
- `nickname` (string): User's display name

**Error Responses**:
- `400 Bad Request`: Invalid input data
  ```json
  {
    "error": "Key: 'loginRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"
  }
  ```
- `401 Unauthorized`: Invalid credentials
  ```json
  {
    "error": "Invalid email or password"
  }
  ```
- `500 Internal Server Error`: Server error during login
  ```json
  {
    "error": "Failed to generate token"
  }
  ```

---

### POST /auth/logout
**Description**: Logout user and clear authentication cookie. Removes the `auth_token` cookie.

**Success Response (200)**:
```json
{
  "message": "Logged out successfully"
}
```

**Error Responses**:
- `500 Internal Server Error`: Server error during logout
  ```json
  {
    "error": "Failed to logout"
  }
  ```

---

### POST /auth/change-password
**Description**: Change user's password (requires login first)

**Request Body**:
```json
{
  "old_password": "oldpassword123",
  "new_password": "newpassword123"
}
```

**Request Body Schema**:
- `old_password` (string, required): Current password (minimum 8 characters)
- `new_password` (string, required): New password (minimum 8 characters)

**Success Response (200)**:
```json
{
  "message": "Password changed successfully"
}
```

**Error Responses**:
- `400 Bad Request`: Invalid input data
  ```json
  {
    "error": "Key: 'changePasswordRequest.OldPassword' Error:Field validation for 'OldPassword' failed on the 'required' tag"
  }
  ```
- `401 Unauthorized`: Missing or invalid token, or incorrect old password
  ```json
  {
    "error": "Unauthorized"
  }
  ```
  ```json
  {
    "error": "Old password is incorrect"
  }
  ```
- `403 Forbidden`: Password change not allowed for this account type
  ```json
  {
    "error": "Password change only allowed for password-based accounts"
  }
  ```
- `500 Internal Server Error`: Server error during password change
  ```json
  {
    "error": "Failed to update password"
  }
  ```

### GET /auth/login-github
Description: Start GitHub OAuth login flow. Optionally accept a `redirect` query param to indicate where the browser should be redirected after a successful login.

Query Parameters:
- `redirect` (string, optional): A full URL to redirect to on success. This value is preserved via OAuth state and used in the callback.

Response:
- 302 Redirect to GitHub authorization URL

Notes:
- The server encodes a nonce and the `redirect` value into the OAuth `state` parameter.

### GET /auth/login-github-callback
Description: OAuth callback endpoint for GitHub. Exchanges the authorization code for a token, creates or finds the user, sets the `auth_token` HTTP-only cookie, and then redirects back to the provided `redirect` URL if present. If no `redirect` is provided, returns JSON.

Reads:
- `state` (from GitHub): contains an encoded object with fields `n` (nonce) and `r` (redirect URL).
- `redirect` (optional query): used only if not present in state.

On success with redirect present:
- 302 Redirect to the `redirect` URL with the following query parameters appended:
  - `login=success`
  - `user_id` (number)
  - `message` (string): "GitHub login successful"
  - `role` (string)
  - `nickname` (string)

On success without redirect:
- 200 JSON:
```json
{
  "message": "GitHub login successful",
  "user_id": 1,
  "role": "user",
  "nickname": "username"
}
```

Error Responses:
- `400 Bad Request`: Missing `code`
- `401 Unauthorized`: Code exchange failed
- `500 Internal Server Error`: OAuth not configured, DB or token errors

---

### GET /auth/login-google
Description: Start Google OAuth login flow. Optionally accept a `redirect` query param to indicate where the browser should be redirected after a successful login.

Query Parameters:
- `redirect` (string, optional): A full URL to redirect to on success. This value is preserved via OAuth state and used in the callback.

Response:
- 302 Redirect to Google authorization URL

Notes:
- The server encodes a nonce and the `redirect` value into the OAuth `state` parameter.

### GET /auth/login-google-callback
Description: OAuth callback endpoint for Google. Exchanges the authorization code, creates or finds the user, sets the `auth_token` HTTP-only cookie, and then redirects back to the provided `redirect` URL if present. If no `redirect` is provided, returns JSON.

Reads:
- `state` (from Google): contains an encoded object with fields `n` (nonce) and `r` (redirect URL).
- `redirect` (optional query): used only if not present in state.

On success with redirect present:
- 302 Redirect to the `redirect` URL with the following query parameters appended:
  - `login=success`
  - `user_id` (number)
  - `message` (string): "Google login successful"
  - `role` (string)
  - `nickname` (string)

On success without redirect:
- 200 JSON similar to GitHub callback.

Error Responses:
- `400 Bad Request`: Missing `code`
- `401 Unauthorized`: Code exchange failed
- `500 Internal Server Error`: OAuth not configured, DB or token errors

---

## Authentication Flow

1. **Register**: Create a new account using `/auth/register`. Or you don't need to do that if you use OAuth.
2. **Login**: Authenticate using `/auth/login` and you don't need to manage any thing about session. Or use `/auth/login-{3rd-platform}` to use OAuth login.
3. **Access Protected Resources**: token will saved in http only cookie
4. **Change Password**: Use `/auth/change-password` with valid authentication

---

## Storage APIs
**Description**:
```
All users can access its own storage and users who have not logged in share a storage.

Don't upload any thing you don't want to share whit admins here. Admins could see all your files.
```

### POST /storage/folder/*folder_path
**Description**: Create a new folder in the storage system

**Path Parameters**:
- `folder_path` (string, required): The folder path to create (supports nested paths, e.g., `/documents/2024/reports`)

**Headers**:
- `Cookie`: auth_token (optional) - Authentication cookie for user identification

**Success Response (200)**:
```json
{
  "message": "Directory created successfully"
}
```

**Error Responses**:
- `500 Internal Server Error`: Failed to create directory
  ```json
  {
    "error": "Failed to create directory"
  }
  ```

**Example**:
```bash
POST /storage/folder/documents/2024/reports
```

---

### GET /storage/folder/*folder_path
**Description**: List contents of a folder

**Path Parameters**:
- `folder_path` (string, required): The folder path to list (supports nested paths, e.g., `/documents/2024`)

**Headers**:
- `Cookie`: auth_token (optional) - Authentication cookie for user identification

**Success Response (200)**:
```json
[
    {
        "is_dir": true,
        "name": "test",
        "size": 4096
    },
    {
        "is_dir": false,
        "name": "test.txt",
        "size": 3
    }
]
```

**Response Schema**:
  - `name` (string): File or folder name
  - `is_dir` (bool): Whether it is a folder
  - `size` (number): File or folder size in bytes

**Error Responses**:
- `500 Internal Server Error`: Failed to list folder contents
  ```json
  {
    "error": "Failed to list folder contents"
  }
  ```

**Example**:
```bash
GET /storage/folder/documents/2024
```

---

### PATCH /storage/folder/*folder_path
**Description**: Update a folder

**Path Parameters**:
- `folder_path` (string, required): The current folder path to update

**Request Body**:
```json
{
  "path": "new folder path"
}
```

**Request Body Schema**:
- `path` (string, optional): The new folder path.


**Headers**:
- `Cookie`: auth_token (optional) - Authentication cookie for user identification

**Success Response (200)**:
```json
{
  "message": "Folder updated successfully"
}
```

**Error Responses**:
- `500 Internal Server Error`: Failed to update folder
  ```json
  {
    "error": "Failed to update folder"
  }
  ```
- `500 Internal Server Error`: Failed to rename folder
  ```json
  {
    "error": "Failed to move folder"
  }
  ```

**Example**:
```bash
PATCH /storage/folder/documents/old_name
```

---

### DELETE /storage/folder/*folder_path
**Description**: Delete a folder and all its contents

**Path Parameters**:
- `folder_path` (string, required): The folder path to delete

**Headers**:
- `Cookie`: auth_token (optional) - Authentication cookie for user identification

**Success Response (200)**:
```json
{
  "message": "Folder deleted successfully"
}
```

**Error Responses**:
- `500 Internal Server Error`: Failed to delete folder
  ```json
  {
    "error": "Failed to delete folder"
  }
  ```

**Example**:
```bash
DELETE /storage/folder/documents/temp
```

---

### GET /storage/file/*file_path
**Description**: Download/retrieve a file from storage

**Path Parameters**:
- `file_path` (string, required): The file path to retrieve (supports nested paths, e.g., `/documents/2024/report.pdf`)

**Headers**:
- `Cookie`: auth_token (optional) - Authentication cookie for user identification

**Success Response (200)**:
- Returns the file content with appropriate Content-Type header
- File is served directly for download or display

**Error Responses**:
- `400 Bad Request`: Invalid file path
  ```json
  {
    "error": "Cannot get file"
  }
  ```
- `404 Not Found`: File not found (served by web server)

**Example**:
```bash
GET /storage/file/documents/2024/report.pdf
```

---

### POST /storage/file/*file_path
**Description**: Upload a file to storage with chunked upload support for large files

**Path Parameters**:
- `file_path` (string, required): The destination file path (supports nested paths, e.g., `/documents/2024/report.pdf`)

**Form Data (multipart/form-data)**:
- `file_id` (string, required): Unique identifier for the file (used for chunked uploads)
- `chunk_index` (integer, required): Index of current chunk (0-based)
- `total_chunks` (integer, required): Total number of chunks for this file
- `chunk_data` (file, required): The file chunk data

**Headers**:
- `Cookie`: auth_token (optional) - Authentication cookie for user identification
- `Content-Type`: multipart/form-data

**Success Response (200)**:
For non-final chunks:
```json
{
  "message": "Chunk uploaded successfully"
}
```

For final chunk (file complete):
```json
{
  "message": "File uploaded successfully"
}
```

**Error Responses**:
- `400 Bad Request`: Invalid request parameters
  ```json
  {
    "error": "Cannot upload file"
  }
  ```
- `400 Bad Request`: Missing chunk data
  ```json
  {
    "error": "Missing chunk_data"
  }
  ```
- `500 Internal Server Error`: Upload failed
  ```json
  {
    "error": "Failed to save file"
  }
  ```

**Chunked Upload Process**:
1. Split large files into chunks (recommended: 1-10MB per chunk)
2. Upload each chunk with the same `file_id` and sequential `chunk_index`
3. Server automatically merges chunks when the final chunk is received
4. Temporary chunks are stored in `tmp/upload_chunks/{file_id}/` during upload

**Example**:
```bash
# Upload chunk 0 of 3
POST /storage/file/documents/large_video.mp4
Content-Type: multipart/form-data

file_id=unique_file_123
chunk_index=0
total_chunks=3
chunk_data=<binary_chunk_0>

# Upload chunk 1 of 3
POST /storage/file/documents/large_video.mp4
Content-Type: multipart/form-data

file_id=unique_file_123
chunk_index=1
total_chunks=3
chunk_data=<binary_chunk_1>

# Upload final chunk 2 of 3 (triggers merge)
POST /storage/file/documents/large_video.mp4
Content-Type: multipart/form-data

file_id=unique_file_123
chunk_index=2
total_chunks=3
chunk_data=<binary_chunk_2>
```

---

### PATCH /storage/file/*file_path
**Description**: Update/move a file to a new location

**Path Parameters**:
- `file_path` (string, required): The current file path to update

**Request Body**:
```json
{
  "path": "/new/location/filename.ext"
}
```

**Request Body Schema**:
- `path` (string, optional): New file path (relative to storage root)

**Headers**:
- `Cookie`: auth_token (optional) - Authentication cookie for user identification

**Success Response (200)**:
```json
{
  "message": "File updated successfully"
}
```

**Error Responses**:
- `400 Bad Request`: Invalid request payload
  ```json
  {
    "error": "Invalid request payload"
  }
  ```
- `400 Bad Request`: Invalid new file path
  ```json
  {
    "error": "Invalid new file path"
  }
  ```
- `500 Internal Server Error`: Failed to update file
  ```json
  {
    "error": "Failed to update file"
  }
  ```
- `500 Internal Server Error`: Failed to move file
  ```json
  {
    "error": "Failed to move file"
  }
  ```

**Example**:
```bash
PATCH /storage/file/documents/old_name.txt
Content-Type: application/json

{
  "path": "/documents/renamed_file.txt"
}
```

---

### DELETE /storage/file/*file_path
**Description**: Delete a file from storage

**Path Parameters**:
- `file_path` (string, required): The file path to delete

**Headers**:
- `Cookie`: auth_token (optional) - Authentication cookie for user identification

**Success Response (200)**:
```json
{
  "message": "File deleted successfully"
}
```

**Error Responses**:
- `400 Bad Request`: Invalid file path
  ```json
  {
    "error": "Cannot delete file"
  }
  ```
- `500 Internal Server Error`: Failed to delete file
  ```json
  {
    "error": "Failed to delete file"
  }
  ```

**Example**:
```bash
DELETE /storage/file/documents/unwanted_file.txt
```

---

## Storage Notes

- All folder and file paths support nested directory structures
- Authentication is optional for storage operations.
- User have their own storage if they logged in and they share a storage with other users if they did not log in.
- Folder and file names are case-sensitive
- The `*folder_path` and `*file_path` parameters capture the entire path after `/folder/` or `/file/`

## Battle Cat APIs

### GET /battle-cat/levels
**Description**: Filter and list Battle Cat levels by stage and up to 3 enemies. Returns an array of level collections. Internally, results may include matches for 3-enemy, 2-enemy, and 1-enemy combinations.

**Query Parameters**:
- `stage` (string, required, max length 3): Stage identifier to filter on
- `enemy` (string, required, repeated, max 3): Enemy names; pass 1 to 3 values. Only give a part of the name is not functional.
  - Examples: `?enemy=臭老兔&enemy=狗仔`
  - Wrong: `?enemy=兔` (This will return no error but you can't find the levels with an enemy named "臭老兔")

**Success Response (200)**:
Returns an array of collections. Each collection contains the requested enemies echo and a list of levels.
```json
[
  {
    "enemies": ["enemyA", "enemyB", "enemyC"],
    "levels": [
      { "level": "001", "name": "Stage 1", "hp": 1200, "enemies": "enemyA, enemyX" },
      { "level": "002", "name": "Stage 2", "hp": 1500, "enemies": "enemyB, enemyY" }
    ]
  }
]
```

**Response Schema**:
- Array of objects with:
  - `enemies` (string[]): The enemies from the query (echoed)
  - `levels` (array): Matched levels for a particular enemy combination
    - `level` (string): Level code/id
    - `name` (string): Level name
    - `hp` (number): Level HP
    - `enemies` (string): Original enemies string from DB for that level

**Error Responses**:
- `401 Unauthorized`: Invalid query parameters (binding/validation failed)
  ```json
  { "error": "Key: 'FilterLevelsRequest.Enemies' Error:Field validation for 'Enemies' failed on the 'max' tag" }
  ```

**Examples**:
```bash
# Single enemy
GET /battle-cat/levels?stage=E01&enemy=dog

# Two enemies (repeat the query parameter)
GET /battle-cat/levels?stage=E01&enemy=dog&enemy=snake

# Three enemies
GET /battle-cat/levels?stage=E01&enemy=dog&enemy=snake&enemy=boss
```

Notes:
- The endpoint may return multiple collections representing matches for three-enemy, pairwise, and single-enemy filters.

## Error Handling

All endpoints return appropriate HTTP status codes:
- `200`: Success
- `400`: Bad Request (validation errors)
- `401`: Unauthorized (authentication required or failed)
- `403`: Forbidden (insufficient permissions)
- `500`: Internal Server Error

## Notes

- All passwords must be at least 8 characters long
- Email addresses must be in valid email format
- Password changes are only allowed for accounts created with email/password (not OAuth accounts)

---

## Reurl APIs

**Description**:
The Reurl APIs allow users to create, manage, and use short URL redirects. Users can create short keys that redirect to target URLs with optional expiration times. Keys are automatically generated using a short base62 encoding. Regular users can only manage their own reurls, while admins can manage all reurls.

### GET /reurl
**Description**: List reurls owned by the authenticated user. Admins see all reurls.

**Headers**:
- `Cookie`: auth_token (required) - Authentication cookie

**Success Response (200)**:
```json
[
  {
    "id": 1,
    "key": "abc123",
    "target_url": "https://example.com",
    "expires_at": "2025-11-17T10:00:00Z",
    "owner_id": 1,
    "created_at": "2025-11-10T10:00:00Z",
    "updated_at": "2025-11-10T10:00:00Z"
  }
]
```

**Response Schema**:
- Array of reurl objects:
  - `id` (uint): Reurl ID
  - `key` (string): Short key for redirect
  - `target_url` (string): Destination URL
  - `expires_at` (string, nullable): Expiration timestamp (ISO 8601)
  - `owner_id` (uint): Owner user ID
  - `created_at` (string): Creation timestamp
  - `updated_at` (string): Last update timestamp

**Error Responses**:
- `401 Unauthorized`: Missing or invalid authentication
  ```json
  {
    "error": "Unauthorized"
  }
  ```

---

### POST /reurl
**Description**: Create a new reurl with auto-generated key

**Request Body**:
```json
{
  "target_url": "https://example.com",
  "expires_in": "7d",
  "key": "123"
}
```

**Request Body Schema**:
- `target_url` (string, required): The URL to redirect to (must be valid URL)
- `expires_in` (string, optional): Expiration duration. Allowed values: "1h", "12h", "1d", "7d", "30d". Defaults to "7d"
- `key` (string, optional): Everyone including who not authrized can access the reurl record by a key. The server will generate a key if not found in request body.

**Headers**:
- `Cookie`: auth_token (required) - Authentication cookie

**Success Response (201)**:
```json
{
  "id": 1,
  "key": "abc123",
  "target_url": "https://example.com",
  "expires_at": "2025-11-17T10:00:00Z",
  "owner_id": 1,
  "created_at": "2025-11-10T10:00:00Z",
  "updated_at": "2025-11-10T10:00:00Z"
}
```

**Error Responses**:
- `400 Bad Request`: Invalid input data
  ```json
  {
    "error": "Key: 'CreateReurlRequest.TargetURL' Error:Field validation for 'TargetURL' failed on the 'url' tag"
  }
  ```
- `400 Bad Request`: Invalid expires_in value
  ```json
  {
    "error": "Invalid expires_in value"
  }
  ```
- `401 Unauthorized`: Missing or invalid authentication
  ```json
  {
    "error": "Unauthorized"
  }
  ```
- `500 Internal Server Error`: Failed to generate key or create reurl
  ```json
  {
    "error": "Failed to create reurl"
  }
  ```

---

### GET /reurl/:id
**Description**: Get a specific reurl by ID. Users can only access their own reurls, admins can access any.

**Path Parameters**:
- `id` (uint, required): Reurl ID

**Headers**:
- `Cookie`: auth_token (required) - Authentication cookie

**Success Response (200)**:
```json
{
  "id": 1,
  "key": "abc123",
  "target_url": "https://example.com",
  "expires_at": "2025-11-17T10:00:00Z",
  "owner_id": 1,
  "created_at": "2025-11-10T10:00:00Z",
  "updated_at": "2025-11-10T10:00:00Z"
}
```

**Error Responses**:
- `401 Unauthorized`: Missing or invalid authentication
  ```json
  {
    "error": "Unauthorized"
  }
  ```
- `403 Forbidden`: Attempting to access another user's reurl (non-admin)
  ```json
  {
    "error": "Forbidden"
  }
  ```
- `404 Not Found`: Reurl not found
  ```json
  {
    "error": "Reurl not found"
  }
  ```

---

### PATCH /reurl/:id
**Description**: Update a reurl. Users can only update their own reurls, admins can update any.

**Path Parameters**:
- `id` (uint, required): Reurl ID

**Request Body**:
```json
{
  "target_url": "https://new-example.com",
  "expires_in": "1d",
  "key": "new"
}
```

**Request Body Schema**:
- `target_url` (string, optional): New target URL (must be valid URL if provided)
- `expires_in` (string, optional): New expiration duration. Allowed values: "1h", "12h", "1d", "7d", "30d"
- `key` (string, optional): Everyone including who not authrized can access the reurl record by a key.

**Headers**:
- `Cookie`: auth_token (required) - Authentication cookie

**Success Response (200)**:
```json
{
  "id": 1,
  "key": "abc123",
  "target_url": "https://new-example.com",
  "expires_at": "2025-11-11T10:00:00Z",
  "owner_id": 1,
  "created_at": "2025-11-10T10:00:00Z",
  "updated_at": "2025-11-10T10:00:00Z"
}
```

**Error Responses**:
- `400 Bad Request`: Invalid input data
  ```json
  {
    "error": "Invalid expires_in value"
  }
  ```
- `401 Unauthorized`: Missing or invalid authentication
  ```json
  {
    "error": "Unauthorized"
  }
  ```
- `403 Forbidden`: Attempting to update another user's reurl (non-admin)
  ```json
  {
    "error": "Forbidden"
  }
  ```
- `404 Not Found`: Reurl not found
  ```json
  {
    "error": "Reurl not found"
  }
  ```
- `500 Internal Server Error`: Failed to update reurl
  ```json
  {
    "error": "Failed to update reurl"
  }
  ```

---

### DELETE /reurl/:id
**Description**: Delete a reurl. Users can only delete their own reurls, admins can delete any.

**Path Parameters**:
- `id` (uint, required): Reurl ID

**Headers**:
- `Cookie`: auth_token (required) - Authentication cookie

**Success Response (200)**:
```json
{
  "message": "Reurl deleted successfully"
}
```

**Error Responses**:
- `401 Unauthorized`: Missing or invalid authentication
  ```json
  {
    "error": "Unauthorized"
  }
  ```
- `403 Forbidden`: Attempting to delete another user's reurl (non-admin)
  ```json
  {
    "error": "Forbidden"
  }
  ```
- `404 Not Found`: Reurl not found
  ```json
  {
    "error": "Reurl not found"
  }
  ```
- `500 Internal Server Error`: Failed to delete reurl
  ```json
  {
    "error": "Failed to delete reurl"
  }
  ```

---

### GET /reurl/redirect/:key
**Description**: Redirect to the target URL using the short key. This endpoint is public and does not require authentication.

**Path Parameters**:
- `key` (string, required): Short key for the reurl

**Success Response (302)**:
- Redirects to the target URL with 302 Found status

**Error Responses**:
- `404 Not Found`: Reurl not found or expired
  ```json
  {
    "error": "Reurl not found or expired"
  }
  ```

**Example**:
```bash
GET /reurl/redirect/abc123
# Redirects to https://example.com
```

---

## Reurl Notes

- Expired reurls are not accessible and will return 404 on redirect
- Users can only manage their own reurls unless they have admin role
- The redirect endpoint is public and can be shared freely
- Expiration times are calculated from creation/update time

---

## Posts / Comments / Reactions APIs

This section describes the post (article) system: creating and managing posts, nested comments, and reaction types (like, love, haha, wow, sad, angry, care).

Notes:
- Authentication: most write endpoints require authentication (the server uses an HTTP-only `auth_token` cookie). Read endpoints are usually `AuthOptional`.
- Visibility: posts support `public` and `private` visibility. Private posts are visible only to the author.

---

### GET /posts
**Description**: List posts with simple filtering and pagination. Authentication optional.

**Query Parameters**:
- `page` (int, optional): page number (default 1)
- `limit` (int, optional): items per page (default 20)
- `author_id` (uint, optional): filter by author
- `tag` (string, optional): filter by tag slug
- `status` (string, optional): `published|draft|archived` (only authors/admins can see drafts)

**Success Response (200)**:
```json
{
  "page": 1,
  "limit": 20,
  "total": 123,
  "posts": [
    {
      "id": 10,
      "author_id": 1,
      "title": "Hello world",
      "summary": "Short summary",
      "cover_image": "/uploads/cover.jpg",
      "status": "published",
      "visibility": "public",
      "view_count": 42,
      "published_at": "2025-11-10T12:00:00Z",
      "tags": ["news","update"]
    }
  ]
}
```

**Error Responses**:
- `400 Bad Request`: invalid query parameters

---

### POST /posts
**Description**: Create a new post. Requires authentication.

**Request Body**:
```json
{
  "title": "My post title",
  "content": "Full HTML or markdown content",
  "summary": "Optional short summary",
  "cover_image": "/uploads/cover.jpg",
  "status": "draft", // draft | published | archived
  "visibility": "public", // public | private
  "tags": ["go","programming"]
}
```

**Request Body Schema**:
- `title` (string, required)
- `content` (string, required)
- `summary` (string, optional)
- `cover_image` (string, optional)
- `status` (string, optional): `draft|published|archived` (default `draft`)
- `visibility` (string, optional): `public|private` (default `public`)
- `tags` (string[], optional)

**Success Response (201)**:
```json
{
  "id": 11,
  "author_id": 1,
  "title": "My post title",
  "status": "draft",
  "visibility": "public",
  "published_at": null
}
```

**Error Responses**:
- `400 Bad Request`: validation error
- `401 Unauthorized`: missing/invalid authentication
- `500 Internal Server Error`: DB error

---

### GET /posts/:id
**Description**: Get a single post by ID. Authentication optional. Private posts return 403 for non-author.

**Path Parameters**:
- `id` (uint, required): Post ID

**Success Response (200)**:
```json
{
  "id": 11,
  "author_id": 1,
  "title": "My post title",
  "content": "Full content",
  "summary": "Optional",
  "cover_image": "/uploads/cover.jpg",
  "status": "published",
  "visibility": "public",
  "view_count": 123,
  "published_at": "2025-11-10T12:00:00Z",
  "tags": ["go","programming"],
  "comments_count": 5,
  "reactions_summary": {"like": 10, "love": 2}
}
```

**Error Responses**:
- `401 Unauthorized` / `403 Forbidden`: access denied for private posts
- `404 Not Found`: post not found

---

### PUT /posts/:id
**Description**: Update a post. Requires authentication and must be the author (or admin).

**Path Parameters**:
- `id` (uint, required): Post ID

**Request Body**: same as `POST /posts` (partial updates allowed)

**Success Response (200)**:
```json
{
  "message": "Post updated",
  "id": 11
}
```

**Error Responses**:
- `400 Bad Request`: validation error
- `401 Unauthorized`: not logged in
- `403 Forbidden`: not the author
- `404 Not Found`: post not found

---

### DELETE /posts/:id
**Description**: Delete a post. Requires authentication and must be the author or an admin. This performs a soft delete.

**Path Parameters**:
- `id` (uint, required): Post ID

**Success Response (200)**:
```json
{
  "message": "Post deleted"
}
```

**Error Responses**:
- `401 Unauthorized`, `403 Forbidden`, `404 Not Found`

---

### GET /posts/:id/comments
**Description**: List comments for a post. Returns nested replies. Authentication optional.

**Path Parameters**:
- `id` (uint, required): Post ID

**Success Response (200)**:
```json
[
  {
    "id": 1,
    "post_id": 11,
    "author_id": 2,
    "content": "Top-level comment",
    "parent_id": null,
    "replies": [
      { "id": 2, "parent_id": 1, "content": "Reply" }
    ],
    "is_deleted": false
  }
]
```

---

### POST /posts/:id/comments
**Description**: Create a comment on a post. Requires authentication.

**Path Parameters**:
- `id` (uint, required): Post ID

**Request Body**:
```json
{
  "content": "Nice post!",
  "parent_id": null
}
```

**Success Response (201)**:
```json
{
  "id": 21,
  "post_id": 11,
  "author_id": 3,
  "content": "Nice post!",
  "parent_id": null
}
```

**Error Responses**:
- `400 Bad Request`, `401 Unauthorized`, `404 Not Found` (post not found)

---

### PUT /comments/:comment_id
**Description**: Update a comment. Requires authentication and comment owner.

**Request Body**:
```json
{
  "content": "Updated content"
}
```

**Success Response (200)**:
```json
{
  "message": "Comment updated",
  "id": 21
}
```

---

### DELETE /comments/:comment_id
**Description**: Soft-delete a comment. Requires authentication and comment owner or post author.

**Success Response (200)**:
```json
{
  "message": "Comment deleted"
}
```

---

### POST /posts/:id/reactions
**Description**: Add or toggle a reaction to a post. Requires authentication.

**Path Parameters**:
- `id` (uint, required): Post ID

**Request Body**:
```json
{ "type": "like" }
```

**Behavior**:
- If the user has no reaction on the target, the reaction is created.
- If the user has the same reaction, it is removed (toggle off).
- If the user has a different reaction, it is updated to the new type.

**Success Response (200)**:
```json
{ "message": "Reaction updated", "summary": {"like": 10, "love": 2} }
```

---

### POST /comments/:comment_id/reactions
**Description**: Add or toggle a reaction to a comment. Requires authentication. Same behavior as post reactions.

---

### GET /posts/:id/reactions
**Description**: Get reaction summary for a post. Authentication optional.

**Success Response (200)**:
```json
{ "like": 10, "love": 2, "haha": 0 }
```

---

### GET /comments/:comment_id/reactions
**Description**: Get reaction summary for a comment.

**Success Response (200)**:
```json
{ "like": 2, "love": 0 }
```

---

## Posts API Notes
- Reactions types: `like`, `love`, `haha`, `wow`, `sad`, `angry`, `care`.
- Comments use soft-delete: deleted comments keep place in thread but marked as deleted.
- Posts with `status: draft` are visible only to the author and admins.
- All write operations require authentication; read endpoints are often available without auth but hide private content.
