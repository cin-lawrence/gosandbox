# syntax=docker/dockerfile:1.0.0-experimental
FROM python:3.8-alpine

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1


RUN mkdir /worker

RUN pip install --upgrade pip
COPY ./build/wkrce/ /worker/

RUN pip install --no-cache-dir -r /worker/requirements.txt

WORKDIR /worker
