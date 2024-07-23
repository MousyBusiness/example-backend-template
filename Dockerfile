# ----verify----
FROM alpine:3.12
ARG SERVICE
RUN test -n "$SERVICE" || (echo "require SERVICE build arg" && exit 1)

# ----builder----
FROM golang:1.19.5 as builder

ARG SERVICE
ARG APP_USER=app
ARG APP_HOME=/home/app

RUN groupadd $APP_USER && useradd -m -g $APP_USER -l $APP_USER
RUN mkdir -p $APP_HOME && chown -R $APP_USER:$APP_USER $APP_HOME
USER $APP_USER

# generate .netrc file for private github access
ARG ACCESS_TOKEN_USR=""
ARG ACCESS_TOKEN_PWD=""
RUN printf "machine github.com\n\
    login ${ACCESS_TOKEN_USR}\n\
    password ${ACCESS_TOKEN_PWD}\n\
    \n\
    machine api.github.com\n\
    login ${ACCESS_TOKEN_USR}\n\
    password ${ACCESS_TOKEN_PWD}\n"\
    >> $APP_HOME/.netrc
RUN chmod 600 $APP_HOME/.netrc

# copy go mod first
COPY --chown=$APP_USER:$APP_USER  go.mod $APP_HOME
COPY --chown=$APP_USER:$APP_USER  go.sum $APP_HOME

# set working directory app home for go mod
WORKDIR $APP_HOME

# get dependencies
RUN go mod download
RUN go mod verify

# copy all local dependencies
COPY --chown=$APP_USER:$APP_USER  . $APP_HOME

# set working directory to specific service
WORKDIR $APP_HOME/cmd/$SERVICE


# build
RUN GOARCH=amd64 GOOS=linux go build -o $SERVICE

## ----runtime----
FROM alpine:3.12

# so CGO_ENABLED=0 isnt required
RUN apk add --no-cache \
    libc6-compat

ARG SERVICE

# set developement stage
ARG STAGE=dev
ENV STAGE $STAGE

# dont use root
ARG APP_USER=app
ARG APP_HOME=/home/app/cmd/$SERVICE

ENV ENTRY="$APP_HOME/$SERVICE"
ENV FIREBASE_CONFIG_FILE="example-$STAGE-firebase-admin-sa.json"

# alpine command slightly different https://stackoverflow.com/questions/49955097/how-do-i-add-a-user-when-im-using-alpine-as-a-base-image
RUN addgroup -S $APP_USER && adduser -S $APP_USER -G $APP_USER

RUN mkdir -p $APP_HOME

WORKDIR $APP_HOME

COPY --chown=$APP_USER:$APP_USER --from=builder $APP_HOME/$SERVICE $APP_HOME
COPY --chown=$APP_USER:$APP_USER $FIREBASE_CONFIG_FILE $APP_HOME

#https://jessicadeen.com/how-to-solve-the-listen-tcp-80-bind-permission-denied-error-in-docker/
RUN apk add libcap && setcap 'cap_net_bind_service=+ep' $ENTRY

EXPOSE 80

USER $APP_USER

CMD "$ENTRY"