# Project Requirements: JSON Document Database

## 1. Databases
- Any number of databases can exist. Each database is an independent logical unit and is identified by a unique name or ID.
- Each database gets its own folder: `/data/{databasename}/`
- Database management (create, delete, list) is done via the server API.

## 2. Tables
- Any number of tables can exist within each database. Tables serve as logical groupings of documents.
- Each table has its own subfolder: `/data/{databasename}/{tablename}/`
- Each table has its own index and its own storage location for documents.
- The documents of a table are stored as individual JSON files: `/data/{databasename}/{tablename}/entries/{document_id}.json`
- Each table has its own metadata file in JSON format: `/data/{databasename}/{tablename}/meta.json`

## 3. Fields and Supported Data Types
- Supported data types for fields:
    - string (text)
    - number (integers and floating-point numbers)
    - boolean (true/false)
    - json (JSON objects)
    - date (ISO 8601, stored as string)
- Fields are described in the metadata file (`meta.json`) as an array of objects with name, type, and optional properties (e.g., required, default, description).
- Optionally, further properties or constraints for values can be defined, e.g.:
    - `minLength`, `maxLength` (for string, number, and json)
    - `nullable` (for all types, allows null values)
    - `pattern` (for string, regex)
    - `enum` (for string, number, boolean)

**Example field definitions with data types and optional properties:**

- string:
  ```json
  { "name": "username", "type": "string", "required": true, "minLength": 3, "maxLength": 20, "nullable": false }
  ```
- number:
  ```json
  { "name": "age", "type": "number", "minLength": 1, "maxLength": 3, "nullable": true }
  ```
- boolean:
  ```json
  { "name": "isActive", "type": "boolean", "default": false, "nullable": false }
  ```
- json:
  ```json
  { "name": "address", "type": "json", "maxLength": 5, "nullable": true }
  ```
- date:
  ```json
  { "name": "created_at", "type": "date", "nullable": false }
  ```

## 4. Index Engine (Table-Based)
Index search is performed in memory (in-memory) for fast access. In addition, the index of each table is regularly persisted to a file so that it can be restored after a restart.

**Index file format and number of index files:**
- Each table has at least one index file for the primary key: `/data/{databasename}/{tablename}/index.jsonl`.
- For each defined secondary index (e.g., on a specific field), an additional index file is created, e.g.: `/data/{databasename}/{tablename}/index_{fieldname}.jsonl`.
- The number of index files per table is: 1 (primary index) + number of secondary indexes.
- Each index file is in `.jsonl` format (JSON Lines), where each line represents an index entry as a JSON object.
- This format is easy to parse, compatible with many tools, and allows efficient sequential reading and writing.
- The index is always table-based, i.e., each table manages its own index files independently of other tables.

## 5. Filter & Query

**Query types:**
- Support for simple filters (e.g., search for field values such as equality, range, substring).
- Support for complex filters with logical operators (AND, OR, NOT).
- Support for joins across multiple tables.
- Comparison operators: =, !=, <, <=, >, >=, IN, NOT IN, LIKE/Pattern.

**Search methods and text search:**
- For structured filters (equality, range, IN, etc.), direct comparisons are made on the field values, supported by indexes if available.
- For LIKE/pattern searches, a substring search is used by default.
- For advanced text search (e.g., search for individual words, word stems, phrases), the introduction of a tokenizer and an inverted index is planned for the future (similar to PostgreSQL or Elasticsearch). In the first version, text search is performed without a tokenizer.
- Full-text search, stemming, and stopword filters are planned as optional extensions and will be documented as future features.
- The supported operators are based on common standards from SQL and document-based databases (e.g., =, !=, <, <=, >, >=, IN, NOT IN, LIKE/Pattern).

**Index usage:**
- Filters on indexed fields always use the corresponding index for the search.
- Filters on non-indexed fields result in a full table scan of all entries.

