import os


class Config:
    DEBUG = False
    TOKEN_EXP = 3600
    SECRET_KEY = os.environ["SECRET_KEY"]
    HOST = os.environ.get("GATEWAY_HOST", "0.0.0.0")
    PORT = os.environ.get("GATEWAY_PORT", 8080)

    USERS_API_BASE_URL = os.environ["USERS_API_BASE_URL"]

    def verify(self):
        if not self.SECRET_KEY:
            print("Config::verify() : $SECRET_KEY is unset")
            exit(1)
