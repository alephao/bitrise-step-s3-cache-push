# S3 Cache Push

A bitrise step to store your cache in a s3 bucket with custom keys.

Should be used with [S3 Cache Pull](https://github.com/alephao/bitrise-step-s3-cache-pull)

### Inputs

- **aws_access_key_id**: Your aws access key id. You can set an environment var `AWS_ACCESS_KEY_ID` instead of using this input.
- **aws_secret_access_key**: Your aws secret access key. You can set an environment var `AWS_SECRET_ACCESS_KEY` instead of using this input. 
- **aws_region**: The region of your S3 bucket. E.g.: `us-east-1 `. You can set an environment var `AWS_S3_REGION` instead of using this input.
- **bucket_name**: The name of your S3 bucket. E.g.: `mybucket`. You can set an environment var `S3_BUCKET_NAME` instead of using this input.
- **path**: The path to the file or folder you want to cache. E.g.: `Carthage`
- **key**: The key that will be used to restore the cache later. E.g.: `carthage-{{ checksum "Cartfile.resolved" }}`. You can use `{{ checksum "path/to/file" }}` to use the checksum of a file as part of the cache key.