**Query syntax (API):**
- Queries are passed to the API as a JSON object, e.g.:
  ```json
  {
    "table": "customers",
    "filter": {
      "age": { "gte": 18, "lte": 65 },
      "isActive": true
    },
    "joins": [
      {
        "table": "orders",
        "on": { "customers.id": "orders.customer_id" },
        "type": "inner",
        "filter": { "status": "open" },
        "joins": [
          {
            "table": "products",
            "on": { "orders.product_id": "products.id" },
            "type": "left",
            "filter": { "category": "digital" }
          }
        ]
      }
    ],
    "sort": [{ "field": "created_at", "direction": "desc" }],
    "limit": 20,
    "offset": 0
  }
  ```
- The `joins` field is recursive: Each join definition can itself contain a `joins` array to allow arbitrarily deep join hierarchies.
- The maximum depth for nested joins is 64 levels (`maxJoinDepth = 64`).
- Each join definition consists of:
    - `table`: Name of the target table
    - `on`: Join condition (field in main table → field in target table)
    - `type`: Join type (e.g., inner, left; optional, default: inner)
    - `filter`: Optional filter on the joined table
    - `joins`: Optional, array of further joins at this level
- Support for sorting, limiting, and pagination.

**Examples:**
- Simple query: All active customers over 30 years old
  ```json
  {
    "table": "customers",
    "filter": {
      "isActive": true,
      "age": { "gt": 30 }
    }
  }
  ```
- Complex query with join: Customers with open orders, sorted by creation date
  ```json
  {
    "table": "customers",
    "joins": [
      {
        "table": "orders",
        "on": { "customers.id": "orders.customer_id" },
        "type": "inner",
        "filter": { "status": "open" }
      }
    ],
    "sort": [{ "field": "created_at", "direction": "desc" }],
    "limit": 20
  }
  ```
- Complex query with nested joins:
  ```json
  {
    "table": "customers",
    "joins": [
      {
        "table": "orders",
        "on": { "customers.id": "orders.customer_id" },
        "type": "inner",
        "filter": { "status": "open" },
        "joins": [
          {
            "table": "products",
            "on": { "orders.product_id": "products.id" },
            "type": "left",
            "filter": { "category": "digital" }
          }
        ]
      }
    ],
    "sort": [{ "field": "created_at", "direction": "desc" }],
    "limit": 20
  }
  ```

## 6. Storage Structure
The actual data contents (documents) are stored persistently only as individual JSON files and are not kept in RAM. Only the index data of the respective tables is kept in memory (in-memory) to enable fast searches and access.

**JSON as the primary input and output format:**
- JSON is the central data format for all server input and output.
- All requests to the API (e.g., to create, update, query, or delete documents) as well as all server responses are in JSON format.
- Validation, storage, and transmission of data are fully based on JSON.

**ACID compliance:**
- The system MUST guarantee the ACID principles (Atomicity, Consistency, Isolation, Durability) for all write and read operations. In particular, this means:
    - Write operations are atomic and either complete fully or not at all.
    - Data integrity is maintained in all operations (consistency).
    - Concurrent accesses do not interfere with each other (isolation).
    - After completion of an operation, the data is permanently stored and survives system crashes (durability).

**Validation and error handling:**
- Before saving or updating a document, comprehensive validation is performed against the schema defined in `meta.json`.
    - All required fields (`required: true`) must be present and valid.
    - Data types must exactly match the definition (e.g., a field of type `string` must not contain a number).
    - Constraints such as `minLength`, `maxLength`, `min`, `max`, `pattern`, `enum`, and `nullable` are strictly checked.
    - Fields not defined in the schema are rejected by default unless explicitly allowed.
    - Default values (`default`) are set if a value is missing and the field is not `required`.
- In case of validation errors, the write operation is aborted and a detailed error message with all violations is returned.
- Error codes and messages are consistent and machine-readable (e.g., `ERR_VALIDATION_FAILED`, `ERR_TYPE_MISMATCH`, `ERR_MISSING_REQUIRED_FIELD`).
- Faulty or incomplete documents are never saved or overwrite existing data.
- All error cases are logged to ensure traceability and debugging.

