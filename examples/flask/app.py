from flask import Flask, session

import hashlib

app = Flask(__name__)

app.secret_key = "supersecret"

app.config["SESSION_COOKIE_NAME"] = "session"

from flask.sessions import SecureCookieSessionInterface
from itsdangerous import URLSafeTimedSerializer


class SHA256SessionInterface(SecureCookieSessionInterface):
    salt = "11"

    def get_signing_serializer(self, app):
        if not app.secret_key:
            return None

        signer_kwargs = {
            "key_derivation": "hmac",
            "digest_method": hashlib.sha256,
        }

        return URLSafeTimedSerializer(
            secret_key=app.secret_key,
            salt=self.salt,
            serializer=self.serializer,
            signer_kwargs=signer_kwargs,
        )


app.session_interface = SHA256SessionInterface()


@app.route("/")
def index():
    session["user"] = "admin"
    session["role"] = "root"

    return "Cookie SHA256 creada"


if __name__ == "__main__":
    app.run(debug=True)