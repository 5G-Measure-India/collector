# 5G Logs Collector

## Getting Started

1. Build docker image

```sh
docker build -t collector:dev .
```

2. Create docker volume to store adb keys

```sh
docker volume create adb_keys
```

3. Setup and auth adb key

```sh
docker run --rm \
  --volume adb_keys:/root/.android \
  --device /dev/bus/usb \
  --entrypoint adb \
  collector:dev devices
```

> _Repeat until you see an authorized device in output_

> _If you receive a pop-up on your phone screen, check "Always allow ..." and click on "Allow" to proceed_

4. Switch to diag mode

```sh
docker run --rm \
  --volume adb_keys:/root/.android \
  --device /dev/bus/usb \
  --entrypoint adb \
  collector:dev shell su -c setprop sys.usb.config diag,adb
```

5. Run

```sh
docker run \
  --rm --detach \
  --volume adb_keys:/root/.android \
  --volume ./data:/app/data \
  --device /dev/bus/usb \
  --device /dev/ttyUSB0 \
  --env TZ=Asia/Kolkata \
  --name collector \
  collector:dev
```

6. Stop

```sh
docker stop collector
```
