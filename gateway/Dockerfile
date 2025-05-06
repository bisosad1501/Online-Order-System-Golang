FROM traefik:v2.9

# Create log directory
RUN mkdir -p /var/log/traefik

COPY ./config/traefik.yml /etc/traefik/traefik.yml
COPY ./config/dynamic_conf.yml /etc/traefik/dynamic_conf.yml

EXPOSE 8000
