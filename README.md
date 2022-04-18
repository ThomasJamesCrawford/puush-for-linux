# Puush-on-linux

A puush clone for linux that uses your own S3 bucket

This will let you take a SS of a selected area, then upload it to S3 and copy a presigned URL (with a 24hr expiry into the clipboard)

`puush-on-linux -b my-bucket-name`

On ubuntu you can bind this to a hotkey easily

`Settings -> Keyboard -> Keyboard Shortcuts -> Custom Shortcuts`

## dependencies

- `slop`
- `notify-send`
- `~.aws/config` access keys setup