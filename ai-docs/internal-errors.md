# Internal Errors

This document describes the internal error handling system and the list of custom error codes used in the paper analyzer.

## Overview

The system uses a structured error handling mechanism defined in `internal/pkg/errors`. Each error is represented by a `CustomError` struct containing a unique 6-digit error code, a static description, and an optional underlying error.

## Error Structure

```go
type CustomError struct {
    Code    int
    Message string
    Err     error
}
```

Errors are formatted as `[Code] Message: Underlying Error` (if an underlying error exists).

## Error Hierarchy

Error codes are 6 digits long and categorized by the first two digits, following a standard convention similar to HTTP status codes.

### Categories

| Range | Category | Description |
| :--- | :--- | :--- |
| **10xxxx** | General / Internal | Unexpected system errors or unimplemented features. |
| **20xxxx** | Authentication | Issues related to user identity and authentication tokens. |
| **30xxxx** | Authorization | Issues related to permissions and access control. |
| **40xxxx** | Validation / Input | Errors caused by invalid client input or missing parameters. |
| **50xxxx** | Infrastructure | Errors related to external systems, databases, networks, or third-party APIs. |

## Error Codes

### General / Internal (10xxxx)

| Code | Variable | Message |
| :--- | :--- | :--- |
| `100001` | `ErrInternalServer` | Internal server error occurred. |
| `100002` | `ErrNotImplemented` | Feature not implemented. |

### Authentication (20xxxx)

| Code | Variable | Message |
| :--- | :--- | :--- |
| `200001` | `ErrUnauthorized` | User is not authorized. |
| `200002` | `ErrTokenInvalid` | Authentication token is invalid. |
| `200003` | `ErrTokenExpired` | Authentication token has expired. |

### Authorization (30xxxx)

| Code | Variable | Message |
| :--- | :--- | :--- |
| `300001` | `ErrForbidden` | Access to resource is forbidden. |
| `300002` | `ErrInsufficientPermissions` | User has insufficient permissions. |

### Validation / Input (40xxxx)

| Code | Variable | Message |
| :--- | :--- | :--- |
| `400001` | `ErrInvalidInput` | Input parameters are invalid. |
| `400002` | `ErrMissingRequiredField` | Required field is missing. |

### Infrastructure (50xxxx)

| Code | Variable | Message |
| :--- | :--- | :--- |
| `500001` | `ErrDatabase` | Database operation failed. |
| `500002` | `ErrRecordNotFound` | Requested record was not found. |
| `500003` | `ErrDuplicateRecord` | Record already exists. |
| `500004` | `ErrNetwork` | Network communication failed. |
| `500005` | `ErrTimeout` | Operation timed out. |
| `500006` | `ErrExternalAPI` | External API returned an error. |
| `500007` | `ErrExternalAPIParsing` | Failed to parse response from external API. |

## Usage

### Creating Errors

Use `errors.New` to create a new custom error:

```go
return errors.New(100001, "Something went wrong")
```

### Wrapping Errors

Use `errors.Wrap` to add context to an existing error while preserving the error code:

```go
if err != nil {
    return errors.Wrap(err, errors.ErrDatabase)
}
```

### Checking Errors

Use `errors.Is` or check the code directly:

```go
if err != nil {
    var customErr *errors.CustomError
    if errors.As(err, &customErr) {
        if customErr.Code == 500002 {
            // Handle record not found
        }
    }
}
```
