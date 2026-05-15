from itsdangerous import URLSafeTimedSerializer
import hashlib

s = URLSafeTimedSerializer(
    "supersecret",
    salt="5",
    signer_kwargs={
        "digest_method": hashlib.sha256
    }
)

cookie = s.dumps({"user": "admin"})
print(cookie)