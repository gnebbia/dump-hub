FROM node AS ng-builder
WORKDIR /opt/dump-hub-fe
COPY . .
RUN npm install
RUN npm install -g @angular/cli
RUN ng build --prod

FROM nginx AS dump-hub-web
WORKDIR /usr/share/nginx/html
COPY --from=ng-builder /opt/dump-hub-fe/dist/dump-hub-fe .
RUN rm 3rdpartylicenses.txt
WORKDIR /etc/nginx
RUN openssl req -new -newkey rsa:4096 -days 365 -nodes -x509 \
    -subj "/C=AU/ST=None/L=None/O=x0e1f/CN=dump-hub.priv" \
    -keyout dump-hub.key  -out dump-hub.crt