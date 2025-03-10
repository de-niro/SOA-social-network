FROM python:3.11-slim
LABEL authors="de-niro"

WORKDIR /app

COPY . .

RUN pip3 install --upgrade pip && pip install --no-cache-dir -r requirements.txt

EXPOSE 8080

CMD ["python", "main.py"]
