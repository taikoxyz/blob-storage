# blob-storage
Repo for BLOB storage (archive and serve data)

## how to run ?

Prerequisite is to have docker engine up and running.

1. Start the mongoDB
```bash
cd local_docker && docker-compose up -d
```

2. Start the `blob-catcher`
```bash
cd cmd/blob_catcher && go run .
```

It uses the config from `internal/logic/config.go`.

## how to test ?

When the `DB` and `blob-catcher` is running, the `blob-catcher` is outputting the `blobHash` to the terminal. Use that `blobHash` (including the 0x) in the python_query.py's `blob_hash` variable and run the script.

```bash
python3 python_query.py
```

## todos
What is still missing is:
- Server listening and serving incoming requests
- Proper DB connection (prod-grade DB with creditentials)
- Proper containerization 