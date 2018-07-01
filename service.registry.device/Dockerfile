FROM python:3.5-jessie

# Upgrade pip
RUN pip install --upgrade pip

# Add the libraries
COPY ./libraries/python /root/.local/lib/python3.5/site-packages/

# Install requirements
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY ./service.registry.device/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy the application files
COPY ./service.registry.device .

CMD python run.py