**Schema versioning:**
- Each table contains a `schemaVersion` field in its meta.json (e.g., `"schemaVersion": 1`).
- Changes to the table schema are tracked by increasing the version number.
- When the schema changes, a migration strategy must be defined (e.g., automatic adjustment of existing documents or explicit migration by the user).

**Handling unknown fields:**
- By default, documents may not contain fields that are not defined in the schema (`allowAdditionalFields: false`).
- Optionally, it can be explicitly allowed in meta.json to store additional fields (`allowAdditionalFields: true`).

**Index options:**
- Index definitions in meta.json can have optional properties such as `unique` (uniqueness) and `sparse` (only for existing values).
- Example: `{ "name": "email", "type": "secondary", "fields": ["email"], "unique": true, "sparse": false }`

**Soft deletes and history:**
- Optionally, a field such as `deleted: true` can be used for soft deletes, so that deleted documents are not physically removed but marked as deleted.
- For change history, an audit trail strategy can be defined (e.g., storing changes with timestamp and user ID).

**Backup and recovery strategy:**
- It is recommended to regularly back up the data and index files.
- Recovery is done by restoring the backup files to the appropriate database structure.

**API error codes and responses:**
- The API returns standardized error codes and messages (e.g., 400 for validation errors, 404 for not found resources, 409 for conflicts).
- Error responses contain machine-readable codes and a human-readable description.

**Security and access control:**
- Optionally, authentication (e.g., API keys) and authorization (e.g., roles, rights at DB/table level) can be implemented.
- Access restrictions can be defined in the server configuration or in meta.json.

**Performance and scalability options:**
- The system should be optimized for large amounts of data and many concurrent accesses (e.g., through efficient indexes, caching, possibly sharding).
- Notes on scaling and recommended hardware requirements can be documented.

**Document limitation:**
- Optionally, a maximum number of documents per table or database can be set.
- If a table exceeds the limit, further write operations are rejected and an appropriate error code is returned.

**Logging and monitoring:**
- The system logs all relevant events (e.g., errors, accesses, changes) in machine-readable log files.
- Monitoring interfaces (e.g., health checks, metrics) can be provided.

**Transaction support:**
- The system can group multiple operations into a transaction, which is either executed completely or not at all (all-or-nothing principle).
- Transactions are explicitly started and completed via the API.

## 7. Event Sourcing
- The system MUST support event sourcing, i.e., every change to a document (creation, update, deletion) is stored as a separate, immutable event.
- All events are stored chronologically and immutably, so that the complete state of a table or document can be restored at any time.
- Advantages of event sourcing:
    - Enables complete traceability and auditability of all changes (audit trail).
    - Increases fault tolerance: After a system failure, the state can be reconstructed exactly by replaying all events.
    - Supports advanced features such as time travel (querying the system state at any point in time), undo/redo, and flexible data migrations.
    - Facilitates integration with external systems via event streams.
- Events are stored persistently and are part of the backup and recovery strategy.

## 8. Secondary Indexes
Additional indexes are maintained for frequently queried fields within a table to speed up targeted searches.

## 9. Server
A server provides an API for creating, deleting, and managing databases and tables, as well as for storing and retrieving JSON documents.

## 10. Installation & Setup

**Database initialization:**
- Before the server is started for the first time, the data directory should be initialized. This can be done via a terminal command:
  ```sh
  ./coco-db init --data-dir=/path/to/datadirectory
  ```
- The command creates the necessary directory structure and checks whether the directory is writable.
- Optionally, a first database can be created directly during initialization (e.g., with `--database=mydb`).
- Initialization is a prerequisite before the server can be started in production.

**Starting the server (example):**
  ```sh
  ./coco-db --data-dir=/path/to/datadirectory
  ```
- The data directory MUST be specified explicitly (see storage structure).

**Notes:**
- Configuration can also be done via environment variables or a configuration file (see documentation).
- For production, it is recommended to run the data directory on a persistent and secure drive.
- Further start parameters and configuration options are described in the technical documentation.

---

**Note:** No code is written before the project plan and roadmap are finalized.