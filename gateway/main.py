import asyncio
from hypercorn.config import Config as HConfig
from hypercorn.asyncio import serve
from asgiref.wsgi import WsgiToAsgi

from gateway import app, conf


def main():
    print("main() : starting...")
    h_conf = HConfig()
    h_conf.bind = [conf.HOST + ":" + str(conf.PORT)]
    app_asgi = WsgiToAsgi(app)
    asyncio.run(serve(app_asgi, h_conf, mode="asgi"))


if __name__ == "__main__":
    main()
