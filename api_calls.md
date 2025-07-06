# FORUM API calls

## 🚀 Curl
### HEALTH CHECK
```sh
curl --request GET \
  --url http://localhost:8080/api/v1/health \
  --header 'Accept: */*' \
  --header 'Accept-Encoding: gzip, deflate, br' \
  --header 'Connection: keep-alive' \
  --header 'Content-Type: application/json'
```