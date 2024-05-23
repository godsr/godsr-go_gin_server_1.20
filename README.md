.env 파일 생성 후

<!-- #DB & REDIS INFO -->

DB="postgres://{DB_ID}:{DB_PW}@{DB_URL}:{DB_PORT}/{DB_NAME}" #DB 정보
PORT=":8080" #포트
REDIS_DSN = "127.0.0.1:6379" #redis 주소

<!-- #ETC -->

HASH_SALT = "SALT" #hash salt 값

<!-- # TOKEN INFO -->

ACCESS_SECRET = "accessToken secret key" #access token secret key
REFRESH_SECRET = "refreshToken secret key" #refresh token secret key

입력
