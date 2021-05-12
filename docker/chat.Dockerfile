FROM build

FROM alpine:latest

COPY --from=build /app/chat /app/

WORKDIR /app

CMD ./chat