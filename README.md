# BackupService

This is a simple solution to push the contents of a specified folder to respective Minio buckets with similar names.

I used the latest version of Golang (1.18 as of writing this), so you may need to refactor accordingly.

### Disclosure
This repo uses some experimental features in Go 1.18. This is the [slices](https://pkg.go.dev/golang.org/x/exp/slices) package using generics on 1.18. For further info, please refer to [https://pkg.go.dev/golang.org/x/exp/slices](https://pkg.go.dev/golang.org/x/exp/slices)

### Set-up

- Clone this repo
    ```bash
      git clone https://github.com/benzerbett/backupservice
    ```
- Install dependencies
    ```bash
        cd backupservice
        go get
    ```
- Set up environment variables. The variables are:
  * MINIO_ENDPOINT
  * MINIO_ACCESSKEY
  * MINIO_SECRET
  * BACKUP_FILES_DIR
  ```bash
      nano .env
  ```
- To test it locally, simply run the ```main.go``` file
 ```bash
    go run main.go
```
- Build & deploy
 ```bash
    go build -o backupservice
```
To deploy, you can schedule this as a cron job or otherwise, depending on your environment.