from flask import Flask, request, Response
from flask_httpauth import HTTPTokenAuth

import aiohttp
import jwt
import config


"""
class AuthManager:
    tokens: dict[str, int] = {}

    def __init__(self):
        pass

    def auth_user(self, token: str, user_id: int):
        self.tokens[token] = user_id

    def get_user(self, token: str):
        return self.tokens.get(token)  # returns None if it doesn't exist
"""


app = Flask(__name__)
conf = config.Config()
app.config.from_object(conf)


class FlaskGateway:
    auth = HTTPTokenAuth(scheme='Bearer')
    client = None

    def __init__(self):
        conf.verify()

    def set_auth_header(self, resp: Response, token: str):
        resp.headers['Authorization'] = 'Bearer ' + token

    def generate_auth_token(self, user_id: int):
        return jwt.encode({'user_id': user_id, 'exp': conf.TOKEN_EXP}, conf.SECRET_KEY, algorithm='HS256')

    async def response_aio2flask(self, resp: aiohttp.ClientResponse):
        # Copy the status
        status = resp.status
        # Copy the headers
        headers = dict(resp.headers)
        # Copy the body
        body = await resp.read()
        # Create a new flask.Response object
        flask_response = Response(response=body, status=status, headers=headers)
        return flask_response

    def generate_401(self, headers=None):
        return Response(status=401, response="NOT_AUTHORIZED", headers=headers)


gateway = FlaskGateway()


@app.before_request
async def post_init():
    app.before_request_funcs[None].remove(post_init)
    gateway.client = aiohttp.ClientSession(conf.USERS_API_BASE_URL)


@app.route('/login', methods=['POST'])
async def login():
    async with gateway.client.post('/login', data=request.data) as resp:
        r = await gateway.response_aio2flask(resp)
        if resp.status == 201:
            resp_json = await resp.json()
            user_id = resp_json["id"]
            token = gateway.generate_auth_token(user_id)

            gateway.set_auth_header(r, token)
        return r


@gateway.auth.login_required
@app.route('/users/<username>')
async def get_user(username: str):
    async with gateway.client.get('/users/' + username, data=request.data) as resp:
        r = await gateway.response_aio2flask(resp)
        resp_json = await resp.json()
        user_id = resp_json["id"]

        # Check if the user id matches
        if not gateway.auth.current_user() == user_id:
            return gateway.generate_401(r.headers)

        return r


@app.route('/register', methods=["POST"])
async def register():
    async with gateway.client.post('/register', data=request.data) as resp:
        r = await gateway.response_aio2flask(resp)
        return r


@gateway.auth.login_required
@app.route('/edit/info', methods=["POST"])
async def edit_info(self):
    user_info = request.json()
    if not self.auth.current_user() == user_info["id"]:
        return self.generate_401()

    async with self.client.post('/edit/info', data=request.data) as resp:
        r = await self.response_aio2flask(resp)
        return r


@gateway.auth.login_required
@app.route('/edit/credentials', methods=["POST"])
async def edit_cred():
    user_info = request.json()
    if not gateway.auth.current_user() == user_info["id"]:
        return gateway.generate_401()

    async with gateway.client.post('/edit/credentials', data=request.data) as resp:
        r = await gateway.response_aio2flask(resp)
        return r

@gateway.auth.verify_token
def verify_token(token):
    try:
        payload = jwt.decode(token, conf.SECRET_KEY, algorithms=['HS256'])

        # Extract user ID from the payload
        return payload.get('user_id')
    except jwt.ExpiredSignatureError or jwt.InvalidTokenError:
        return None
