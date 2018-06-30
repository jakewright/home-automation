FROM s1lvester/docker-ola

# Override the entrypoint because we want to wrap start.sh 
# in a script that also starts the DMX python service.
ENTRYPOINT []

# Install wget
RUN  apt-get update \
  && apt-get install -y wget \
  && rm -rf /var/lib/apt/lists/*

# Install pip
RUN wget -O get-pip.py 'https://bootstrap.pypa.io/get-pip.py' && python get-pip.py

# Add the libraries
COPY ./libraries/python /root/.local/lib/python2.7/site-packages/

# Install the requirements
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY ./service.controller.dmx/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Add the startup scripts
COPY ./service.controller.dmx/process_manager.sh /
COPY ./service.controller.dmx/wait-for-it.sh /
RUN ["chmod", "+x", "/process_manager.sh"]
RUN ["chmod", "+x", "/wait-for-it.sh"]

# Add the DMX service
COPY ./service.controller.dmx .

CMD ["/process_manager.sh"]