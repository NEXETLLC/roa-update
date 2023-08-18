# Auto-update Prefix filter for Juniper using bgpq3.


This is a simple program writing in go to upload 'as-set-v4.juniper.update' to s3 bucket.

## How to use

1. Create a s3 bucket
2. Create a IAM user with s3 write permission
3. write a .env file with the following content:

```bash
MINIO_ENDPOINT=
MINIO_ACCESS_KEY=
MINIO_SECRET_KEY=
MINIO_BUCKET_NAME=
```


4. write a `update.lists` file with the following content:



```bash
as-set
as65535
as65536
```

Output should be like this:

```bash
go run . 
Successfully uploaded as-set-v4.juniper.update to `as-set-v4.juniper.update` of bucket `your-bucket-name`
```


## How to update in Juniper

```bash
load replace https://s3.amazonaws.com/your-bucket-name/as-set-v4.juniper.update
```