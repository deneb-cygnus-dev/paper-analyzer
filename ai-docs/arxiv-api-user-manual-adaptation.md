# ArXiv API Adaptation Details

This document details how the `ArxivFetcher` adapts the generic `FetchConfig` to the specific requirements of the arXiv API.

## API Reference

The arXiv API User Manual is available at [https://arxiv.org/help/api/user-manual](https://arxiv.org/help/api/user-manual).

## Configuration Mapping

The `ArxivFetcher` translates the `entities.FetchConfig` into an arXiv API query URL as follows:

### Base URL

The base URL for all queries is:
`http://export.arxiv.org/api/query?`

### Query Parameters

| FetchConfig Field | arXiv API Parameter | Description |
| :--- | :--- | :--- |
| `Category` | `search_query` | Mapped to `cat:<Category>`. Required. |
| `Keywords` | `search_query` | Appended to `search_query` as `AND all:<Keyword>` for each keyword. |
| `TimeSpan` | `search_query` | If set (e.g., "last_N_days"), calculates the start date and appends `AND submittedDate:[YYYYMMDDHHMM TO *]` to `search_query`. |
| `MaxResults` | `max_results` | Mapped directly to `max_results`. |
| N/A | `sortBy` | Always set to `submittedDate`. |
| N/A | `sortOrder` | Always set to `descending`. |

### Example

**Config:**

```go
entities.FetchConfig{
    Category:   "cs.SE",
    MaxResults: 10,
    Keywords:   []string{"fuzzing"},
    TimeSpan:   "last_5_days", // Assuming today is 2025-11-24
}
```

**Generated URL (Conceptual):**
`http://export.arxiv.org/api/query?search_query=cat:cs.SE+AND+all:fuzzing+AND+submittedDate:[202511190000+TO+*]&sortBy=submittedDate&sortOrder=descending&max_results=10`

## Implementation Details

### Validation

The `Fetch` method enforces the following validation rules:

1. **Category is Required**: The `Category` field must not be empty. Returns `ErrMissingRequiredField` (400002).
1. **Limit is Required**: Either `TimeSpan` or `MaxResults` (or both) must be specified to prevent fetching excessive data. Returns `ErrInvalidInput` (400001).

### URL Encoding

The implementation uses `net/url` to ensure all query parameters are properly encoded, handling special characters in keywords or categories correctly.

### Date Handling

For `TimeSpan` (e.g., "last_N_days"), the start date is calculated relative to the current server time (`time.Now()`). The time component is set to `0000` (midnight) for the start date.
