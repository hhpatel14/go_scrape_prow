FROM python:3.9

COPY requirements.txt /
RUN /usr/local/bin/python -m pip install --upgrade pip
RUN pip install -r requirements.txt

COPY *.json *.py *.sh /
ENTRYPOINT ["./entrypoint.sh"]
CMD [""]
